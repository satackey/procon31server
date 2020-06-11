// chromeとかで開いたページはGETとして認識される

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	// "../apisec/apisec.go"

	"sync"
)

func main() {

	log.Println("localhost:3000 でサーバを起動しました")
	api := rest.NewApi()

	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/ping", PingsThings),

		rest.Get("/matches", GetAllMatch),

		// 実際には降ってくるけど, 練習環境ではそれがないのでpostを用意しておく
		// rest.Post("mathces", PostMatch),
		// PathExp must start with /
		rest.Post("/matches", PostMatch),

		// CRUDのC Post or Put
		rest.Get("/matches/:id", GetUpdateAction),
		// rest.Delete("/matches/:id", DeleteMatch),
		rest.Post("/matches/:id/advanceaction", PostAction),
		// rest.Get("/matches/:id/action"),
		rest.Post("/matches/:id/action", PostUpdate),
		// // rest.Post("/matches/:id/action", NewPostUpdate),
		// rest.Post("/matchesnew/:id", GetFieldStatus),
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

// こういう構造体をapisec.goから引っ張ってきたい
// type Status struct {
//     Id int `json:id`
// }

type Match struct {
	Id             int    `json:"id"`             // 試合のid
	IntervalMillis int    `json:"intervalMillis"` // ターンとターンの間の時間
	MatchTo        string `json:"matchTo"`        // 対戦相手の名前
	TeamID         int    `json:"teamID"`         // 自分のteamid
	TurnMillis     int    `json:"turnMillis"`     // 1ターンの(制限?)時間
	Turns          int    `json:"turns"`          // 試合のターン数
}

// var store = map[int]*Match{}

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
	// lock.Lock()
	// store[matchstruct.Id] = &matchstruct
	// lock.Unlock()
	// w.WriteJson(&matchstruct)
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

func GetAllMatch(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	// manymatches := make([]Match, len(store)) // mapを初期化してる storeの数だけ場所を確保している
	// i := 0
	// for _, matchstruct := range store {
	// 	manymatches[i] = *matchstruct
	// 	i++
	// }
	// x := unsafe.Sizeof(manymatches)
	// fmt.Println(x)

	// lock.RUnlock() // ここに書かないとpostができない
	// // X := make([]int, 5, 5)
	// // 全部24じゃん
	// // fmt.Println(unsafe.Sizeof(X))
	// // fmt.Println(unsafe.Sizeof(manymatches))
	// fmt.Println(len(store)) // 0??????
	// fmt.Println("len:", len(manymatches))
	// // zero := len(manymatches)
	// if len(manymatches) == 0 {
	// 	MyResponse(w, "値一個もないぞ〜〜〜", 401)
	// 	return
	// }

	// fmt.Printf("%T %T", manymatches, &manymatches)

	// w.WriteJson(&manymatches)
	MyResponse(w, "OK", 202)
	return
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
	// lock.Lock()
	// store[matchstruct.Id] = &matchstruct
	// lock.Unlock()
	// w.WriteJson(&matchstruct)
	MyResponse(w, "OK", 202)
	return
}

type UpdateAction struct {
	AgentID int    `json:"agentID"`
	DX      int    `json:"dx"`
	DY      int    `json:"dy"`
	Type    string `json:"type"`
	Turn    int    `json:"turn"`
	// Array   []Agent `json:"agents"`
}

// TODO:変数名をちゃんと考える
// var store_2 = map[int]*UpdateAction{}

func PostUpdate(w rest.ResponseWriter, r *rest.Request) {

	// id_val := r.PathParam("id")
	// // code はstr型
	// var id int
	// id, _ = strconv.Atoi(id_val)

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

	// lock.Lock()
	// FieldStatusActionStruct := FieldStatusAction{}
	// デコードした値を別の構造体型mapに代入しなきゃいけない
	// 別の型に代入する方法(ActionStoreにNewActionStructのアドレスを代入)を調べるか左辺を細かく指定する(値を一個ずつ ActionStoreのn番目のidというkeyに NewActionStruct.Idを代入みたいな)方法を調べる
	// matchstruct := Match{}

	// store_2[id] = &NewActionStruct
	// lock.Unlock()
	// w.WriteJson(&NewActionStruct)
	MyResponse(w, "OK", 202)
	return
}

type FieldStatusAction struct {
	AgentID int    `json:"agentID"`
	DX      int    `json:"dx"`
	DY      int    `json:"dy"`
	Type    string `json:"type"`
	Apply   int    `json:"apply"`
	Turn    int    `json:"turn"`
}

var ActionStore = map[int]*FieldStatusAction{}

// Matchのidを判別しつつUpdateActionをGetする
func GetUpdateAction(w rest.ResponseWriter, r *rest.Request) {
	// id_val := r.PathParam("id")
	// // code はstr型
	// var id int
	// id, _ = strconv.Atoi(id_val)
	// fmt.Println(id) // urlのパラメータ

	// lock.RLock()

	// var updateactionstruct *UpdateAction
	// var otherstruct *Y
	// newaction := UpdateAction{}
	// agentid := store_2.agentID

	// if store_2[id] != nil {
	// 	// 右辺でアドレスを指定することで型を流し込んでいる
	// 	// matchstruct = &Match{}
	// 	updateactionstruct = &UpdateAction{}

	// 	// otherstruct = &Y{}
	// 	// store[id]に格納されている値をpointerで呼び出して代入
	// 	// *matchstruct = *store[id]
	// 	*updateactionstruct = *store_2[id]
	// 	// *otherstruct = *store_2[id]

	// 	// fmt.Println(*country)
	// 	// fmt.Println("x")
	// }

	// updateactionstruct

	// lock.RUnlock()
	// if country == nil { // ここの左辺がなぜstore[code]じゃだめかわからない 動いたからヨシはだめなので後で考える

	// if store[id] == nil {

	// 	rest.NotFound(w, r)
	// 	return
	// }
	// w.WriteJson(updateactionstruct)
	MyResponse(w, "OK", 202)
	return
}

// }
