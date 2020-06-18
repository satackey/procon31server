package gamemaster

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/satackey/procon31server/pkg/apispec"
)

func TestMain(m *testing.M) {
	println("before all...")

	code := m.Run()

	println("after all...")
	err := deleteAllCreatedMatch()
	if err != nil {
		fmt.Println("%w", err)
		os.Exit(1)
		return
	}

	os.Exit(code)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString1(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}

// 0 ～ 配列の要素数までのランダム値取得
func choice(s []string) string {
	rand.Seed(time.Now().UnixNano()) // 乱数のシードとして現在時刻のナノ秒を渡す
	i := rand.Intn(len(s))           // 0 ～ 配列の要素数までのランダム値取得
	return s[i]
}

func TestConnectDB(t *testing.T) {
	gm := &GameMaster{}
	err := gm.ConnectDB()

	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}
}

// TestConnectDBを関数化
func createGameMasterInstanceConnectedDB(tb testing.TB) *GameMaster {
	gm := &GameMaster{}
	err := gm.ConnectDB()
	if err != nil {
		tb.Fatalf("connect 失敗: %s", err)
		return nil
	}
	return gm
}

var registerTeams []string = []string{}

func TestRegisterTeam(t *testing.T) {
	gm := createGameMasterInstanceConnectedDB(t)
	words := []string{"苫小牧", "旭川", "函館", "釧路", "帯広"}

	name := choice(words)
	globalid := RandString1(10)
	err := gm.RegisterTeam(globalid, name)
	// "学校名"だけだとおもしろくないから適当な地名にしてみた、動くかな？
	if err != nil {
		t.Fatalf("チーム登録 失敗: %s", err)
		return
	}
	// Todo: チーム削除
	registerTeams = append(registerTeams, globalid)
}

// TestRegisterTeamを関数化
func registerTeam(tb testing.TB) (string, error) {
	gm := createGameMasterInstanceConnectedDB(tb)
	words := []string{"苫小牧", "旭川", "函館", "釧路", "帯広"}

	name := choice(words)
	globalid := RandString1(10)
	err := gm.RegisterTeam(globalid, name)
	// "学校名"だけだとおもしろくないから適当な地名にしてみた、動くかな？
	if err != nil {
		tb.Fatalf("チーム登録 失敗: %s", err)
		return "", nil
	}

	// Todo: チーム削除
	registerTeams = append(registerTeams, globalid)
	return globalid, nil
}

func TestTeamExistsAri(t *testing.T) {
	gm := createGameMasterInstanceConnectedDB(t)

	globalid, err := registerTeam(t)
	if err != nil {
		t.Fatalf("チーム登録の時点で失敗: %s", err)
		return
	}
	// ランダムIDでチーム登録

	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}

	a1, err := gm.TeamExists(globalid)
	if err != nil {
		t.Fatalf("チーム存在確認 失敗: %s", err)
		return
	}

	if !a1 {
		t.Fatal("チームが存在しません")
	}
}

func TestTeamExistsNasi(t *testing.T) {
	gm := createGameMasterInstanceConnectedDB(t)

	globalid := RandString1(10)

	a1, err := gm.TeamExists(globalid)
	if err != nil {
		t.Fatalf("チーム存在確認 失敗: %s", err)
		return
	}

	if a1 {
		t.Fatal("チームが存在します")
	}
}

func TestCreareMatch(t *testing.T) {
	createMatch(t)
}

var createdMatchIDs []int = []int{}

func createMatch(tb testing.TB) int {
	gm := createGameMasterInstanceConnectedDB(tb)

	cell := apispec.Cell{
		Status: "free",
	}

	TestCase := apispec.FieldStatus{
		Width:             2,
		Height:            2,
		Points:            [][]int{{1, 1}, {1, 1}},
		StartedAtUnixtime: 0,
		Turn:              0,
		Cells:             [][]apispec.Cell{{cell, cell}, {cell, cell}},
		Teams:             []apispec.Team{},
		Actions:           []apispec.FieldStatusAction{},
	}

	globalid1, err := registerTeam(tb)
	globalid2, err := registerTeam(tb)

	matchID, err := gm.CreateMatch(&TestCase, 1599066568, 15000, 2000, 15, globalid1, globalid2)
	if err != nil {
		tb.Fatalf("マッチ登録 失敗: %s", err)
		return 0
	}

	createdMatchIDs = append(createdMatchIDs, matchID)

	return matchID
}

// MatchとTeamを削除
func deleteAllCreatedMatch() error {
	for _, id := range createdMatchIDs {
		match, err := GetMatch(id)
		if err != nil {
			return err
		}
		err = match.Delete()
		if err != nil {
			return err
		}
	}

	for _, id := range registerTeams {
		team, err := GetTeam(id)
		if err != nil {
			return err
		}
		err = team.Delete()
		if err != nil {
			return err
		}
	}

	return nil
}

func TestGetMatch(t *testing.T) {
	matchID := createMatch(t)

	_, err := GetMatch(matchID)
	if err != nil {
		t.Fatalf("失敗: %s", err)
		return
	}
	return
}

func TestGetRemainingMSecToTheTransitionOnTurn(t *testing.T) {
	matchID := createMatch(t)

	m, err := GetMatch(matchID)
	if err != nil {
		t.Fatalf("失敗: %s", err)
		return
	}

	// startatが1599066568
	atTime := time.Unix(1599066568, 0)
	sum, err := m.GetRemainingMSecToTheTransitionOnTurn(2, atTime)
	if err != nil {
		t.Fatalf("計算失敗: %s", err)
		return
	}
	t.Log(sum / 1000)
	// endtimeはミリ秒だから1000で割ってみる
	// 結果は要検証
	// この方法を使うときは -v オプションを付けないと出力されないぞ
	return
}

// -count=1 オプションはキャッシュなしでtestしてくれるぞ

// func TestStartAutoUpdate(t *testing.T) {
// 	gm := &GameMaster{}
// 	err := gm.ConnectDB()
// 	if err != nil {
// 		t.Fatalf("connect 失敗: %s", err)
// 		return
// 	}

// 	err = gm.ConnectDB()
// 	if err != nil {
// 		t.Fatalf("connect 失敗: %s", err)
// 		return
// 	}

// 	m := gm.DB
// 	err = m.StartAutoUpdate()
// 	if err != nil {
// 		t.Fatalf("失敗: %s", err)
// 		return
// 	}
// 	return
// }

func TestPostAgentActions(t *testing.T) {
	gm := createGameMasterInstanceConnectedDB(t)

	TestCase1 := &apispec.UpdateAction{
		AgentID: 2,
		DX:      2,
		DY:      2,
		Type:    "hoge",
		X:       2,
		Y:       2,
	}

	TestCase2 := &apispec.UpdateAction{
		AgentID: 2,
		DX:      2,
		DY:      2,
		Type:    "hoge",
		X:       2,
		Y:       2,
	}

	testCase := []*apispec.UpdateAction{TestCase1, TestCase2}

	err := gm.PostAgentActions(2, testCase)
	if err != nil {
		t.Fatalf("PostAgentActions 失敗: %s", err)
		return
	}
}

func TestUpdateTurn(t *testing.T) {
	matchID := createMatch(t)

	m, err := GetMatch(matchID)
	if err != nil {
		t.Fatalf("失敗: %s", err)
		return
	}

	err = m.UpdateTurn()
	if err != nil {
		t.Fatalf("UpdateTurn 失敗: %s", err)
		return
	}
}
