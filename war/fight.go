package war

import (
	"context"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	war_pb "github.com/hiank/thinkend/war/proto"
)

//Fight 战斗
type Fight struct {
	doReq   chan *war_pb.S_War_Do
	mapID   int32             //NOTE: 地图id, 用于标识不同的障碍物
	cmdReq  chan *war_pb.Tick //NOTE: 用于驱动服务端战斗计算
	teamHub []*Team
}

//NewFight new a fight
//idecode content WarType FightID
func NewFight(idecode IDecode) *Fight {

	var teamCnt int
	switch idecode.WarType() {

	case war_pb.War_Type_SOCCER: //NOTE: 足球模式
		teamCnt = 2
	case war_pb.War_Type_ONEWINNER: //NOTE: 吃鸡模式
		teamCnt = 12
	case war_pb.War_Type_HUNTING: //NOTE: 猎杀模式[统计标记]
		teamCnt = 2
	case war_pb.War_Type_TEAMWINNER: //NOTE: 3v3 战斗
		teamCnt = 2
	default:
		glog.Warningln("战斗类型错误") //NOTE: 未识别的类型
	}

	return &Fight{
		teamHub: make([]*Team, 0, teamCnt),
		cmdReq:  make(chan *war_pb.Tick),
		doReq:   make(chan *war_pb.S_War_Do),
	}
}

//Start 启动战斗
//complete 战斗结束后，将fid 发送到此chan
func (f *Fight) Start(ctx context.Context, complete chan<- uint32) {

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		complete <- ctx.Value(CtxKeyIDecode).(IDecode).FightID()
		cancel()
	}()

	f.init(ctx) //NOTE: 初始化，主要是对各个team设置新的IDecode

	if !f.waiting(ctx) {
		return
	}
	go f.doing(ctx) //NOTE: 服务端演算

	ticker, tickIdx, cmds := time.NewTicker(time.Millisecond*100), int32(0), []*any.Any{}
	defer ticker.Stop()
L:
	for {
		select {
		case <-ctx.Done():
			glog.Infoln("oh fight done")
			break L
		case req := <-f.cmdReq:
			cmds = append(cmds, req.GetActions()...)
		case <-ticker.C:
			//此处生成一个Ticker,将缓存中的指令发给客户端
			tick := &war_pb.Tick{
				Index:   tickIdx,
				Actions: cmds,
			}
			cmds = []*any.Any{}
			tickIdx++

			f.cmdReq <- tick
			anyMsg, _ := ptypes.MarshalAny(tick)
			for _, team := range f.teamHub {
				team.Post(anyMsg) //NOTE: 向所有客户端发送tick
			}
		case <-time.After(time.Minute * 10): //NOTE: 10分钟后退出fight
			glog.Infoln("Fight Timeout")
			break L
		}
	}
}

func (f *Fight) init(ctx context.Context) {

	idecode := ctx.Value(CtxKeyIDecode).(IDecode)
	for idx, team := range f.teamHub {
		team.SetIDecode(EncodeTeamID(idecode, uint8(idx)))
	}
}

//doing 服务端演算，用于验证或结算同步
func (f *Fight) doing(ctx context.Context) {

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case tick := <-f.cmdReq:
			f.calculate(tick)
		}
	}
}

//waiting 匹配并等待
func (f *Fight) waiting(ctx context.Context) (success bool) {

	select {
	case <-ctx.Done():
		return false
	case <-time.After(time.Millisecond * 300): //NOTE: 等待0.3s 向全体客户端发送匹配完成消息
	}

	idecode, teams := ctx.Value(CtxKeyIDecode).(IDecode), make([]*war_pb.Team, len(f.teamHub))
	for i, t := range f.teamHub {
		t.SetIDecode(EncodeTeamID(idecode, uint8(i)))
		teams[i] = t.ProtoTeam()
	}
	matched := &war_pb.War_Match{
		MapId:  f.mapID,
		League: teams,
	}
	for _, team := range f.teamHub {
		team.PostMatched(matched)
	}

	select {
	case <-ctx.Done():
		return false
	case <-time.After(time.Second * 3): //NOTE: 3s 后，战斗开始
	}
	return true
}

func (f *Fight) optMove(msg *war_pb.Move) {

}

func (f *Fight) optShoot(msg *war_pb.Shoot) {

}

//calculate 执行下一次tick 演算
func (f *Fight) calculate(tick *war_pb.Tick) {

	for _, opt := range tick.GetActions() {
		name, err := ptypes.AnyMessageName(opt)
		if err != nil {
			continue //NOTE: 应该不会出现这种情况吧
		}
		switch name {
		case "Move": //NOTE: 移动操作
			msg := &war_pb.Move{}
			if err := ptypes.UnmarshalAny(opt, msg); err != nil {
				glog.Warningln("Fight doing Move : ", err)
				continue
			}
			f.optMove(msg)
		case "Shoot": //NOTE: 射击操作
			msg := &war_pb.Shoot{}
			if err := ptypes.UnmarshalAny(opt, msg); err != nil {
				glog.Warningln("Fight doing Shoot : ", err)
				continue
			}
			f.optShoot(msg)
		}
	}
}

//Pit team pit left num
func (f *Fight) Pit() int {

	return cap(f.teamHub) - len(f.teamHub)
}

//Join add team to Fight
func (f *Fight) Join(team *Team) {

	f.teamHub = append(f.teamHub, team)
}

//DoReq 战斗操作请求
func (f *Fight) DoReq() chan<- *war_pb.S_War_Do {

	return f.doReq
}
