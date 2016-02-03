package game

import (
	"database/sql"
	"db"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"
)

// Game is the class that represents all of a T9 game data

// Metadata about the game
// Pretty much everything besides the Board
type GameInfo struct {
	GameID      uint
	Players     [2]uint
	Turn        uint // 0-19
	MoveHistory MoveHistory
	Started     time.Time
	Modified    time.Time
}

// GameInfo and a Board
type Game struct {
	GameInfo
	Board Board //0 - 4^10-1
}

// Returns a symbol based on the uint representation
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

// Returns a dbgame representation of the Game
func (g *Game) dbgame() *dbgame {
	comprBoard := g.Board.Compress()
	m1, m2 := g.MoveHistory.Compress()
	return &dbgame{g.GameID, g.Players[0], g.Players[1], g.Turn, comprBoard[0], comprBoard[1], comprBoard[2], comprBoard[3], comprBoard[4], comprBoard[5], comprBoard[6], comprBoard[7], comprBoard[8], m1, m2, g.Started, g.Modified}
}

// Updates the database version to equal the current Game
func (g *Game) Update() (sql.Result, error) {
	dg := g.dbgame()
	return dg.update()
}

// Validates and then makes a move
func (g *Game) MakeMove(player, box, square uint) error {
	if g.Board.Box().CheckOwned() != 0 {
		return errors.New("Game already finished")
	}

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

	g.CheckVictor()

	return nil

}

func (g *Game) CheckVictor() uint {

	victor := g.Board.Box().CheckOwned()
	if victor != 0 {
		g.Turn = 20 + victor
	}

	return g.Turn
}

// Returns the info of the game
func (g *Game) Info() GameInfo {
	return g.GameInfo
}

// Prints info about the game and the board
func (g *Game) Print() {
	log.Println("GameID:", g.GameID)
	log.Println("Players:", g.Players)
	log.Println("Turn:", g.Turn)
	log.Println("Started:", g.Started)
	log.Println("Modified:", g.Modified)
	g.Board.Print()
	log.Println(g.MoveHistory)
}

// Gets a Game frome the database
func GetGame(id uint) (*Game, error) {
	err := db.Db.Ping()
	if err != nil {
		return nil, err
	}

	var game dbgame

	var started, modified string

	//TODO: handle NULLS
	err = db.Db.QueryRow("SELECT gameid, player0, player1, turn, box0, box1, box2, box3, box4, box5, box6, box7, box8, movehistory0, movehistory1, started, modified FROM games WHERE gameid=?", id).Scan(&game.gameid, &game.player0, &game.player1, &game.turn, &game.box0, &game.box1, &game.box2, &game.box3, &game.box4, &game.box5, &game.box6, &game.box7, &game.box8, &game.movehistory0, &game.movehistory1, &started, &modified)
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

// Creates a new game and uploads it to the database
func MakeGame(player0, player1 uint) (*Game, error) {
	if player0 == player1 {
		return nil, errors.New("player1 cannot be the same as player2 in creating game")
	}

	err := db.Db.Ping()
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

	_, err = g.dbgame().upload()

	if err != nil {
		return nil, err
	}

	return &g, nil

}

// Used for sorting game id's by time modified
type idModified struct {
	id       uint
	modified time.Time
}

type idModifiedSlice []*idModified

//sorts by most recent first
func (a idModifiedSlice) Len() int           { return len(a) }
func (a idModifiedSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a idModifiedSlice) Less(i, j int) bool { return a[i].modified.After(a[j].modified) }

// Returns all games in the database sorted by most recently modified
func GetAllGames() ([]uint, error) {
	err := db.Db.Ping()
	if err != nil {
		return nil, err
	} //TODO: handle NULLS

	var games idModifiedSlice

	rows, err := db.Db.Query("SELECT gameid, modified FROM games")
	defer rows.Close()
	for rows.Next() {
		var tempgameactual idModified
		tempgame := &tempgameactual
		var modified string
		if err := rows.Scan(&(tempgame.id), &modified); err != nil {
			return nil, err
		}
		//golang constant thingy
		//reference time is "Mon Jan 2 15:04:05 -0700 MST 2006"
		const sqlForm = "2006-01-02 15:04:05"

		tempgame.modified, err = time.Parse(sqlForm, modified)
		if err != nil {
			return nil, err
		}
		games = append(games, tempgame)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Sort(games)

	var ids []uint

	for _, tempgame := range games {
		ids = append(ids, tempgame.id)
	}

	return ids, nil

}
func GetUserGames(userID uint) ([]uint, error) {
	err := db.Db.Ping()
	if err != nil {
		return nil, err
	} //TODO: handle NULLS

	var ids []uint

	rows, err := db.Db.Query("SELECT gameid FROM games WHERE player0=? OR player1=?", userID, userID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
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
