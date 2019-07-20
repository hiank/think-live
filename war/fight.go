package war

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/any"
	// "github.com/hiank/think/net/k8s"
	"github.com/hiank/think/pb"
	master_pb "github.com/hiank/thinkend/master/proto"
	war_pb "github.com/hiank/thinkend/war/proto"
)

//Team 队伍
type Team struct {
	// id        uint32   //NOTE: Team id
	idecode   IDecode  //NOTE: 保存 War_Type FightID TeamID信息
	gamers    []*Gamer //NOTE: 队伍中的玩家
	maxpitnum int      //NOTE: 最大坑位数量
}

//NewTeam 创建一个队伍
func NewTeam(idecode IDecode) *Team {

	var num int
	switch idecode.WarType() {

	case war_pb.War_Type_HUNTING:
		num = 3
	case war_pb.War_Type_ONEWINNER:
		num = 1
	case war_pb.War_Type_SOCCER:
		num = 3
	case war_pb.War_Type_TEAMWINNER:
		num = 3
	default:
		glog.Warningln("未知战斗类型")
	}

	team := &Team{
		gamers:    make([]*Gamer, 0, num),
		maxpitnum: num,
		idecode:   idecode,
	}
	return team
}

//Pit 获取队伍剩余坑位
func (t *Team) Pit() int {

	return t.maxpitnum - len(t.gamers)
}

//Join 将gamer加入到team中
func (t *Team) Join(j *join) (err error) {

	gamer := j.gamer
	if len(t.gamers) >= t.maxpitnum {

		glog.Infoln("team join : ", len(t.gamers), t.maxpitnum)
		err = errors.New("team is full")
		return
	}
	t.gamers = append(t.gamers, gamer)
	//此处需要将当前team匹配进度发送给客户端
	anyMsg, err := ptypes.MarshalAny(t.ProtoTeam())
	glog.Infoln(anyMsg, err)
	t.Post(anyMsg)
	return
}

//Post 向Team中的所有玩家发送消息
func (t *Team) Post(anyMsg *any.Any) {

	// glog.Infoln("gamers num : ", len(t.gamers))
	for _, gamer := range t.gamers {

		// glog.Infoln("Team Post ", gamer.GetToken())
		GetNetPool().Post(&pb.Message{

			Key:   gamer.GetKey(),
			Token: gamer.GetToken(),
			Data:  anyMsg,
		})
	}
}

//PostMatched 发送匹配完成消息，因为匹配消息需要对不同玩家发送一个唯一的id，所以需要单独处理，无法直接使用通用的post方法
func (t *Team) PostMatched(msg *war_pb.War_Match) {

	for idx, gamer := range t.gamers {

		msg.Id = uint64(EncodeRoleID(t.idecode, uint32(idx+1)))
		anyMsg, _ := ptypes.MarshalAny(msg)
		GetNetPool().Post(&pb.Message{

			Key:   gamer.GetKey(),
			Token: gamer.GetToken(),
			Data:  anyMsg,
		})
	}
}

//ProtoTeam 转换为war_pb.Team
func (t *Team) ProtoTeam() *war_pb.Team {

	roles := make([]*master_pb.Role, len(t.gamers))
	for i, gamer := range t.gamers {

		roles[i] = gamer.Role
	}
	return &war_pb.Team{
		Roles: roles,
	}
}

//Fight 战斗
type Fight struct {
	*list.Element

	ctx   context.Context    //NOTE:
	Close context.CancelFunc //NOTE:

	idecode IDecode //NOTE: 保存War_Type fightID 信息
	// t 			war_pb.War_Type 		//NOTE: 战斗类型
	cup int32 //NOTE: 基础奖杯数，第一个Gamer加入进来时赋值
	// id 			uint32					//NOTE: 战斗id
	mapID    int32         //NOTE: 地图id, 用于标识不同的障碍物
	tickIdx  int32         //NOTE: 当前tick编号, 用于生成下一个tick
	tick     *list.List    //NOTE: 保存tick信息, 只保存有指令的tick
	lastTick *list.Element //NOTE: 保存最后一次完全确认的tick

	cmds []*any.Any        //NOTE: 命令组, 下一个tick 使用
	cmd  chan *war_pb.Tick //NOTE: 用于驱动服务端战斗计算

	teams []*Team

	back chan *list.Element //NOTE: 匹配完成后写入此chan

	mtx sync.Mutex //NOTE: 读写锁
}

//NewFight 开一场新的战斗
func NewFight(ctx context.Context, idecode IDecode, back chan *list.Element) *Fight {

	ctx, cancel := context.WithCancel(ctx)

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

	f := &Fight{
		ctx:     ctx,
		Close:   cancel,
		idecode: idecode,
		teams:   make([]*Team, 0, teamCnt),
		back:    back,
		cup:     -1,
		cmd: 	 make(chan *war_pb.Tick),
	}
	return f
}

