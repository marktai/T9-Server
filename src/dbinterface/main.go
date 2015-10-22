package dbinterface

import (
	"database/sql"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"log"
)

type Person struct {
	Name  string
	Phone string
}

type Game struct {
	ID         int
	comprBoard []int //0 - 4^10-1
	Turn       int   // 0-19
	Players    []int
}

type dbgame struct {
	gameid  int
	player0 int
	player1 int
	turn    int
	box0    int
	box1    int
	box2    int
	box3    int
	box4    int
	box5    int
	box6    int
	box7    int
	box8    int
}

// Turn
// 0_ -> player 1 turn
// _0-8 -> box 0-8, 9 -> anywhere

type Box struct {
	Squares []int
	Owned   int
}

type Board struct {
	Boxes []Box
}

func (g dbgame) game() Game {
	players := []int{g.player0, g.player1}
	comprBoard := []int{g.box0, g.box1, g.box2, g.box3, g.box4, g.box5, g.box6, g.box7, g.box8}
	return Game{g.gameid, comprBoard, g.turn, players}
}

func (g Game) GetBoard() (Board, error) {
	return Board{nil}, nil
}

func Tester() {
	db, err := sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/TT2")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		// do something here
		log.Fatal(err)
	}

	var id int
	err = db.QueryRow("select gameid from games").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(id)
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
