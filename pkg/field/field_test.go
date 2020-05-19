package field

import (
	"testing"

	"github.com/satackey/procon31server/pkg/apispec"
)

func TestCelSelectedTimesCount(t *testing.T) {
	f := New()

	// 本来は 12^2 以上 24^2 以下のサイズである
	points := [][]int{
		[]int{ 3,  2,  1,  0,  2},
		[]int{-1, -3,  3,  1, -2},
		[]int{ 3, -1,  1, -3,  3},
		[]int{ 0, -2,  1,  2,  0},
		
	}
	cells := [][]apispec.Cell{}
	for i, cellLine := range cells {
		for j := range cellLine {
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
		[]int{0, 0, 0, 0, 0},
		[]int{2, 1, 0, 1, 0},
		[]int{0, 1, 1, 0, 1},
		[]int{0, 0, 1, 0, 0},
		
	}

	// result の配列のサイズは正しいか？
	if len(result) != len(expected) {
		t.Fatalf("len(result): ", len(result), "\nlen(expected): ", len(expected))
	}
	for i := range result {
		if len(result[i]) != len(expected[i]) {
			t.Errorf("i: ", i, "\nlen(result[", i, "]): ", len(result[i]), "\nlen(expected[", i, "]): ", len(expected[i]))
		}
	}

	// 各マスの数値は正しいか？
	if t.Failed() == false {
		for i, resultLine := range result {
			for j := resultLine {
				if result[i][j] != expected[i][j] {
					t.Errorf("i: ", i, ", j: ", j, "\nresult[", i, "][", j, "]: ", result[i][j], "\nexpected[", i, "][", j, "]: ", expected[i][j])
				}
			}
		}
	}

	t.Log("Test is finished!")
}
