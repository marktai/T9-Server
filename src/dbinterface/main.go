package dbinterface

import (
	"database/sql"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"fmt"
	"log"
)

var (
	db *sql.DB
)

type Person struct {
	Name  string
	Phone string
}

type Game struct {
	ID         uint
	Players    [2]uint
	Turn       uint    // 0-19
	comprBoard [9]uint //0 - 4^10-1
}

type dbgame struct {
	gameid  uint
	player0 uint
	player1 uint
	turn    uint
	box0    uint
	box1    uint
	box2    uint
	box3    uint
	box4    uint
	box5    uint
	box6    uint
	box7    uint
	box8    uint
}

// Turn
// 0_ -> player 1 turn
// _0-8 -> box 0-8, 9 -> anywhere

type Box struct {
	Owned   uint
	Squares [9]uint
}

type Board struct {
	Boxes [9]*Box
}

func (g dbgame) game() *Game {
	players := [2]uint{g.player0, g.player1}
	comprBoard := [9]uint{g.box0, g.box1, g.box2, g.box3, g.box4, g.box5, g.box6, g.box7, g.box8}
	return &Game{g.gameid, players, g.turn, comprBoard}
}

func (g *Game) GetBoard() (Board, error) {

	var board Board

	for i, comprBox := range g.comprBoard {
		board.Boxes[i].Decompress(comprBox)
	}
	return board, nil
}

func (b *Box) Decompress(compressed uint) error {

	i := 9
	for i >= 0 {
		num := compressed & 0x03
		if i == 9 {
			b.Owned = num
		} else {
			b.Squares[i] = num
		}
		compressed >>= 2
		i--
	}
	log.Println(b.Squares)
	return nil
}

func (b *Box) Compress() uint {

	var total uint

	for _, square := range b.Squares {
		total += square
		total *= 4
	}
	total += b.Owned

	return total
}

func (b *Box) Print() {

	out := "\n"

	for i := 0; i < 3; i++ {
		if i != 0 {
			out += fmt.Sprintln("---------")
		}

		out += fmt.Sprintf("%d | %d | %d\n", b.Squares[3*i], b.Squares[3*i+1], b.Squares[3*i+2])
	}
	// 	1 | 2 | 1
	// 	---------
	// 	2 | 1 | 2
	// 	---------
	// 	0 | 2 | 3
	log.Println(out)
}

func GetGame(id uint) (*Game, error) {
	err := db.Ping()
	if err != nil {
		// do something here
		return nil, err
	}

	var game dbgame

	err = db.QueryRow("select * from games where gameid=?", id).Scan(&game.gameid, &game.player0, &game.player1, &game.turn, &game.box0, &game.box1, &game.box2, &game.box3, &game.box4, &game.box5, &game.box6, &game.box7, &game.box8)
	if err != nil {
		log.Fatal(err)
	}
	return game.game(), nil
}

func Server() {
	var err error
	db, err = sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/TT2")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	game, err := GetGame(0)
	if err != nil {
		log.Println(err)
	}

	log.Println(game)
	// game.GetBoard()

	b := &Box{0, [9]uint{0, 1, 1, 0, 2, 2, 2, 1, 0}}
	b.Print()
	compressed := b.Compress()
	log.Printf("Compressed is %d\n", compressed)

	c := &Box{}
	c.Decompress(compressed)
	c.Print()
}

func Tester() {
}

// func Tester() {
// 	session, err := mgo.Dial("localhost:27017")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()

// 	// Optional. Switch the session to a monotonic behavior.
// 	session.SetMode(mgo.Monotonic, true)

// 	c := session.DB("test").C("people")
// 	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
// 		&Person{"Cla", "+55 53 8402 8510"})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	result := Person{}
// 	err = c.Find(bson.M{"name": "Ale"}).One(&result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Phone:", result.Phone)
// }
