package field

import (
	"fmt"
	"math"

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
	TeamID int
}

type UpdateActionWithoutType struct {
	AgentID int `json:"agentID"`
}

type UpdateActionAbsoluteOnly struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type UpdateActionAbsolute struct {
	UpdateActionAbsoluteOnly
	UpdateActionWithoutType
}

type UpdateActionRelativeOnly struct {
	DX int `json:"dx"`
	DY int `json:"dy"`
}

type UpdateActionRelative struct {
	UpdateActionWithoutType
	UpdateActionRelativeOnly
}

// UpdateAction は行動情報 を表します
type UpdateAction struct {
	UpdateActionWithoutType
	UpdateActionAbsoluteOnly
	UpdateActionRelativeOnly
	Type string `json:"type"`
}

type UpdateAction3 interface {
	Act(*Field)
	IsValid(*Field) bool
	GetUpdateAction() *UpdateAction2
}

type PutUpdateAction struct {
	UpdateAction2
	UpdateActionAbsolute
}

type MoveUpdateAction struct {
	UpdateAction2
	UpdateActionRelative
}

type StayUpdateAction struct {
	UpdateAction2
	UpdateActionRelative
}

type RemoveUpdateAction struct {
	UpdateAction2
	UpdateActionRelative
}

// Act は
func (p *PutUpdateAction) Act(field *Field) {
	// 移動先のx, y座標を取得する
	x := p.X
	y := p.Y
	// 配置される新しいエージェントの情報を作り、その情報をフィールドに保存する
	// newAgentID の決め方を考えよう
	// newAgentID は 現在存在するIDをインクリメントしていくとき存在してなかったIDにする
	newAgentID := 1
	for {
		_, keyExist := field.Agents[newAgentID]
		if !keyExist {
			break
		}
		newAgentID++
	}

	field.Agents[newAgentID] = &Agent{
		ID:     newAgentID,
		TeamID: p.TeamID,
		X:      x,
		Y:      y,
		field:  field,
	}
}

// GetUpdateAction は
func (p *PutUpdateAction) GetUpdateAction() *UpdateAction2 {
	return &p.UpdateAction2
}

// IsValid は
func (p *PutUpdateAction) IsValid(field *Field) bool {
	// 移動先のセルは範囲外か(nextX, nextY)
	if field.IsOutsideField(nextX, nextY) {
		return false
	}
	// 移動先は敵の城壁か(nextX, nextY, f.Agents[updateAction.AgentID].TeamID)
	if field.IsOpponentWall(nextX, nextY, field.Agents[p.AgentID].TeamID) {
		return false
	}

	return true
}

// Act は
func (m *MoveUpdateAction) Act(field *Field) {
	// 移動先のx, y座標を取得する
	x := field.Agents[m.AgentID].X + m.DX
	y := field.Agents[m.AgentID].Y + m.DY
	// エージェントの座標を変える
	field.Agents[m.AgentID].X = x
	field.Agents[m.AgentID].Y = y
	// 移動先の座標を自陣の城壁に変える
	field.Cells[y][x].TiledBy = field.Agents[m.AgentID].TeamID
	field.Cells[y][x].Status = "wall"
}

// GetUpdateAction は
func (m *MoveUpdateAction) GetUpdateAction() *UpdateAction2 {
	return &m.UpdateAction2
}

// IsValid は
func (m *MoveUpdateAction) IsValid(field *Field) bool {
	// DX, DYの値は正常か(updateAction.DX, updateAction.DY)
	if !field.IsDXDYValidValue(m.DX, m.DY) {
		return false
	}
	// 移動先のセルは範囲外か(nextX, nextY)
	if field.IsOutsideField(nextX, nextY) {
		return false
	}
	// 移動先は敵の城壁か(nextX, nextY, f.Agents[updateAction.AgentID].TeamID)
	if field.IsOpponentWall(nextX, nextY, field.Agents[m.AgentID].TeamID) {
		return false
	}

	return true
}

// Act は
func (s *StayUpdateAction) Act(field *Field) {
	field.ActStay(s.GetUpdateAction())
}

// GetUpdateAction は
func (s *StayUpdateAction) GetUpdateAction() *UpdateAction2 {
	return &s.UpdateAction2
}

// IsValid は
func (s *StayUpdateAction) IsValid(field *Field) bool {
	if !field.IsZero(nextX, nextY) {
		return false
	}

	return true
}

