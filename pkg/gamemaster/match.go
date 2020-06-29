package gamemaster

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/satackey/procon31server/pkg/apispec"
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

	fJSON, err := json.Marshal(fieldStatus)
	if err != nil {
		return 0, fmt.Errorf("jsonの読み込みに失敗しました: %w", err)
	}

	stmt, err := g.DB.Prepare("INSERT INTO `matches` (`id`, `start_at`, `turn_ms`, `interval_ms`, `turn_num`, `field`) VALUES (NULL, ?, ?, ?, ?, ?)")
	creatematch, err := stmt.Exec(startsAt, turnMillis, intervalMillis, turns, fJSON)
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

	sql := fmt.Sprintf("INSERT INTO `match_teams` (`match_id`, `local_team_id`, `global_team_id`, `update_actions`) VALUES ('%d', NULL, '%s', 'null'), ('%d', NULL, '%s', 'null')", matchID, globalTeamID1, matchID, globalTeamID2)
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

// RegisterTeam は チームを登録します
func (g *GameMaster) RegisterTeam(globalTeamID string, name string) error {
	sql := fmt.Sprintf("INSERT INTO `teams` (`global_id`, `name`) VALUES ('%s', '%s')", globalTeamID, name)
	_, err := g.DB.Query(sql)
	if err != nil {
		return fmt.Errorf("チーム登録に失敗しました: %w", err)
	}
	return nil
}
