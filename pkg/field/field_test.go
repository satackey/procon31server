package field

import (
	"testing"

	"github.com/satackey/procon31server/pkg/apispec"
)

func TestCellSelectedTimesCount(t *testing.T) {
	f := New()

	// 本来は 12^2 以上 24^2 以下のサイズである
	width := 5
	height := 4

	points := [][]int{
		{ 3,  2,  1,  0,  2},
		{-1, -3,  3,  1, -2},
		{ 3, -1,  1, -3,  3},
		{ 0, -2,  1,  2,  0},
		
	}
	cells := make([][]apispec.Cell, height)
	for i, cellRow := range cells {
		cells[i] = make([]apispec.Cell, width)
		for j := range cellRow {
			cells[i][j].Status = "free"
			cells[i][j].TeamID = 0
		}
	}
	teams := []apispec.Team{
		apispec.Team{
			TeamID: 3,
			Agents: []apispec.Agent{
				apispec.Agent{
					AgentID: 303,
					X: 1,
					Y: 0,
				},
				apispec.Agent{
					AgentID: 304,
					X: 1,
					Y: 1,
				},
				apispec.Agent{
					AgentID: 305,
					X: 3,
					Y: 1,
				},
				apispec.Agent{
					AgentID: 306,
					X: 0,
					Y: 2,
				},
				
			},
			WallPoint: 0,
			AreaPoint: 0,
		},
		apispec.Team{
			TeamID: 4,
			Agents: []apispec.Agent{
				apispec.Agent{
					AgentID: 403,
					X: 1,
					Y: 2,
				},
				apispec.Agent{
					AgentID: 404,
					X: 2,
					Y: 2,
				},
				apispec.Agent{
					AgentID: 405,
					X: 3,
					Y: 2,
				},
				apispec.Agent{
					AgentID: 406,
					X: 2,
					Y: 3,
				},
				
			},
			WallPoint: 0,
			AreaPoint: 0,
		},
		
	}
	actions := []apispec.FieldStatusAction{}

	fieldStatus := &apispec.FieldStatus{
		Width:             5,
		Height:            4,
		Points:            points,
		StartedAtUnixtime: 1576800000,
		Turn:              0,
		Cells:             cells,
		Teams:             teams,
		Actions:           actions,
	}

	f.InitField(fieldStatus)

	agentCount := 8
	isValid := make([]bool, agentCount)
	for i := range isValid {
		isValid[i] = true
	}

	updateActions := []*apispec.UpdateAction{
		&apispec.UpdateAction{
			AgentID: 303,
			DX:      -1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 304,
			DX:      -1,
			DY:      0,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 305,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 306,
			DX:      1,
			DY:      -1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 403,
			DX:      1,
			DY:      0,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 404,
			DX:      0,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 405,
			DX:      0,
			DY:      -1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 406,
			DX:      -1,
			DY:      -1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
	}

	result := f.CellSelectedTimesCount(isValid, updateActions)
	
	expected := [][]int{
		{0, 0, 0, 0, 0},
		{2, 1, 0, 1, 0},
		{0, 1, 1, 0, 1},
		{0, 0, 1, 0, 0},
		
	}

	// result の配列のサイズは正しいか？
	if len(result) != len(expected) {
		t.Fatalf("\nlen(result): %d\nlen(expected): %d", len(result), len(expected))
	}
	for i := range result {
		if len(result[i]) != len(expected[i]) {
			t.Errorf("\ni: %d\nlen(result[%d]): %d\nlen(expected[%d]): %d", i, i, len(result[i]), i, len(expected))
		}
	}

	// 各マスの数値は正しいか？
	if t.Failed() == false {
		for i, resultLine := range result {
			for j := range resultLine {
				if result[i][j] != expected[i][j] {
					t.Errorf("\ni: %d, j: %d\nresult[%d][%d]: %d\nexpected[%d][%d]: %d", i, j, i, j, result[i][j], i, j, expected[i][j])
				}
			}
		}
	}

	// テストが成功しているなら褒める
	if t.Failed() == false {
		t.Logf("The test is successful!!!")
	}

	t.Log("CellSelectedTimesCount() Test is finished.")
}
