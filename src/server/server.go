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

func Run(port uint16) {
	//start := time.Now()
	r := mux.NewRouter()
	r.HandleFunc("/game", makeGame).Methods("POST")
	r.HandleFunc("/game/{ID}", getGame).Methods("GET")
	r.HandleFunc("/game/{ID}/info", getGame).Methods("GET")
	r.HandleFunc("/game/{ID}", makeGameMove).Methods("POST")
	r.HandleFunc("/game/{ID}/string", getGameString).Methods("GET")
	r.HandleFunc("/game/{ID}/ws", ws.ServeWs).Methods("GET")
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