// Act は
func (r *RemoveUpdateAction) Act(field *Field) {
	// 移動先のx, y座標を取得する
	x := field.Agents[r.AgentID].X + r.DX
	y := field.Agents[r.AgentID].Y + r.DY
	// 城壁 (wall) を除去する、つまりfreeに…
	// そうはいかないわ！私は怪人ジンチー。除去されたセルが囲われている場合、陣地にするわ！
	// 後回し！！！！！！！！！！！！！

	// 考察した結果、除去されたセルは仮にfreeとしておき、全ての行動を適用した後にareaになるかどうか計算すればいい！！！
	// 怪人ジンチー、死亡…😢

	// wallを除去し、freeにする
	field.Cells[y][x].TiledBy = 0
	field.Cells[y][x].Status = "free"
	// field.ActRemove(p.GetUpdateAction())
}

// GetUpdateAction は
func (r *RemoveUpdateAction) GetUpdateAction() *UpdateAction2 {
	return &r.UpdateAction2
}

// IsValid は
func (r *RemoveUpdateAction) IsValid(field *Field) bool {
	if !field.IsDXDYValidValue(r.DX, r.DY) {
		return false
	}
	// 移動先のセルは範囲外か(nextX, nextY)
	if field.IsOutsideField(nextX, nextY) {
		return false
	}
	// 移動先は城壁か(nextX, nextY)
	if !field.IsWall(nextX, nextY) {
		return false
	}

	return true
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
// かならず ActAgents() を実行後に実行すること
func (f *Field) CalcPoint(teamID int) int {
	sum := 0
	for y := range f.Cells {
		for x := range f.Cells[y] {
			if f.Cells[y][x].TiledBy == teamID {
				if f.Cells[y][x].Status == "wall" {
					sum += f.Cells[y][x].Point
				} else if f.Cells[y][x].Status == "position" {
					sum += int(math.Abs(float64(f.Cells[y][x].Point)))
				}
			}
		}
	}
	return sum
}

// MakeUpdateAction2s は updateActions と updateActionIDs をまとめて返します
func (f *Field) MakeUpdateAction2s(updateActions []*apispec.UpdateAction, updateActionIDs []int) []UpdateAction3 {
	result := make([]UpdateAction3, len(updateActions))
	for i, action := range updateActions {
		result[i] = CreateUpdateAction(action, updateActionIDs[i])
	}
	return result
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

	// Todo: Interface を使う
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
func (f *Field) ActuallyActAgent(action UpdateAction3) {
	action.Act(f)
}

// ActMove は type = "move" のとき ActuallyActAgent により実行されます
func (f *Field) ActMove(updateAction2 *UpdateAction2) {
}

// ActRemove は type = "remove" のとき ActuallyActAgent により実行されます
// func (f *Field) ActRemove(updateAction2 *UpdateAction2) {
// 	// RemoveUpdateAction に移動された　ActRemove は無職。
// }

// ActStay は type = "stay" のとき ActuallyActAgent により実行されます
func (f *Field) ActStay(updateAction2 *UpdateAction2) {
	// 特に判定することもない
	// Q.何故関数化した？ A.見栄えがいいから
}

// ActPut は type = "put" のとき ActuallyActAgent により実行されます
func (f *Field) ActPut(updateAction2 *UpdateAction2) {

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
	for i := 0; i < f.Height; i++ {
		seen[i] = make([]bool, f.Width)
		for j := 0; j < f.Width; j++ {
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
		for i := 0; i < 8; i++ {
			if f.IsOutsideField(x+dx[i], y+dy[i]) {
				return nil, false
			}
			if f.IsWallByteamIDOrSeen(x+dx[i], y+dy[i], teamID, seen) {
				continue
			}
			st.Push([]int{x + dx[i], y + dy[i]})
			seen[y+dy[i]][x+dx[i]] = true
		}
	}
	return seen, true
}

// FinalCheckByDFS はenclosedBy[1]で囲まれた中でenclosedBy[0]によって囲めるか？をDFSで判定し、
// (x, y)をより内側で囲っているteamのIDを返します
func (f *Field) FinalCheckByDFS(enclosedBy []int, startX int, startY int, isAreaBy map[int][][]bool) int {
	teamID := enclosedBy[0]
	dx := []int{1, 1, 0, -1, -1, -1, 0, 1}
	dy := []int{0, 1, 1, 1, 0, -1, -1, -1}
	seen := make([][]bool, f.Height)
	for i := 0; i < f.Height; i++ {
		seen[i] = make([]bool, f.Width)
		for j := 0; j < f.Width; j++ {
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
		for i := 0; i < 8; i++ {
			if f.IsOutsideField(x+dx[i], y+dy[i]) || !isAreaBy[enclosedBy[1]][y+dy[i]][x+dx[i]] {
				// [1]の内側に[0]があるならこの条件を満たさないﾊｽﾞ。
				// が、満たしたのだから、仮定が偽
				return enclosedBy[1]
			}
			if f.IsWallByteamIDOrSeen(x+dx[i], y+dy[i], teamID, seen) {
				continue
			}
			st.Push([]int{x + dx[i], y + dy[i]})
			seen[y+dy[i]][x+dx[i]] = true
		}
	}

	return enclosedBy[0]
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

		for i := 0; i < 8; i++ {
			if f.IsOutsideField(x+dx[i], y+dy[i]) || f.IsWallOrSeen(x+dx[i], y+dy[i], *seen) {
				continue
			}
			st.Push([]int{x + dx[i], y + dy[i]})
			(*seen)[y+dy[i]][x+dx[i]] = true
		}
	}
}

// EnclosedByWoHenkou はセル(x, y)を囲っているチームIDと、実際にどこのセルが囲まれているかを記録して返します
func (f *Field) EnclosedByWoHenkou(startX int, startY int) ([]int, map[int][][]bool) {
	// isAreaBy[ID][Y][X] := 座標 (X, Y) が TeamID による城壁で囲まれたエリアか
	isAreaBy := map[int][][]bool{}
	// (x, y) を囲んでいる城壁のteamIDのスライス
	enclosedBy := []int{}

	for _, team := range f.Teams {
		isAreaBy[team.ID] = make([][]bool, f.Height)
		for y := 0; y < f.Height; y++ {
			isAreaBy[team.ID][y] = make([]bool, f.Width)
		}
		seen, ok := f.CheckAreaByDFS(team.ID, startX, startY)
		if ok {
			// team.IDによって(x, y)は囲われている！
			isAreaBy[team.ID] = seen
			enclosedBy = append(enclosedBy, team.ID)
		}
	}

	return enclosedBy, isAreaBy
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
			enclosedBy, isAreaBy := f.EnclosedByWoHenkou(x, y)

			var teamID int

			// どちらのチームにも囲まれていないのなら変更しない。
			// 1チームにしか囲われていないのならそのチームの陣地である。
			// 2チームともに囲まれているのなら片方の領域内で囲めるかどうかで決める。

			switch len(enclosedBy) {
			case 0:
				teamID = 0
			case 1:
				teamID = enclosedBy[0]
			case 2:
				teamID = f.FinalCheckByDFS(enclosedBy, x, y, isAreaBy)
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
func (f *Field) CheckIfAgentInfoIsValid(updateActions []*apispec.UpdateAction) (res []bool) {
	res = make([]bool, len(updateActions))

	for index, updateAction := range updateActions {
		// 行動情報が正しいものと仮定する
		res[index] = true

		// AgentIDが存在するか
		if _, ok := f.Agents[updateAction.AgentID]; !ok {
			// AgentIDは存在しなかった
			res[index] = false
			continue
		}

		ua := CreateUpdateAction(updateAction, f.Agents[updateAction.AgentID].TeamID)
		if ua == nil {
			// Typeの文字列が不正
			res[index] = false
			continue
		}

		// 移動先の座標
		nextX, nextY := f.CalcAgentDestination(updateAction)

		// switch updateAction.Type {
		// case "stay":
		// 	updateAction.DX = 0
		// 	updateAction.DY = 0
		// 	updateAction.X = 0
		// 	updateAction.Y = 0
		// case "move":
		// 	updateAction.X = 0
		// 	updateAction.Y = 0
		// 	// DX, DYの値は正常か(updateAction.DX, updateAction.DY)
		// 	if !f.IsDXDYValidValue(updateAction.DX, updateAction.DY) {
		// 		res[index] = false
		// 		continue
		// 	}
		// 	// 移動先のセルは範囲外か(nextX, nextY)
		// 	if f.IsOutsideField(nextX, nextY) {
		// 		res[index] = false
		// 		continue
		// 	}
		// 	// 移動先は敵の城壁か(nextX, nextY, f.Agents[updateAction.AgentID].TeamID)
		// 	if f.IsOpponentWall(nextX, nextY, f.Agents[updateAction.AgentID].TeamID) {
		// 		res[index] = false
		// 		continue
		// 	}
		// case "remove":
		// 	updateAction.X = 0
		// 	updateAction.Y = 0
		// 	// DX, DYの値は正常か(updateAction.DX, updateAction.DY)
		// 	if !f.IsDXDYValidValue(updateAction.DX, updateAction.DY) {
		// 		res[index] = false
		// 		continue
		// 	}
		// 	// 移動先のセルは範囲外か(nextX, nextY)
		// 	if f.IsOutsideField(nextX, nextY) {
		// 		res[index] = false
		// 		continue
		// 	}
		// 	// 移動先は城壁か(nextX, nextY)
		// 	if !f.IsWall(nextX, nextY) {
		// 		res[index] = false
		// 		continue
		// 	}
		// case "put":
		// 	updateAction.DX = 0
		// 	updateAction.DY = 0
		// 	// 移動先のセルは範囲外か(nextX, nextY)
		// 	if f.IsOutsideField(nextX, nextY) {
		// 		res[index] = false
		// 		continue
		// 	}
		// 	// 移動先は敵の城壁か(nextX, nextY, f.Agents[updateAction.AgentID].TeamID)
		// 	if f.IsOpponentWall(nextX, nextY, f.Agents[updateAction.AgentID].TeamID) {
		// 		res[index] = false
		// 		continue
		// 	}
		// default:
		// 	// Typeの文字列がおかしい
		// 	res[index] = false
		// }

		res[index] = updateActionHoge.IsValid(f)

	}
	return
}

// CalcAgentDestination は行動情報が指し示す移動先の座標を返します
func (f *Field) CalcAgentDestination(updateAction *apispec.UpdateAction) (x int, y int) {
	x = f.Agents[updateAction.AgentID].X + updateAction.DX
	y = f.Agents[updateAction.AgentID].Y + updateAction.DY
	if updateAction.Type == "put" {
		x = updateAction.X
		y = updateAction.Y
	}
	return
}

// IsDXDYValidValue はDXとDYが有効な値であるならtrueを返します
func (f *Field) IsDXDYValidValue(DX int, DY int) bool {
	if DX != -1 && DX != 0 && DX != 1 {
		return false
	} else if DY != -1 && DY != 0 && DY != 1 {
		return false
	}
	return true
}

// IsOpponentWall は与えられたセルが自分以外の城壁ならtrueを返します
// 第3引数は「移動するエージェント自身のTeamID」である！「敵のTeamID」ではない！！！
func (f *Field) IsOpponentWall(x int, y int, myTeamID int) bool {
	if f.Cells[y][x].Status == "wall" && f.Cells[y][x].TiledBy != myTeamID {
		return true
	}
	return false
}

// IsWall は与えられたセルが城壁ならtrueを返します
func (f *Field) IsWall(x int, y int) bool {
	if f.Cells[y][x].Status == "wall" {
		return true
	}
	return false
}

// IsDXDYZero はDXDYが0であるならtrueを返します
func (f *Field) IsDXDYZero(dx int, dy int) bool {
	if dx == 0 && dy == 0 {
		return true
	}
	return false
}

func CreateUpdateAction(action *apispec.UpdateAction, teamID int) UpdateAction3 {
	switch action.Type {
	case "move":
		result := &MoveUpdateAction{}
		result.TeamID = teamID
		result.DX = action.DX
		result.DY = action.DY
		result.AgentID = action.AgentID
		return result

	case "remove":
		result := &RemoveUpdateAction{}
		result.TeamID = teamID
		result.DX = action.DX
		result.DY = action.DY
		result.AgentID = action.AgentID
		return result

	case "stay":
		result := &StayUpdateAction{}
		result.TeamID = teamID
		result.DX = action.DX
		result.DY = action.DY
		result.AgentID = action.AgentID
		return result

	case "put":
		result := &PutUpdateAction{}
		result.TeamID = teamID
		result.X = action.X
		result.Y = action.Y
		result.AgentID = action.AgentID
		return result
	}

	return nil
}
