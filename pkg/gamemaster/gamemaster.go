package gamemaster

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	// "strings"
	"time"

	// MySQLを使うため
	_ "github.com/go-sql-driver/mysql"

	"github.com/satackey/procon31server/pkg/apispec"
	"github.com/satackey/procon31server/pkg/db"
	"github.com/satackey/procon31server/pkg/field"
)

// Match は
type Match struct {
	id int
	DB *sql.DB
}

// Team は
type Team struct {
	JoinedMatchesByLocalTeamID map[int]*joinedMatch // Key: LocalTeamID と Value: MatchID の紐付けをする
	id                         string
	DB                         *sql.DB
}

type joinedMatch struct {
	ID int
	// LocalTeamID int
	UpdateActions []*apispec.UpdateAction
}

// GameMaster は
type GameMaster struct {
	Matches map[int]*Match
	// Teams   map[string]*Team
	// LocalTeamIDs map[int]int
	// GlobalTeamIDsByLocalTeamID map[int]string
	DB *sql.DB
}

// GetMatch は
func GetMatch(id int) (*Match, error) {
	// matchが存在するか調べる
	db, err := db.GetDBConnection()
	if err != nil {
		return nil, fmt.Errorf("GetDBConnection 失敗: %w", err)
	}

	sql := fmt.Sprintf("SELECT id FROM `matches` WHERE `id` = '%d'", id)
	match, err := db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("データベースに接続できませんでした: %w", err)
	}

	for match.Next() {
		var queriedid int
		if err := match.Scan(&queriedid); err != nil {
			return nil, fmt.Errorf("情報の抽出に失敗しました: %w", err)
		}

		if id == queriedid {
			result := &Match{
				id: id,
				DB: db,
			}

			return result, nil
		}
	}

	message := fmt.Sprintf("id: %dが存在しません", id)
	return &Match{}, errors.New(message)
}

// GetTeam は
func GetTeam(id string) (*Team, error) {
	// teamが存在するか調べる
	db, err := db.GetDBConnection()
	if err != nil {
		return nil, fmt.Errorf("GetDBConnection 失敗: %w", err)
	}

	sql := fmt.Sprintf("SELECT global_id FROM `teams` WHERE `global_id` = '%s'", id)
	team, err := db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("データベースに接続できませんでした: %w", err)
	}

	for team.Next() {
		var queriedid string
		if err := team.Scan(&queriedid); err != nil {
			return nil, fmt.Errorf("情報の抽出に失敗しました: %w", err)
		}

		if id == queriedid {
			result := &Team{
				id: id,
				DB: db,
			}

			return result, nil
		}
	}

	message := fmt.Sprintf("id: %sが存在しません", id)
	return &Team{}, errors.New(message)
}

