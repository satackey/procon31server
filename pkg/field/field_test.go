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
		Width: 4,
		Height: 3,
		Points: points,
		StartedAtUnixtime: 1576800000,
		Turn: 0,
		Cells: cells,
		Teams: teams,
		Actions: actions,
	}

	f.InitField(fieldStatus)
	
	isValid := []bool{true, true, true, true, true, true, true}
	updateActions := []*apispec.UpdateAction{
		&apispec.UpdateAction{
			AgentID: 330,
			DX: 1,
			DY: 1,
			Type: "move",
			X: 0,
			Y: 0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX: 1,
			DY: 1,
			Type: "move",
			X: 0,
			Y: 0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX: 1,
			DY: 1,
			Type: "move",
			X: 0,
			Y: 0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX: 1,
			DY: 1,
			Type: "move",
			X: 0,
			Y: 0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX: 1,
			DY: 1,
			Type: "move",
			X: 0,
			Y: 0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX: 1,
			DY: 1,
			Type: "move",
			X: 0,
			Y: 0,
		},
		&apispec.UpdateAction{
			AgentID: 330,
			DX: 1,
			DY: 1,
			Type: "move",
			X: 0,
			Y: 0,
		},
		
	}
	result := f.CellSelectedTimesCount(isValid, updateActions)
	expext := [][]int{}
	if result == expext {
		// error‚ð‹L˜^
		// t.Error("\nresult: ", result, "\nexpext: ", expext)
	}

	t.Log("Test is finished!")
}

func (f *Field) TestActAgents(t *testing.T) {
	result := f.ActAgents(hoge)
	expext := fuga
	if result != expext {
		// error‚ð‹L˜^
	}
	t.Log("Test is finished")
}