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
	return f.CalcTilePoint(teamID) + f.CalcAreaPoint(teamID)
}

func (f *Field) CalcTilePoint(teamID int) int {
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

//ある座標がわかっているときにそのtiledbyがわかる関数が欲しい！！
func (f *Field) SearchTiledBy(x, y int) int {
	var res int
	res = f.Cells[y][x].TiledBy
	// いい感じの処理、一行で書けちゃった（笑）
	return res
}

//左上から調べていく…？
//外側でも探索してNGなセルをみつけておく
//エリア　中身囲まれている
//たいる　自陣のセル
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
				if seen[y][x] == true {
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

						// seen[x+dx][y+dy] == true なら既にみているからcontinue
						if seen[x+dx[index]][y+dy[index]] == true {
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
			if IsAreaPoint == true {
				Sum += SumKari
			}

		}
	}
	return Sum
}

// CellSelectedTimesCount は各セルが行動先に選ばれた回数を返します
func (f *Field) CellSelectedTimesCount(isValid []bool, updateActions []*apispec.UpdateAction) [][]int {
	var distinationCount [][]int

	for i, updateAction := range updateActions {
		if isValid[i] == true {
			var x, y int
			if updateActions[i].Type == "put" {
				x = updateAction.X
				y = updateAction.Y
			} else {
				x = f.Agents[updateAction.AgentID].X + updateAction.DX
				y = f.Agents[updateAction.AgentID].Y + updateAction.DY
			}
			distinationCount[x][y]++
		}
	}

	return distinationCount
}

// ConvertIntoHistory は エージェント1体の行動情報を、行動履歴に変換します。
func (f *Field) ConvertIntoHistory(isValid bool, updateAction *apispec.UpdateAction, distinationCount [][]int) AgentActionHistory {
	agentActionHistory := AgentActionHistory{
		AgentID: updateAction.AgentID,
		Type:    updateAction.Type,
		Turn:    f.Turn,
	}

	if updateAction.Type == "put" {
		agentActionHistory.X = updateAction.X
		agentActionHistory.Y = updateAction.Y
	} else {
		agentActionHistory.DX = updateAction.DX
		agentActionHistory.DY = updateAction.DY
	}

	if isValid == true {
		var x, y int
		if updateAction.Type == "put" {
			x = updateAction.X
			y = updateAction.Y
		} else {
			x = f.Agents[updateAction.AgentID].X + updateAction.DX
			y = f.Agents[updateAction.AgentID].Y + updateAction.DY
		}
		if distinationCount[y][x] == 1 {
			agentActionHistory.Apply = 1
		} else {
			agentActionHistory.Apply = 0
		}
	} else {
		agentActionHistory.Apply = -1
	}

	return agentActionHistory
}

// ActuallyActAgent は マジで行動情報に基づいてフィールド情報を更新します
func (f *Field) ActuallyActAgent(updateAction *apispec.UpdateAction) {
	switch updateAction.Type {
		case "move" :
			f.ActMove(updateAction);
		case "remove" :
			f.ActRemove(updateAction);
		case "stay" :
			f.ActStay(updateAction);
		case "put" :
			f.ActPut(updateAction);
	}
}

// ActMove は type = "move" のとき ActuallyActAgent により実行されます
func (f *Field) ActMove(updateAction *apispec.UpdateAction) {
	// 移動先のx, y座標を取得する
	x := f.Agents[updateAction.AgentID].X + updateAction.DX;
	y := f.Agents[updateAction.AgentID].Y + updateAction.DY;
	// エージェントの座標を変える
	f.Agents[updateAction.AgentID].X = x;
	f.Agents[updateAction.AgentID].Y = y;
	// 移動先の座標を自陣の城壁に変える
	f.Cells[y][x].TiledBy = f.Agents[updateAction.AgentID].TeamID;
	f.Cells[y][x].Status = "wall";
}

// ActRemove は type = "remove" のとき ActuallyActAgent により実行されます
func (f *Field) ActRemove(updateAction *apispec.UpdateAction) {
	// 移動先のx, y座標を取得する
	// x := f.Agents[updateAction.AgentID].X + updateAction.DX;
	// y := f.Agents[updateAction.AgentID].Y + updateAction.DY;
	// 城壁 (wall) を除去する、つまりfreeに…
	// そうはいかないわ！私は怪人ジンチー。除去されたセルが囲われている場合、陣地にするわ！
	// 後回し！！！！！！！！！！！！！
}

