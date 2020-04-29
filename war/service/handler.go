package main

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/token"
	master_pb "github.com/hiank/thinkend/master/proto"
	"github.com/hiank/thinkend/war"
	war_pb "github.com/hiank/thinkend/war/proto"
)

//Handler 处理k8s 请求
type Handler struct {
	k8s.IgnoreGet  //NOTE: 忽略Get 方法处理
	k8s.IgnorePost //NOTE: 忽略Post 方法处理

	w *war.War //NOTE: 战争
}

//NewHandler 新建一个Handler
func NewHandler(ctx context.Context) *Handler {

	return &Handler{
		w: war.NewWar(ctx),
	}
}

//Handle 处理读到的消息
func (h *Handler) Handle(msg *pool.Message) (err error) {

	var name string
	if name, err = ptypes.AnyMessageName(msg.GetData()); err != nil {
		glog.Warningln("get message name error : ", err)
		return
	}
	name = name[strings.LastIndexByte(name, '_')+1:]
	glog.Infoln("message name : ", name)

	switch name {
	case "Want": //"S_War_Want":
		err = h.handleWant(msg)
	case "Do": //"S_War_Do":
		err = h.handleDo(msg)
	default:
		err = errors.New("undefined operate for message named " + name)
		glog.Warningln(err)
	}
	return
}

//handleWant operate type 'Want' message
func (h *Handler) handleWant(msg *pool.Message) error {

	uid, err := strconv.ParseUint(msg.ToString(), 10, 64)
	if err != nil {
		glog.Warningln("parse token to uint64 error : ", err)
		return err
	}

	role := &master_pb.Role{
		Uid:     uid,
		ModelId: 1,
		Cup:     1000,
		Uname:   msg.ToString() + "Tank",
	}

	var want war_pb.S_War_Want
	if err = ptypes.UnmarshalAny(msg.GetData(), &want); err != nil {
		glog.Warningln("unmarshal want message error : ", err)
		return err
	}
	tok, _ := token.GetBuilder().Get(msg.GetToken())
	h.w.Join(want.GetType(), war.NewGamer(tok, role))
	return nil
}

//handleDo operate type 'Do' message
func (h *Handler) handleDo(msg *pool.Message) (err error) {

	var doMsg war_pb.S_War_Do
	if err = ptypes.UnmarshalAny(msg.GetData(), &doMsg); err != nil {
		//NOTE: 数据解析出错
		glog.Warningln("Do : ", err)
		return
	}
	h.w.Do(&doMsg)
	return
}
