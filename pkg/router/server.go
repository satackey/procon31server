// chromeとかで開いたページはGETとして認識される

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	// "time"
	// "../apisec/apisec.go"
	"strconv"
	"sync"
	"unsafe"
)

func main() {

	// users := Users{
	// 	Store: map[string]*User{},
	// }

	log.Println("localhost:3000 でサーバを起動しました")
	api := rest.NewApi()

	// https://github.com/ant0ine/go-json-rest#streamingに使われていたapi
	// api.Use(&rest.AccessLogApacheMiddleware{})
	// api.Use(rest.DefaultCommonStack...)

	// 基本的なやつ?
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/stream", StreamThings),
		// rest.Get("/matches", MatchesThings),
		// rest.Get("/matches/:id", FieldStatus),
		// rest.Get("/matches/:id", users.FieldStatus),
		// rest.Post("/matches/:id/action", Action),s
		rest.Get("/ping", PingsThings),

		rest.Get("/matches", GetAllMatch),

		// 実際には降ってくるけど, 練習環境ではそれがないのでpostを用意しておく
		// rest.Post("mathces", PostMatch),
		// PathExp must start with /
		rest.Post("/matches", PostMatch),

		// CRUDのC Post or Put
		rest.Get("/matches/:id", GetUpdateAction),
		// 必要かまだわからない
		rest.Delete("/matches/:id", DeleteMatch),

		rest.Post("/matches/:id/advanceaction", PostAction),
		// rest.Get("/matches/:id/action"),

		rest.Post("/matches/:id/action", PostUpdate),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":3000", api.MakeHandler()))
}

// Thing は
type Thing struct {
	Name string
}

// StreamThings は
func StreamThings(w rest.ResponseWriter, r *rest.Request) {
	// cpt := 0
	// for {
	// 	cpt++

	// w.WriteJson(
	// 	&Thing{
	// 		Name: fmt.Sprintf("thing #%d", cpt),
	// 	},
	// )
	w.(http.ResponseWriter).Write([]byte("\n"))
	// Flush the buffer to client
	w.(http.Flusher).Flush()
	// wait 3 seconds

	// time.Sleep(time.Duration(3) * time.Second)
	// }
}

// Ping は
type Ping struct {
	Status string `json:"status"`
}

// PingsThings は
func PingsThings(w rest.ResponseWriter, r *rest.Request) {

	// ping := Ping{"OK"}
	// err := r.DecodeJsonPayload(&ping)
	// if err != nil {
	// 	rest.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.(http.ResponseWriter).Write([]byte("\n"))
	// 	// Flush the buffer to client
	// w.(http.Flusher).Flush()
	// // w.WriteHeader(http.StatusOK)

	// w.WriteJson(&ping)

	// ping := Ping{"OK"}

	// // res, err := json.Marshal(ping)

	// // res, err := EncodeJson(ping)
	// // res, err := rest.EncodeJson(ping)
	// // res, err := r.EncodeJson(ping)

	// if {
	// 	rest.Error(w, "")
	// }

	// res, err := w.EncodeJson(ping)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")

	// x := r.PathParam()
	// if x == "ping" {
	// 	w.WriteJson(http.StatusOK)
	// }

	// if x != "ping" {
	// 	w.WriteJson(http.StatusUnauthorized)
	// }

	// w.Write(res)

	MyResponse(w, "OK", 202)
	return

}

// Match は
// type Match struct {
// 	 string
// }

// type User struct {
// 	Id string
// 	Name string
// }

// type Users struct {
// 	sync.RWMutex
// 	Store map[string]*User
// }

// MatchesThings は
func MatchesThings(w rest.ResponseWriter, r *rest.Request) {

	w.(http.ResponseWriter).Write([]byte("\n"))
	// Flush the buffer to client
	w.(http.Flusher).Flush()
	// wait 3 seconds
	// time.Sleep(time.Duration(3) * time.Second)
}

// こういう構造体をapisec.goから引っ張ってきたい
// type Status struct {
//     Id int `json:id`
// }

// FieldStatus は
// func (u *Users) FieldStatus(w rest.ResponseWriter, r *rest.Request) {
func FieldStatus(w rest.ResponseWriter, r *rest.Request) {

	// id := r.PathParam("id")
	// fmt.Print(id)
	// u.RLock()  // グローバル変数にアクセスするのに必要?
	// var user *User
	// if u.Store[id] != nil {
	// 	user = &User{}
	// 	*user = *u.Store[id]
	// }
	// u.RUnlock()
	// if user == nil {
	// 	rest.NotFound(w, r)
	// 	return
	// }
	// w.WriteJson(user)

	w.(http.ResponseWriter).Write([]byte("\n"))
	// Flush the buffer to client
	w.(http.Flusher).Flush()
	// wait 3 seconds
	// time.Sleep(time.Duration(3) * time.Second)
}

