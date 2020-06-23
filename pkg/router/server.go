// chromeとかで開いたページはGETとして認識される

package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/satackey/procon31server/pkg/apispec"
)

func main() {

	// apispec.Match.ID = 12
	log.Println("localhost:3000 でサーバを起動しました")
	api := rest.NewApi()

	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/ping", PingsThings),
		rest.Get("/matches", GetAllMatch),
		rest.Get("/matches/:id", GetUpdateAction),
		rest.Post("/matches/:id/action", PostAction),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":3000", api.MakeHandler()))
}

// Ping は
type Ping struct {
	Status string `json:"status"`
}

// PingsThings は
func PingsThings(w rest.ResponseWriter, r *rest.Request) {

	MyResponse(w, "OK", 202)
	return

}

// https://github.com/ant0ine/go-json-rest/blob/ebb33769ae013bd5f518a8bac348c310dea768b8/rest/response.go

type ResponseWriter interface {

	// Identical to the http.ResponseWriter interface
	Header() http.Header
	WriteJson(v interface{}) error

	EncodeJson(v interface{}) ([]byte, error)

	WriteHeader(int)
}

var ErrorFieldName = "status"

func MyResponse(w ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	err := w.WriteJson(map[string]string{ErrorFieldName: error})
	if err != nil {
		panic(err)
	}
}

var StartAtUnixTime = "startAtUnixTime"
var TwoError = map[int]*TwoResponse{}

func DoubleResponse(w ResponseWriter, error string, unixtime int, code int) {
	// TwoError := TwoResponse{}
	w.WriteHeader(code)
	// err := w.WriteJson(map[string]string{ErrorFieldName: error, StartAtUnixTime: unixtime})
	err := w.WriteJson(TwoResponse{Status: error, UnixTime: unixtime})
	if err != nil {
		panic(err)
	}
}

type TwoResponse struct {
	UnixTime int    `json:"startAtUnixTime"`
	Status   string `json:"status"`
}

func GetAllMatch(w rest.ResponseWriter, r *rest.Request) {
	allmatch := &apispec.Match{}
	allmatch.Turns = 9999

	if allmatch.ID == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if allmatch.IntervalMillis == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if allmatch.MatchTo == "" {
		rest.Error(w, "id required", 400)
		return
	}
	if allmatch.TeamID == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if allmatch.TurnMillis == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if allmatch.Turns == 9999 {
		rest.Error(w, "id required", 400)
		return
	}

	MyResponse(w, "OK", 202)
	return
}

func PostAction(w rest.ResponseWriter, r *rest.Request) {

	actions := &apispec.UpdateAction{}
	// actions := &apispec.Update{}
	// action.Turns = 9999
	actions.DX = 9999
	actions.DY = 9999
	err := r.DecodeJsonPayload(&actions)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}

	if actions.AgentID == 0 {
		rest.Error(w, "agentid required", 400)
		return
	}
	if actions.DX == 9999 {
		rest.Error(w, "DX ないぞ", 400)
		return
	}
	if actions.DX < -1 || actions.DX > 1 {
		rest.Error(w, "DXは-1から1の範囲にしましょう", 400)
		return
	}
	if actions.DY == 9999 {
		rest.Error(w, "DY ないぞ", 400)
		return
	}
	if actions.DY < -1 || actions.DY > 1 {
		// rest.Error(w, "dy required", 400)
		DoubleResponse(w, "invalid", 0, 400)
		return
	}
	// こっちは必要ないかも 必要だった
	if actions.Type == "" {
		rest.Error(w, "type ないぞ", 400)
		return
	}
	if actions.Type != "move" || actions.Type != "remove" || actions.Type != "stay" {
		rest.Error(w, "typeはmoveかremoveかstayで入力しよう", 400)
		return
	}
	// if actions.Turns == 9999 {
	// 	rest.Error(w, "turns required", 400)
	// 	return
	// }
	// else {

	//辺を越えようとしてきた場合のエラー
	// 4辺の値 width height を事前情報から引っ張ってくる

	MyResponse(w, "OK", 202)
	return
}

var ActionStore = map[int]*apispec.FieldStatusAction{}

// Matchのidを判別しつつUpdateActionをGetする
func GetUpdateAction(w rest.ResponseWriter, r *rest.Request) {

	// TODO ここちょっと怪しい(値自体を読み込めていない気がする)
	NewAction := &apispec.FieldStatus{}
	NewAction.StartedAtUnixtime = 9999
	// NewAction.Turn

	if NewAction.Width == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if NewAction.Height == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if len(NewAction.Points) == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if NewAction.StartedAtUnixtime == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if NewAction.Turn == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if len(NewAction.Cells) == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if len(NewAction.Teams) == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if len(NewAction.Actions) == 0 {
		rest.Error(w, "id required", 400)
		return
	}

	// Teamの中身の確認
	// TODO ここもちょっと怪しい
	TeamStatus := &apispec.Team{}

	if TeamStatus.TeamID == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if len(TeamStatus.Agents) == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	// tilepointは存在してた
	// TODO tilepointかcellpoint どちらかに統一
	if TeamStatus.TilePoint == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if TeamStatus.AreaPoint == 0 {
		rest.Error(w, "id required", 400)
		return
	}

	// Agent
	// AgentStatus := &apispec.Agent{}

	// if AgentStatus.AgentID == 0 {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }
	// if AgentStatus.X == 0 {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }
	// if AgentStatus.Y == 0 {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }

	// // FieldStatusのAction
	// FieldAction := &apispec.FieldStatusAction{}

	// if FieldAction.AgentID == 0 {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }
	// if FieldAction.DX > 3 {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }
	// if FieldAction.DY == 0 {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }
	// if FieldAction.Type == "" {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }
	// if FieldAction.Apply == 0 {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }
	// if FieldAction.Turn == 0 {
	// 	rest.Error(w, "id required", 400)
	// 	return
	// }

	MyResponse(w, "OK", 202)
	return
}

// 初期値が0じゃだめな変数
// startAtUnixTime
// turn
// areapoint

// turnとturnsの違い
// startAtUnixTimeとstartedAtUnixTimeの違い
// どちらも0が代入される可能性があるかどうか