// CreateMatch は 新しい試合を作ります　戻り値 は作られた試合のIDです
func (g *GameMaster) CreateMatch(fieldStatus *apispec.FieldStatus, startsAt int64, turnMillis int, intervalMillis int, turns int, globalTeamID1 string, globalTeamID2 string) (int, error) {
	now := time.Now()
	if now.Unix() > int64(startsAt) {
		return 0, errors.New("startsAtが今の時刻より前です")
	}

	exists, err := g.TeamExists(globalTeamID1)
	if err != nil {
		return 0, err
	}
	if !exists {
		message := fmt.Sprintf("globalTeamID1: %sが存在しません", globalTeamID1)
		return 0, errors.New(message)
	}

	exists, err = g.TeamExists(globalTeamID2)
	if err != nil {
		return 0, err
	}
	if !exists {
		message := fmt.Sprintf("globalTeamID2: %sが存在しません", globalTeamID2)
		return 0, errors.New(message)
	}
	// 渡されたglobalTeamIDたちが存在するかの判定、存在しない場合はその旨をエラーで表す

	fieldJSON := "{}"

	sql := fmt.Sprintf("")

	// creatematch, err := g.DB.QueryRow(sql)
	stmt, err := g.DB.Prepare("INSERT INTO `matches` (`id`, `start_at`, `turn_ms`, `interval_ms`, `turn_num`, `field`) VALUES (NULL, ?, ?, ?, ?, ?)")
	creatematch, err := stmt.Exec(startsAt, turnMillis, intervalMillis, turns, fieldJSON)
	if err != nil {
		return 0, fmt.Errorf("データベースに接続できませんでした1: %w", err)
	}

	// var matchID int
	matchID64, err := creatematch.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("情報の抽出に失敗しました: %w", err)
	}
	matchID := int(matchID64)
	// for creatematch.Next() {
	// 	if err := creatematch.Scan(&matchID); err != nil {
	// 		return 0, fmt.Errorf("情報の抽出に失敗しました: %w", err)
	// 	}
	// }

	sql = fmt.Sprintf("INSERT INTO `match_teams` (`match_id`, `local_team_id`, `global_team_id`, `update_actions`) VALUES ('%d', NULL, '%s', 'null'), ('%d', NULL, '%s', 'null')", matchID, globalTeamID1, matchID, globalTeamID2)
	_, err = g.DB.Query(sql)
	if err != nil {
		// return 0, err // 取得に失敗
		return 0, fmt.Errorf("データベースに接続できませんでした2: %w", err)
	}

	field := &field.Field{}
	field.InitField(fieldStatus)
	// fieldStatusをfieldに、、

	// match.StartAutoTurnUpdate()

	return matchID, nil
	// matchIDを関数の戻り値にする
}

// GetRemainingMSecToTheTransitionOnTurn は nターン終了時までの時間を計算する関数
func (m *Match) GetRemainingMSecToTheTransitionOnTurn(n int, atTime time.Time) (int, error) {
	sql := fmt.Sprintf("SELECT `start_at`, `turn_ms`, `interval_ms` FROM `matches` WHERE `id` = '%d'", m.id)
	matches, err := m.DB.Query(sql)
	if err != nil {
		return 0, fmt.Errorf("取得に失敗しました: %w", err)
	}

	var StartsAt, TurnMillis, IntervalMillis int
	for matches.Next() {
		if err := matches.Scan(&StartsAt, &TurnMillis, &IntervalMillis); err != nil {
			return 0, fmt.Errorf("取得に失敗しました: %w", err)
		}
	}
	// StartsAt = time.Now().Add(time.Duration(1) * time.Minute)
	turnSec := TurnMillis / 1000
	intervalSec := IntervalMillis / 1000
	// timeパッケージにはミリ秒が無いので求め、m.StartsAtをミリ秒にする

	endtime := (StartsAt + turnSec*n + intervalSec*(n-1)) - int(atTime.Unix())
	return int(endtime), nil
}

// StartAutoTurnUpdate は 各ターン終了の時間に点数計算をする
// func (m *Match) StartAutoTurnUpdate() {
// 	// sql := fmt.Sprintf("SELECT `start_at`, `turn_ms`, `interval_ms` FROM `matches` WHERE `id` = '%d'", m.id)
// 	// matches, err := m.DB.Query(sql)
// 	// if err != nil {
// 	// 	return 0, fmt.Errorf("取得に失敗しました: %w", err)
// 	// }
// 	// var StartsAt, TurnMillis, IntervalMillis int
// 	// 	for matches.Next() {
// 	// 	if err := matches.Scan(&StartsAt, &TurnMillis, &IntervalMillis); err != nil {
// 	// 		return 0, fmt.Errorf("取得に失敗しました: %w", err)
// 	// 	}
// 	// }
// 	// startsAtMillis := int64(StartsAt) * 1000
// 	// n64 := int64(n)
// 	// // timeパッケージにはミリ秒が無いので求め、m.StartsAtをミリ秒にする

