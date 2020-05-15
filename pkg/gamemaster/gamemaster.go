package gamemaster

import (
	"errors"

	"github.com/satackey/procon31server/pkg/apispec"
)

type Match struct {
	FieldStatus    apispec.FieldStatus
	TurnMillis     int
	IntervalMillis int
}

type GameMaster struct {
	Matches map[int]*Match
}

func (g *GameMaster) CreateMatch(fieldStatus, turnMillis, intervalMillis) {
	matchID := len(g.Matches)
	// マップの数(=Matchの数)をmatchIDに

	match := &Match{
		FieldStatus:    fieldStatus,
		TurnMillis:     turnMillis,
		IntervalMillis: intervalMillis,
		// 型: 値,
	}
	// FieldStatus, 各ターンの時間, ターン数をGameMasterで保管

	g.Matches[matchID] = match
	// matchIDをkeyにしてmap(Matches)に値(match)をセット

	// todo 各ターン終了の時間に点数計算をするようにする

	return macthID
	// matchIDを関数の戻り値にする
}

func (g *GameMaster) GetFieldByID(matchID) (apispec.FieldStatus, error) {
	match, ok := g.Matches[matchID]

	if !ok {
		return nil, errors.New("エラーを表す文字列")
	}
	// 受け取ったmatchIDが存在するかの判定、存在しない場合はその旨をエラーで表す

	return match.fieldStatus
	// match(g.Matches[key])のなかのFieldStatus
}

func (g *GameMaster) PostAgentActions(teamID, FieldStatusAction) {

}
