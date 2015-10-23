package game

import (
	"database/sql"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// TODO: move to queue
type MoveHistory [18]uint

type Game struct {
	GameID      uint
	Players     [2]uint
	Turn        uint  // 0-19
	Board       Board //0 - 4^10-1
	MoveHistory MoveHistory
	Started     time.Time
	Modified    time.Time
}

func (m *MoveHistory) Decompress(a, b uint64) {
	compressed := [2]uint64{a, b}

	for i := 0; i < 2; i++ {
		for moveIndex := 0; moveIndex < 9; moveIndex++ {
			move := compressed[i] & 0x7F
			m[moveIndex+9*i] = uint(move)
			compressed[i] >>= 7
		}
	}
}

func (m *MoveHistory) Compress() (uint64, uint64) {
	var compressed [2]uint64
	j := 1
	for i := 17; i >= 0; i-- {
		if i == 8 {
			j -= 1
		}
		compressed[j] *= 128
		compressed[j] += uint64(m[i])
	}
	return compressed[0], compressed[1]
}

func (m *MoveHistory) AddMove(move uint) {
	for i := 17; i >= 1; i-- {
		m[i] = m[i-1]
	}
	m[0] = move
}

func symbol(input uint, translate bool) string {
	if !translate {
		return fmt.Sprintf("%d", input)
	}
	switch input {
	case 0:
		return " "
	case 1:
		return "x"
	case 2:
		return "o"
	case 3:
		return "/"
	default:
		return "?"
	}
}

func (g *Game) dbgame() *dbgame {
	comprBoard := g.Board.Compress()
	m1, m2 := g.MoveHistory.Compress()
	return &dbgame{g.GameID, g.Players[0], g.Players[1], g.Turn, comprBoard[0], comprBoard[1], comprBoard[2], comprBoard[3], comprBoard[4], comprBoard[5], comprBoard[6], comprBoard[7], comprBoard[8], m1, m2, g.Started, g.Modified}
}

func (g *Game) Update() (sql.Result, error) {
	dg := g.dbgame()
	return dg.update()
}

func (g *Game) MakeMove(player, box, square uint) error {
	playerTurn := g.Turn / 10 % 2
	if player != g.Players[playerTurn] {
		return errors.New("Not player's turn")
	}

	moveBox := g.Turn % 10

	if moveBox != 9 && box != moveBox {
		return errors.New("Not correct box")
	}

	if box > 8 {
		return errors.New("Box out of range")
	}

	if g.Board[box].Owned != 0 {
		return errors.New("Box already taken")
	}

	if square > 8 {
		return errors.New("Square out of range")
	}

	err := g.Board[box].MakeMove(playerTurn+1, square)
	if err != nil {
		return err
	}

	g.MoveHistory.AddMove(9*box + square)
	g.Modified = time.Now().UTC()

	g.Turn = (1 - playerTurn) * 10
	if g.Board[square].Owned != 0 {
		g.Turn += 9
	} else {
		g.Turn += square
	}

	return nil

}

func (g *Game) Print() {
	log.Println("GameID:", g.GameID)
	log.Println("Players:", g.Players)
	log.Println("Turn:", g.Turn)
	log.Println("Started:", g.Started)
	log.Println("Modified:", g.Modified)
	g.Board.Print()
	log.Println(g.MoveHistory)
}

func GetGame(id uint) (*Game, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	var game dbgame

	var started, modified string

	//TODO: handle NULLS
	err = db.QueryRow("SELECT gameid, player0, player1, turn, box0, box1, box2, box3, box4, box5, box6, box7, box8, movehistory0, movehistory1, started, modified FROM games WHERE gameid=?", id).Scan(&game.gameid, &game.player0, &game.player1, &game.turn, &game.box0, &game.box1, &game.box2, &game.box3, &game.box4, &game.box5, &game.box6, &game.box7, &game.box8, &game.movehistory0, &game.movehistory1, &started, &modified)
	if err != nil {
		return nil, err
	}
	//golang constant thingy
	//reference time is "Mon Jan 2 15:04:05 -0700 MST 2006"
	const sqlForm = "2006-01-02 15:04:05"

	game.started, err = time.Parse(sqlForm, started)
	if err != nil {
		return nil, err
	}
	game.modified, err = time.Parse(sqlForm, modified)
	if err != nil {
		return nil, err
	}

	return game.game(), nil
}

func MakeGame(player0, player1 uint) (*Game, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	var g Game

	rand.Seed(time.Now().Unix())

	collision := 1
	times := 0
	for collision != 0 {

		g.GameID = uint(rand.Int31n(65536))
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM games WHERE gameid=?)", g.GameID).Scan(&collision)
		if err != nil {
			return nil, err
		}
		times++
		if times > 20 {
			return nil, errors.New("Too many attempts to find a unique game ID")
		}
	}
	g.Players = [2]uint{player0, player1}
	g.Turn = 9
	g.Started = time.Now().UTC()
	g.Modified = time.Now().UTC()
	g.MoveHistory.AddMove(127)

	addGame, err := db.Prepare("INSERT INTO games VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return nil, err
	}

	dg := g.dbgame()

	_, err = addGame.Exec(dg.gameid, dg.player0, dg.player1, dg.turn, dg.box0, dg.box1, dg.box2, dg.box3, dg.box4, dg.box5, dg.box6, dg.box7, dg.box8, dg.movehistory0, dg.movehistory1, dg.started, dg.modified)

	if err != nil {
		return nil, err
	}

	return &g, nil

}
