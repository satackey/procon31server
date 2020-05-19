package field

import (
	"testing"

	"github.com/satackey/procon31server/pkg/apispec"
)

func TestCelSelectedTimesCount(t *testing.T) {
	f := New()

	points := [][]int{}
	cells := [][]apispec.Cell{}
	teams := []apispec.Team{}
	actions := []apispec.FieldStatusAction{}

	fieldStatus := &apispec.FieldStatus{
		Width:             4,
		Height:            3,
		Points:            points,
		StartedAtUnixtime: 1576800000,
		Turn:              0,
		Cells:             cells,
		Teams:             teams,
		Actions:           actions,
	}

	f.InitField(fieldStatus)

	isValid := []bool{true, true, true, true, true, true, true}
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
