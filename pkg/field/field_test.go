package field

import (
	"testing"

	"github.com/satackey/procon31server/pkg/apispec"
)

func SetDataForTest() (*Field, []bool, []*apispec.UpdateAction){
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
		{
			TeamID: 3,
			Agents: []apispec.Agent{
				{
					AgentID: 303,
					X: 1,
					Y: 0,
				},
				{
					AgentID: 304,
					X: 1,
					Y: 1,
				},
				{
					AgentID: 305,
					X: 3,
					Y: 1,
				},
				{
					AgentID: 306,
					X: 0,
					Y: 2,
				},
				
			},
			WallPoint: 0,
			AreaPoint: 0,
		},
		{
			TeamID: 4,
			Agents: []apispec.Agent{
				{
					AgentID: 403,
					X: 1,
					Y: 2,
				},
				{
					AgentID: 404,
					X: 2,
					Y: 2,
				},
				{
					AgentID: 405,
					X: 3,
					Y: 2,
				},
				{
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
		{
			AgentID: 303,
			DX:      -1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 304,
			DX:      -1,
			DY:      0,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 305,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 306,
			DX:      1,
			DY:      -1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 403,
			DX:      1,
			DY:      0,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 404,
			DX:      0,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 405,
			DX:      0,
			DY:      -1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 406,
			DX:      -1,
			DY:      -1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
	}

	return f, isValid, updateActions
}

func TestConvertIntoHistory(t *testing.T){
	f, isValid, updateActions := SetDataForTest()

	selectedAgentsIndex := f.RecordCellSelectedAgents(isValid, updateActions)
	isApply := make([]int, len(updateActions))
	f.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex, &isApply)

	result0 := f.ConvertIntoHistory(isValid[0], updateActions[0], isApply[0])
	expected0 := AgentActionHistory{
		AgentID: 303,
		DX: -1, 
		DY: 1,
		X: 0,
		Y: 0,
		Type: "move",
		Turn: 1,
		Apply: 0,
	}

	result4 := f.ConvertIntoHistory(isValid[4], updateActions[4], isApply[4])
	expected4 := AgentActionHistory{
		AgentID: 403,
		DX: 1,
		DY: 0,
		X: 0,
		Y: 0,
		Type: "move",
		Turn: 1,
		Apply: 1,
	}

	if result0 != expected0 {
		t.Fatalf("\nresult0: %+v\nexpected0: %+v\n", result0, expected0)
	}

	if result4 != expected4 {
		t.Fatalf("\nresult4: %+v\nexpected4: %+v\n", result4, expected4)
	}

	// テストが成功しているなら褒める
	if t.Failed() == false {
		t.Logf("ConvertIntoHistory() is correct!!!")
	}

	t.Log("Test is finished.")

}

func TestDetermineIfApplied(t *testing.T) {
	f, isValid, updateActions := SetDataForTest()

	selectedAgentsIndex := f.RecordCellSelectedAgents(isValid, updateActions)
	result := make([]int, len(updateActions))

	f.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex, &result)
	expected := []int{0, 0, 1, 0, 1, 1, 1, 1}

	// サイズは正しいか
	if len(result) != len(expected) {
		t.Fatalf("len(result): %d\nlen(expected): %d\n", len(result), len(expected))
	}
	// 値は正しいか
	for i := range result {
		if result[i] != expected[i] {
			t.Fatalf("i: %d\nresult[i]: %d\nexpected[i]: %d\n", i, result[i], expected[i])
		}
	}

	// テストが成功しているなら褒める
	if t.Failed() == false {
		t.Logf("DetermineIfApplied() is correct!!!")
	}

	t.Log("Test is finished.")

}

func TestCellSelectedTimesCount(t *testing.T) {
	f, isValid, updateActions := SetDataForTest()

	result := f.RecordCellSelectedAgents(isValid, updateActions)
	
	expected := make([][][]int, f.Height)
	for i := range expected {
		expected[i] = make([][]int, f.Width)
		for j := range expected[i] {
			expected[i][j] = make([]int, 0)
		}
	}
	expected[1][0] = []int{0, 1}
	expected[1][1] = []int{3}
	expected[1][3] = []int{6}
	expected[2][1] = []int{7}
	expected[2][2] = []int{4}
	expected[2][4] = []int{2}
	expected[3][2] = []int{5}

	// result の配列のサイズは正しいか？
	if len(result) != len(expected) {
		t.Fatalf("\nlen(result): %d\nlen(expected): %d\n", len(result), len(expected))
	}
	for i := range result {
		if len(result[i]) != len(expected[i]) {
			t.Fatalf("\ni: %d\nlen(result[i]): %d\nlen(expected[i]): %d\n", i, len(result[i]), len(expected[i]))
		}
		for j := range result[i] {
			if len(result[i][j]) != len(expected[i][j]) {
				t.Fatalf("\ni: %d, j: %d\nlen(result[i][j]): %d\nlen(expected[i][j]): %d\n", i, j, len(result[i][j]), len(expected[i][j]))
			}
			// 保存されているindexの値は正しいか？
			for k := range result[i][j] {
				if result[i][j][k] != expected[i][j][k] {
					t.Fatalf("\ni: %d, j: %d, k: %d\nresult[i][j][k]: %d\nexpected[i][j][k]: %d\n", i, j, k, result[i][j][k], expected[i][j][k])
				}
				t.Log(result[i][j][k], " ")
			}
		}
	}

	// テストが成功しているなら褒める
	if t.Failed() == false {
		t.Logf("RecordCellSelectedAgents() is correct!!!")
	}

	t.Log("Test is finished.")
}
