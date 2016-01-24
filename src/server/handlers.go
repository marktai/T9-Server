package server

import (
	"auth"
	"encoding/json"
	"fmt"
	"game"
	"github.com/gorilla/mux"
	// "io/ioutil"
	"encoding/base64"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"ws"
)

func makeGame(w http.ResponseWriter, r *http.Request) {
	player1, err := stringtoUint(r.FormValue("Player1"))
	if err != nil {
		WriteError(w, errors.New("Error parsing Player1 form value"), 400)
		return
	}

	player2, err := stringtoUint(r.FormValue("Player2"))
	if err != nil {
		WriteError(w, errors.New("Error parsing Player2 form value"), 400)
		return
	}

	starter, err := stringtoUint(r.FormValue("Starter"))
	// 0 means random
	// 1 means player1 starts
	// 2 means player2 starts
	if err != nil {
		starter = 0
	}

	if starter > 3 {
		WriteError(w, errors.New("Starter must be 0-2"), 400)
	}

	if requireAuth {
		authed, err := auth.AuthRequest(r, player1)
		if err != nil || !authed {
			if err != nil {
				log.Println(err)
			}
			WriteErrorString(w, "Not Authorized Request", 401)
			return
		}
	}

	// 0 means that 50% chance of switching
	switch starter {
	case 0:
		if rand.Intn(2) == 1 {
			break
		}
		fallthrough
	case 2:
		player1, player2 = player2, player1
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
		WriteError(w, err, 400)
		return
	}

	player, err := stringtoUint(r.FormValue("Player"))
	if err != nil {
		WriteError(w, errors.New("Error parsing Player form value"), 400)
		return
	}

	box, err := stringtoUint(r.FormValue("Box"))
	if err != nil {
		WriteError(w, errors.New("Error parsing Box form value"), 400)
		return
	}

	square, err := stringtoUint(r.FormValue("Square"))
	if err != nil {
		WriteError(w, errors.New("Error parsing Square form value"), 400)
		return
	}

	if requireAuth {
		authed, err := auth.AuthRequest(r, player)
		if err != nil || !authed {
			if err != nil {
				log.Println(err)
			}
			WriteErrorString(w, "Not Authorized Request", 401)
			return
		}
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

	var parsedJson map[string]string
	if encodedAuth := r.Header.Get("Authorization"); encodedAuth != "" {
		authBytes, err := base64.StdEncoding.DecodeString(encodedAuth)
		if err != nil {
			//ERROR
		}
		auth := string(authBytes[:])
		if strings.Count(auth, ":") != 1 {
			// ERROR
		}
		authSlice := strings.Split(auth, ":")

		parsedJson = make(map[string]string)
		parsedJson["User"] = authSlice[0]
		parsedJson["Password"] = authSlice[1]
	} else {

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&parsedJson)
		if err != nil {
			WriteErrorString(w, err.Error()+" in parsing POST body (JSON)", 400)
			return
		}
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
		WriteErrorString(w, "User and password combination incorrect", 401)
		return
	}

	retMap := map[string]string{"UserID": fmt.Sprintf("%d", userID), "Secret": secret.Base64(), "Expiration": secret.ExpirationUTC()}
	WriteJson(w, retMap)
}

func verifySecret(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var parsedJson map[string]string
	err := decoder.Decode(&parsedJson)
	if err != nil {
		WriteErrorString(w, err.Error()+" in parsing POST body (JSON)", 400)
		return
	}

	user, ok := parsedJson["User"]
	if !ok {
		WriteErrorString(w, "No 'User' set in POST body", 400)
		return
	}
	inpSecret, ok := parsedJson["Secret"]
	if !ok {
		WriteErrorString(w, "No 'Secret' set in POST body", 400)
		return
	}

	userID, secret, err := auth.VerifySecret(user, inpSecret)
	if err != nil {
		// hides details about server from login attempts"
		log.Println(err)
		WriteErrorString(w, "User and secret combination incorrect", 400)
		return
	}

	retMap := map[string]string{"UserID": fmt.Sprintf("%d", userID), "Secret": secret.Base64(), "Expiration": secret.ExpirationUTC()}
	WriteJson(w, retMap)
}

func getUserGames(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userID, err := stringtoUint(vars["userID"])
	if err != nil {
		WriteError(w, err, 400)
	}

	if requireAuth {

		authed, err := auth.AuthRequest(r, userID)
		if err != nil || !authed {

			if err != nil {
				log.Println(err)
			}
			WriteErrorString(w, "Not Authorized Request", 401)
			return
		}
	}

	games, err := game.GetUserGames(userID)

	if err != nil {
		WriteError(w, err, 400)
		return
	}

	WriteJson(w, genMap("Games", games))
}
