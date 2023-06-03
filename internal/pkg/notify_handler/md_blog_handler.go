package notify_handler

import (
	"log"
	"net/http"

	"github.com/walkerdu/wecom-backend/pkg/wecom"
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

	// 上传
	uploader := HandlerInst().GetUploader()
	mediaID, err := uploader(wecom.MessageTypeFile, msg.ContentAttr, []byte(msg.Content))
	if err != nil {
		log.Printf("[ERROR]HandleMessage|uploader media %s failed, err=%s", msg.ContentAttr, err)
		textMsgRsp.RetCode = http.StatusInternalServerError
		textMsgRsp.RetMsg = err.Error()
		return textMsgRsp, nil
	}

	userID := "walkerdu"
	pusher := HandlerInst().GetPusher()
	if err := pusher(userID, mediaID); err != nil {
		log.Printf("[ERROR]HandleMessage|push to user=%s failed, err=%s", userID, err)
		textMsgRsp.RetCode = http.StatusInternalServerError
		textMsgRsp.RetMsg = err.Error()
	}

	return textMsgRsp, nil
}
