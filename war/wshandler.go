package war

import (
	"github.com/hiank/think/pb"
	"context"
)

//WSHandler 处理k8s 请求
type WSHandler struct {

	handler *Handler
	// w 			*War 			//NOTE: 战争
	// ctx 		context.Context			//NOTE:
	// Close 		context.CancelFunc		//NOTE:
}

//NewWSHandler 新建一个Handler
func NewWSHandler(ctx context.Context) *WSHandler {

	return &WSHandler {
		handler : NewHandler(ctx),
	}
}

//Handle 处理读到的消息
func (h *WSHandler) Handle(msg *pb.Message) (err error) {

	msg.Key = "ws"
	return h.handler.Handle(msg)
}


func (h *WSHandler) Close() {

	h.handler.Close()
}
