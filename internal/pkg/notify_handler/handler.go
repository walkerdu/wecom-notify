package notify_handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/walkerdu/wecom-notify/pkg/message"
)

var once sync.Once
var handler *Handler

type LogicHandler interface {
	GetHandlerType() message.MessageType
	HandleMessage(message.MessageReq) (message.MessageRsp, error)
}

// 所有周边来源notify的Handler管理
type Handler struct {
	logicHandlerMap map[message.MessageType]LogicHandler

	// 目前所有LogicHandler复用一个pusher，后面可以根据需要收敛到LogicHandler内部
	pusher func(string, string) error
}

// NewHandler 返回一个新的Handler实例
func HandlerInst() *Handler {
	once.Do(func() {
		handler = &Handler{
			logicHandlerMap: make(map[message.MessageType]LogicHandler),
		}
	})

	return handler
}

func (h *Handler) RegisterLogicHandler(msgType message.MessageType, logicHandler LogicHandler) {
	h.logicHandlerMap[msgType] = logicHandler
}

func (h *Handler) GetLogicHandlerMap() map[message.MessageType]LogicHandler {
	return h.logicHandlerMap
}

func (h *Handler) RegisterPusher(pusher func(string, string) error) {
	h.pusher = pusher
}

func (h *Handler) GetPusher() func(string, string) error {
	return h.pusher
}

// ServeHTTP 实现http.Handler接口
func (h *Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Printf("[DEBUG]ServeHttp|recv request URL:%s, Method:%s", req.URL, req.Method)

	if req.Method == http.MethodGet {
		log.Printf("[WARN]ServeHttp|do not support http get Method")
	} else if req.Method == http.MethodPost {
		// 读取HTTP请求体
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			err = fmt.Errorf("Failed to read request Body:%s", err)
			log.Printf("[DEBUG]%s", err)

			http.Error(wr, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[DEBUG]ServeHTTP|recv request Body:%s", body)

		var msg message.MessageReq
		if err := json.Unmarshal(body, &msg); err != nil {
			err = fmt.Errorf("Failed to unmarshal request Body, err:%s", err)
			log.Printf("[ERROR]ServeHTTP=%s", err)

			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}

		logicHandler, exist := h.logicHandlerMap[msg.MsgType]
		if !exist {
			err = fmt.Errorf("unsuport request message type=%s", msg.MsgType)
			log.Printf("[ERROR]ServeHTTP=%s", err)

			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}

		logicRsp, err := logicHandler.HandleMessage(msg)
		if err != nil {
			log.Printf("[ERROR]ServeHTTP=%s", err)
			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}

		logicRspBytes, err := json.Marshal(logicRsp)
		if err != nil {
			log.Printf("[ERROR]ServeHTTP|json Marshal failed, err:%s", err)
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(wr, string(logicRspBytes))
	}
}
