package war_test

import (
	"testing"

	"github.com/hiank/thinkend/war"
	war_pb "github.com/hiank/thinkend/war/proto"
	"gotest.tools/v3/assert"
)

func TestFightPit(t *testing.T) {

	check := func (t *testing.T, warType war_pb.War_Type, cnt int)  {
		
		fight := war.NewFight(war.EncodeWarType(0, warType))
		for i:=1; i<=cnt; i++ {
			fight.Join(war.NewTeam(warType))
			assert.Equal(t, fight.Pit(), cnt-i)
		}
	}

	var i war_pb.War_Type
	for i=1; i<5; i++ {
		var cnt int
		switch i {
		case war_pb.War_Type_HUNTING:
			cnt = 2
		case war_pb.War_Type_ONEWINNER:
			cnt = 12
		case war_pb.War_Type_SOCCER:
			cnt = 2
		case war_pb.War_Type_TEAMWINNER:
			cnt = 2
		}
		check(t, i, cnt)
	}
}


func TestFightJoin(t *testing.T) {


}