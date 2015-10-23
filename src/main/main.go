package main

import (
	"flag"
	"game"
	"log"
	"server"
)

func main() {

	game.Open()
	defer game.Close()

	Tester()
	server.Run(8080)

}

func Tester() {

	var boxNumber uint
	var squareNumber uint
	var player uint

	flag.UintVar(&boxNumber, "boxNumber", 0, "boxNumber")
	flag.UintVar(&boxNumber, "b", 0, "boxNumber")
	flag.UintVar(&squareNumber, "squareNumber", 0, "squareNumber")
	flag.UintVar(&squareNumber, "s", 0, "squareNumber")
	flag.UintVar(&player, "player", 0, "player")
	flag.UintVar(&player, "p", 0, "player")
	flag.Parse()

	game, err := game.GetGame(63714)

	// game, err := game.MakeGame(0, 1)

	if err != nil {
		log.Println(err)
		return
	}
	game.Print()
	err = game.MakeMove(player, boxNumber, squareNumber)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = game.Update()
	if err != nil {
		log.Println(err)
		return
	}

	game.Print()
}
