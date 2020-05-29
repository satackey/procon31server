package gamemaster

import (
	"testing"
)

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

	err = gm.RegisterTeam("hoge", "釧路")
	if err != nil {
		t.Fatalf("チーム登録 失敗: %s", err)
		return
	}
}

func TestTeamExistsAri(t *testing.T) {
	gm := &GameMaster{}
	err := gm.ConnectDB()

	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}

	a1, err := gm.TeamExists("hoge")
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

	if err != nil {
		t.Fatalf("connect 失敗: %s", err)
		return
	}

	a1, err := gm.TeamExists("neko")
	if err != nil {
		t.Fatalf("チーム存在確認 失敗: %s", err)
		return
	}

	if a1 {
		t.Fatal("チームが存在します")
	}
}
