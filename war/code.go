package war

import (
	war_pb "github.com/hiank/thinkend/war/proto"
)

//IDecode 操作单位id 解码
// site [0, 7], War_Type: max 2^7-1 diffent war, one type own one Battle
// site [8, 39], FightID: max 2^32
// site [40, 47], TeamID: max 2^8
// site [48, 55], GamerID: max 2^8
// site [56, 63], reserve
type IDecode uint64

//WarType 获得War_Type
func (id IDecode) WarType() war_pb.War_Type {

	return war_pb.War_Type((1<<8 - 1) & id)
}

//FightID 获得战斗id
func (id IDecode) FightID() uint32 {

	return uint32((1<<40 - 1) & (id >> 8))
}

//TeamID 获得队伍id
func (id IDecode) TeamID() uint8 {

	return uint8((1<<48 - 1) & (id >> 40))
}

//GamerID 获得角色id
func (id IDecode) GamerID() uint8 {

	return uint8((1<<56 - 1) & (id >> 48))
}

//EncodeWarType 设置 War_Type，返回新的id
func EncodeWarType(id IDecode, wt war_pb.War_Type) IDecode {

	return (id & ^IDecode(1<<8 - 1)) | IDecode(wt)
}

//EncodeFightID 设置 fightID 并返回新的id
func EncodeFightID(id IDecode, fid uint32) IDecode {

	return (id & (^(IDecode(1<<32 - 1)<<8))) | (IDecode(fid)<<8)
}

//EncodeTeamID 设置 teadID 并返回新的id
func EncodeTeamID(id IDecode, tid uint8) IDecode {

	return (id & ^(IDecode(1<<8 - 1)<<40)) | (IDecode(tid)<<40)
}

//EncodeGamerID 设置 gamerID 并返回新的id
func EncodeGamerID(id IDecode, gid uint8) IDecode {

	return (id & ^(IDecode(1<<8 - 1)<<48) | (IDecode(gid)<<48))
}
