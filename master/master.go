package master

import (
	"strconv"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
	master_pb "github.com/hiank/thinkend/master/proto"
)

//Handler 处理消息
type Handler struct {
}

//Handle 处理stream消息
func (h *Handler) Handle(msg *pool.Message) error {

	return nil
}

//HandleGet 处理get 消息
func (h *Handler) HandleGet(msg *pb.Message) (res *pb.Message, err error) {

	req := &master_pb.G_Master_Role{}
	if err = ptypes.UnmarshalAny(msg.GetData(), req); err != nil {

		glog.Warningf("unmarshal G_Master_Role error %v\n", err)
		//此处处理收到的数据异常的bug
		return
	}

	var uid uint64
	if uid, err = strconv.ParseUint(msg.GetToken(), 10, 64); err != nil {
		glog.Warningln("parse token to uint64 error : ", err)
		return
	}

	role := &master_pb.Role{

		Uid:     uid,
		ModelId: 1,
		Cup:     1000,
		Uname:   msg.GetToken() + "Tank",
	}

	var anyMsg *any.Any
	if anyMsg, err = ptypes.MarshalAny(role); err != nil {
		glog.Warningln("message to any error : ", err)
		return
	}
	res = &pb.Message{
		Token: msg.GetToken(),
		Data:  anyMsg,
	}
	return
}

//HandlePost 处理post 消息
func (h *Handler) HandlePost(msg *pb.Message) error {

	return nil
}
