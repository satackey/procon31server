package field

import (
	"fmt"
	"testing"

	"github.com/strvworks/procon30/pkg/facilitator"
)

var exampleFieldStatus *facilitator.FieldStatus = &facilitator.FieldStatus{
	Width:  6,
	Height: 4,
	Points: [][]int{
		{
			2, -4, 0, 0, -4, 2,
		},
		{
			5, -1, 1, 1, -1, 5,
		},
		{
			5, -1, 1, 1, -1, 5,
		},
		{
			2, -4, 0, 0, -4, 2,
		},
	},
	StartedAtUnixtime: 0,
	Turn:              0,
	Tiled: [][]int{
		{
			0, 0, 6, 0, 0, 0,
		},
		{
			5, 0, 0, 0, 0, 0,
		},
		{
			0, 0, 0, 0, 0, 6,
		},
		{
			0, 0, 0, 5, 0, 0,
		},
	},
	Teams: []facilitator.Team{
		facilitator.Team{
			TeamID: 5,
			Agents: []facilitator.Agent{
				facilitator.Agent{
					AgentID: 9,
					X:       1,
					Y:       2,
				},
				facilitator.Agent{
					AgentID: 10,
					X:       4,
					Y:       4,
				},
			},
			TilePoint: 5,
			AreaPoint: 0,
		},
		facilitator.Team{
			TeamID: 6,
			Agents: []facilitator.Agent{
				facilitator.Agent{
					AgentID: 11,
					X:       3,
					Y:       1,
				},
				facilitator.Agent{
					AgentID: 12,
					X:       6,
					Y:       3,
				},
			},
			TilePoint: 5,
			AreaPoint: 0,
		},
	},

	Actions: []facilitator.FieldStatusAction{
		/*
			//フィールド情報_turn1.json の Actionsのなかみ
			&facilitator.FieldStatusAction{
				AgentID: 9
				DX: 1
				DY: 1
				Type: "move"
				Apply: 1
				Turn: 1
			},
			&facilitator.FieldStatusAction{
				AgentID: 10
				DX: -1
				DY: -1
				Type: "move"
				Apply: 1
				Turn: 1
			},
			&facilitator.FieldStatusAction{
				AgentID: 11
				DX: 1
				DY: 0
				Type: "move"
				Apply: 1
				Turn: 1
			},
			&facilitator.FieldStatusAction{
				AgentID: 12
				DX: 0
				DY: -1
				Type:"move"
				Apply: 1
				Turn: 1
			}
		*/
	},
}

func TestField(t *testing.T) {
	field := New()
	field.InitField(exampleFieldStatus)
	fmt.Fprintf("%+v\n", field.Cells)
}
