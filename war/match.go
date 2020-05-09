package war

import (
	"context"
	"math"
)

//Matcher 匹配Gamer 到fight
type Matcher struct {
	idecode    IDecode
	joinReq    chan *Gamer
	matchedRes chan *Fight      //NOTE: 每个匹配成功的fight 都会通过这个通知出去
	teamHub    map[int32]*Team  //NOTE: 匹配中的team
	fightHub   map[int32]*Fight //NOTE: 匹配中的fight，每个team匹配完成后，加到对应奖杯数的fight中
}

//NewMatcher 构建匹配器
func NewMatcher(ctx context.Context) *Matcher {

	matcher := &Matcher{
		joinReq:    make(chan *Gamer),
		matchedRes: make(chan *Fight),
		idecode:    ctx.Value(CtxKeyIDecode).(IDecode),
		teamHub:    make(map[int32]*Team),
		fightHub:   make(map[int32]*Fight),
	}
	go matcher.loop(ctx)
	return matcher
}

func (matcher *Matcher) loop(ctx context.Context) {

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case gamer := <-matcher.joinReq:
			key := int32(math.Floor(float64(gamer.GetCup()/30) + 0.5))
			if team := matcher.joinTeam(key, gamer); team.Pit() == 0 {
				delete(matcher.teamHub, key)
				if fight := matcher.joinFight(key, team); fight.Pit() == 0 {
					delete(matcher.fightHub, key)
					matcher.matchedRes <- fight
				}
			}
		}
	}
}

//JoinReq 加入游戏请求
func (matcher *Matcher) JoinReq() chan<- *Gamer {

	return matcher.joinReq
}

//MatchedRes 匹配成功后通知
func (matcher *Matcher) MatchedRes() <-chan *Fight {

	return matcher.matchedRes
}

func (matcher *Matcher) joinTeam(key int32, gamer *Gamer) *Team {

	team, ok := matcher.teamHub[key]
	if !ok {
		team = NewTeam(matcher.idecode.WarType())
		matcher.teamHub[key] = team
	}
	team.Join(gamer)
	return team
}

//team完成匹配后，加入到对应奖杯数的fight中
func (matcher *Matcher) joinFight(key int32, team *Team) *Fight {

	fight, ok := matcher.fightHub[key]
	if !ok {
		fight = NewFight(matcher.idecode)
		matcher.fightHub[key] = fight
	}
	fight.Join(team)
	return fight
}
