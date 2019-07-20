package war

import (
	"sync"
	"context"
	"github.com/golang/glog"
	war_pb "github.com/hiank/thinkend/war/proto"
	"github.com/hiank/think/pool"
	master_pb "github.com/hiank/thinkend/master/proto"
)

//Gamer 玩家信息
type Gamer struct {

	key, token string
	*master_pb.Role					//NOTE: 玩家信息
}

//GetKey 获得Gamer 连接的key
func (gamer *Gamer) GetKey() string {

	return gamer.key
}

//GetToken 获得Gamer的token
func (gamer *Gamer) GetToken() string {

	return gamer.token
}

type join struct {

	gamer 		*Gamer							//NOTE: 玩家信息
}


//War 战争，根据War_Type 保存各种Battle
type War struct {

	ctx 	context.Context 					//NOTE: 根Context
	m 		map[war_pb.War_Type]*Battle			//NOTE: map[War_Type]*Battle
	Close 	context.CancelFunc 					//NOTE: 关闭此War中所有
}


//NewWar 创建一个War
func NewWar(ctx context.Context) *War {

	ctx, cancel := context.WithCancel(ctx)
	w := &War{
		ctx 	: ctx,
		m 		: make(map[war_pb.War_Type]*Battle),
		Close 	: cancel,
	}
	return w
}


//Join 加入战斗，排队
func (w *War) Join(t war_pb.War_Type, j *join) {

	b, ok := w.m[t]
	if !ok {
		b = NewBattle(w.ctx, t)
		w.m[t] = b
	}
	b.Join(j)
}

//Do 处理操作命令
func (w *War) Do(d *war_pb.S_War_Do) {

	idecode := IDecode(d.GetId())
	if b, ok := w.m[idecode.WarType()]; ok {

		b.Do(d)
	} else {

		glog.Warningf("can't find battle typed : %v\n", idecode.WarType())
	}
}


type PoolGetter func() *pool.Pool

func (pg PoolGetter) Get() *pool.Pool {

	return pg()
}


var poolGetter PoolGetter
var mtx sync.RWMutex

func SetPoolGetter(pg PoolGetter) {

	mtx.Lock()
	defer mtx.Unlock()

	if poolGetter == nil {
		poolGetter = pg
	}
}

func GetNetPool() *pool.Pool {

	mtx.RLock()
	defer mtx.RUnlock()

	if poolGetter == nil {
		glog.Fatalln("poolGetter should be setted before get")
	}
	return poolGetter.Get()
}



// var netPool *pool.Pool
// var mtx sync.RWMutex

// //SetNetPool 设置netPool
// func SetNetPool(p *pool.Pool) {

// 	mtx.Lock()
// 	defer mtx.Unlock()

// 	if netPool == nil {
// 		netPool = p
// 	}
// }

// //GetNetPool 获得netPool
// func GetNetPool() *pool.Pool {

// 	mtx.RLock()
// 	defer mtx.RUnlock()

// 	if netPool == nil {
// 		glog.Fatalln("netPool should be setted before get")
// 	}
// 	return netPool
// }
