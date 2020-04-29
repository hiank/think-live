package war

import (
	"github.com/hiank/think/token"

	master_pb "github.com/hiank/thinkend/master/proto"
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
