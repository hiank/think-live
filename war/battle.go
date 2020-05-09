package war

import (
	"container/list"
	"context"

	"github.com/golang/glog"
	war_pb "github.com/hiank/thinkend/war/proto"
)

//ContextKey for context.WithValue
type ContextKey int

//CtxKeyIDecode key for context value [IDecode]
var CtxKeyIDecode = ContextKey(0)

//FidMaker make fid for Fight
type FidMaker struct {
	max  uint32
	free *list.List
}

//newFidMaker new FidMaker
func newFidMaker() *FidMaker {

	return &FidMaker{
		free: list.New(),
	}
}

//Make make a new fid
func (fm *FidMaker) Make() uint32 {

	if fm.free.Len() > 0 {
		return fm.free.Remove(fm.free.Front()).(uint32)
	}
	fm.max++
	return fm.max - 1
}

//Free free a fid
func (fm *FidMaker) Free(fid uint32) {

	fm.free.PushBack(fid)
}

//Battle 战役，包含一组战斗
type Battle struct {
	matcher *Matcher
	doReq   chan *war_pb.S_War_Do
}

//NewBattle 创建一场战役
func NewBattle(ctx context.Context, t war_pb.War_Type) *Battle {

	ctx = context.WithValue(ctx, CtxKeyIDecode, EncodeWarType(0, t))
	b := &Battle{
		matcher: NewMatcher(ctx),
		doReq:   make(chan *war_pb.S_War_Do),
	}
	go b.loop(ctx)
	return b
}

func (b *Battle) loop(ctx context.Context) {

	fightHub, fidMaker, freeReq := make(map[uint32]*Fight), newFidMaker(), make(chan uint32)
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case fight := <-b.matcher.MatchedRes(): //NOTE: 匹配成功后，处理
			fid := fidMaker.Make()
			fightHub[fid] = fight
			go fight.Start(context.WithValue(ctx, CtxKeyIDecode, EncodeFightID(ctx.Value(CtxKeyIDecode).(IDecode), fid)), freeReq)
		case fid := <-freeReq:
			fidMaker.Free(fid)
			delete(fightHub, fid)
		case req := <-b.doReq:
			b.optDo(req, fightHub)
		}
	}
}

//JoinReq 获得加入战斗请求chan
func (b *Battle) JoinReq() chan<- *Gamer {

	return b.matcher.JoinReq()
}

//DoReq 获得执行操作请求chan
func (b *Battle) DoReq() chan<- *war_pb.S_War_Do {

	return b.doReq
}

//optDo 处理收到的消息
func (b *Battle) optDo(warDo *war_pb.S_War_Do, fightHub map[uint32]*Fight) {

	fid := IDecode(warDo.GetId()).FightID()
	if fight, ok := fightHub[fid]; ok {
		fight.DoReq() <- warDo
	} else {
		glog.Warningf("can't find fight %d\n", fid)
	}
}
