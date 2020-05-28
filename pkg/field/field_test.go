package field

import (
	"testing"

	"github.com/satackey/procon31server/pkg/apispec"
)

func SetDataForTest() (*Field, []bool, []*apispec.UpdateAction) {
	f := New()

	// 本来は 12^2 以上 24^2 以下のサイズである
	width := 6
	height := 8
	points := [][]int{
		{4, -5, 0, 0, -1, 4},
		{-4, -1, 0, -4, -3, -2},
		{-1, 3, -2, 4, -1, 3},
		{4, 0, 0, 1, 1, 1},
		{-2, 3, 2, 5, -2, 0},
		{-1, 1, 4, 2, -3, 1},
		{-1, 3, 5, 3, -4, 0},
		{-4, 0, 2, 1, 2, 2},
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
					X:       0,
					Y:       0,
				},
				{
					AgentID: 304,
					X:       0,
					Y:       2,
				},
				{
					AgentID: 305,
					X:       4,
					Y:       2,
				},
				{
					AgentID: 306,
					X:       5,
					Y:       4,
				},
				{
					AgentID: 307,
					X:       1,
					Y:       6,
				},
				{
					AgentID: 308,
					X:       1,
					Y:       3,
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
					X:       2,
					Y:       1,
				},
				{
					AgentID: 404,
					X:       3,
					Y:       4,
				},
				{
					AgentID: 405,
					X:       0,
					Y:       7,
				},
				{
					AgentID: 406,
					X:       1,
					Y:       7,
				},
				{
					AgentID: 407,
					X:       5,
					Y:       7,
				},
				{
					AgentID: 408,
					X:       2,
					Y:       3,
				},
			},
			WallPoint: 0,
			AreaPoint: 0,
		},
	}
	actions := []apispec.FieldStatusAction{}

	fieldStatus := &apispec.FieldStatus{
		Width:             width,
		Height:            height,
		Points:            points,
		StartedAtUnixtime: 1576800000,
		Turn:              0,
		Cells:             cells,
		Teams:             teams,
		Actions:           actions,
	}

	f.InitField(fieldStatus)

	updateActions := []*apispec.UpdateAction{
		{
			AgentID: 303,
			DX:      1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 304,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 305,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       3,
			Y:       4,
		},
		{
			AgentID: 306,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       4,
			Y:       2,
		},
		{
			AgentID: 307,
			DX:      0,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 308,
			DX:      -1,
			DY:      1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 403,
			DX:      -1,
			DY:      0,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 404,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       5,
			Y:       4,
		},
		{
			AgentID: 405,
			DX:      1,
			DY:      -1,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 406,
			DX:      -1,
			DY:      0,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 407,
			DX:      1,
			DY:      0,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 408,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       3,
			Y:       7,
		},
	}

	agentCount := len(updateActions)
	isValid := make([]bool, agentCount)
	for i := range isValid {
		isValid[i] = true
	}
	isValid[10] = false

	return f, isValid, updateActions
}