// ActStay は type = "stay" のとき ActuallyActAgent により実行されます
func (f *Field) ActStay(updateAction *apispec.UpdateAction) {
	// 特に判定することもない
	// Q.何故関数化した？ A.見栄えがいいから
}

// ActPut は type = "put" のとき ActuallyActAgent により実行されます
func (f *Field) ActPut(updateAction *apispec.UpdateAction) {
	
}

// ActAgents はエージェントの行動に基づいてフィールドを変更し、履歴を保存します。
func (f *Field) ActAgents(isValid []bool, updateActions []*apispec.UpdateAction) {

	// 行動を精査します
	// もうやったのでIsValidは信用していいデータらしい。

	// セルが行動先に選ばれた回数をカウントします
	distinationCount := f.CellSelectedTimesCount(isValid, updateActions)
	// 競合しているセルと、そのセルを選んでいるエージェントがいるセルには行けません
	// 二重ループで実装できそう。
	// 競合しているセルをqueueに突っ込む
	// queueから出したセルを行動先に選んでいるセル


	// 一時的にコメントアウト

	// セルに保存された回数が1回なら、実行できます(Apply:1)
	// Apply := 1
	// for i, updateAction := range updateActions {
	// 	if IsValid[i] == false {
	// 		var slise []AgentActionHistory
	// 		f.ActionHistories[i].AgentActionHistories = append(f.ActionHistories[i].AgentActionHistories, slise...) // sliseの中身を追加する
	// 		Apply = -1
	// 		// f.Turn
	// 		continue
	// 	}
	// 	x := f.Agents[updateAction.AgentID].X + updateAction.DX
	// 	y := f.Agents[updateAction.AgentID].Y + updateAction.DY
	// 	if DistinationCount[x][y] == 1 {
	// 		// 動ける、動かす
	// 		Apply = 1
	// 	} else {
	// 		// 競合して動けない、stayにする
	// 		Apply = 0
	// 	}
	// }

	// DistinationCount と IsValid に基づいて apply が決定

	// []AgentActionHistoryつくる
	var agentActionHistories []AgentActionHistory
	// 各updateActionに対して
	for i, updateAction := range updateActions {
		// updateaction -> []AgentActionHistry に変換して代入
		agentActionHistories[i] = f.ConvertIntoHistory(isValid[i], updateAction, distinationCount)
		// apply == 1 なら実際に動かす
		if agentActionHistories[i].Apply == 1 {
			f.ActuallyActAgent(updateAction)
		}
	}
	// f.ActionHistories[i].AgentActionHistories に agentActionHistories を代入
	// 何番目？0-indexedかなどうかな
	f.ActionHistories[0].AgentActionHistories = agentActionHistories
}


// 旧ActAgents (削除予定)

// ActAgents はエージェントの行動を指定し、フィールドを変更します。
// func (f *Field) ActAgents(updateActions []*apispec.UpdateAction) error {
// 	// この座標を行動先に選んだエージェントの数
// 	var DistinationCount [][]int

// 	for i := 0; i < len(updateActions); i++ {
// 		for j := 0; j < len(f.Agents); j++ {
// 			if updateActions[i].AgentID == f.Agents[j].ID {
// 				x := updateActions[i].DX + f.Agents[j].X
// 				y := updateActions[i].DY + f.Agents[j].Y
// 				if DistinationCount[y][x] == 1 {
// 					// 動かす
// 					if updateActions[i].Type == "move" {
// 						if f.Cells[y][x].TiledBy != updateActions[i].AgentID && f.Cells[y][x].TiledBy != 0 {
// 							continue
// 						}
// 						f.Cells[y][x].TiledBy = updateActions[i].AgentID
// 						f.Agents[j].X += updateActions[i].DX
// 						f.Agents[j].Y += updateActions[i].DY
// 					} else if updateActions[i].Type == "remove" {
// 						if f.Cells[y][x].TiledBy == updateActions[i].AgentID || f.Cells[y][x].TiledBy == 0 {
// 							continue
// 						}
// 						f.Cells[y][x].TiledBy = updateActions[i].AgentID
// 					} else if updateActions[i].Type == "stay" {
// 						continue
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// 範囲外にアクセスしようとしたとき err(error型) を返すようにしてほしい
// 	return nil
// }

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
