package gamemaster

import (
	"math/rand"
	"testing"
	"time"

	"github.com/satackey/procon31server/pkg/apispec"
)

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

func TestConnectDB(t *testing.T) {
	gm := &GameMaster{}
	err := gm.ConnectDB()

	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}
}

func TestRegisterTeam(t *testing.T) {
	gm := &GameMaster{}
	err := gm.ConnectDB()

	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}
	globalid := RandString1(10)
	err = gm.RegisterTeam(globalid, "学校名")
	if err != nil {
		t.Fatalf("チーム登録 失敗: %s", err)
		return
	}
	// Todo: チーム削除
}

func TestTeamExistsAri(t *testing.T) {
	gm := &GameMaster{}
	err := gm.ConnectDB()

	globalid := RandString1(10)
	err = gm.RegisterTeam(globalid, "学校名")
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
	gm := &GameMaster{}
	err := gm.ConnectDB()

	globalid := RandString1(10)
	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}

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
	gm := &GameMaster{}
	err := gm.ConnectDB()
	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}

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

	_, err = gm.CreateMatch(&TestCase, 1599066568, 15000, 2000, 15, "7r64phsgztwm2n4wr27du7nmxnxgaemt3wnnzwxaxc53dw7yt3", "haae42hngzahwewty5azjnnpgaxbibnfyfugpbhd7hmrds2sy7")
	if err != nil {
		t.Fatalf("マッチ登録 失敗: %s", err)
		return
	}
}

func TestGetMatch(t *testing.T) {
	gm := &GameMaster{}
	err := gm.ConnectDB()
	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}

	err = gm.ConnectDB()
	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}

	_, err = GetMatch(gm.DB, 6)
	if err != nil {
		t.Fatalf("失敗: %s", err)
		return
	}
	return
}

func TestGetRemainingMSecToTheTransitionOnTurn(t *testing.T) {
	gm := &GameMaster{}
	err := gm.ConnectDB()
	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}

	m, err := GetMatch(gm.DB, 6)
	if err != nil {
		t.Fatalf("失敗: %s", err)
		return
	}
	sum, err := m.GetRemainingMSecToTheTransitionOnTurn(2)
	if err != nil {
		t.Fatalf("計算失敗: %s", err)
		return
	}
	t.Log(sum)
	// 結果は要検証
	// この方法を使うときは -v オプションを付けないと出力されないぞ
	return
}

// -count=1 オプションはキャッシュなしでtestしてくれるぞ
