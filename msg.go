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

type ClientMsg struct {
	Id        string          `json:"id"`
	Type      ClientMsgType   `json:"t"`
	Ack       int8            `json:"ack,omitempty"`
	ChannelId string          `json:"cid,omitempty"`
	Data      json.RawMessage `json:"d,omitempty"`
}

func (msg *ClientMsg) Validate() error {
	if msg == nil {
		return fmt.Errorf("ClientMsg is null")
	}

	if !msg.Type.checkValid() {
		return fmt.Errorf("unknown ClientMsg type: %v", msg.Type)
	}

	if msg.Type.dataMust() && len(msg.Data) == 0 {
		return fmt.Errorf("ClientMsg type %v must have data", msg.Type)
	}

	switch msg.Type {
	case ClientMsgTypeBiz:
		if msg.ChannelId == "" {
			return errors.New("empty cid for channel biz msg")
		}
	default:
		//
	}

	return nil
}

type InitMsgData struct {
	SessionId       string `json:"session_id,omitempty"`
	ContentEncoding string `json:"content_encoding,omitempty"`
	AcceptEncoding  string `json:"accept_encoding,omitempty"`
	PingInterval    int    `json:"ping_interval"`
	Code            string `json:"code"`
}
