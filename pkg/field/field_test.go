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
	teams := []apispec.Team{
		apispac.Team{
			TeamID: 3,
			Agents: []Agent{
				ID: 303,
				TeamID: 3,
				X: 
				Y:
				field: &f,
			}
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
			AgentID: 330,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
	}
	result := f.CellSelectedTimesCount(isValid, updateActions)
	expected := [][]int{}
	if result == expected {
		// errorを記録
		// t.Error("\nresult: ", result, "\nexpected: ", expected)
	}

	t.Log("Test is finished!")
}

func (f *Field) TestActAgents(t *testing.T) {
	result := f.ActAgents(hoge)
	expected := fuga
	if result != expected {
		// errorを記録
	}
	t.Log("Test is finished")
}
