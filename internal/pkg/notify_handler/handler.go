package notify_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/walkerdu/wecom-backend/pkg/wecom"
	"github.com/walkerdu/wecom-notify/pkg/message"
)

var once sync.Once
var handler *Handler

type LogicHandler interface {
	GetHandlerType() message.MessageType
	HandleMessage(message.MessageReq) (message.MessageRsp, error)
}

type UploaderFunc func(wecom.MessageType, string, []byte) (string, error)

// 所有周边来源notify的Handler管理
type Handler struct {
	logicHandlerMap map[message.MessageType]LogicHandler

	// TODO 目前所有LogicHandler复用一个pusher，后面可以根据需要收敛到LogicHandler内部
	pusher func(string, string) error

	uploader UploaderFunc
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

func (h *Handler) RegisterUploader(uploader UploaderFunc) {
	h.uploader = uploader
}

func (h *Handler) GetUploader() UploaderFunc {
	return h.uploader
}

// ServeHTTP 实现http.Handler接口
func (h *Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Printf("[DEBUG]ServeHttp|recv request URL:%s, Method:%s", req.URL, req.Method)

	if req.Method == http.MethodGet {
		err := errors.New("do not support HTTP GET Method")
		log.Printf("[WARN]ServeHttp|%s", err)
		http.Error(wr, err.Error(), http.StatusBadRequest)

		return
	} else if req.Method == http.MethodPost {
		var msg message.MessageReq

		contentType := req.Header.Get("Content-Type")
		if contentType == "application/json" {
			// 1.http请求体body的content为json格式
			// 读取HTTP请求体
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				err = fmt.Errorf("Failed to read request Body:%s", err)
				log.Printf("[DEBUG]%s", err)

				http.Error(wr, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Printf("[DEBUG]ServeHTTP|recv request Body:%s", body)

			if err := json.Unmarshal(body, &msg); err != nil {
				err = fmt.Errorf("Failed to unmarshal request Body, err:%s", err)
				log.Printf("[ERROR]ServeHTTP|%s", err)

				http.Error(wr, err.Error(), http.StatusBadRequest)
				return
			}
		} else if contentType == "multipart/form-data" {
			// 2.http请求体body的content为form表单数据

			// 解析 multipart/form-data 请求体
			// 最大1MB
			err := req.ParseMultipartForm(1 * 1024 * 1024)
			if err != nil {
				log.Printf("[ERROR]ServeHTTP|ParseMultipartForm failed, err=%s", err)
				http.Error(wr, err.Error(), http.StatusBadRequest)
				return
			}

			// 获取上传的文件
			mpFile, mpFileHeader, err := req.FormFile("file")
			if err != nil {
				log.Printf("[ERROR]ServeHTTP|FormFile failed, err=%s", err)
				http.Error(wr, err.Error(), http.StatusBadRequest)
				return
			}
			defer mpFile.Close()

			body, err := ioutil.ReadAll(mpFile)
			if err != nil {
				log.Printf("[ERROR]ServeHTTP|ReadAll failed, err=%s", err)
				http.Error(wr, err.Error(), http.StatusInternalServerError)
				return
			}

			msg.MsgType = message.MessageType_MDBlog
			msg.Content = string(body)
			msg.ContentAttr = mpFileHeader.Filename
		} else {
			err := fmt.Errorf("HTTP POST Method: unkown content-type:%s", contentType)
			log.Printf("[WARN]ServeHttp|%s", err)
			http.Error(wr, err.Error(), http.StatusBadRequest)

			return
		}

		logicHandler, exist := h.logicHandlerMap[msg.MsgType]
		if !exist {
			err := fmt.Errorf("unsuport request message type=%s", msg.MsgType)
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
