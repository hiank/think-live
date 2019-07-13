package war

import (
	"container/list"
	"context"
	"sync"

	"github.com/golang/glog"
	war_pb "github.com/hiank/thinkend/war/proto"
)

//Battle 战役，包含一组战斗
type Battle struct {
	ctx   context.Context    //NOTE:
	Close context.CancelFunc //NOTE:

	// t war_pb.War_Type //NOTE: 战斗类型
	idecode IDecode //NOTE: 保存WarType 信息

	active  map[uint32]*Fight  //NOTE: map[fightId]*Fight 已经匹配完成的战斗
	waiting *list.List         //NOTE: 匹配中的战斗
	free    *list.List         //NOTE: 空闲id 列表
	max     uint32             //NOTE: 当前最大id值
	matched chan *list.Element //NOTE: 当匹配完成是调用这个chan

	mtx sync.RWMutex //NOTE: 战斗缓存需要读写需要上锁
}


//NewBattle 创建一场战役
func NewBattle(ctx context.Context, t war_pb.War_Type) *Battle {

	ctx, cancel := context.WithCancel(ctx)

	b := &Battle{
		ctx:     ctx,
		Close:   cancel,
		idecode: EncodeWarType(0, t),
		active:  make(map[uint32]*Fight),
		waiting: list.New(),
		free:    list.New(),
		matched: make(chan *list.Element),
	}
	go b.loop()
	return b
}

func (b *Battle) loop() {

L: for {

		select {
		case <-b.ctx.Done():
			glog.Infoln("Battle loop Done")
			break L
		case ele := <-b.matched:
			b.mtx.Lock()
			f := b.waiting.Remove(ele).(*Fight)
			b.active[f.GetID()] = f
			b.mtx.Unlock()
		}
	}
}

//Join 玩家加入战役，排队等待进入战斗
func (b *Battle) Join(j *join) {

	b.mtx.Lock()
	defer b.mtx.Unlock()

	joined := false
	for element := b.waiting.Front(); element != nil; element = element.Next() {

		f := element.Value.(*Fight)
		if f.Join(j) {
			joined = true
			break
		}
	}

	if !joined {

		var id uint32
		if ele := b.free.Front(); ele != nil {

			id = b.free.Remove(ele).(uint32)
		} else {

			b.max++
			id = b.max
		}

		f := NewFight(b.ctx, EncodeFightID(b.idecode, id), b.matched)
		f.Element = b.waiting.PushBack(f)
		f.Join(j)
	}
}

//Do 处理收到的消息
func (b *Battle) Do(d *war_pb.S_War_Do) {

	b.mtx.RLock()
	defer b.mtx.RUnlock()

	idecode := IDecode(d.GetId())
	if f, ok := b.active[idecode.FightID()]; ok {

		f.Do(d)
	} else {

		glog.Warningf("can't find fight %d\n", idecode.FightID())
	}
}
