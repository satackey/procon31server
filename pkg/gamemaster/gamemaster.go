package gamemaster

import (
	"errors"
	"strings"
	"time"

	"github.com/satackey/procon31server/pkg/apispec"
	"github.com/satackey/procon31server/pkg/field"
)

// Match は
type Match struct {
	FieldStatus    *apispec.FieldStatus
	Field          *field.Field
	TurnMillis     int
	IntervalMillis int
	StartsAt       int
	Turns  int
}

// Team は
type Team struct {
	JoinedMatchesByLocalTeamID map[int]*joinedMatch // Key: LocalTeamID と Value: MatchID の紐付けをする
}

type joinedMatch struct {
	ID int
	// LocalTeamID int
	UpdateActions apispec.UpdateActions
}

// GameMaster は
type GameMaster struct {
	Matches map[int]*Match
	Teams   map[string]*Team
	// LocalTeamIDs map[int]int
	GlobalTeamIDsByLocalTeamID map[int]string
}

// CreateMatch は 新しい試合を作ります　戻り値 は作られた試合のIDです
func (g *GameMaster) CreateMatch(fieldStatus *apispec.FieldStatus, startsAt int, turnMillis int, intervalMillis int, turns int, globalTeamID1 string, globalTeamID2 string) (int, error) {
	now := time.Now()
	if now.Unix() > int64(startsAt) {
		return 0, errors.New("startsAtが今の時刻より前です")
	}

	globalTeam1, ok1 := g.Teams[globalTeamID1]
	globalTeam2, ok2 := g.Teams[globalTeamID2]
	// _, ok := マップ[キー]
	// マップ内にキーが存在するかどうか調べるときはこうやって書く
	if !ok1 {
		return 0, errors.New(strings.Join("globalTeamID: ", globalTeamID1, "が存在しません"))
	}
	if !ok2 {
		return 0, errors.New(strings.Join("globalTeamID: ", globalTeamID2, "が存在しません"))
	}
	// 渡されたglobalTeamIDたちが存在するかの判定、存在しない場合はその旨をエラーで表す

	matchID := len(g.Matches)
	// マップの数(=Matchの数)をmatchIDに

	localTeamID1 := matchID * 2
	localTeamID2 := localTeamID1 + 1

	globalTeam1.JoinedMatchesByLocalTeamID[localTeamID1] = &joinedMatch{
		ID: matchID,
	}
	globalTeam2.JoinedMatchesByLocalTeamID[localTeamID2] = &joinedMatch{
		ID: matchID,
	}

	field := &field.Field{}
	field.initField(fieldStatus)
	// fieldStatusをfieldに、、

	match := &Match{
		// FieldStatus:    fieldStatus,
		Field:          field,
		TurnMillis:     turnMillis,
		IntervalMillis: intervalMillis,
		StartsAt:       startsAt,
		Turns:			turns,
		// 型: 値,
	}
	// FieldStatus, 各ターンの時間, ターン数をGameMasterで保管

	g.Matches[matchID] = match
	// matchIDをkeyにしてmap(Matches)に値(match)をセット

	match.StartAutoTurnUpdate()

	return matchID, nil
	// matchIDを関数の戻り値にする
}

// GetRemainingMSecToTheTransitionOnTurn は nターン終了時までの時間を計算する関数
func (m *Match) GetRemainingMSecToTheTransitionOnTurn(n int) int {
	now := time.Now()
	nowMillis := now.UnixNano() / int64(time.Millisecond)
	startsAtMillis := int64(m.StartsAt) * 1000
	n64 := int64(n)
	// timeパッケージにはミリ秒が無いので求め、m.StartsAtをミリ秒にする

	endtime := (startsAtMillis + int64(m.TurnMillis)*n64 + int64(m.IntervalMillis)*(n64-1)) - nowMillis
	// 求めたいもの = (m.StartsAt /* s */ + m.TurnMillis /* ms */) - いま /* s */
	return int(endtime)
}

// StartAutoTurnUpdate は 各ターン終了の時間に点数計算をする
func (m *Match) StartAutoTurnUpdate() {

	go func() {
		time.Sleep(time.Duration(m.GetRemainingMSecToTheTransitionOnTurn(1)) * time.Millisecond)
		// 時間を計算する関数を呼び出す
		// endtime秒後にfield.ActAgents()をしたい
		field.ActAgents()
		// 渡す中身は後で、、、
	}()
}

// GetFieldByID は 指定された試合IDの保管しているFieldStatusを返します
func (g *GameMaster) GetFieldByID(matchID int) (*apispec.FieldStatus, error) {
	match, ok := g.Matches[matchID]

	if !ok {
		return &apispec.FieldStatus{}, errors.New("試合のIDが存在しません")
	}
	// 受け取ったmatchIDが存在するかの判定、存在しない場合はその旨をエラーで表す

	return match.FieldStatus, nil
	// match(g.Matches[key])のなかのFieldStatus
}

// PostAgentActions は 各チームのエージェントの行動情報を受け取ります
func (g *GameMaster) PostAgentActions(localTeamID int, FieldStatusAction *apispec.FieldStatusAction) {
	// localTeamID → globalTeamID
	// g.GlobalTeamIDsByLocalTeamID[localTeamID]
	// g.Teams[key].JoinedMatchesByLocalTeamID[]
	// globalTeamID →　Team
	
	// Team → localTeamID → joinedMatches
	// joinedMatches.UpdateActions ←　代入！！
}

// RegisterTeam は チームを登録します
func (g *GameMaster) RegisterTeam(globalTeamID string) error {
	_, exists := g.Teams[globalTeamID]
	if exists {
		return errors.New(strings.Join("globalTeamIDはすでに登録されています"))
		// エラー
	}
	// 同じチームIDを登録しようとしていたらエラー

	g.Teams[globalTeamID].JoinedMatchesByLocalTeamID = make(map[int]*joinedMatch)
	return nil
}

// GetMatchesByGlobalTeamID は 参加する試合の情報を取得します
func (g *GameMaster) GetMatchesByGlobalTeamID(globaTeamID string) (apispec.Matches, error) {
	team, exists := g.Teams[globalTeamID]
	if !exists {
		return &apispec.Matches{}, errors.New(strings.Join("globaslTeamIDが存在しません"))
		// エラー
	}
	// 存在しないチームIDを取得したらエラー

	result := make(apispec.Matches, 0)
	// team.JoinedMatchesByLocalTeamID[n].ID
	for localTeamID, joinedMatchOfTeam := range team.JoinedMatchesByLocalTeamID {
		// joinedMatchOfTeam.ID が MatchID

		result := append(result, &apispec.Match{
			ID:				joinedMatchOfTeam;
			IntervalMillis:	g.Matches[joinedMatchOfTeam].IntervalMillis;
			// MatchTo:
			TeamID:			localTeamID;
			TurnMillis:		g.Matches[joinedMatchOfTeam].TurnMillis;
			Turns:			g.Matches[joinedMatchOfTeam].Turns;
		})
		// resultを埋めていく
	}
	return result, nil
}

// gm := &GameMaster{}

// `gm.RegisterTeam("tomakomai")
// gm.RegisterTeam("asahikawa")`
// // ..
// gm.CreateMatch(... "tomakomai", "asahikawa")
// gm.CreateMatch(... "tokyo", "tomakomai")
// globalTeamID1 = tomakomai, localTeamID1 = 1
// globalTeamID1 = asahikawa, localTeamID1 = 2
