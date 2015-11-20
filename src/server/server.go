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

func Run(port int) {
	//start := time.Now()
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/users", makeUser).Methods("POST")
	r.HandleFunc("/games", getAllGames).Methods("GET")
	r.HandleFunc("/games", makeGame).Methods("POST")
	r.HandleFunc("/games/{ID}", getGame).Methods("GET")
	r.HandleFunc("/games/{ID}/info", getGame).Methods("GET")
	r.HandleFunc("/games/{ID}", makeGameMove).Methods("POST")
	r.HandleFunc("/games/{ID}/string", getGameString).Methods("GET")
	r.HandleFunc("/games/{ID}/ws", ws.ServeWs).Methods("GET")
	//log.Println("Took %s", time.Now().Sub(start))
	//log.Println(post)

	// r.HandleFunc("/posts", getPostList)
	// r.HandleFunc("/posts/{Title}", getPost)
	// r.HandleFunc("/posts/{Title}/paragraph/{id:[0-9]+}", getParagraph).Methods("GET")
	// r.HandleFunc("/posts/{Title}/info", getInfo).Methods("GET")
	// r.HandleFunc("/desktopIP", getIP).Methods("GET")
	// r.HandleFunc("/desktopIP", setIP).Methods("POST")
	// r.HandleFunc("/desktopIP", clearIP).Methods("DELETE")

	for {
		log.Printf("Running at 0.0.0.0:%d\n", port)
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), r))
		time.Sleep(1 * time.Second)
	}
}
