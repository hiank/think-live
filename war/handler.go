package war

import (
	"strconv"
	"errors"
	"github.com/golang/glog"
	"strings"
	"github.com/golang/protobuf/ptypes"
	"github.com/hiank/think/pb"
	"context"
	"github.com/hiank/think/net/k8s"
	"github.com/hiank/thinkend/war/proto"
	master_pb "github.com/hiank/thinkend/master/proto"
)

//Handler 处理k8s 请求
type Handler struct {

	k8s.IgnoreGet				//NOTE: 忽略Get 方法处理
	k8s.IgnorePost 				//NOTE: 忽略Post 方法处理

	w 			*War 			//NOTE: 战争
	ctx 		context.Context			//NOTE:
	Close 		context.CancelFunc		//NOTE:
}

//NewHandler 新建一个Handler
func NewHandler(ctx context.Context) *Handler {

	ctx, cancel := context.WithCancel(ctx)
	return &Handler {
		w 		: NewWar(ctx),
		ctx 	: ctx,
		Close 	: cancel,
	}
}

//Handle 处理读到的消息
func (h *Handler) Handle(msg *pb.Message) (err error) {

	var name string
	if name, err = ptypes.AnyMessageName(msg.GetData()); err != nil {

		glog.Warningln("get message name error : ", err)
		return
	}
	name = name[strings.LastIndexByte(name, '_') + 1:]
	glog.Infoln("message name : ", name)
	glog.Infoln("handle key : ", msg.GetKey())

	switch name {

	case "Want":		//"S_War_Want":
		h.handleWant(msg)
	case "Do": 		//"S_War_Do":
		doMsg := &war_pb.S_War_Do{}
		if err = ptypes.UnmarshalAny(msg.GetData(), doMsg); err != nil {
			//NOTE: 数据解析出错
			glog.Warningln("Do : ", err)
			return 
		}
		h.w.Do(doMsg)
	default:
		err = errors.New("undefined operate for message named " + name)
		glog.Warningln(err.Error())
	}
	return
}

func (h *Handler) handleWant(msg *pb.Message) {

	var uid uint64
	var err error
	if uid, err = strconv.ParseUint(msg.GetToken(), 10, 64); err != nil {
		glog.Warningln("parse token to uint64 error : ", err)
		return
	}
	
	role := &master_pb.Role {

		Uid : uid,
		ModelId : 1,
		Cup : 1000,
		Uname : msg.GetToken() + "Tank",
	}

	// role := &master_pb.Role{}
	// if err = ptypes.UnmarshalAny(r.GetData(), role); err != nil {

	// 	glog.Warningf("unmarshal role error %v\n", err)
	// 	//此处处理收到的数据异常的bug
	// 	return
	// }

	var want war_pb.S_War_Want
	if err = ptypes.UnmarshalAny(msg.GetData(), &want); err != nil {

		glog.Warningln("unmarshal want message error : ", err)
		return
	}

	gamer := &Gamer{msg.GetKey(), msg.GetToken(), role}
	j := &join{gamer: gamer}
	h.w.Join(want.GetType(), j)
}