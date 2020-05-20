// chromeとかで開いたページはGETとして認識される

package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
    // "time"
	// "../apisec/apisec.go"
	"sync"
	"strconv"
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
		rest.Get("/matches/:code", GetCountry),
		// 必要かまだわからない
		// rest.Delete("/countries/:code", DeleteCountry),

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
	Status string
}

// PingsThings は
func PingsThings(w rest.ResponseWriter, r *rest.Request) {

	w.(http.ResponseWriter).Write([]byte("\n"))
		// Flush the buffer to client
	w.(http.Flusher).Flush()

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

type Country struct {
	Code int `json:"code"`
	Name string `json:"name"`
	Hoge string `json:"fuga"`
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
var store = map[int]*Country{}
// var store = []*Country{}

var lock = sync.RWMutex{}

func PostMatch(w rest.ResponseWriter, r *rest.Request) {
	// country := new(Country)
	country := Country{}
	// fmt.Print(country) { }が入ってた 初期化されている

	// リクエストを読み取り, json.UnmarshalでJSONをデコードしてる
	// &countryはJSONが入っている
	err := r.DecodeJsonPayload(&country)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError) // == 500
		return
	}
	if country.Code == 0 {
		rest.Error(w, "country code required", 400)
		return
	}
	if country.Name == "" {
		rest.Error(w, "country name required", 400)
		return
	}
	if country.Hoge == "" {
		rest.Error(w, "hoge required", 400)
		return
	}
	// if country.Hoge == "" {
	// 	rest.Error(w, "hoge required", 400)
	// 	return
	// }
	// if country.Hoge == "" {
	// 	rest.Error(w, "hoge required", 400)
	// 	return
	// }
	lock.Lock()
	store[country.Code] = &country
	lock.Unlock()
	w.WriteJson(&country)
}

func GetAllMatch(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	countries := make([]Country, len(store)) // mapを初期化してる storeの数だけ場所を確保している
	i := 0
	for _, country := range store {
		countries[i] = *country
		i++
	}
	lock.RUnlock()
	w.WriteJson(&countries)
}

func GetCountry(w rest.ResponseWriter, r *rest.Request) {
	code_val := r.PathParam("code")
	// fmt.Println("aaaa")
	// code はstr型
	var code int
	code, _ = strconv.Atoi(code_val)
	fmt.Println(code) // urlのパラメータ
	lock.RLock()
	var country *Country
	// fmt.Printf("%T\n", country) // nil
	// fmt.Println("%T", store) //map

	// fmt.Printf("%T %T\n", store, code) // map[int]*main.Countryと string
	// if store[code] == nil {
	// 	fmt.Printf("nil")
	// }
    if store[code] != nil {
        country = &Country{}
		*country = *store[code]
		// fmt.Println(*country)
		// fmt.Println("x")
    }
    lock.RUnlock()

	

    // if country == nil { // ここの左辺がなぜstore[code]じゃだめかわからない 動いたからヨシはだめなので後で考える
	if store[code] == nil {
		rest.NotFound(w, r)
		// fmt.Println("a")
        return
    }
    w.WriteJson(country)
}


// Routes sharing a common placeholder MUST name it consistently: id != code
// /matches/:[] ここがかぶるとだめ