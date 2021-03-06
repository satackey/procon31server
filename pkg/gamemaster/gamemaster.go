package gamemaster

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	// MySQLを使うため
	_ "github.com/go-sql-driver/mysql"

	"github.com/satackey/procon31server/pkg/apispec"
	"github.com/satackey/procon31server/pkg/db"
)

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

// GetRemainingMSecToTheTransitionOnTurn は nターン終了時までの時間を計算する関数
func (m *Match) GetRemainingMSecToTheTransitionOnTurn(n int, atTime time.Time) (int, error) {
	sql := fmt.Sprintf("SELECT `start_at`, `turn_ms`, `interval_ms` FROM `matches` WHERE `id` = '%d'", m.id)
	matches, err := m.DB.Query(sql)
	if err != nil {
		return 0, fmt.Errorf("取得に失敗しました: %w", err)
	}

	var startsAt int64
	var turnMillis, intervalMillis int
	for matches.Next() {
		if err := matches.Scan(&startsAt, &turnMillis, &intervalMillis); err != nil {
			return 0, fmt.Errorf("取得に失敗しました: %w", err)
		}
	}

	endtime := int(startsAt-atTime.Unix())*1000 + turnMillis*(n+1) + intervalMillis*n
	// endtime is Msec
	return endtime, nil
}

// StartAutoTurnUpdate は 各ターン終了の時間に点数計算をする
func (m *Match) StartAutoTurnUpdate() error {
	for {
		turn, err := m.GetTurn(time.Now())
		if err != nil {
			return fmt.Errorf("現在のターン数取得に失敗しました (AutoTurnUpdate は終了されます): %w", err)
		}

		result, err := m.GetRemainingMSecToTheTransitionOnTurn(turn, time.Now())
		if err != nil {
			return fmt.Errorf("現在のターンの終了時刻までの時間の取得に失敗しました (AutoTurnUpdate は終了されます): %w", err)
		}

		// Todo: ログを Println する
		// matchID: ターン3の更新を1, 15 秒後にします
		time.Sleep(time.Duration(result) * time.Millisecond)

		// 途中で AutoTurnUpdate を終了できるようにする

		err = m.UpdateTurn()

		if err != nil {
			return fmt.Errorf("ターンの更新に失敗しました (AutoTurnUpdate は終了されます): %w", err)
		}

		turn, err = m.GetTurn(time.Now())
		if err != nil {
			return fmt.Errorf("現在のターン数取得に失敗しました (AutoTurnUpdate は終了されます): %w", err)
		}

		// 最終ターンになったら終わりにする。
		lastTurn, err := m.GetLastTurn()
		if err != nil {
			return fmt.Errorf("ターン総数の取得に失敗しました: %w", err)
		}
		if turn == lastTurn {
			return nil
		}
	}
}

// ターン総数を調べる関数
func (m *Match) GetLastTurn() (int, error) {
	sql := fmt.Sprintf("SELECT `turn_num` FROM `matches` WHERE `id` = '%d'", m.id)
	matches, err := m.DB.Query(sql)
	if err != nil {
		return -1, fmt.Errorf("取得に失敗しました: %w", err)
	}

	var turnnum int
	for matches.Next() {
		if err := matches.Scan(&turnnum); err != nil {
			return -1, fmt.Errorf("取得に失敗しました: %w", err)
		}
	}
	return turnnum, nil
}

// 今のターンを調べる関数
func (m *Match) GetTurn(atTime time.Time) (int, error) {
	sql := fmt.Sprintf("SELECT `start_at`, `turn_ms`, `interval_ms` FROM `matches` WHERE `id` = '%d'", m.id)
	matches, err := m.DB.Query(sql)
	if err != nil {
		return 0, fmt.Errorf("取得に失敗しました: %w", err)
	}

	var startsAt int64
	var turnMillis, intervalMillis int
	for matches.Next() {
		if err := matches.Scan(&startsAt, &turnMillis, &intervalMillis); err != nil {
			return 0, fmt.Errorf("取得に失敗しました: %w", err)
		}
	}

	if atTime.Unix() < int64(startsAt) {
		fmt.Printf("atTime.Unix() < int64(startsAt)")
		return -1, nil
	}

	// fmt.Printf("%d, %d\n", (int(atTime.Unix()-startsAt) * 1000), (turnMillis + intervalMillis))
	turn := (int(atTime.Unix()-startsAt) * 1000) / (turnMillis + intervalMillis)

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

	i := 0
	for table.Next() || i > 3 {
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

func (m *Match) UpdateTurnkari(Updatefield *apispec.FieldStatus) error {
	sql := fmt.Sprintf("SELECT `field` FROM `matches` WHERE `id` = '%d'", m.id)
	_, err := m.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("取得に失敗しました: %w", err)
	}

	uf, err := json.Marshal(Updatefield)
	if err != nil {
		return fmt.Errorf("jsonの読み込みに失敗しました: %w", err)
	}

	sql = fmt.Sprintf("UPDATE `matches` SET `field` = '%s' WHERE `matches`.`id`= %d", uf, m.id)
	_, err = m.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("書き込みに失敗しました: %w", err)
	}
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
