package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("localhost:3000 でサーバを起動しました")
	api := rest.NewApi()
	api.Use(&rest.AccessLogApacheMiddleware{})
	api.Use(rest.DefaultCommonStack...)
	router, err := rest.MakeRouter(
		rest.Get("/stream", StreamThings),
		rest.Get("/matches", MatchesThings),
		rest.Get("/matches/:id", FieldStatus),
		// rest.Post("/matches/:id/action", Action),s
		rest.Get("/ping", PingsThings),


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
	cpt := 0
	for {
		cpt++
		w.WriteJson(
			&Thing{
				Name: fmt.Sprintf("thing #%d", cpt),
			},
		)
		w.(http.ResponseWriter).Write([]byte("\n"))
		// Flush the buffer to client
		w.(http.Flusher).Flush()
		// wait 3 seconds
		time.Sleep(time.Duration(3) * time.Second)
	}
}

// Ping は
type Ping struct {
	Status string
}

// PingsThings は
func PingsThings(w rest.ResponseWriter, r *rest.Request) {
	// for {
	// 	cpt++
	// 	w.WriteJson(
	// 		&Thing{
	// 			Name: fmt.Sprintf("thing #%d", cpt),
	// 		},
	// 	)

	// ping := Ping{
	// 	status: "OK",
	// }
	// w.WriteJson(&ping)
	w.(http.ResponseWriter).Write([]byte("\n"))
		// Flush the buffer to client
	w.(http.Flusher).Flush()
		// wait 3 seconds
	time.Sleep(time.Duration(3) * time.Second)
	// }
}


// Match は
// type Match struct {
// 	 string
// }

// MatchesThings は
func MatchesThings(w rest.ResponseWriter, r *rest.Request) {

	w.(http.ResponseWriter).Write([]byte("\n"))
		// Flush the buffer to client
	w.(http.Flusher).Flush()
		// wait 3 seconds
	time.Sleep(time.Duration(3) * time.Second)
}

// FieldStatus は
func FieldStatus(w rest.ResponseWriter, r *rest.Request) {

	w.(http.ResponseWriter).Write([]byte("\n"))
		// Flush the buffer to client
	w.(http.Flusher).Flush()
		// wait 3 seconds
	time.Sleep(time.Duration(3) * time.Second)
}


// Action は 行動情報を更新してpostします
// func Action(w rest.ResponseWriter, r *rest.Request) {

// 	w.(http.ResponseWriter).Write([]byte("\n"))
// 		// Flush the buffer to client
// 	w.(http.Flusher).Flush()
// }