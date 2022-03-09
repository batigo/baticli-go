package baticli

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ClientMsgType int8

const (
	ClientMsgTypeInit     ClientMsgType = 1
	ClientMsgTypeInitResp ClientMsgType = 2
	ClientMsgTypeBiz      ClientMsgType = 3
	ClientMsgTypeAck      ClientMsgType = 4
	ClientMsgTypeEcho     ClientMsgType = 100
)

func (m ClientMsgType) dataMust() bool {
	return m == ClientMsgTypeInit
}

func (m ClientMsgType) checkValid() bool {
	return m == ClientMsgTypeBiz ||
		m == ClientMsgTypeAck ||
		m == ClientMsgTypeEcho ||
		m == ClientMsgTypeInit ||
		m == ClientMsgTypeInitResp
}

func (m ClientMsgType) str() string {
	switch m {
	case ClientMsgTypeInit:
		return "init"
	case ClientMsgTypeInitResp:
		return "init_resp"
	case ClientMsgTypeBiz:
		return "biz"
	case ClientMsgTypeAck:
		return "ack"
	case ClientMsgTypeEcho:
		return "echo"
	}

	return "unknown"
}

type ClientMsgSend struct {
	Id        string        `json:"id"`
	Type      ClientMsgType `json:"t"`
	Ack       int8          `json:"ack,omitempty"`
	ServiceId string        `json:"sid,omitempty"`
	Data      interface{}   `json:"d,omitempty"`
}

type ClientMsgRecv struct {
	Id        string          `json:"id"`
	Type      ClientMsgType   `json:"t"`
	Ack       int8            `json:"ack,omitempty"`
	ServiceId string          `json:"sid,omitempty"`
	Data      json.RawMessage `json:"d,omitempty"`
}

func (msg *ClientMsgRecv) decode(bs []byte) error {
	return json.Unmarshal(bs, msg)
}

func (msg *ClientMsgRecv) Validate() error {
	if msg == nil {
		return fmt.Errorf("ClientMsgRecv is null")
	}

	if !msg.Type.checkValid() {
		return fmt.Errorf("unknown ClientMsgRecv type: %v", msg.Type)
	}

	if msg.Type.dataMust() && len(msg.Data) == 0 {
		return fmt.Errorf("ClientMsgRecv type %v must have data", msg.Type)
	}

	switch msg.Type {
	case ClientMsgTypeBiz:
		if msg.ServiceId == "" {
			return errors.New("empty cid for channel biz msg")
		}
	default:
		//
	}

	return nil
}

type InitMsgData struct {
	AcceptEncoding CompressorType `json:"accept_encoding,omitempty"`
	PingInterval   int            `json:"ping_interval,omitempty"`
}

func (d *InitMsgData) decode(bs []byte) error {
	return json.Unmarshal(bs, d)
}
