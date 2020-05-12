package field

// Cell はフィールドのますを表します
type Cell struct {
	Point int
	// 下の変数のいい名前は無いかな～
	TiledBy int
	// state + IDで、どのチームのどんなマスか
	status string
	teamID int
	// というのはどうかしら？

	x     int
	y     int
	field *Field
}

func newCell(point int, tiledBy int, x int, y int, field *Field) *Cell {
	return &Cell{point, tiledBy, x, y, field}
}

// GetAgentID はますにいるエージェントのIDを返します。
// エージェントがますにいない時は、-1を返します。
func (c Cell) GetAgentID() int {
	// c.field を使い、Agentがこのマスにいれば、そのIDを返す
	// c.field.Agents[].x
	/*
		for i, agentID := range agentIDs {
			c.field.Agents[agentID]
		}
	*/
	return 0
}
