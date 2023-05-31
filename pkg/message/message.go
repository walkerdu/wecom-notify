package message

// 推送消息类型
type MessageType string

const (
	MessageType_MDBlog MessageType = "md_blog" // blog文档信息，直接推送
)

type MessageReq struct {
	MsgType MessageType `json:"msg_type"`
	Content string      `json:"content"`
}

type MessageRsp struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
}
