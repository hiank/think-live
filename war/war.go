package war

import (
	"context"

	war_pb "github.com/hiank/thinkend/war/proto"
)

type warReq struct {
	warType war_pb.War_Type
	data    interface{}
}

//War 战争，根据War_Type 保存各种Battle
type War struct {
	joinReq chan *warReq
	doReq   chan *warReq
}

//NewWar 创建一个War
func NewWar(ctx context.Context) *War {

	w := &War{
		joinReq: make(chan *warReq),
		doReq:   make(chan *warReq),
	}
	go w.loop(ctx)
	return w
}

func (w *War) loop(ctx context.Context) {

	dic := make(map[war_pb.War_Type]*Battle)
	getBattle := func(warType war_pb.War_Type) *Battle {
		battle, ok := dic[warType]
		if !ok {
			battle = NewBattle(ctx, warType)
			dic[warType] = battle
		}
		return battle
	}

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case joinReq := <-w.joinReq:
			getBattle(joinReq.warType).Join(joinReq.data.(*Gamer))
		case doReq := <-w.doReq:
			getBattle(doReq.warType).Do(doReq.data.(*war_pb.S_War_Do))
		}
	}
}

//Join 加入战斗，排队
func (w *War) Join(t war_pb.War_Type, gamer *Gamer) {

	w.joinReq <- &warReq{warType: t, data: gamer}
}

//Do 处理操作命令
func (w *War) Do(d *war_pb.S_War_Do) {

	w.doReq <- &warReq{warType: IDecode(d.GetId()).WarType(), data: d}
}
