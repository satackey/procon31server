package field

import (
	"fmt"

	"github.com/golang-collections/collections/stack"
	"github.com/satackey/procon31server/pkg/apispec"
)

// Field はフィールド情報を表します
type Field struct {
	Width           int
	Height          int
	Turn            int
	Cells           [][]*Cell
	Agents          map[int]*Agent
	Teams           []*Team
	ActionHistories []ActionHistory
}

// ActionHistory は各エージェントの行動履歴を表します
type ActionHistory struct {
	AgentActionHistories []AgentActionHistory
}

// AgentActionHistory はエージェントの行動履歴を表します
type AgentActionHistory struct {
	AgentID int
	DX      int
	DY      int
	X       int
	Y       int
	Type    string
	Turn    int
	Apply   int
}

// UpdateAction2 は行動情報がどのチームによるものなのかを表します
type UpdateAction2 struct {
	*apispec.UpdateAction
	TeamID int
}

// New は初期化された Field を返します
func New() *Field {
	return &Field{
		Agents: map[int]*Agent{},
	}
}

// InitField はフィールド情報を登録します
func (f *Field) InitField(fieldStatus *apispec.FieldStatus) {
	f.Width = fieldStatus.Width
	f.Height = fieldStatus.Height
	f.Turn = fieldStatus.Turn

	f.Cells = make([][]*Cell, f.Height)
	for y, fieldRow := range fieldStatus.Points {
		f.Cells[y] = make([]*Cell, f.Width)
		for x, fieldColumn := range fieldRow {
			cell := newCell(fieldColumn, fieldStatus.Cells[y][x].TeamID, fieldStatus.Cells[y][x].Status, x, y, f)
			f.Cells[y][x] = cell
		}
	}

	for _, team := range fieldStatus.Teams {
		fmt.Printf("\nTeams %+v\n", team)
		teamData := &Team{
			ID: team.TeamID,
		}
		f.Teams = append(f.Teams, teamData)
		for _, agent := range team.Agents {
			fmt.Printf("\nAgent: %d\n", agent.AgentID)
			f.Agents[agent.AgentID] = &Agent{
				ID:     agent.AgentID,
				TeamID: team.TeamID,
				X:      agent.X,
				Y:      agent.Y,
				field:  f,
			}
		}
	}

}

// CalcPoint は現在のフィールドの指定されたチームIDの得点を計算します
func (f *Field) CalcPoint(teamID int) int {
	return f.CalcWallPoint(teamID) + f.CalcAreaPoint(teamID)
}

// CalcWallPoint は現在のフィールドの指定されたチームIDの城壁の点数を計算します
func (f *Field) CalcWallPoint(teamID int) int {
	sum := 0
	for _ /*y*/, fieldRow := range f.Cells {
		for _ /*x*/, cell := range fieldRow {
			if cell.TiledBy == teamID {
				sum += cell.Point
			}
		}
	}
	return sum
}

/*
//stackを使うなら再帰関数は使わないほうが安全…
func (f Field) dfs(vx int, vy int) {
	var seen [f.Width][f.Height]bool{{}}
	seen[vx][vy] = true

	for {
		if (seen[nextVx][nextVy]) {
			continue
		}

	}
}
*/

// SearchTiledBy は指定された座標の tiledby を返します
// ある座標がわかっているときにそのtiledbyがわかる関数が欲しい！！
func (f *Field) SearchTiledBy(x, y int) int {
	var res int
	res = f.Cells[y][x].TiledBy
	// いい感じの処理、一行で書けちゃった（笑）
	return res
}