func TestActAgents(t *testing.T) {
	result, isValid, updateActions := SetDataForTest()
	expected, _, _ := SetDataForTest()


	result.ActAgents(isValid, updateActions)
	// todo: expectedを変更する
	expected.ActionHistories = []ActionHistory{
		{},
		{
			AgentActionHistories: []AgentActionHistory{
				{
					AgentID: 303,
					DX: 1,
					DY: 1,
					X: 0,
					Y: 0,
					Type: "move",
					Turn: 1,
					Apply: 0,
				},
				{
					AgentID: 304,
					DX: 0,
					DY: 0,
					X: 0,
					Y: 0,
					Type: "put",
					Turn: 1,
					Apply: 0,
				},
				{
					AgentID: 305,
					DX: 0,
					DY: 0,
					X: 3,
					Y: 4,
					Type: "put",
					Turn: 1,
					Apply: 1,
				},
				{
					AgentID: 306,
					DX: 0,
					DY: 0,
					X: 4,
					Y: 2,
					Type: "put",
					Turn: 1,
					Apply: 1,
				},
				{
					AgentID: 307,
					DX: 0,
					DY: 1,
					X: 0,
					Y: 0,
					Type: "move",
					Turn: 1,
					Apply: 1,
				},
				{
					AgentID: 308,
					DX: -1,
					DY: 1,
					X: 0,
					Y: 0,
					Type: "move",
					Turn: 1,
					Apply: 1,
				},
				{
					AgentID: 403,
					DX: -1,
					DY: 0,
					X: 0,
					Y: 0,
					Type: "move",
					Turn: 1,
					Apply: 0,
				},
				{
					AgentID: 404,
					DX: 0,
					DY: 0,
					X: 5,
					Y: 4,
					Type: "put",
					Turn: 1,
					Apply: 1,
				},
				{
					AgentID: 405,
					DX: 1,
					DY: -1,
					X: 0,
					Y: 0,
					Type: "move",
					Turn: 1,
					Apply: 1,
				},
				{
					AgentID: 406,
					DX: -1,
					DY: 0,
					X: 0,
					Y: 0,
					Type: "move",
					Turn: 1,
					Apply: 1,
				},
				{
					AgentID: 407,
					DX: 1,
					DY: 0,
					X: 0,
					Y: 0,
					Type: "move",
					Turn: 1,
					Apply: -1,
				},
				{
					AgentID: 408,
					DX: 0,
					DY: 0,
					X: 3,
					Y: 7,
					Type: "move",
					Turn: 1,
					Apply: 1,
				},
			},
		},
	}
	// 明日の僕へ: 競合するようなput行動をするエージェントがいるのでテストケースを変更してね
	expected.Cell[hoge][fuga] = hogehoge

}

func TestConvertIntoHistory(t *testing.T) {
	f, isValid, updateActions := SetDataForTest()

	selectedAgentsIndex := f.RecordCellSelectedAgents(isValid, updateActions)
	isApply := make([]int, len(updateActions))
	f.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex, &isApply)

	result0 := f.ConvertIntoHistory(isValid[0], updateActions[0], isApply[0])
	expected0 := AgentActionHistory{
		AgentID: 303,
		DX:      1,
		DY:      1,
		X:       0,
		Y:       0,
		Type:    "move",
		Turn:    1,
		Apply:   0,
	}

	result3 := f.ConvertIntoHistory(isValid[3], updateActions[3], isApply[3])
	expected3 := AgentActionHistory{
		AgentID: 306,
		DX:      0,
		DY:      0,
		X:       4,
		Y:       2,
		Type:    "put",
		Turn:    1,
		Apply:   1,
	}

	result10 := f.ConvertIntoHistory(isValid[10], updateActions[10], isApply[10])
	expected10 := AgentActionHistory{
		AgentID: 407,
		DX:      1,
		DY:      0,
		X:       0,
		Y:       0,
		Type:    "move",
		Turn:    1,
		Apply:   -1,
	}

	if result0 != expected0 {
		t.Fatalf("\nresult0: %+v\nexpected0: %+v\n", result0, expected0)
	}

	if result3 != expected3 {
		t.Fatalf("\nresult3: %+v\nexpected3: %+v\n", result3, expected3)
	}

	if result10 != expected10 {
		t.Fatalf("\nresult10: %+v\nexpected10: %+v\n", result10, expected10)
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
	expected := []int{0, 0, 1, 1, 1, 1, 0, 1, 1, 1, -1, 1}

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
	expected[0][0] = []int{1}
	expected[1][1] = []int{0, 6}
	expected[2][4] = []int{3}
	expected[4][0] = []int{5}
	expected[4][3] = []int{2}
	expected[4][5] = []int{7}
	expected[6][1] = []int{8}
	expected[7][0] = []int{9}
	expected[7][1] = []int{4}
	expected[7][3] = []int{11}

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
