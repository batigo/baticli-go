package baticli

import "C"
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"sync"
	"time"
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

func NewConn(ctx context.Context, conf ConnConfig) (conn *Conn, err error) {
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

type Conn struct {
	inited         bool
	conn           *websocket.Conn
	compressorType CompressorType
	compressor     Compressor
	msgType        int

	lock sync.RWMutex
}

func (c *Conn) Init() (err error) {
	if c.isInited() {
		return
	}

	err = c.sendInitMsg()
	if err != nil {
		return err
	}

	return
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
	err = json.Unmarshal(msg.Data, &initData)
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

func (c *Conn) isInited() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inited
}

func (c *Conn) startRead() {}