// Action は 行動情報を更新してpostします
// func Action(w rest.ResponseWriter, r *rest.Request) {

// 	w.(http.ResponseWriter).Write([]byte("\n"))
// 		// Flush the buffer to client
// 	w.(http.Flusher).Flush()
// }

// type Country struct {
// 	Code string
// 	Name string
// }

type Match struct {
	Id             int    `json:"id"`             // 試合のid
	IntervalMillis int    `json:"intervalMillis"` // ターンとターンの間の時間
	MatchTo        string `json:"matchTo"`        // 対戦相手の名前
	TeamID         int    `json:"teamID"`         // 自分のteamid
	TurnMillis     int    `json:"turnMillis"`     // 1ターンの(制限?)時間
	Turns          int    `json:"turns"`          // 試合のターン数
}

// type Match struct {
// 	ID             int    `json:"id"`
// 	IntervalMillis int    `json:"intervalMillis"`
// 	// MatchTo        string `json:"matchTo"`
// 	// TeamID         int    `json:"teamID"`
// 	// TurnMillis     int    `json:"turnMillis"`
// 	// Turns          int    `json:"turns"`
// }

// var store = map[string]*Country{}
var store = map[int]*Match{}

// var store = []*Country{}

var lock = sync.RWMutex{}

func PostMatch(w rest.ResponseWriter, r *rest.Request) {
	// country := new(Country)
	matchstruct := Match{}
	// fmt.Print(country) { }が入ってた 初期化されている

	// リクエストを読み取り, json.UnmarshalでJSONをデコードしてる
	// &countryはJSONが入っている
	err := r.DecodeJsonPayload(&matchstruct)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}
	if matchstruct.Id == 0 {
		rest.Error(w, "country id required", 400)
		return
	}
	if matchstruct.IntervalMillis == 0 {
		rest.Error(w, "country intervalmills ないぞ", 400)
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
	if matchstruct.Turns == 0 {
		rest.Error(w, "turns required", 400)
		return
	}
	lock.Lock()
	store[matchstruct.Id] = &matchstruct
	lock.Unlock()
	w.WriteJson(&matchstruct)
}

// type ZERO struct {
// 	Id int `json:"id"`
// 	IntervalMillis string `json:"intervalMillis"`
// 	MatchTo string `json:"matchTo"`
// 	TeamID         int    `json:"teamID"`
// 	TurnMillis     int    `json:"turnMillis"`
// 	Turns          int    `json:"turns"`
// }

// https://github.com/ant0ine/go-json-rest/blob/ebb33769ae013bd5f518a8bac348c310dea768b8/rest/response.go