// CalcAreaPoint は 現在のフィールドの指定されたチームIDの城壁の点数を計算します
func (f *Field) CalcAreaPoint(teamID int) int {
	Sum := 0

	seen := [][]bool{{}}
	todo := stack.New()

	//右、下、左、上
	dx := [4]int{1, 0, -1, 0}
	dy := [4]int{0, 1, 0, -1}

	for y, fieldRow := range f.Cells {
		for x, cell := range fieldRow {
			IsAreaPoint := true // 今見ている連結成分がエリアポイントかどうか。falseになった時点でどこかのセルが外側のセルである
			SumKari := 0        // 今見ているセルを含む連結成分のPointの合計値を保存する変数

			if cell.TiledBy == teamID {
				//タイルが自陣か否か
			} else {
				if seen[y][x] {
					// ここに入ったらDFSはしない、みたことある
					continue
				}
				// ここに入ったらDFSをする、みたことない
				todo.Push([2]int{x, y})

				for todo.Len() != 0 { // スタックの中身があるならDFSを続ける　0にならない限り
					v := todo.Pop().([2]int)
					SumKari += f.Cells[v[0]][v[1]].Point

					//外か内かの判定
					if v[0] == 0 || v[0] == f.Width-1 || v[1] == 0 || v[1] == f.Height-1 {
						IsAreaPoint = false
					}

					for index := 0; index < 4; index++ {
						// x+dxが範囲外かどうか(範囲外ならcontinue), yについても同様
						if x+dx[index] < 0 || x+dx[index] >= f.Width {
							continue
						}
						if y+dy[index] < 0 || y+dy[index] >= f.Height {
							continue
						}

						// 移動先が相手のタイルならcontinue
						if f.SearchTiledBy(x+dx[index], y+dy[index]) != teamID || f.SearchTiledBy(x+dx[index], y+dy[index]) != -1 {
							//移動先.TiledBy は 敵のIDになる
							continue
						}

						// seen[x+dx][y+dy] が true なら既にみているからcontinue
						if seen[x+dx[index]][y+dy[index]] {
							continue
						}

						// seen[x+dx][y+dy] に true を代入
						seen[x+dx[index]][y+dy[index]] = true

						//v[0]:x  v[1]:y todoの配列
						todo.Push([2]int{v[0] + dx[index], v[1] + dy[index]})
					}
				}
			}
			// SumKari の値が確定！！
			if IsAreaPoint {
				Sum += SumKari
			}

		}
	}
	return Sum
}

// MakeUpdateAction2s は updateActions と updateActionIDs をまとめて返します
func (f *Field) MakeUpdateAction2s(updateActions []*apispec.UpdateAction, updateActionIDs []int) []*UpdateAction2 {
	updateAction2s := make([]*UpdateAction2, len(updateActions))
	for i := range updateActions {
		updateAction2s[i] = &UpdateAction2{
			TeamID:       updateActionIDs[i],
			UpdateAction: updateActions[i],
		}
	}
	return updateAction2s
}

// RecordCellSelectedAgents は各セルを行動先に選んでいるような行動情報の要素番号を記録します
func (f *Field) RecordCellSelectedAgents(isValid []bool, updateActions []*apispec.UpdateAction) [][][]int {

	selectedAgents := make([][][]int, f.Height)
	for i := range selectedAgents {
		selectedAgents[i] = make([][]int, f.Width)
	}

	for i, updateAction := range updateActions {
		if isValid[i] {
			var x, y int
			if updateActions[i].Type == "put" {
				x = updateAction.X
				y = updateAction.Y
			} else {
				x = f.Agents[updateAction.AgentID].X + updateAction.DX
				y = f.Agents[updateAction.AgentID].Y + updateAction.DY
			}
			selectedAgents[y][x] = append(selectedAgents[y][x], i)
		}

	}

	return selectedAgents
}

// GiveNewStack は競合しているマスが入れられたスタックのポインタを返します
func (f *Field) GiveNewStack(selectedAgentsIndex [][][]int) *stack.Stack {
	stk := stack.New()
	for y := range selectedAgentsIndex {
		for x := range selectedAgentsIndex[y] {
			if len(selectedAgentsIndex[y][x]) >= 2 {
				for i := range selectedAgentsIndex[y][x] {
					stk.Push(selectedAgentsIndex[y][x][i])
				}
			}
		}
	}
	return stk
}

// DetermineIfApplied は 行動情報が競合か許容か不正かを判定して isApplicable に保存します
func (f *Field) DetermineIfApplied(isValid []bool, updateActions []*apispec.UpdateAction, selectedAgentsIndex [][][]int) []int {
	// isApplicable を初期化
	isApplicable := make([]int, len(updateActions))
	for i := range isApplicable {
		isApplicable[i] = 1
	}

	// 競合しているセルと、そのセルを選んでいるエージェントがいるセルには行けません
	// 競合しているセルをstackに突っ込む
	stk := f.GiveNewStack(selectedAgentsIndex)

	// 不正行動は先にはじいておく
	for i := range isValid {
		if isValid[i] == false {
			isApplicable[i] = -1
		}
	}

	// stackから出したセルを行動先に選んでいるセル
	for stk.Len() != 0 {
		updateActionIndex := stk.Pop().(int)
		if isApplicable[updateActionIndex] == -1 {
			fmt.Printf("yabeeeeeeee\n")
		}
		isApplicable[updateActionIndex] = 0
		if updateActions[updateActionIndex].Type != "put" {
			x := f.Agents[updateActions[updateActionIndex].AgentID].X
			y := f.Agents[updateActions[updateActionIndex].AgentID].Y
			for i := range selectedAgentsIndex[y][x] {
				stk.Push(selectedAgentsIndex[y][x][i])
			}
		}
	}

	return isApplicable
}

