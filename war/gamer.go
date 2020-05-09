package war

import (
	"errors"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hiank/think/token"

	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
	master_pb "github.com/hiank/thinkend/master/proto"
	war_pb "github.com/hiank/thinkend/war/proto"
)

//Gamer 玩家信息
type Gamer struct {
	*token.Token
	*master_pb.Role //NOTE: 玩家信息
}

//NewGamer create new gamer
func NewGamer(tok *token.Token, pbRole *master_pb.Role) *Gamer {

	return &Gamer{tok, pbRole}
}

//Team 队伍
type Team struct {
	idecode   IDecode  //NOTE: 保存 War_Type FightID TeamID信息
	gamers    []*Gamer //NOTE: 队伍中的玩家
	maxpitnum int      //NOTE: 最大坑位数量
}

//NewTeam 创建一个队伍
func NewTeam(warType war_pb.War_Type) *Team {

	var num int
	switch warType {

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
	}
	return team
}

//Pit 获取队伍剩余坑位
func (t *Team) Pit() int {

	return t.maxpitnum - len(t.gamers)
}

//Join 将gamer加入到team中
func (t *Team) Join(gamer *Gamer) (err error) {

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

//SetIDecode 加入Fight 时，需要设置 IDecode
func (t *Team) SetIDecode(idecode IDecode) {

	t.idecode = idecode
}

//Post 向Team中的所有玩家发送消息
func (t *Team) Post(anyMsg *any.Any) {

	for _, gamer := range t.gamers {

		t.post(gamer, anyMsg)
	}
}

//PostMatched 发送匹配完成消息，因为匹配消息需要对不同玩家发送一个唯一的id，所以需要单独处理，无法直接使用通用的post方法
func (t *Team) PostMatched(msg *war_pb.War_Match) {

	for idx, gamer := range t.gamers {

		msg.Id = uint64(EncodeGamerID(t.idecode, uint8(idx+1)))
		anyMsg, _ := ptypes.MarshalAny(msg)
		t.post(gamer, anyMsg)
	}
}

func (t *Team) post(gamer *Gamer, anyMsg *any.Any) {

	pbMsg := &pb.Message{
		Token: gamer.ToString(),
		Data:  anyMsg,
	}
	var writer k8s.Writer
	writer.Handle(pool.NewMessage(pbMsg, gamer.Derive()))
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
