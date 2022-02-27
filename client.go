package baticli

import (
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

func NewConn(conf ConnConfig) (conn *Conn, err error) {
	err = conf.validate()
	if err != nil {
		return
	}

	params := url.Values{}
	params.Add("uid", conf.Uid)
	params.Add("did", conf.Did)
	params.Add("dt", fmt.Sprintf("%v", conf.Dt))
	urlStr := fmt.Sprintf("%s?%s", conf.Url, params.Encode())

	_conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		return
	}

	conn = &Conn{
		conn:       _conn,
		compressor: newCompressor(conf.Compressor),
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
}

func (conf *ConnConfig) validate() error {
	if conf.Dt.validate() {
		return fmt.Errorf("unknown device-type: %v", conf.Dt)
	}
	if conf.Did == "" || len(conf.Did) > 64 {
		return fmt.Errorf("device-id empty or too long(max length=64): %s")
	}
	if conf.Uid == "" || len(conf.Uid) > 64 {
		return fmt.Errorf("user-id empty or too long(max length=64): %s")
	}

	return nil
}

type Conn struct {
	inited     bool
	conn       *websocket.Conn
	compressor Compressor

	lock sync.RWMutex
}

func (c *Conn) Init() (err error) {
	if c.isInited() {
		return
	}

	return
}

func (c *Conn) isInited() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inited
}