// ConvertIntoHistory は エージェント1体の行動情報を行動履歴に変換します
func (f *Field) ConvertIntoHistory(isValid bool, updateAction *apispec.UpdateAction, isApplicable int) AgentActionHistory {
	agentActionHistory := AgentActionHistory{
		AgentID: updateAction.AgentID,
		Type:    updateAction.Type,
		Turn:    f.Turn + 1,
		Apply:   isApplicable,
	}
	if updateAction.Type == "put" {
		agentActionHistory.X = updateAction.X
		agentActionHistory.Y = updateAction.Y
	} else {
		agentActionHistory.DX = updateAction.DX
		agentActionHistory.DY = updateAction.DY
	}

	return agentActionHistory
}

// ActuallyActAgent は マジで行動情報に基づいてフィールド情報を更新します
func (f *Field) ActuallyActAgent(updateAction2 *UpdateAction2) {
	switch updateAction2.Type {
	case "move":
		f.ActMove(updateAction2)
	case "remove":
		f.ActRemove(updateAction2)
	case "stay":
		f.ActStay(updateAction2)
	case "put":
		f.ActPut(updateAction2)
	}
}

// ActMove は type = "move" のとき ActuallyActAgent により実行されます
func (f *Field) ActMove(updateAction2 *UpdateAction2) {
	// 移動先のx, y座標を取得する
	x := f.Agents[updateAction2.AgentID].X + updateAction2.DX
	y := f.Agents[updateAction2.AgentID].Y + updateAction2.DY
	// エージェントの座標を変える
	f.Agents[updateAction2.AgentID].X = x
	f.Agents[updateAction2.AgentID].Y = y
	// 移動先の座標を自陣の城壁に変える
	f.Cells[y][x].TiledBy = f.Agents[updateAction2.AgentID].TeamID
	f.Cells[y][x].Status = "wall"
}

// ActRemove は type = "remove" のとき ActuallyActAgent により実行されます
func (f *Field) ActRemove(updateAction2 *UpdateAction2) {
	// 移動先のx, y座標を取得する
	x := f.Agents[updateAction2.AgentID].X + updateAction2.DX
	y := f.Agents[updateAction2.AgentID].Y + updateAction2.DY
	// 城壁 (wall) を除去する、つまりfreeに…
	// そうはいかないわ！私は怪人ジンチー。除去されたセルが囲われている場合、陣地にするわ！
	// 後回し！！！！！！！！！！！！！

	// 考察した結果、除去されたセルは仮にfreeとしておき、全ての行動を適用した後にareaになるかどうか計算すればいい！！！
	// 怪人ジンチー、死亡…😢

	// wallを除去し、freeにする
	f.Cells[y][x].TiledBy = 0
	f.Cells[y][x].Status = "free"
}

// ActStay は type = "stay" のとき ActuallyActAgent により実行されます
func (f *Field) ActStay(updateAction2 *UpdateAction2) {
	// 特に判定することもない
	// Q.何故関数化した？ A.見栄えがいいから
}

// ActPut は type = "put" のとき ActuallyActAgent により実行されます
func (f *Field) ActPut(updateAction2 *UpdateAction2) {
	// 移動先のx, y座標を取得する
	x := updateAction2.X
	y := updateAction2.Y
	// 配置される新しいエージェントの情報を作り、その情報をフィールドに保存する
	// newAgentID の決め方を考えよう
	// newAgentID は 現在存在するIDをインクリメントしていくとき存在してなかったIDにする
	newAgentID := 1
	for {
		_, keyExist := f.Agents[newAgentID]
		if !keyExist {
			break
		}
		newAgentID++
	}

	f.Agents[newAgentID] = &Agent{
		ID:     newAgentID,
		TeamID: updateAction2.TeamID,
		X:      x,
		Y:      y,
		field:  f,
	}
}

// IsOutsideField は与えられた座標がフィールドの外側にあるならtrueを返します
func (f *Field) IsOutsideField(x int, y int) bool {
	if x < 0 || x >= f.Width {
		return true
	}
	if y < 0 || y >= f.Height {
		return true
	}
	return false
}

