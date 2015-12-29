package server

import (
	"fmt"
	"game"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"ws"
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

	if err != nil {
		WriteError(w, err, 400)
		return
	}

	WriteJson(w, genMap("ID", game.GameID))
}

func getAllGames(w http.ResponseWriter, r *http.Request) {
	games, err := game.GetAllGames()

	if err != nil {
		WriteError(w, err, 400)
		return
	}

	WriteJson(w, genMap("Games", games))
}

func getGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := stringtoUint(vars["ID"])
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	game, err := game.GetGame(id)

	if err != nil {
		WriteError(w, err, 400)
		return
	}

	WriteJson(w, genMap("Game", game.GameInfo))
}

func getBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := stringtoUint(vars["ID"])
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	game, err := game.GetGame(id)

	if err != nil {
		WriteError(w, err, 400)
		return
	}

	WriteJson(w, genMap("Board", game.Board))
}

func getGameString(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := stringtoUint(vars["ID"])
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	game, err := game.GetGame(id)

	if err != nil {
		WriteError(w, err, 400)
		return
	}

	WriteJson(w, genMap("Board", game.Board.StringArray(true)))
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

	if err != nil {
		WriteError(w, err, 400)
		return
	}

	err = game.MakeMove(player, box, square)
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	_, err = game.Update()
	WriteOutputError(w, genMap("Output", "Successful"), err)

	if err == nil {
		err = ws.Broadcast(id, []byte(fmt.Sprintf("Changed %d, %d", box, square)))
		if err != nil {
			log.Println(err)
		}

	}
}
