package main

import (
	"dbinterface"
	"flag"
	"log"
)

func main() {

	dbinterface.Open()
	defer dbinterface.Close()

	Tester()

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

	game, err := dbinterface.GetGame(63714)

	// game, err := dbinterface.MakeGame(0, 1)

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
