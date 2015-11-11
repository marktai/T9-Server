package game

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type GameInfo struct {
	GameID      uint
	Players     [2]uint
	Turn        uint // 0-19
	MoveHistory MoveHistory
	Started     time.Time
	Modified    time.Time
}

type Game struct {
	GameInfo
	Board Board //0 - 4^10-1
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

func (g *Game) Info() GameInfo {
	return g.GameInfo
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
		if err == sql.ErrNoRows {
			return nil, errors.New("Game not found")
		}
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
	g.GameID, err = getUniqueID()
	if err != nil {
		return nil, err
	}
	g.Players = [2]uint{player0, player1}
	g.Turn = 9
	g.Started = time.Now().UTC()
	g.Modified = time.Now().UTC()
	g.MoveHistory.AddMove(127)

	if err != nil {
		return nil, err
	}

	g.dbgame().upload()

	if err != nil {
		return nil, err
	}

	return &g, nil

}

func GetAllGames() ([]uint, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	} //TODO: handle NULLS

	var ids []uint

	rows, err := db.Query("SELECT gameid FROM games")
	defer rows.Close()
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return ids, err
	}

	return ids, nil

}
