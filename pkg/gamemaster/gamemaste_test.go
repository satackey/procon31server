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
