// chromeとかで開いたページはGETとして認識される

package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/satackey/procon31server/pkg/apispec"

	"sync"
)

func main() {

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

var lock = sync.RWMutex{}

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
var googlemap = map[int]*TwoResponse{}

func DoubleResponse(w ResponseWriter, error string, unixtime int, code int) {
	// googlemap := TwoResponse{}
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
	lock.RLock()
	allmatch := apispec.Match{}
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

	matchstruct := apispec.Match{}
	matchstruct.Turns = 9999
	err := r.DecodeJsonPayload(&matchstruct)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}
	if matchstruct.ID == 0 {
		rest.Error(w, "id required", 400)
		return
	}
	if matchstruct.IntervalMillis == 0 {
		rest.Error(w, "intervalmills ないぞ", 400)
		return
	}
	if matchstruct.MatchTo == "" {
		rest.Error(w, "Matchto required", 400)
		return
	}
	if matchstruct.TeamID == 0 {
		rest.Error(w, "teamid required", 400)
		return
	}
	if matchstruct.TurnMillis == 0 {
		rest.Error(w, "turnmillis required", 400)
		return
	}
	if matchstruct.Turns == 9999 {
		rest.Error(w, "turns required", 400)
		return
	}
	// lock.Lock()
	// store[matchstruct.Id] = &matchstruct
	// lock.Unlock()
	// w.WriteJson(&matchstruct)
	MyResponse(w, "OK", 202)
	return
}

var ActionStore = map[int]*apispec.FieldStatusAction{}

// Matchのidを判別しつつUpdateActionをGetする
func GetUpdateAction(w rest.ResponseWriter, r *rest.Request) {

	// TODO ここちょっと怪しい(値自体を読み込めていない気がする)
	NewAction := apispec.FieldStatus{}
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
	TeamStatus := apispec.Team{}

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
