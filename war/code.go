package war

import (
	war_pb "github.com/hiank/thinkend/war/proto"
)

//IDecode 操作单位id 解码
type IDecode uint64

//WarType 获得War_Type
func (id IDecode) WarType() war_pb.War_Type {

	t := (id >> 58) & (1<<4 - 1)
	return war_pb.War_Type(t)
}

//FightID 获得战斗id
func (id IDecode) FightID() uint32 {

	fightID := (id >> 26) & (1<<32 - 1)
	return uint32(fightID)
}

//TeamID 获得队伍id
func (id IDecode) TeamID() uint32 {

	teamID := (id >> 20) & (1<<6 - 1)
	return uint32(teamID)
}

//RoleID 获得角色id
func (id IDecode) RoleID() uint32 {

	roleID := (id >> 16) & (1<<4 - 1)
	return uint32(roleID)
}

// //IDEncode 生成id
// func IDEncode(wt war_pb.War_Type, fid uint32, tid uint32, rid uint32) uint64 {

// 	id := (uint64(wt) << 58) & (uint64(fid) << 26) & (uint64(tid) << 20) & (uint64(rid) << 16)
// 	return id
// }

// type IDEncode uint64

// func (id IDEncode) SetWarType(uint32) {

// 	id =
// }

//EncodeWarType 设置 War_Type，返回新的id
func EncodeWarType(id IDecode, wt war_pb.War_Type) IDecode {

	// uint64(wt) << 58
	var val IDecode = 1 << 4 - 1
	val = (id & (^(val << 58))) | (IDecode(wt) << 58)
	return val
}

//EncodeFightID 设置 fightID 并返回新的id
func EncodeFightID(id IDecode, fid uint32) IDecode {

	var val IDecode = 1 << 32 - 1
	val = (id & (^(val << 26))) | (IDecode(fid) << 26)
	return val
}

//EncodeTeamID 设置 teadID 并返回新的id
func EncodeTeamID(id IDecode, tid uint32) IDecode {

	var val IDecode = 1 << 6 - 1
	val = (id & (^(val << 20))) | (IDecode(tid) << 20)
	return val
}

//EncodeRoleID 设置 roleID 并返回新的id
func EncodeRoleID(id IDecode, rid uint32) IDecode {

	var val IDecode = 1 << 4 - 1
	val = (id & (^(val << 16))) | (IDecode(rid) << 16)
	return val
}
