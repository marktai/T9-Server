package server

import (
	"auth"
	"encoding/json"
	"fmt"
	"game"
	"github.com/gorilla/mux"
	// "io/ioutil"
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

	authed, err := auth.AuthRequest(r)
	if err != nil || !authed {

		if err != nil {
			log.Println(err)
		}
		WriteErrorString(w, "Not Authorized Request", 400)
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
		ws.Broadcast(id, []byte("Changed"))
	}
}

func makeUser(w http.ResponseWriter, r *http.Request) {

	secret := r.FormValue("Secret")
	if secret != "thisisatotallysecuresecret" {
		WriteErrorString(w, "Sorry, you can't make a user now", 500)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var parsedJson map[string]string
	err := decoder.Decode(&parsedJson)
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	user, ok := parsedJson["User"]
	if !ok {
		WriteErrorString(w, "No 'User' set in POST body", 400)
		return
	}
	pass, ok := parsedJson["Password"]
	if !ok {
		WriteErrorString(w, "No 'Password' set in POST body", 400)
		return
	}

	userID, err := auth.MakeUser(user, pass)
	if err != nil {
		WriteError(w, err, 500)
		return
	}

	WriteJson(w, genMap("UserID", userID))

}

func login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var parsedJson map[string]string
	err := decoder.Decode(&parsedJson)
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	user, ok := parsedJson["User"]
	if !ok {
		WriteErrorString(w, "No 'User' set in POST body", 400)
		return
	}
	pass, ok := parsedJson["Password"]
	if !ok {
		WriteErrorString(w, "No 'Password' set in POST body", 400)
		return
	}

	userID, secret, err := auth.Login(user, pass)
	if err != nil {
		// hides details about server from login attempts"
		log.Println(err)
		WriteErrorString(w, "Username and password combination incorrect", 500)
		return
	}

	retMap := map[string]string{"UserID": fmt.Sprintf("%d", userID), "Secret": secret.Base64()}
	WriteJson(w, retMap)
}
