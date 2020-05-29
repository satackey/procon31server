package gamemaster

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	// MySQLを使うため
	_ "github.com/go-sql-driver/mysql"

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
	Turns          int
}

// Team は
type Team struct {
	JoinedMatchesByLocalTeamID map[int]*joinedMatch // Key: LocalTeamID と Value: MatchID の紐付けをする
}

type joinedMatch struct {
	ID int
	// LocalTeamID int
	UpdateActions []*apispec.UpdateAction
}

// GameMaster は
type GameMaster struct {
	Matches map[int]*Match
	Teams   map[string]*Team
	// LocalTeamIDs map[int]int
	GlobalTeamIDsByLocalTeamID map[int]string
	DB                         *sql.DB
}

// CreateMatch は 新しい試合を作ります　戻り値 は作られた試合のIDです
// func (g *GameMaster) CreateMatch(fieldStatus *apispec.FieldStatus, startsAt int, turnMillis int, intervalMillis int, turns int, globalTeamID1 string, globalTeamID2 string) (int, error) {
// 	now := time.Now()
// 	if now.Unix() > int64(startsAt) {
// 		return 0, errors.New("startsAtが今の時刻より前です")
// 	}

// 	// sql := fmt.Sprintf("", globalTeamID1)
// 	// globalTeam1, ok1 := g.DB.Query(sql)
// 	// sql = fmt.Sprintf("", globalTeamID2)
// 	// globalTeam2, ok2 := g.DB.Query(sql)
// 	// globalTeam12はどう扱う？、globalteamIDが存在しない場合のエラーとsqlのエラーはこの場合どうする？

// 	// globalTeam1, ok1 := g.Teams[globalTeamID1]
// 	// globalTeam2, ok2 := g.Teams[globalTeamID2]
// 	// _, ok := マップ[キー]
// 	// マップ内にキーが存在するかどうか調べるときはこうやって書く
// 	// if !ok1 {
// 	// 	return 0, errors.New(strings.Join([]string{"globalTeamID: ", globalTeamID1, "が存在しません"}, ""))
// 	// }
// 	// if !ok2 {
// 	// 	return 0, errors.New(strings.Join([]string{"globalTeamID: ", globalTeamID2, "が存在しません"}, ""))
// 	// }
// 	// 渡されたglobalTeamIDたちが存在するかの判定、存在しない場合はその旨をエラーで表す

// 	// sql = fmt.Sprintf("", %d)
// 	// matchID := len(g.Matches)
// 	// // マップの数(=Matchの数)をmatchIDに

// 	// localTeamID1 := matchID * 2
// 	// localTeamID2 := localTeamID1 + 1

// 	// globalTeam1.JoinedMatchesByLocalTeamID[localTeamID1] = &joinedMatch{
// 	// 	ID: matchID,
// 	// }
// 	// globalTeam2.JoinedMatchesByLocalTeamID[localTeamID2] = &joinedMatch{
// 	// 	ID: matchID,
// 	// }

// 	field := &field.Field{}
// 	field.InitField(fieldStatus)
// 	// fieldStatusをfieldに、、

// 	match := &Match{
// 		// FieldStatus:    fieldStatus,
// 		Field:          field,
// 		TurnMillis:     turnMillis,
// 		IntervalMillis: intervalMillis,
// 		StartsAt:       startsAt,
// 		Turns:          turns,
// 		// 型: 値,
// 	}
// 	// FieldStatus, 各ターンの時間, ターン数をGameMasterで保管

// 	g.Matches[matchID] = match
// 	// matchIDをkeyにしてmap(Matches)に値(match)をセット

// 	match.StartAutoTurnUpdate()