type ResponseWriter interface {

	// Identical to the http.ResponseWriter interface
	Header() http.Header

	// Use EncodeJson to generate the payload, write the headers with http.StatusOK if
	// they are not already written, then write the payload.
	// The Content-Type header is set to "application/json", unless already specified.
	WriteJson(v interface{}) error

	// Encode the data structure to JSON, mainly used to wrap ResponseWriter in
	// middlewares.
	EncodeJson(v interface{}) ([]byte, error)

	// Similar to the http.ResponseWriter interface, with additional JSON related
	// headers set.
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

func GetAllMatch(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	manymatches := make([]Match, len(store)) // mapを初期化してる storeの数だけ場所を確保している
	i := 0
	for _, matchstruct := range store {
		manymatches[i] = *matchstruct
		i++
	}
	x := unsafe.Sizeof(manymatches)
	fmt.Println(x)

	// y := Match{}
	// if unsafe.Sizeof(manymatches) == 24 {
	// 	rest.Error(w, "値一個もないぞ", 401)
	// 	return
	// }

	lock.RUnlock() // ここに書かないとpostができない
	// X := make([]int, 5, 5)
	// 全部24じゃん
	// fmt.Println(unsafe.Sizeof(X))
	// fmt.Println(unsafe.Sizeof(manymatches))
	fmt.Println(len(store)) // 0??????
	fmt.Println("len:", len(manymatches))
	// zero := len(manymatches)
	if len(manymatches) == 0 {
		MyResponse(w, "値一個もないぞ〜〜〜", 401)
		return
	}

	fmt.Printf("%T %T", manymatches, &manymatches)

	w.WriteJson(&manymatches)
}

// func GetMatch(w rest.ResponseWriter, r *rest.Request) {
// 	id_val := r.PathParam("id")
// 	// fmt.Println("aaaa")
// 	// code はstr型
// 	var id int
// 	id, _ = strconv.Atoi(id_val)
// 	// fmt.Println(id) // urlのパラメータ
// 	lock.RLock()
// 	var matchstruct *Match
// 	// fmt.Printf("%T\n", country) // nil
// 	// fmt.Println("%T", store) //map

// 	// fmt.Printf("%T %T\n", store, code) // map[int]*main.Countryと string
// 	// if store[code] == nil {
// 	// 	fmt.Printf("nil")
// 	// }
// 	if store[id] != nil {
// 		matchstruct = &Match{}
// 		*matchstruct = *store[id]
// 		// fmt.Println(*country)
// 		// fmt.Println("x")
// 	}
// 	lock.RUnlock()

// 	// if country == nil { // ここの左辺がなぜstore[code]じゃだめかわからない 動いたからヨシはだめなので後で考える
// 	if store[id] == nil {
// 		rest.NotFound(w, r)
// 		// fmt.Println("a")
// 		return
// 	}
// 	w.WriteJson(matchstruct)
// }

// Routes sharing a common placeholder MUST name it consistently: id != code
// /matches/:[] ここがかぶるとだめ

func DeleteMatch(w rest.ResponseWriter, r *rest.Request) {
	id_val := r.PathParam("id")
	var id int
	id, _ = strconv.Atoi(id_val)
	lock.Lock()
	delete(store, id)
	lock.Unlock()
	w.WriteHeader(http.StatusOK)
}

// このpostをmatches:id に反映させる? させた(/actionの使い方があっているか確認) 違う
func PostAction(w rest.ResponseWriter, r *rest.Request) {

	// actionstruct := UpdateAction{}
	matchstruct := Match{}

	err := r.DecodeJsonPayload(&matchstruct)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}
	if matchstruct.Id == 0 {
		rest.Error(w, "country id required", 400)
		return
	}
	if matchstruct.IntervalMillis == 0 {
		rest.Error(w, "country intervalmills ないぞ", 400)
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
	if matchstruct.Turns == 0 {
		rest.Error(w, "turns required", 400)
		return
	}
	lock.Lock()
	store[matchstruct.Id] = &matchstruct
	lock.Unlock()
	w.WriteJson(&matchstruct)
}

type UpdateAction struct {
	AgentID int    `json:"agentID"`
	DX      int    `json:"dx"`
	DY      int    `json:"dy"`
	Type    string `json:"type"`
	Turn    int    `json:"turn"`
	// Array   []Agent `json:"agents"`
}

// type Agent struct {
// 	AgentID int `json:"agentID"`
// 	X       int `json:"x"`
// 	Y       int `json:"y"`
// }

// TODO:変数名をちゃんと考える
var store_2 = map[int]*UpdateAction{}

func PostUpdate(w rest.ResponseWriter, r *rest.Request) {

	id_val := r.PathParam("id")
	// code はstr型
	var id int
	id, _ = strconv.Atoi(id_val)

	NewActionStruct := UpdateAction{}

	// turnが変更されないか監視するための変数(decodeされる前なのでpostされる値はまだ入っていない(に等しい))
	NewActionStruct.Turn = 5
	CheckTurnVal := NewActionStruct.Turn
	fmt.Println(CheckTurnVal)
	err := r.DecodeJsonPayload(&NewActionStruct)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}
	if NewActionStruct.AgentID == 0 {
		rest.Error(w, "agent id required", 400)
		return
	}
	if NewActionStruct.DX == 0 {
		rest.Error(w, "dx ないぞ", 400)
		return
	}
	if NewActionStruct.DY == 0 {
		rest.Error(w, "dy required", 400)
		return
	}
	if NewActionStruct.Type == "" {
		rest.Error(w, "type required", 400)
		return
	}

	// turnが書き換えられていないか確認している
	if CheckTurnVal != NewActionStruct.Turn {
		NewActionStruct.Turn = CheckTurnVal
	}

	lock.Lock()

	// matchstruct := Match{}

	store_2[id] = &NewActionStruct
	lock.Unlock()
	w.WriteJson(&NewActionStruct)
}

// Matchのidを判別しつつUpdateActionをGetする
func GetUpdateAction(w rest.ResponseWriter, r *rest.Request) {
	id_val := r.PathParam("id")
	// code はstr型
	var id int
	id, _ = strconv.Atoi(id_val)
	// fmt.Println(id) // urlのパラメータ

	lock.RLock()

	// var matchstruct *Match
	var updateactionstruct *UpdateAction

	if store_2[id] != nil {
		// 右辺でアドレスを指定することで型を流し込んでいる
		// matchstruct = &Match{}
		updateactionstruct = &UpdateAction{}
		// store[id]に格納されている値をpointerで呼び出して代入
		// *matchstruct = *store[id]
		*updateactionstruct = *store_2[id]

		// fmt.Println(*country)
		// fmt.Println("x")
	}
	lock.RUnlock()
	// if country == nil { // ここの左辺がなぜstore[code]じゃだめかわからない 動いたからヨシはだめなので後で考える

	// if store[id] == nil {

	// 	rest.NotFound(w, r)
	// 	return
	// }
	w.WriteJson(updateactionstruct)
	// w.WriteJson(matchstruct)
}

// type Turn struct {
// 	[]
// }
