package server

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"ws"
)

var requireAuth bool

func Run(port int, disableAuth bool) {
	//start := time.Now()
	r := mux.NewRouter()
	requireAuth = !disableAuth

	// user requests
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/verifySecret", verifySecret).Methods("POST")
	r.HandleFunc("/users", makeUser).Methods("POST")

	// unauthorized requests
	r.HandleFunc("/games", getAllGames).Methods("GET")
	r.HandleFunc("/games", makeGame).Methods("POST")
	r.HandleFunc("/games/{ID}", getGame).Methods("GET")
	r.HandleFunc("/games/{ID}/info", getGame).Methods("GET")
	r.HandleFunc("/games/{ID}/board", getBoard).Methods("GET")
	r.HandleFunc("/games/{ID}/string", getGameString).Methods("GET")
	r.HandleFunc("/games/{ID}/ws", ws.ServeWs).Methods("GET")

	// authorized requests
	r.HandleFunc("/games/{ID}", makeGameMove).Methods("POST")
	r.HandleFunc("/games/{ID}/move", makeGameMove).Methods("POST")
	r.HandleFunc("/users/{userID}/games", getUserGames).Methods("GET")

	for {
		log.Printf("Running at 0.0.0.0:%d\n", port)
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), r))
		time.Sleep(1 * time.Second)
	}
}
