package baticli

import (
	"errors"
	"fmt"
)

func (m ClientMsgType) checkValid() bool {
	return m == ClientMsgType_Biz ||
		m == ClientMsgType_Ack ||
		m == ClientMsgType_Echo ||
		m == ClientMsgType_Init ||
		m == ClientMsgType_InitResp
}

func (m ClientMsgType) str() string {
	switch m {
	case ClientMsgType_Init:
		return "init"
	case ClientMsgType_InitResp:
		return "init_resp"
	case ClientMsgType_Biz:
		return "biz"
	case ClientMsgType_Ack:
		return "ack"
	case ClientMsgType_Echo:
		return "echo"
	}

	return "unknown"
}

func (msg *ClientMsg) Validate() error {
	if msg == nil {
		return fmt.Errorf("ClientMsgRecv is null")
	}

	if !msg.Type.checkValid() {
		return fmt.Errorf("unknown ClientMsgRecv type: %v", msg.Type)
	}

	switch msg.Type {
	case ClientMsgType_Biz:
		if msg.GetServiceId() == "" {
			return errors.New("service-id missing for biz msg")
		}
		if msg.GetBizData() == nil {
			return errors.New("biz-data missing for biz msg")
		}
	case ClientMsgType_Init, ClientMsgType_InitResp:
		if msg.GetInitData() == nil {
			return errors.New("init-data missing for init msg")
		}
	default:
		//
	}

	return nil
}
