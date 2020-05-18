package gamemaster

import (
	"errors"
	"time"

	"github.com/satackey/procon31server/pkg/apispec"
)

type Match struct {
	FieldStatus    *apispec.FieldStatus
	TurnMillis     int
	IntervalMillis int
	StartsAt       int
}

type GameMaster struct {
	Matches map[int]*Match
}

// CreateMatch は 新しい試合を作ります　戻り値 は作られた試合のIDです
func (g *GameMaster) CreateMatch(fieldStatus *apispec.FieldStatus, startsAt int, turnMillis int, intervalMillis int) (int, error) {
	now := time.Now()
	if now.Unix() > int64(startsAt) {
		return 0, errors.New("エラーを表す文字列")
	}

	matchID := len(g.Matches)
	// マップの数(=Matchの数)をmatchIDに

	match := &Match{
		FieldStatus:    fieldStatus,
		TurnMillis:     turnMillis,
		IntervalMillis: intervalMillis,
		StartsAt:       startsAt,
		// 型: 値,
	}
	// FieldStatus, 各ターンの時間, ターン数をGameMasterで保管

	g.Matches[matchID] = match
	// matchIDをkeyにしてmap(Matches)に値(match)をセット

	match.Turnendcalc(1)
	// 時間を計算する関数を呼び出す
	// 各ターン終了の時間に点数計算をするようにする

	return matchID, nil
	// matchIDを関数の戻り値にする
}

// Turnendcalc は nターン終了時までの時間を計算する関数
func (m *Match) Turnendcalc(n int) int {
	now := time.Now()
	nowMillis := now.UnixNano() / int64(time.Millisecond)
	startsAtMillis := int64(m.StartsAt) * 1000
	n64 := int64(n)
	// timeパッケージにはミリ秒が無いので求め、m.StartsAtをミリ秒にする

	endtime := (startsAtMillis + int64(m.TurnMillis)*n64 + int64(m.IntervalMillis)*(n64-1)) - nowMillis
	// 求めたいもの = (m.StartsAt /* s */ + m.TurnMillis /* ms */) - いま /* s */
	return int(endtime)
}

// GetFieldByID は 指定された試合IDの保管しているFieldStatusを返します
func (g *GameMaster) GetFieldByID(matchID int) (*apispec.FieldStatus, error) {
	match, ok := g.Matches[matchID]

	if !ok {
		return &apispec.FieldStatus{}, errors.New("エラーを表す文字列")
	}
	// 受け取ったmatchIDが存在するかの判定、存在しない場合はその旨をエラーで表す

	return match.FieldStatus, nil
	// match(g.Matches[key])のなかのFieldStatus
}

// PostAgentActions は
func (g *GameMaster) PostAgentActions(teamID int, FieldStatusAction int) {

}
