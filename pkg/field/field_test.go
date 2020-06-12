package field

import (
	"testing"

	"github.com/satackey/procon31server/pkg/apispec"
)

func GetTestCase01() (*Field, []bool, []*apispec.UpdateAction, []int) {
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
	for y := range cells {
		cells[y] = make([]apispec.Cell, width)
		for x := range cells[y] {
			cells[y][x].Status = "free"
			cells[y][x].TeamID = 0
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

	for y := range f.Cells {
		for x := range f.Cells[y] {
			f.Cells[y][x].x = x
			f.Cells[y][x].y = y
			f.Cells[y][x].field = f
			f.Cells[y][x].Point = points[y][x]
		}
	}

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
			AgentID: 0,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       2,
			Y:       1,
		},
		{
			AgentID: 0,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       4,
			Y:       2,
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
			AgentID: 0,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       1,
			Y:       3,
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
			AgentID: 403,
			DX:      -1,
			DY:      0,
			Type:    "move",
			X:       0,
			Y:       0,
		},
		{
			AgentID: 0,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       4,
			Y:       2,
		},
		{
			AgentID: 0,
			DX:      0,
			DY:      0,
			Type:    "put",
			X:       3,
			Y:       7,
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
	}

	agentCount := len(updateActions)
	isValid := make([]bool, agentCount)
	for i := range isValid {
		isValid[i] = true
	}
	isValid[11] = false

	updateActionIDs := []int{3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4}

	return f, isValid, updateActions, updateActionIDs
}

func GetTestCase02() (*Field, []bool, []*apispec.UpdateAction, []int) {
	f := New()

	width := 12
	height := 13

	points := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	w := apispec.Cell{
		Status: "wall",
		TeamID: 3,
	}
	e := apispec.Cell{
		Status: "free",
		TeamID: 0,
	}
	r := apispec.Cell{
		Status: "wall",
		TeamID: 4,
	}
	cells := [][]apispec.Cell{
		{w,w,w,w,w,e,e,e,r,r,r,r,},
		{w,e,e,e,w,w,e,r,r,e,e,r,},
		{w,w,w,w,w,w,w,w,w,r,r,r,},
		{w,r,r,r,r,e,e,e,w,w,w,w,},
		{w,e,e,e,r,r,r,r,r,r,r,w,},
		{w,r,r,r,r,e,e,w,e,e,r,w,},
		{w,r,e,e,r,e,w,w,w,e,r,w,},
		{w,r,r,r,r,r,w,e,w,w,r,w,},
		{w,w,e,w,r,e,w,w,w,e,r,w,},
		{w,r,r,r,r,e,e,r,e,e,r,w,},
		{w,r,e,e,r,r,r,r,r,r,r,w,},
		{w,r,r,r,e,e,e,e,e,e,e,w,},
		{w,w,w,w,w,w,w,w,w,w,w,w,},
	}
	teams := []apispec.Team{
		{
			TeamID: 3,
			Agents: []apispec.Agent{
				{
					AgentID: 303,
					X: 7,
					Y: 0,
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
					X: 2,
					Y: 8,
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
	for y := range f.Cells {
		for x := range f.Cells[y] {
			f.Cells[y][x].x = x
			f.Cells[y][x].y = y
			f.Cells[y][x].field = f
			f.Cells[y][x].Point = points[y][x]
		}
	}

	updateActions := []*apispec.UpdateAction{
		{
			AgentID: 303,
			DX: 0,
			DY: 1,
			Type: "remove",
			X: 0,
			Y: 0,
		},
		{
			AgentID: 403,
			DX: 1,
			DY: 0,
			Type: "remove",
			X: 0,
			Y: 0,
		},
	}

	isValid := []bool{true, true}

	updateActionIDs := []int{3, 4}

	return f, isValid, updateActions, updateActionIDs
}

func TestFinalCheckByDFS(t *testing.T) {
	f, isValid, updateActions, updateActionIDs := GetTestCase02()
	updateAction2s := f.MakeUpdateAction2s(updateActions, updateActionIDs)
	selectedAgentsIndex := f.RecordCellSelectedAgents(isValid, updateActions)
	isApply := f.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex)
	agentActionHistories := make([]AgentActionHistory, len(updateActions))
	for i, updateAction := range updateActions {
		agentActionHistories[i] = f.ConvertIntoHistory(isValid[i], updateAction, isApply[i])
		if agentActionHistories[i].Apply == 1 {
			f.ActuallyActAgent(updateAction2s[i])
		}
	}

	x := []int{6, 9, 7}
	y := []int{2, 5, 7}
	expected := []bool{false, false, true}

	for i := 0; i < 3; i ++ {
		isAreaBy := map[int][][]bool{
			3: {},
			4: {},
		}
		for teamID := range isAreaBy {
			isAreaBy[teamID] = make([][]bool, f.Height)
			for y := range isAreaBy[teamID] {
				isAreaBy[teamID][y] = make([]bool, f.Width)
			}
		}
		if f.CheckAreaByDFS(3, x[i], y[i], &isAreaBy) != true {
			t.Fatalf("\nerror\n")
		}
		if f.CheckAreaByDFS(4, x[i], y[i], &isAreaBy) != true {
			t.Fatalf("\nerror\n")
		}
		if f.FinalCheckByDFS(3, x[i], y[i], isAreaBy[4]) != expected[i] {
			t.Fatalf("\ni: %d\nf.FinalCheckByDFS(~): %v\nexpected: %v\n", i, f.FinalCheckByDFS(3, x[i], y[i], isAreaBy[4]), expected[i])
		}
	}

	// テストが成功しているなら褒める
	if t.Failed() == false {
		t.Logf("FinalCheckByDFS() is correct!!!")
	}
}

func TestCheckAreaByDFS(t *testing.T) {
	f, isValid, updateActions, updateActionIDs := GetTestCase02()
	updateAction2s := f.MakeUpdateAction2s(updateActions, updateActionIDs)
	selectedAgentsIndex := f.RecordCellSelectedAgents(isValid, updateActions)
	isApply := f.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex)
	agentActionHistories := make([]AgentActionHistory, len(updateActions))
	for i, updateAction := range updateActions {
		agentActionHistories[i] = f.ConvertIntoHistory(isValid[i], updateAction, isApply[i])
		if agentActionHistories[i].Apply == 1 {
			f.ActuallyActAgent(updateAction2s[i])
		}
	}

	x, y := 7, 7

	result := map[int][][]bool{
		3: {},
		4: {},
	}
	for teamID := range result {
		result[teamID] = make([][]bool, f.Height)
		for y := range result[teamID] {
			result[teamID][y] = make([]bool, f.Width)
		}
	}

	if f.CheckAreaByDFS(3, x, y, &result) != true {
		t.Fatalf("error")
	}
	if f.CheckAreaByDFS(4, x, y, &result) != true {
		t.Fatalf("error")
	}

	expected := map[int][][]bool{}
	expected[3] = [][]bool{
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false, true,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
	}
	expected[4] = [][]bool{
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false, true, true, true, true, true,false,false,},
		{false,false,false,false,false, true, true, true, true, true,false,false,},
		{false,false,false,false,false,false, true, true, true, true,false,false,},
		{false,false,false,false,false, true, true, true, true, true,false,false,},
		{false,false,false,false,false, true, true,false, true, true,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
		{false,false,false,false,false,false,false,false,false,false,false,false,},
	}
	if len(result) != len(expected) {
		t.Fatalf("\nlen(result): %d\nlen(expected): %d\n", len(result), len(expected))
	}
	for i := range result {
		_, ok := expected[i]
		if ok == false {
			t.Fatalf("\n%d does not exist\n", i)
		}
		if len(result[i]) != len(expected[i]) {
			t.Fatalf("\ni: %d\nlen(result[i]): %d\nlen(expected[i]): %d\n", i, len(result[i]), len(expected[i]))
		}
		for y := range expected[i] {
			if len(result[i][y]) != len(expected[i][y]) {
				t.Fatalf("\ni: %d, y: %d\nlen(result[i][y]): %d\nlen(expected[i][y]): %d\n", i, y, len(result[i][y]), len(expected[i][y]))
			}
			for x := range result[i][y] {
				if result[i][y][x] != expected[i][y][x] {
					t.Fatalf("\ni: %d, y: %d, x: %d\nresult[i][y][x]: %v\nexpected[i][y][x]: %v\n", i, y, x, result[i][y][x], expected[i][y][x])
				}
			}
		}
	}

	// テストが成功しているなら褒める
	if t.Failed() == false {
		t.Logf("CheckAreaByDFS() is correct!!!")
	}
}

func TestActAgents(t *testing.T) {
	result, isValid, updateActions, updateActionIDs := GetTestCase01()
	expected, _, _, _ := GetTestCase01()

	result.ActAgents(isValid, updateActions, updateActionIDs)
	expected.ActionHistories = []ActionHistory{
		{},
		{
			AgentActionHistories: []AgentActionHistory{
				{
					AgentID: 303,
					DX:      1,
					DY:      1,
					X:       0,
					Y:       0,
					Type:    "move",
					Turn:    1,
					Apply:   0,
				},
				{
					AgentID: 0,
					DX:      0,
					DY:      0,
					X:       2,
					Y:       1,
					Type:    "put",
					Turn:    1,
					Apply:   0,
				},
				{
					AgentID: 0,
					DX:      0,
					DY:      0,
					X:       4,
					Y:       2,
					Type:    "put",
					Turn:    1,
					Apply:   0,
				},
				{
					AgentID: 308,
					DX:      -1,
					DY:      1,
					X:       0,
					Y:       0,
					Type:    "move",
					Turn:    1,
					Apply:   1,
				},
				{
					AgentID: 0,
					DX:      0,
					DY:      0,
					X:       1,
					Y:       3,
					Type:    "put",
					Turn:    1,
					Apply:   1,
				},
				{
					AgentID: 307,
					DX:      0,
					DY:      1,
					X:       0,
					Y:       0,
					Type:    "move",
					Turn:    1,
					Apply:   1,
				},

				{
					AgentID: 403,
					DX:      -1,
					DY:      0,
					X:       0,
					Y:       0,
					Type:    "move",
					Turn:    1,
					Apply:   0,
				},
				{
					AgentID: 0,
					DX:      0,
					DY:      0,
					X:       4,
					Y:       2,
					Type:    "put",
					Turn:    1,
					Apply:   0,
				},
				{
					AgentID: 0,
					DX:      0,
					DY:      0,
					X:       3,
					Y:       7,
					Type:    "put",
					Turn:    1,
					Apply:   1,
				},
				{
					AgentID: 405,
					DX:      1,
					DY:      -1,
					X:       0,
					Y:       0,
					Type:    "move",
					Turn:    1,
					Apply:   1,
				},
				{
					AgentID: 406,
					DX:      -1,
					DY:      0,
					X:       0,
					Y:       0,
					Type:    "move",
					Turn:    1,
					Apply:   1,
				},
				{
					AgentID: 407,
					DX:      1,
					DY:      0,
					X:       0,
					Y:       0,
					Type:    "move",
					Turn:    1,
					Apply:   -1,
				},
			},
		},
	}
	expected.Cells[4][0].TiledBy = 3
	expected.Cells[4][0].Status = "wall"
	expected.Cells[6][1].TiledBy = 4
	expected.Cells[6][1].Status = "wall"
	expected.Cells[7][0].TiledBy = 4
	expected.Cells[7][0].Status = "wall"
	expected.Cells[7][1].TiledBy = 3
	expected.Cells[7][1].Status = "wall"

	expected.Agents[1] = &Agent{
		ID:     1,
		TeamID: 3,
		X:      1,
		Y:      3,
		field:  expected,
	}
	expected.Agents[2] = &Agent{
		ID:     2,
		TeamID: 4,
		X:      3,
		Y:      7,
		field:  expected,
	}
	expected.Agents[308].X = 0
	expected.Agents[308].Y = 4
	expected.Agents[307].X = 1
	expected.Agents[307].Y = 7
	expected.Agents[405].X = 1
	expected.Agents[405].Y = 6
	expected.Agents[406].X = 0
	expected.Agents[406].Y = 7

	// Test用 ポインタを合わせる
	for i := range expected.Agents {
		expected.Agents[i].field = result
	}

	// 各要素に対して判定する
	if result.Width != expected.Width {
		t.Fatalf("\nresult.Width: %d\nexpected.Width: %d\n", result.Width, expected.Width)
	}
	if result.Height != expected.Height {
		t.Fatalf("\nresult.Height: %d\nexpected.Height: %d\n", result.Height, expected.Height)
	}
	if result.Turn != expected.Turn {
		t.Fatalf("\nresult.Turn: %d\nexpected.Turn: %d\n", result.Turn, expected.Turn)
	}

	if len(result.Cells) != len(expected.Cells) {
		t.Fatalf("\nlen(result.Cells): %d\nlen(expected.Cells): %d\n", len(result.Cells), len(expected.Cells))
	}
	for i := range result.Cells {
		if len(result.Cells[i]) != len(expected.Cells[i]) {
			t.Fatalf("\ni: %d\nlen(result.Cells[i]): %d\nlen(expected.Cells[i]): %d\n", i, len(result.Cells[i]), len(expected.Cells[i]))
		}
		for j := range result.Cells[i] {
			if result.Cells[i][j] == expected.Cells[i][j] {
				t.Fatalf("\ni: %d, j: %d\nresult.Cells[i][j]: %+v\nexpected.Cells[i][j]: %+v\n", i, j, result.Cells[i][j], expected.Cells[i][j])
			}
		}
	}
	if len(result.Agents) != len(expected.Agents) {
		for i := range result.Agents {
			t.Logf("%+v\n", result.Agents[i])
		}
		t.Fatalf("\nlen(result.Agents): %d\nlen(expected.Agents): %d\n", len(result.Agents), len(expected.Agents))
	}
	for resultKey, resultValue := range result.Agents {
		expectedValue, isExist := expected.Agents[resultKey]
		if isExist == false {
			t.Fatalf("\nkey: %d\nresult.Agents[key] is exist, But expected.Agents[key] is not exist.\n", resultKey)
		}
		if *resultValue != *expectedValue {
			t.Fatalf("\nkey: %d\nresult.Agents[key]: %+v\nexpected.Agents[key]: %+v\n", resultKey, resultValue, expectedValue)
		}
	}
	if len(result.Teams) != len(expected.Teams) {
		t.Fatalf("\nlen(result.Teams): %d\nlen(expected.Teams): %d\n", len(result.Teams), len(expected.Teams))
	}
	for i := range result.Teams {
		if *result.Teams[i] != *expected.Teams[i] {
			t.Fatalf("\ni: %d\nresult.Teams[i]: %+v\nexpected.Teams[i]: %+v\n", i, *result.Teams[i], *expected.Teams[i])
		}
	}
	if len(result.ActionHistories) != len(expected.ActionHistories) {
		t.Fatalf("\nlen(result.ActionHistories): %d\nlen(expected.ActionHistories): %d\n", len(result.ActionHistories), len(expected.ActionHistories))
	}
	for i := range result.ActionHistories {
		if len(result.ActionHistories[i].AgentActionHistories) != len(expected.ActionHistories[i].AgentActionHistories) {
			t.Fatalf("\nlen(result.ActionHistories[i].AgentActionHistories): %d\nlen(expected.ActionHistories[i].AgentActionHistories): %d\n", len(result.ActionHistories[i].AgentActionHistories), len(expected.ActionHistories[i].AgentActionHistories))
		}
		for j := range result.ActionHistories[i].AgentActionHistories {
			if result.ActionHistories[i].AgentActionHistories[j] != expected.ActionHistories[i].AgentActionHistories[j] {
				t.Fatalf("\ni: %d, j: %d\nresult.ActionHistories[i].AgentActionHistories[j]: %+v\nexpected.ActionHistories[i].AgentActionHistories[j]: %+v\n", i, j, result.ActionHistories[i].AgentActionHistories[j], expected.ActionHistories[i].AgentActionHistories[j])
			}
		}
	}


	// テストが成功しているなら褒める
	if t.Failed() == false {
		t.Logf("ActAgents() is correct!!!")
	}
}

func TestActuallyActAgent(t *testing.T) {
	result, isValid, updateActions, updateActionIDs := GetTestCase01()
	expected, _, _, _ := GetTestCase01()

	t.Logf("%+v\n", len(result.Agents))

	selectedAgentsIndex := result.RecordCellSelectedAgents(isValid, updateActions)
	isApply := result.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex)
	agentActionHistories := make([]AgentActionHistory, len(updateActions))
	for i := range updateActions {
		agentActionHistories[i] = result.ConvertIntoHistory(isValid[i], updateActions[i], isApply[i])
	}
	updateAction2s := result.MakeUpdateAction2s(updateActions, updateActionIDs)

	result.ActuallyActAgent(updateAction2s[3])
	result.ActuallyActAgent(updateAction2s[4])
	result.ActuallyActAgent(updateAction2s[5])
	result.ActuallyActAgent(updateAction2s[8])
	result.ActuallyActAgent(updateAction2s[9])
	result.ActuallyActAgent(updateAction2s[10])

	expected.Agents[308].X = 0
	expected.Agents[308].Y = 4
	expected.Cells[4][0].TiledBy = 3
	expected.Cells[4][0].Status = "wall"
	expected.Agents[307].X = 1
	expected.Agents[307].Y = 7
	expected.Cells[7][1].TiledBy = 3
	expected.Cells[7][1].Status = "wall"
	expected.Agents[405].X = 1
	expected.Agents[405].Y = 6
	expected.Cells[6][1].TiledBy = 4
	expected.Cells[6][1].Status = "wall"
	expected.Agents[406].X = 0
	expected.Agents[406].Y = 7
	expected.Cells[7][0].TiledBy = 4
	expected.Cells[7][0].Status = "wall"

	expected.Agents[1] = &Agent{
		ID:     1,
		TeamID: 3,
		X:      1,
		Y:      3,
		field:  expected,
	}
	expected.Agents[2] = &Agent{
		ID:     2,
		TeamID: 4,
		X:      3,
		Y:      7,
		field:  expected,
	}

	// テスト用に Field.Agents[].field のアドレスを合わせる
	for i := range expected.Agents {
		expected.Agents[i].field = result
	}

	for y := range result.Cells {
		for x := range result.Cells[y] {
			if result.Cells[y][x].TiledBy != expected.Cells[y][x].TiledBy {
				t.Fatalf("\ny: %d, x:%d\nresult.Cells[y][x].TiledBy: %d\nexpected.Cells[y][x].TiledBy: %d\n", y, x, result.Cells[y][x].TiledBy, expected.Cells[y][x].TiledBy)
			}
			if result.Cells[y][x].Status != expected.Cells[y][x].Status {
				t.Fatalf("\ny: %d, x:%d\nresult.Cells[y][x].Status: %s\nexpected.Cells[y][x].Status: %s\n", y, x, result.Cells[y][x].Status, expected.Cells[y][x].Status)
			}
		}
	}
	if len(result.Agents) != len(expected.Agents) {
		t.Fatalf("\nlen(result.Agents): %d\nlen(expected.Agents): %d\n", len(result.Agents), len(expected.Agents))
	}
	for key := range result.Agents {
		if expected.Agents[key] == nil {
			t.Fatalf("\nkey: %d\nexpected.Agents[key] is not exist\n", key)
		}
		if *result.Agents[key] != *expected.Agents[key] {
			t.Fatalf("\nkey: %d\nresult.Agents[key]: %+v\nexpected.Agents[key]: %+v\n", key, result.Agents[key], expected.Agents[key])
		}
	}

}

func TestConvertIntoHistory(t *testing.T) {
	f, isValid, updateActions, _ := GetTestCase01()

	selectedAgentsIndex := f.RecordCellSelectedAgents(isValid, updateActions)
	isApply := f.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex)

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

	result4 := f.ConvertIntoHistory(isValid[4], updateActions[4], isApply[4])
	expected4 := AgentActionHistory{
		AgentID: 0,
		DX:      0,
		DY:      0,
		X:       1,
		Y:       3,
		Type:    "put",
		Turn:    1,
		Apply:   1,
	}

	result11 := f.ConvertIntoHistory(isValid[11], updateActions[11], isApply[11])
	expected11 := AgentActionHistory{
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

	if result4 != expected4 {
		t.Fatalf("\nresult4: %+v\nexpected4: %+v\n", result4, expected4)
	}

	if result11 != expected11 {
		t.Fatalf("\nresult11: %+v\nexpected11: %+v\n", result11, expected11)
	}

	// テストが成功しているなら褒める
	if t.Failed() == false {
		t.Logf("ConvertIntoHistory() is correct!!!")
	}

}

func TestDetermineIfApplied(t *testing.T) {
	f, isValid, updateActions, _ := GetTestCase01()

	selectedAgentsIndex := f.RecordCellSelectedAgents(isValid, updateActions)

	result := f.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex)
	expected := []int{0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, -1}

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

}

func TestCellSelectedTimesCount(t *testing.T) {
	f, isValid, updateActions, _ := GetTestCase01()

	result := f.RecordCellSelectedAgents(isValid, updateActions)

	expected := make([][][]int, f.Height)
	for i := range expected {
		expected[i] = make([][]int, f.Width)
		for j := range expected[i] {
			expected[i][j] = make([]int, 0)
		}
	}
	expected[1][1] = []int{0, 6}
	expected[1][2] = []int{1}
	expected[2][4] = []int{2, 7}
	expected[3][1] = []int{4}
	expected[4][0] = []int{3}
	expected[6][1] = []int{9}
	expected[7][0] = []int{10}
	expected[7][1] = []int{5}
	expected[7][3] = []int{8}

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
}
