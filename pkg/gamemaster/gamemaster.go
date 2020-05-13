package gamemaster

import (
	"github.com/satackey/procon31server/pkg/apispec"
)

type GameMaster struct {
	FieldStatus    apispec.FieldStatus
	turnMillis     int
	intervalMillis int
}

func (g *GameMaster) createMatch(FieldStatus, turnMillis, intervalMillis) {

}

func (g *GameMaster) getFieldByID(matchID) {

}

func (g *GameMaster) postAgentActions(teamID, FieldStatusAction) {

}
