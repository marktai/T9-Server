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

	var boxNumber int
	var squareNumber int
	var newValue uint

	flag.IntVar(&boxNumber, "boxNumber", 0, "boxNumber")
	flag.IntVar(&boxNumber, "b", 0, "boxNumber")
	flag.IntVar(&squareNumber, "squareNumber", 0, "squareNumber")
	flag.IntVar(&squareNumber, "s", 0, "squareNumber")
	flag.UintVar(&newValue, "newValue", 0, "newValue")
	flag.UintVar(&newValue, "n", 0, "newValue")
	flag.Parse()

	game, err := dbinterface.GetGame(0)
	if err != nil {
		log.Println(err)
		return
	}

	game.Board[boxNumber].Squares[squareNumber] = newValue

	_, err = game.Update()
	if err != nil {
		log.Println(err)
		return
	}

	game.Board.Print()
}