func (f *Fight) start() {

	<-time.After(time.Millisecond * 300)			//NOTE: 等待0.3s 向全体客户端发送匹配完成消息

	teams := make([]*war_pb.Team, len(f.teams))
	for i, t := range f.teams {

		teams[i] = t.ProtoTeam()
	}

	matched := &war_pb.War_Match{

		MapId:  f.mapID,
		League: teams,
	}
	for _, team := range f.teams {

		team.PostMatched(matched)
	}

	<-time.After(time.Second * 3)		//NOTE: 3s 后，战斗开始

	go f.doing()

	glog.Infoln("fight will loop")
	ticker := time.NewTicker(time.Millisecond * 100)
	afterChan := time.After(time.Minute * 10)
L: 	for {

		select {
		case <-f.ctx.Done():
			glog.Infoln("oh fight done")
			break L
		case <-ticker.C:
			// glog.Infoln("hello ticker")
			//此处生成一个Ticker,将缓存中的指令发给客户端
			f.mtx.Lock()
			tick := &war_pb.Tick{

				Index:   f.tickIdx,
				Actions: f.cmds,
			}
			f.cmds = nil
			f.tickIdx++
			f.mtx.Unlock()

			f.cmd <- tick
			anyMsg, _ := ptypes.MarshalAny(tick)
			for _, team := range f.teams {

				team.Post(anyMsg) //NOTE: 向所有客户端发送tick
			}
		case <-afterChan:		//NOTE: 10分钟后退出fight
			f.Close()
			glog.Infoln("Fight Timeout")
			break L
		}
	}
}

func (f *Fight) doing() {

	glog.Infoln("Fight doing")
L: 	for {

		select {
		case <-f.ctx.Done():
			break L
		case tick := <-f.cmd:
			//演算战斗
			// glog.Infoln("doing tick : ", tick)
			A: for _, opt := range tick.GetActions() {

				name, err := ptypes.AnyMessageName(opt)
				if err != nil {

					continue A//NOTE: 应该不会出现这种情况吧
				}
				switch name {
				case "Move": //NOTE: 移动操作
					msg := &war_pb.Move{}
					if err := ptypes.UnmarshalAny(opt, msg); err != nil {
						glog.Warningln("Fight doing Move : ", err)
						continue A
					}
					f.optMove(msg)
				case "Shoot": //NOTE: 射击操作
					msg := &war_pb.Shoot{}
					if err := ptypes.UnmarshalAny(opt, msg); err != nil {
						glog.Warningln("Fight doing Shoot : ", err)
						continue A
					}
					f.optShoot(msg)
				}
			}
			f.calculate()
		}
	}
}

func (f *Fight) optMove(msg *war_pb.Move) {

}

func (f *Fight) optShoot(msg *war_pb.Shoot) {

}

//calculate 执行下一次tick 演算
func (f *Fight) calculate() {

	
}

//GetID 得到战斗id
func (f *Fight) GetID() uint32 {

	return f.idecode.FightID()
}

//match 判断teams 是否都已满员
func (f *Fight) match() {

	if len(f.teams) != cap(f.teams) {
		return
	}
	for _, team := range f.teams {

		if team.Pit() > 0 {
			glog.Infoln("pit ", team.Pit())
			return
		}
	}
	f.back <- f.Element	
	go f.start()		//NOTE: 执行战斗开始
}

//Join 加入战斗，如果匹配失败，返回false，匹配成功，返回true
func (f *Fight) Join(j *join) bool {

	f.mtx.Lock()
	defer f.mtx.Unlock()

	var team *Team
	gamer, num := j.gamer, len(f.teams)

	glog.Infoln("gamer cup : ", gamer.GetCup())
	if f.cup == -1 {

		f.cup = gamer.GetCup()
	} else {

		fix := f.cup - gamer.GetCup()
		if fix > 80 || fix < -80 {
			//NOTE: 匹配失败
			return false
		}
	}
	if num == cap(f.teams) {

		team = f.teams[0]
		for i := 1; i < num; i++ {

			t := f.teams[i]
			if team.Pit() < t.Pit() {
				team = t
			}
		}
	} else {

		team = NewTeam(EncodeTeamID(f.idecode, uint32(num)))
		f.teams = append(f.teams, team)
	}
	team.Join(j)
	f.match()
	return true
}

//Do 处理操作指令
func (f *Fight) Do(d *war_pb.S_War_Do) {

	f.mtx.Lock()

	glog.Infoln("Fight Do : ", d)
	if f.cmds == nil {

		f.cmds = []*any.Any{d.GetAction()}
	} else {

		f.cmds = append(f.cmds, d.GetAction())
	}

	f.mtx.Unlock()
}