// IsWallByteamIDOrSeen はteamIDの壁の中、もしくは既に見ているマスならtrueを返します
func (f *Field) IsWallByteamIDOrSeen(x int, y int, teamID int, seen [][]bool) bool {
	if f.Cells[y][x].Status == "wall" && f.Cells[y][x].TiledBy == teamID {
		return true
	}
	if seen[y][x] {
		return true
	}
	return false
}

// CheckAreaByDFS はteamIDの城壁の内部にあるセルを記録して返し、囲われていればtrueを返します
func (f *Field) CheckAreaByDFS(teamID int, startX int, startY int) ([][]bool, bool) {
	seen := make([][]bool, f.Height)
	for i := 0; i < f.Height; i ++ {
		seen[i] = make([]bool, f.Width)
		for j := 0; j < f.Width; j ++ {
			seen[i][j] = false
		}
	}
	dx := []int{1, 1, 0, -1, -1, -1, 0, 1}
	dy := []int{0, 1, 1, 1, 0, -1, -1, -1}
	st := stack.New()
	st.Push([]int{startX, startY})
	seen[startY][startX] = true
	for st.Len() != 0 {
		xy := st.Pop().([]int)
		x := xy[0]
		y := xy[1]
		for i := 0; i < 8; i ++ {
			if f.IsOutsideField(x + dx[i], y + dy[i]) {
				return nil, false
			}
			if f.IsWallByteamIDOrSeen(x+dx[i], y+dy[i], teamID, seen) {
				continue
			}
			st.Push([]int{x+dx[i], y+dy[i]})
			seen[y+dy[i]][x+dx[i]] = true
		}
	}
	return seen, true
}

// FinalCheckByDFS はSurroundedBy[1]で囲まれた中でSurroundedBy[0]によって囲めるか？をDFSで判定し、
// (x, y)をより内側で囲っているteamのIDを返します
func (f *Field) FinalCheckByDFS(surroundedBy []int, startX int, startY int, isAreaBy map[int][][]bool) int {
	teamID := surroundedBy[0]
	dx := []int{1, 1, 0, -1, -1, -1, 0, 1}
	dy := []int{0, 1, 1, 1, 0, -1, -1, -1}
	seen := make([][]bool, f.Height)
	for i := 0; i < f.Height; i ++ {
		seen[i] = make([]bool, f.Width)
		for j := 0; j < f.Width; j ++ {
			seen[i][j] = false
		}
	}

	st := stack.New()
	st.Push([]int{startX, startY})
	seen[startY][startX] = true
	for st.Len() != 0 {
		xy := st.Pop().([]int)
		x := xy[0]
		y := xy[1]
		for i := 0; i < 8; i ++ {
			if  f.IsOutsideField(x + dx[i], y + dy[i]) || !isAreaBy[surroundedBy[1]][y+dy[i]][x+dx[i]]  {
				// [1]の内側に[0]があるならこの条件を満たさないﾊｽﾞ。
				// が、満たしたのだから、仮定が偽
				return surroundedBy[1]
			}
			if f.IsWallByteamIDOrSeen(x+dx[i], y+dy[i], teamID, seen) {
				continue
			}
			st.Push([]int{x+dx[i], y+dy[i]})
			seen[y+dy[i]][x+dx[i]] = true
		}
	}

	return surroundedBy[0]
}

// IsWallOrSeen は壁の中、もしくは既に見ているマスならtrueを返します
func (f *Field) IsWallOrSeen(x int, y int, seen [][]bool) bool {
	if f.Cells[y][x].Status == "wall" {
		return true
	}
	if seen[y][x] {
		return true
	}
	return false
}

// ChangeCellToPositionByDFS は城壁を超えない範囲の連結成分をすべてteamIDの陣地にし、seenを更新します
// teamID == -1 ならなにもしません
func (f *Field) ChangeCellToPositionByDFS(teamID int, startX int, startY int, seen *[][]bool) {
	dx := []int{1, 1, 0, -1, -1, -1, 0, 1}
	dy := []int{0, 1, 1, 1, 0, -1, -1, -1}

	st := stack.New()
	st.Push([]int{startX, startY})
	(*seen)[startY][startX] = true
	for st.Len() != 0 {
		xy := st.Pop().([]int)
		x := xy[0]
		y := xy[1]

		f.Cells[y][x].TiledBy = teamID
		if f.Cells[y][x].TiledBy != 0 {
			f.Cells[y][x].Status = "position"
		}
		
		for i := 0; i < 8; i ++ {
			if f.IsOutsideField(x + dx[i], y + dy[i]) || f.IsWallOrSeen(x+dx[i], y+dy[i], *seen) {
				continue
			}
			st.Push([]int{x+dx[i], y+dy[i]})
			(*seen)[y+dy[i]][x+dx[i]] = true
		}
	}
}

