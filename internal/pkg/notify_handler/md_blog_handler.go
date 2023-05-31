package notify_handler

import (
	"github.com/walkerdu/wecom-notify/pkg/message"
)

func init() {
	handler := &MDBlogMessageHandler{}

	HandlerInst().RegisterLogicHandler(message.MessageType_MDBlog, handler)
}

type MDBlogMessageHandler struct {
}

func (t *MDBlogMessageHandler) GetHandlerType() message.MessageType {
	return message.MessageType_MDBlog
}

func (t *MDBlogMessageHandler) HandleMessage(msg message.MessageReq) (message.MessageRsp, error) {
	textMsgRsp := message.MessageRsp{
		RetCode: 200,
		RetMsg:  "success",
	}

	return textMsgRsp, nil
}
