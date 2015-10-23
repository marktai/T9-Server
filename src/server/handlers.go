package server

import (
	"game"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func makeGame(w http.ResponseWriter, r *http.Request) {
	player1, err := stringtoUint(r.FormValue("Player1"))
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	player2, err := stringtoUint(r.FormValue("Player2"))
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	game, err := game.MakeGame(player1, player2)
	WriteOutputError(w, map[string]uint{"ID": game.GameID}, err)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := stringtoUint(vars["ID"])
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	game, err := game.GetGame(id)
	WriteOutputError(w, map[string]interface{}{"Game": game.Info()}, err)
}

func makeGameMove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := stringtoUint(vars["ID"])
	if err != nil {
		log.Println(id)
		WriteError(w, err, 400)
		return
	}

	player, err := stringtoUint(r.FormValue("Player"))
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	box, err := stringtoUint(r.FormValue("Box"))
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	square, err := stringtoUint(r.FormValue("Square"))
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	game, err := game.GetGame(id)

	err = game.MakeMove(player, box, square)
	if err != nil {
		WriteError(w, err, 500)
		return
	}

	_, err = game.Update()
	WriteOutputError(w, "Successful", err)
}
