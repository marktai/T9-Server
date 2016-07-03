package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"recaptcha"
	"time"
	"ws"
)

var requireAuth bool

func Log(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func Run(port int, disableAuth bool) {
	//start := time.Now()

	rand.Seed(time.Now().UTC().UnixNano())
	r := mux.NewRouter()
	requireAuth = !disableAuth
	recaptcha.ReadSecret("./creds/recaptcha.json", "www.marktai.com")

	// user requests
	// r.HandleFunc("/login", Log(login)).Methods("POST")
	// r.HandleFunc("/verifySecret", Log(verifySecret)).Methods("POST")
	// r.HandleFunc("/users", Log(makeUser)).Methods("POST")

	// unauthorized requests
	// r.HandleFunc("/games", getAllGames).Methods("GET")
	r.HandleFunc("/games/{ID:[0-9]+}", getGame).Methods("GET") // only for backwards compatibility
	r.HandleFunc("/games/{ID:[0-9]+}/info", getGame).Methods("GET")
	r.HandleFunc("/games/{ID:[0-9]+}/board", getBoard).Methods("GET")
	r.HandleFunc("/games/{ID:[0-9]+}/string", getGameString).Methods("GET")
	r.HandleFunc("/games/{ID:[0-9]+}/ws", Log(ws.ServeWs)).Methods("GET")

	// authorized requests
	r.HandleFunc("/games", Log(makeGame)).Methods("POST")
	r.HandleFunc("/games/{ID:[0-9]+}", Log(makeGameMove)).Methods("POST") // only for backwards compatibility
	r.HandleFunc("/games/{ID:[0-9]+}/move", Log(makeGameMove)).Methods("POST")
	r.HandleFunc("/users/{userID:[0-9]+}/games", getUserGames).Methods("GET")

	for {
		log.Printf("Running at 0.0.0.0:%d\n", port)
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), r))
		time.Sleep(1 * time.Second)
	}
}
