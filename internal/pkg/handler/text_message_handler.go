package handler

import (
	"github.com/walkerdu/wecom-backend/pkg/wecom"
)

const WeChatTimeOutSecs = 5

func init() {
	handler := &TextMessageHandler{}

	HandlerInst().RegisterLogicHandler(wecom.MessageTypeText, handler)
}

type TextMessageHandler struct {
}

func (t *TextMessageHandler) GetHandlerType() wecom.MessageType {
	return wecom.MessageTypeText
}

func (t *TextMessageHandler) HandleMessage(msg wecom.MessageIF) (wecom.MessageIF, error) {
	//textMsg := msg.(*wecom.TextMessageReq)

	textMsgRsp := wecom.TextMessageRsp{
		Content: "臣妾来了，有何吩咐!",
	}

	return &textMsgRsp, nil
}
