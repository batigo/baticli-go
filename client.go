package baticli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type DeviceType uint8

func (t DeviceType) validate() bool {
	return t == DeviceTypeAndroid || t == DeviceTypeIos || t == DeviceTypeWeb || t == DeviceTypeOthers
}

const (
	DeviceTypeAndroid DeviceType = 0
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
	conn.msgSendChan = make(chan *ClientMsg, 32)
	conn.msgRecvChan = make(chan *ClientMsg, 32)
	conn.stopChan = make(chan interface{})
	sendFunc = func(msg *ClientMsg) {
		conn.msgSendChan <- msg
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
	if !conf.Dt.validate() {
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
	errMsgTypeAbnormal = errors.New("msg type abnormal")
	errMsgInvalid      = errors.New("invalid msg")
	errMsgDecodeFail   = errors.New("msg decode fail")
	errClientRecv      = errors.New("failed to recv msg")
)

type SendMsgFunc func(msg *ClientMsg)
type RecvMsgHandler func(msg *ClientMsg)
type ConnCloseHandler func()

type Conn struct {
	conn             *websocket.Conn
	compressorType   CompressorType
	compressor       Compressorr
	msgType          int
	msgSendChan      chan *ClientMsg
	msgRecvChan      chan *ClientMsg
	stopChan         chan interface{}
	msgHandler       RecvMsgHandler
	connClosehandler ConnCloseHandler
	lock             sync.RWMutex
	hbInterval       time.Duration
}

func (c *Conn) Init() (err error) {
	if c.connClosehandler == nil {
		err = fmt.Errorf("conn close handler required")
		return
	}

	if c.msgHandler == nil {
		err = fmt.Errorf("recv msg handler required")
		return
	}

	err = c.sendInitMsg()
	if err != nil {
		return
	}

	data, err := c.waitInitResp()
	if err != nil {
		return
	}

	c.compressorType = data.AcceptCompressor
	c.compressor = newCompressor(data.AcceptCompressor)
	c.hbInterval = time.Second * time.Duration(data.PingInterval)
	c.start()
	return
}

func (c *Conn) Close() {
	close(c.stopChan)
	close(c.msgSendChan)
}

func (c *Conn) SetRecvMsgHandler(handler RecvMsgHandler) {
	c.msgHandler = handler
}

func (c *Conn) SetConnCloseHanler(handler ConnCloseHandler) {
	c.connClosehandler = handler
}

func (c *Conn) sendInitMsg() (err error) {
	data := InitData{}
	if c.compressorType != CompressorType_Null {
		data.AcceptCompressor = c.compressorType
	}
	msg := ClientMsg{
		Id:       Genmsgid(),
		Type:     ClientMsgType_Init,
		Ack:      0,
		InitData: &data,
	}
	bs, err := proto.Marshal(&msg)
	if err != nil {
		return
	}
	return c.conn.WriteMessage(c.msgType, bs)
}

func (c *Conn) waitInitResp() (data *InitData, err error) {
	msg, err := c.recvMsg()
	if err != nil {
		return
	}

	if msg.Type != ClientMsgType_InitResp {
		err = errMsgTypeAbnormal
		return
	}

	data = msg.GetInitData()
	return
}

func (c *Conn) recvMsg() (msg ClientMsg, err error) {
	_, bs, err := c.conn.ReadMessage()
	if err != nil {
		log.Printf("failed to recv msg: %s", err.Error())
		err = errClientRecv
		return
	}

	bs, err = c.compressor.Uncompress(bs)
	if err != nil {
		log.Printf("failed to uncomress msg: %s", err.Error())
		return
	}

	log.Printf("recv msg: %s", bs)
	err = proto.Unmarshal(bs, &msg)
	if err != nil {
		log.Printf("failed to decode msg: %s - %s", bs, err.Error())
		err = errMsgDecodeFail
		return
	}

	err = msg.Validate()
	if err != nil {
		log.Printf("recv invalid msg: %s - %s", msg.Id, err.Error())
		err = errMsgInvalid
		return
	}

	return
}

func (c *Conn) start() {
	go func() {
		for {
			msg, err := c.recvMsg()
			if err == errClientRecv {
				return
			}
			if err != nil {
				continue
			}

			select {
			case <-c.stopChan:
				return
			case c.msgRecvChan <- &msg:
				continue
			}

		}
	}()

	go func() {
		defer close(c.stopChan)
		for {
			ticker := time.NewTicker(c.hbInterval).C
			select {
			case <-c.stopChan:
				return
			case msg := <-c.msgRecvChan:
				c.msgHandler(msg)
			case <-ticker:
				err := c.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5))
				if err != nil {
					break
				}
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
			case msg := <-c.msgSendChan:
				bs, err := proto.Marshal(msg)
				bs, err = c.compressor.Compress(bs)
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