// SurroundedByWoHenkou はセル(x, y)を囲っているチームIDと、実際にどこのセルが囲まれているかを記録して返します
func (f *Field) SurroundedByWoHenkou(startX int, startY int) ([]int, map[int][][]bool) {
	// isAreaBy[ID][Y][X] := 座標 (X, Y) が TeamID による城壁で囲まれたエリアか
	isAreaBy := map[int][][]bool{}
	// (x, y) を囲んでいる城壁のteamIDのスライス
	surroundedBy := []int{}

	for _, team := range f.Teams {
		isAreaBy[team.ID] = make([][]bool, f.Height)
		for y := 0; y < f.Height; y ++ {
			isAreaBy[team.ID][y] = make([]bool, f.Width)
		}
		seen, ok := f.CheckAreaByDFS(team.ID, startX, startY);
		if ok {
			// team.IDによって(x, y)は囲われている！
			isAreaBy[team.ID] = seen
			surroundedBy = append(surroundedBy, team.ID)
		}
	}

	return surroundedBy, isAreaBy
}

// CleanUpCellsFormerlyWall は以前は壁だった細胞を片付けます
// (f.Cells[y][x].Status が free になるか position になるか決めます)
func (f *Field) CleanUpCellsFormerlyWall() {
	seen := make([][]bool, f.Height)
	for y := range seen {
		seen[y] = make([]bool, f.Width)
		for x := range seen[y] {
			seen[y][x] = false
		}
	}

	for y := range seen {
		for x := range seen[y] {
			if f.IsWallOrSeen(x, y, seen) {
				continue
			}
			// 各チームの城壁で囲まれているかチェック
			surroundedBy, isAreaBy := f.SurroundedByWoHenkou(x, y)

			var teamID int

			// どちらのチームにも囲まれていないのなら変更しない。
			// 1チームにしか囲われていないのならそのチームの陣地である。
			// 2チームともに囲まれているのなら片方の領域内で囲めるかどうかで決める。

			switch len(surroundedBy) {
			case 0:
				teamID = 0
			case 1:
				teamID = surroundedBy[0]
			case 2:
				teamID = f.FinalCheckByDFS(surroundedBy, x, y, isAreaBy)
			}
			// 連結成分をすべてtrueにする、上記の通り変更する
			f.ChangeCellToPositionByDFS(teamID, x, y, &seen)
		}
	}
}

// ActAgents はエージェントの行動に基づいてフィールドを変更し、履歴を保存します。
func (f *Field) ActAgents(isValid []bool, updateActions []*apispec.UpdateAction, updateActionIDs []int) {
	// updateActionIDs []int をもらって
	// updateAction2s []*UpdateAction2 をつくる
	updateAction2s := f.MakeUpdateAction2s(updateActions, updateActionIDs)

	// 行動を精査します
	// もうやったのでIsValidは信用していいデータらしい。

	// セルが選ばれた回数　ではなく、そのセルを選んでいる行動情報の要素番号をスライスにして保存
	// selectedAgentsIndex[y][x] := (x, y)を移動先に選んでいる行動情報の要素番号の配列
	selectedAgentsIndex := f.RecordCellSelectedAgents(isValid, updateActions)

	// DistinationCount と IsValid に基づいて apply が決定
	// i番目のupdateActionが許容か競合か不正か
	isApplicable := f.DetermineIfApplied(isValid, updateActions, selectedAgentsIndex)
	// []AgentActionHistoryつくる
	agentActionHistories := make([]AgentActionHistory, len(updateActions))
	// 各updateActionに対して
	for i, updateAction := range updateActions {
		// updateaction -> []AgentActionHistry に変換して代入
		agentActionHistories[i] = f.ConvertIntoHistory(isValid[i], updateAction, isApplicable[i])
		// apply == 1 なら実際に動かす
		if agentActionHistories[i].Apply == 1 {
			f.ActuallyActAgent(updateAction2s[i])
		}
	}
	// removeなどによって除去された城壁があるので、エリアを再計算して正しいセルの状態にする
	f.CleanUpCellsFormerlyWall()
	// f.ActionHistories[i].AgentActionHistories に agentActionHistories を代入
	// もし0ターン目ならActionHistories[0]は使わないので空けておく
	if f.Turn == 0 && len(f.ActionHistories) == 0 {
		f.ActionHistories = append(f.ActionHistories, ActionHistory{})
	}
	// f.Turn+1 ターン目の行動情報を記録する
	f.ActionHistories = append(f.ActionHistories, ActionHistory{AgentActionHistories: agentActionHistories})
}

