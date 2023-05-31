package notify_handler

import (
	"log"
	"net/http"

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

	userID := "walkerdu"
	pusher := HandlerInst().GetPusher()
	if err := pusher(userID, msg.Content); err != nil {
		log.Printf("[ERROR]HandleMessage|push to user=%s failed, err=%s", userID, err)
		textMsgRsp.RetCode = http.StatusInternalServerError
		textMsgRsp.RetMsg = err.Error()
	}

	return textMsgRsp, nil
}
