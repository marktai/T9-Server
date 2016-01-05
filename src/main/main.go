package main

import (
	"auth"
	"crypto/hmac"
	"crypto/sha256"
	"db"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"server"
)

func main() {
	// // makeUser()
	// testHMAC()

	// _, secret, err := auth.Login("me", "password")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.Printf("%s", secret.String())

	runServer()
}

func makeUser() {
	_, err := auth.MakeUser("me", "password")
	if err != nil {
		log.Println(err)
	}
}

func testHMAC() {
	id, secret, err := auth.Login("me", "password")
	if err != nil {
		log.Println(err)
		return
	}
	r, err := http.NewRequest("GET", fmt.Sprintf("localhost/here?ID=%d", id), nil)
	if err != nil {
		log.Println(err)
		return
	}

	time := "Naow1"

	r.Header.Add("Time-Sent", time)

	message := append([]byte(time), []byte(fmt.Sprintf("localhost/here?ID=%d", id))...)

	mac := hmac.New(sha256.New, secret.Bytes())
	mac.Write(message)
	hmacBytes := mac.Sum(nil)

	hmacString := base64.StdEncoding.EncodeToString(hmacBytes)
	r.Header.Add("HMAC", hmacString)

	// log.Println(r)

	authed, err := auth.AuthRequest(r, id)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(authed)

}

func runServer() {
	var port int

	flag.IntVar(&port, "Port", 8081, "Port the server listens to")

	flag.Parse()

	db.Open()
	defer db.Db.Close()
	server.Run(port)
}
