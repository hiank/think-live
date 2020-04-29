package war_test

import (
	"testing"

	"github.com/hiank/thinkend/war"
	"gotest.tools/assert"

	war_pb "github.com/hiank/thinkend/war/proto"
)

var testCode war.IDecode = 43093660299234562 //NOTE: 1001 1001 0001 1001 0111 0100 1010 1010 0011 0101 0000 1101 0000 0010‬

func TestDecodeWarType(t *testing.T) {

	assert.Equal(t, testCode.WarType(), war_pb.War_Type_ONEWINNER)
}

func TestDecodeFightID(t *testing.T) {

	assert.Equal(t, testCode.FightID(), uint32(1957311757))
}

func TestDecodeTeamID(t *testing.T) {

	assert.Equal(t, testCode.TeamID(), uint8(25))
}

func TestDecodeGamerID(t *testing.T) {

	assert.Equal(t, testCode.GamerID(), uint8(153))
}

func TestEncodeWarType(t *testing.T) {

	code := war.EncodeWarType(testCode, war_pb.War_Type_SOCCER)
	assert.Equal(t, code.WarType(), war_pb.War_Type_SOCCER)
}

func TestEncodeFightID(t *testing.T) {

	code := war.EncodeFightID(testCode, 1)
	assert.Equal(t, code.FightID(), uint32(1))
}

func TestEncodeTeamID(t *testing.T) {

	code := war.EncodeTeamID(testCode, 10)
	assert.Equal(t, code.TeamID(), uint8(10))
}

func TestEncodeGamerID(t *testing.T) {

	code := war.EncodeGamerID(testCode, 11)
	assert.Equal(t, code.GamerID(), uint8(11))
}