// GetFieldEasyToSee は見やすいフィールド情報を返します
func (f Field) GetFieldEasyToSee() [][]string {
	resField := make([][]string, len(f.Cells))

	for iRow, rowCells := range f.Cells {
		resField[iRow] = make([]string, len(f.Cells[iRow]))
		for iCol, cell := range rowCells {
			agent, err := f.GetAgent(cell.x, cell.y)
			if err == nil {
				if agent.TeamID == f.Teams[0].ID {
					resField[iRow][iCol] = fmt.Sprintf(" (%3d) ", cell.Point)
				} else if agent.TeamID == f.Teams[1].ID {
					resField[iRow][iCol] = fmt.Sprintf(" [%3d] ", cell.Point)
				} else {
					resField[iRow][iCol] = fmt.Sprintf("  %3d  ", cell.Point)
				}
			}
		}
	}
	return resField
}

// GetAgent は指定された座標のAgentを取得します
func (f Field) GetAgent(x, y int) (*Agent, error) {

	for _, agent := range f.Agents {
		if agent.X == x && agent.Y == y {
			return agent, nil
		}
	}

	return &Agent{}, fmt.Errorf("エージェントがみつかりません")
}

// CheckIfAgentsInfoIsValid は行動情報全体が有効かどうか返します
func (f *Field) CheckIfAgentsInfoIsValid(updateActions []*apispec.UpdateAction) bool {
	if len(f.Agents) != len(updateActions) {
		// 送信されてきたデータ長が違う　出直してこい
		return false
	}

	for _, updateAction := range updateActions {
		if _, ok := f.Agents[updateAction.AgentID]; ok == false {
			// AgentID が不正
			return false
		}
	}

	for _, updateAction1 := range updateActions {
		for _, updateAction2 := range updateActions {
			if updateAction1.AgentID == updateAction2.AgentID {
				// ID がかぶっている
				return false
			}
		}
	}

	// ここまで来れば完璧です by さたけ
	return true
}

// CheckIfAgentInfoIsValid は行動情報一つ一つが有効か判定します
func (f *Field) CheckIfAgentInfoIsValid(teamID int, updateActions []*apispec.UpdateAction) (res []bool) {
	res = make([]bool, len(updateActions))

	for index, updateAction := range updateActions {
		NextX := f.Agents[updateAction.AgentID].X + updateAction.DX
		NextY := f.Agents[updateAction.AgentID].Y + updateAction.DY

		if _, ok := f.Agents[updateAction.AgentID]; !ok {
			// 指定された AgentID は存在しない　想定外　論外
			res[index] = false
			continue
		} else if updateAction.DX != -1 && updateAction.DX != 0 && updateAction.DX != 1 {
			// DX の値が不正　瞬間移動はできない。
			res[index] = false
			continue
		} else if updateAction.DY != -1 && updateAction.DY != 0 && updateAction.DY != 1 {
			// DY の値が不正　瞬間移動はできない。
			res[index] = false
			continue
		} else if NextX < 0 || NextX >= f.Width || NextY < 0 || NextY >= f.Height {
			// 移動先に指定した場所は範囲外　異世界に飛ぶ気か？
			res[index] = false
			continue
		} else if updateAction.Type == "move" {
			if f.Cells[NextX][NextY].TiledBy != teamID && f.Cells[NextX][NextY].TiledBy != 0 {
				// 移動先に指定したセルに敵のタイルがあって動けない！！！
				res[index] = false
				continue
			}
		} else if updateAction.Type == "remove" {
			if f.Cells[NextX][NextY].TiledBy == teamID || f.Cells[NextX][NextY].TiledBy == 0 {
				// 移動先に指定したセルに敵のタイルはない！！！
				res[index] = false
				continue
			}
		} else if updateAction.Type == "stay" {
			// "stay" で衝突する場合はないので true
		} else {
			// upgateAction.Type の文字列が意味不明　そんなものは存在しない
			res[index] = false
			continue
		}

		// ここまで到達したデータに不正はないので true
		res[index] = true
	}
	return
}