// 	return matchID, nil
// 	// matchIDを関数の戻り値にする
// }

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
		// field.ActAgents()
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
func (g *GameMaster) PostAgentActions(localTeamID int, UpdateActions []*apispec.UpdateAction) error {
	globalTeamID, exists := g.GlobalTeamIDsByLocalTeamID[localTeamID]
	if !exists {
		return errors.New(strings.Join([]string{"localTeamID: ", strconv.Itoa(localTeamID), "が存在しません"}, ""))
	}
	// 存在しないlocalTeamIDを渡されたらエラー、存在していたらlocalTeamID → globalTeamID

	g.Teams[globalTeamID].JoinedMatchesByLocalTeamID[localTeamID].UpdateActions = UpdateActions
	// globalTeamID →　Team
	// Team → localTeamID → joinedMatches
	// joinedMatches.UpdateActions ←　代入！！

	return nil
}

// RegisterTeam は チームを登録します
func (g *GameMaster) RegisterTeam(globalTeamID string, name string) error {
	_, exists := g.Teams[globalTeamID]
	if exists {
		return errors.New(strings.Join([]string{"globalTeamID: ", globalTeamID, "はすでに登録されています"}, ""))
		// エラー
	}
	// 同じチームIDを登録しようとしていたらエラー
	sql := fmt.Sprintf("INSERT INTO `teams` (`global_id`, `name`) VALUES ('%s', '%s')", globalTeamID, name)
	_, err := g.DB.Query(sql)
	return err
}

func (g *GameMaster) TeamExists(globalTeamID string) (bool, error) {
	sql := fmt.Sprintf("SELECT global_id FROM `teams` WHERE `global_id` = '%s'", globalTeamID)
	teams, err := g.DB.Query(sql)
	if err != nil {
		return false, err // 取得に失敗しましたとか情報を加える
	}
	// dbに存在を問い合わせる

	for teams.Next() {
		var queriedGlobalTeamID string
		if err := teams.Scan(&queriedGlobalTeamID); err != nil {
			return false, err // 情報の抽出に失敗しました
		}

		if globalTeamID == queriedGlobalTeamID {
			return true, nil
		}
	}

	return false, nil
}

// GetMatchesByGlobalTeamID は 参加する試合の情報を取得します
func (g *GameMaster) GetMatchesByGlobalTeamID(globalTeamID string) (*apispec.Matches, error) {
	team, exists := g.Teams[globalTeamID]
	if !exists {
		return &apispec.Matches{}, errors.New(strings.Join([]string{"globalTeamID: ", globalTeamID, "が存在しません"}, ""))
		// エラー
	}
	// 存在しないチームIDを取得したらエラー

	result := make(apispec.Matches, 0)
	// team.JoinedMatchesByLocalTeamID[n].ID
	for localTeamID, joinedMatchOfTeam := range team.JoinedMatchesByLocalTeamID {
		// joinedMatchOfTeam.ID が MatchID

		result = append(result, &apispec.Match{
			ID:             joinedMatchOfTeam.ID,
			IntervalMillis: g.Matches[joinedMatchOfTeam.ID].IntervalMillis,
			TeamID:         localTeamID,
			TurnMillis:     g.Matches[joinedMatchOfTeam.ID].TurnMillis,
			Turns:          g.Matches[joinedMatchOfTeam.ID].Turns,
			// MatchTo:
		})
		// resultを埋めていく
	}
	return &result, nil
}

// ConnectDB はデータベースに接続します
func (g *GameMaster) ConnectDB() error {
	db, err := sql.Open("mysql", "procon31server:password@tcp(mysql:3306)/procon31")
	if err != nil {
		return fmt.Errorf("データベースに接続できませんでした: %s", err)
	}
	// defer db.Close()
	g.DB = db

	return nil
}

// gm := &GameMaster{}

// `gm.RegisterTeam("tomakomai")
// gm.RegisterTeam("asahikawa")`
// // ..
// gm.CreateMatch(... "tomakomai", "asahikawa")
// gm.CreateMatch(... "tokyo", "tomakomai")
// globalTeamID1 = tomakomai, localTeamID1 = 1
// globalTeamID1 = asahikawa, localTeamID1 = 2
