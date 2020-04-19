package field

// Agent はエージェントを表します
type Agent struct {
	ID     int
	TeamID int
	X      int
	Y      int
	field  *Field
}

func newAgent(id int, teamID int, x int, y int, field *Field) *Agent {
	return &Agent{id, teamID, x, y, field}
}