// 	go func() {
// 		// 今何ターン目かを調べる関数がほしい
// 		startsAtMillis + int64(TurnMillis)*n64
// 		sum, err := m.GetRemainingMSecToTheTransitionOnTurn(1)
// 		if err != nil {
// 			fmt.Println("error:%w", err)
// 		}
// 		time.Sleep(time.Duration(sum) * time.Millisecond)
// 		// 時間を計算する関数を呼び出す
// 		// endtime秒後に
// 	}()
// }

// 今のターンを調べる関数
func (m *Match) GetTurn(n int, atTime time.Time) (int, error) {
	nowMillis := atTime.Unix() * 1000

	sql := fmt.Sprintf("SELECT `start_at`, `turn_ms`, `interval_ms` FROM `matches` WHERE `id` = '%d'", m.id)
	matches, err := m.DB.Query(sql)
	if err != nil {
		return 0, fmt.Errorf("取得に失敗しました: %w", err)
	}

	var StartsAt, TurnMillis, IntervalMillis int
	for matches.Next() {
		if err := matches.Scan(&StartsAt, &TurnMillis, &IntervalMillis); err != nil {
			return 0, fmt.Errorf("取得に失敗しました: %w", err)
		}
	}

	startsAtMillis := int64(StartsAt) * 1000
	// timeパッケージにはミリ秒が無いので求め、m.StartsAtをミリ秒にする

	// hoge := (startsAtMillis + int64(TurnMillis)*n64 + int64(IntervalMillis)*(n64-1))
	// startat turmms1 intervalms1 turnms2 intevalms2
	// 今-(s+1)
	// 今-(s+1)-2
	var hoge int
	var turn int
	for turn = 1; hoge > 0; turn++ {
		hoge = int(nowMillis-startsAtMillis) + TurnMillis*n + IntervalMillis*(n-1)
	}

	return turn, nil
}

// UpdateTurn は 盤面を更新する
func (m *Match) UpdateTurn() error {
	sql := fmt.Sprintf("SELECT update_actions FROM `match_teams` WHERE `match_id` = %d", m.id)
	table, err := m.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("取得に失敗しました: %w", err)
	}
	// dbに存在を問い合わせる

	// Team1UpdateActions
	// Team2UpdateActions
	fmt.Printf("ぷりん:%d", m.id)

	i := 0
	for table.Next() || i > 3 {
		fmt.Printf("ぷりん")
		var queriedUpdateActions []byte
		if err := table.Scan(&queriedUpdateActions); err != nil {
			return fmt.Errorf("情報の抽出に失敗しました: %w", err)
		}

		var updateactions []*apispec.UpdateAction
		if err := json.Unmarshal(queriedUpdateActions, &updateactions); err != nil {
			return fmt.Errorf("Unmarshal 失敗しました: %w", err)
		}
		fmt.Printf("%+v", updateactions)

		i++
	}

	if i != 2 {
		return fmt.Errorf("match_idが2つではない okasii")
	}

	// field.ActAgents()
	// 渡す中身は後で、、
	return nil
}

// GetFieldByID は 指定された試合IDの保管しているFieldStatusを返します
// func (g *GameMaster) GetFieldByID(matchID int) (*apispec.FieldStatus, error) {
// 	match, ok := g.Matches[matchID]

// 	if !ok {
// 		return &apispec.FieldStatus{}, errors.New("試合のIDが存在しません")
// 	}
// 	// 受け取ったmatchIDが存在するかの判定、存在しない場合はその旨をエラーで表す

// 	return match.FieldStatus, nil
// 	// match(g.Matches[key])のなかのFieldStatus
// }

// PostAgentActions は 各チームのエージェントの行動情報を受け取ります
func (g *GameMaster) PostAgentActions(localTeamID int, UpdateActions []*apispec.UpdateAction) error {
	sql := fmt.Sprintf("SELECT local_team_id FROM `match_teams` WHERE `local_team_id` = %d", localTeamID)
	_, err := g.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("localTeamIDが存在しません: %w", err)
	}
	// 存在しないlocalTeamIDを渡されたらエラー

	ua, err := json.Marshal(UpdateActions)
	if err != nil {
		return fmt.Errorf("jsonの読み込みに失敗しました: %w", err)
	}

	sql = fmt.Sprintf("UPDATE `match_teams` SET `update_actions` = '%s' WHERE `match_teams`.`local_team_id` = %d", ua, localTeamID)
	_, err = g.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("書き込みに失敗しました: %w", err)
	}

	return nil
}

