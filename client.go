package baticli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type DeviceType uint8

func (t DeviceType) validate() bool {
	return t == DeviceTypeAndriod || t == DeviceTypeIos || t == DeviceTypeWeb || t == DeviceTypeOthers
}

const (
	DeviceTypeAndriod DeviceType = 0
	DeviceTypeIos     DeviceType = 1
	DeviceTypeWeb     DeviceType = 2
	DeviceTypeOthers  DeviceType = 3
)

func NewConn(ctx context.Context, conf ConnConfig) (conn *Conn, sendFunc SendMsgFunc, err error) {
	err = conf.validate()
	if err != nil {
		return
	}

	params := url.Values{}
	params.Add("uid", conf.Uid)
	params.Add("did", conf.Did)
	params.Add("dt", fmt.Sprintf("%v", conf.Dt))
	urlStr := fmt.Sprintf("%s?%s", conf.Url, params.Encode())

	_conn, _, err := websocket.DefaultDialer.DialContext(ctx, urlStr, nil)
	if err != nil {
		return
	}

	conn = &Conn{
		conn:           _conn,
		compressor:     NullCompressor{},
		compressorType: conf.Compressor,
		msgType:        websocket.TextMessage,
	}
	conn.msgsendChan = make(chan ClientMsgSend, 32)
	conn.msgrecvChan = make(chan ClientMsgRecv, 32)
	conn.stopChan = make(chan interface{})
	sendFunc = func(msg ClientMsgSend) {
		conn.msgsendChan <- msg
	}
	conn.conn.SetReadLimit(1024 * 1024)
	if conf.BinaryMsg {
		conn.msgType = websocket.BinaryMessage
	}
	return
}

type ConnConfig struct {
	id         string
	Url        string
	Uid        string
	Did        string
	Dt         DeviceType
	Timeout    time.Duration
	HeartBeat  time.Duration
	Compressor CompressorType
	BinaryMsg  bool
}

func (conf *ConnConfig) validate() error {
	if conf.Dt.validate() {
		return fmt.Errorf("unknown device-type: %v", conf.Dt)
	}
	if conf.Did == "" || len(conf.Did) > 64 {
		return fmt.Errorf("device-id empty or too long(max length=64): %s", conf.Did)
	}
	if conf.Uid == "" || len(conf.Uid) > 64 {
		return fmt.Errorf("user-id empty or too long(max length=64): %s", conf.Uid)
	}

	return nil
}

var (
	errMsgTypeAbnormal   = errors.New("msg type abnormal")
	errMsgDecodeFail     = errors.New("msg decode fail")
	errMsgDataDecodeFail = errors.New("msg data decode fail")
)

type SendMsgFunc func(msg ClientMsgSend)
type RecvMsgHandler func(msg ClientMsgRecv)
type ConnCloseHandler func()

type Conn struct {
	conn             *websocket.Conn
	compressorType   CompressorType
	compressor       Compressor
	msgType          int
	msgsendChan      chan ClientMsgSend
	msgrecvChan      chan ClientMsgRecv
	stopChan         chan interface{}
	msgHandler       RecvMsgHandler
	connClosehandler ConnCloseHandler
	lock             sync.RWMutex
}

func (c *Conn) Init() (sendFunc SendMsgFunc, err error) {
	if c.connClosehandler == nil {
		return nil, fmt.Errorf("conn close handler required")
	}

	if c.msgHandler == nil {
		return nil, fmt.Errorf("recv msg handler required")
	}

	err = c.sendInitMsg()
	if err != nil {
		return
	}

	data, err := c.waitInitResp()
	if err != nil {
		return
	}

	c.compressor = newCompressor(data.AcceptEncoding)
	c.start()
	return
}

func (c *Conn) Close() {
	close(c.stopChan)
	close(c.msgsendChan)
}

func (c *Conn) SetRecvMsgHandler(handler RecvMsgHandler) {
	c.msgHandler = handler
}

func (c *Conn) SetConnCloseHanler(handler ConnCloseHandler) {
	c.connClosehandler = handler
}

func (c *Conn) sendInitMsg() (err error) {
	data := InitMsgData{}
	if c.compressorType != CompressorTypeNull {
		data.AcceptEncoding = c.compressorType
	}
	msg := ClientMsgSend{
		Id:   Genmsgid(),
		Type: ClientMsgTypeInit,
		Ack:  1,
		Data: data,
	}
	bs, _ := json.Marshal(msg)
	return c.conn.WriteMessage(c.msgType, bs)
}

func (c *Conn) waitInitResp() (data InitMsgData, err error) {
	msg, err := c.recvMsg()
	if err != nil {
		return
	}

	if msg.Type != ClientMsgTypeInitResp {
		err = errMsgTypeAbnormal
		return
	}

	var initData InitMsgData
	err = initData.decode(msg.Data)
	if err != nil {
		err = errMsgDataDecodeFail
		return
	}

	return
}

func (c *Conn) recvMsg() (msg ClientMsgRecv, err error) {
	_, bs, err := c.conn.ReadMessage()
	if err != nil {
		return
	}

	err = msg.decode(bs)
	if err != nil {
		err = errMsgDecodeFail
		return
	}
	return
}

func (c *Conn) start() {
	go func() {
		defer close(c.msgrecvChan)
		for {
			_, bs, err := c.conn.ReadMessage()
			if err != nil {
				return
			}
			var msg ClientMsgRecv
			err = msg.decode(bs)
			if err != nil {
				continue
			}
			select {
			case <-c.stopChan:
				return
			case c.msgrecvChan <- msg:
				continue
			}

		}
	}()

	go func() {
		for {
			select {
			case <-c.stopChan:
				return
			case msg := <-c.msgrecvChan:
				c.msgHandler(msg)
			}
		}
	}()

	go func() {
		defer func() {
			c.conn.Close()
			c.connClosehandler()
		}()
		for {
			select {
			case <-c.stopChan:
				return
			case msg := <-c.msgsendChan:
				bs, err := msg.encode()
				if err != nil {
					continue
				}
				err = c.conn.WriteMessage(c.msgType, bs)
				if err != nil {
					return
				}
			}
		}
	}()
}
