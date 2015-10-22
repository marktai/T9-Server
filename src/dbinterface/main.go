package dbinterface

import (
	"database/sql"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"fmt"
	"log"
	"strings"
)

var (
	db *sql.DB
)

type Person struct {
	Name  string
	Phone string
}

type Game struct {
	ID      uint
	Players [2]uint
	Turn    uint  // 0-19
	Board   Board //0 - 4^10-1
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

type Board [9]Box

func (g *dbgame) game() *Game {
	players := [2]uint{g.player0, g.player1}
	comprBoard := [9]uint{g.box0, g.box1, g.box2, g.box3, g.box4, g.box5, g.box6, g.box7, g.box8}
	var b Board
	b.Decompress(comprBoard)
	return &Game{g.gameid, players, g.turn, b}
}

func (g *dbgame) update() (sql.Result, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	updateGame, err := db.Prepare("UPDATE games SET player0=?, player1=?, turn=?, box0=?, box1=?, box2=?, box3=?, box4=?, box5=?, box6=?, box7=?, box8=? WHERE gameid=?")

	if err != nil {
		return nil, err
	}

	res, err := updateGame.Exec(g.player0, g.player1, g.turn, g.box0, g.box1, g.box2, g.box3, g.box4, g.box5, g.box6, g.box7, g.box8, g.gameid)

	if err != nil {
		return nil, err
	}
	return res, nil
}
func (g *Game) dbgame() *dbgame {
	comprBoard := g.Board.Compress()
	return &dbgame{g.ID, g.Players[0], g.Players[1], g.Turn, comprBoard[0], comprBoard[1], comprBoard[2], comprBoard[3], comprBoard[4], comprBoard[5], comprBoard[6], comprBoard[7], comprBoard[8]}
}

func (g *Game) Update() (sql.Result, error) {
	dg := g.dbgame()
	return dg.update()
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

	out := fmt.Sprintf("Owned = %d\n", b.Owned)
	out += b.String()

	// 	1 | 2 | 1
	// 	---------
	// 	2 | 1 | 2
	// 	---------
	// 	0 | 2 | 3
	log.Println(out)
}

func (b *Box) String() string {
	out := ""
	for i := 0; i < 3; i++ {
		if i != 0 {
			out += fmt.Sprintln(strings.Repeat("-", 9))
		}

		out += fmt.Sprintf("%d | %d | %d\n", b.Squares[3*i], b.Squares[3*i+1], b.Squares[3*i+2])
	}
	return out
}

func (b *Board) Compress() [9]uint {
	var comprBoard [9]uint
	for i, _ := range b {
		comprBoard[i] = b[i].Compress()
	}
	return comprBoard
}

func (b *Board) Decompress(compressed [9]uint) {
	for i, _ := range compressed {
		b[i].Decompress(compressed[i])
	}
}

func (b *Board) Print() {
	log.Println("\n" + b.String())
}

func (b *Board) String() string {
	out := ""
	for i := 0; i < 3; i++ {
		if i != 0 {
			out += "\n" + fmt.Sprintln(strings.Repeat("-", 37)) + "\n"
		}
		for j := 0; j < 3; j++ {
			if j != 0 {
				out += fmt.Sprint(strings.Repeat("-", 9)) + "     " + fmt.Sprint(strings.Repeat("-", 9)) + "     " + fmt.Sprint(strings.Repeat("-", 9)) + "\n"
			}

			line := ""

			for k := 0; k < 3; k++ {
				if k != 0 {
					line += "  |  "
				}
				line += fmt.Sprintf("%d | %d | %d", b[3*i+k].Squares[3*j], b[3*i+k].Squares[3*j+1], b[3*i+k].Squares[3*j+2])
			}
			out += line + "\n"

		}
	}
	return out
}

func GetGame(id uint) (*Game, error) {
	err := db.Ping()
	if err != nil {
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
		return
	}

	_, err = game.Update()
	if err != nil {
		log.Println(err)
		return
	}

	game.Board.Print()
}

func Tester() {
}