// RegisterTeam は チームを登録します
func (g *GameMaster) RegisterTeam(globalTeamID string, name string) error {
	sql := fmt.Sprintf("INSERT INTO `teams` (`global_id`, `name`) VALUES ('%s', '%s')", globalTeamID, name)
	_, err := g.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("チーム登録に失敗しました: %w", err)
	}
	return nil
}

// TeamExists は globalTeamIDが存在するか確認する
func (g *GameMaster) TeamExists(globalTeamID string) (bool, error) {
	sql := fmt.Sprintf("SELECT global_id FROM `teams` WHERE `global_id` = '%s'", globalTeamID)
	teams, err := g.DB.Query(sql)
	if err != nil {
		return false, fmt.Errorf("取得に失敗しました: %w", err)
	}
	// dbに存在を問い合わせる

	for teams.Next() {
		var queriedGlobalTeamID string
		if err := teams.Scan(&queriedGlobalTeamID); err != nil {
			return false, fmt.Errorf("情報の抽出に失敗しました: %w", err)
		}

		if globalTeamID == queriedGlobalTeamID {
			return true, nil
		}
	}

	return false, nil
}

// GetMatchesByGlobalTeamID は 参加する試合の情報を取得します
// func (g *GameMaster) GetMatchesByGlobalTeamID(globalTeamID string) (*apispec.Matches, error) {
// 	team, exists := g.Teams[globalTeamID]
// 	if !exists {
// 		return &apispec.Matches{}, errors.New(strings.Join([]string{"globalTeamID: ", globalTeamID, "が存在しません"}, ""))
// 		// エラー
// 	}
// 	// 存在しないチームIDを取得したらエラー

// 	result := make(apispec.Matches, 0)
// 	team.JoinedMatchesByLocalTeamID[n].ID
// 	for localTeamID, joinedMatchOfTeam := range team.JoinedMatchesByLocalTeamID {
// 		// joinedMatchOfTeam.ID が MatchID

// 		result = append(result, &apispec.Match{
// 			ID:             joinedMatchOfTeam.ID,
// 			IntervalMillis: g.Matches[joinedMatchOfTeam.ID].IntervalMillis,
// 			TeamID:         localTeamID,
// 			TurnMillis:     g.Matches[joinedMatchOfTeam.ID].TurnMillis,
// 			Turns:          g.Matches[joinedMatchOfTeam.ID].Turns,
// 			// MatchTo:
// 		})
// 		// resultを埋めていく
// 	}
// 	return &result, nil
// }

// ConnectDB は DBに接続する
func (g *GameMaster) ConnectDB() error {
	db, err := db.GetDBConnection()
	if err != nil {
		return fmt.Errorf("GetDBConnection 失敗: %w", err)
	}
	g.DB = db
	return nil
}

// Delete は matchesとmatch_teamsのidを削除する
func (m *Match) Delete() error {
	sql := fmt.Sprintf("DELETE FROM `matches` WHERE `id` = '%d'", m.id)
	_, err := m.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("Delete 失敗: %w", err)
	}

	sql = fmt.Sprintf("DELETE FROM `match_teams` WHERE `match_id` = '%d'", m.id)
	_, err = m.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("Delete 失敗: %w", err)
	}
	return nil
}

// Delete は teamsのglobal_idを削除する
func (t *Team) Delete() error {
	sql := fmt.Sprintf("DELETE FROM `teams` WHERE `global_id` = '%s'", t.id)
	_, err := t.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("Delete 失敗: %w", err)
	}
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
