package game

import (
	"database/sql"
	"db"
	_ "github.com/go-sql-driver/mysql"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	//	"errors"
	//	"math/rand"
	"time"
)

//exactly how it is in the database
type dbgame struct {
	gameid       uint
	player0      uint
	player1      uint
	turn         uint
	box0         uint
	box1         uint
	box2         uint
	box3         uint
	box4         uint
	box5         uint
	box6         uint
	box7         uint
	box8         uint
	movehistory0 uint64
	movehistory1 uint64
	started      time.Time
	modified     time.Time
}

// checks if a id already exists in the database
func checkIDConflict(id uint) (bool, error) {
	collision := 1
	err := db.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM games WHERE gameid=?)", id).Scan(&collision)
	return collision != 0, err
}

// gets a unique id for a new game
func getUniqueID() (uint, error) {
	var count uint
	var scale uint
	var addConst uint

	var newID uint

	conflict := true
	err := db.Db.QueryRow("SELECT count, scale, addConst FROM count WHERE type='games'").Scan(&count, &scale, &addConst)
	if err != nil {
		return 0, err
	}

	for conflict {
		count += 1
		newID = (count*scale + addConst) % 65536
		conflict, err = checkIDConflict(newID)
		if newID == 0 {
			continue
		}
		if err != nil {
			return 0, err
		}
	}

	updateCount, err := db.Db.Prepare("UPDATE count SET count=? WHERE type='games'")
	if err != nil {
		return newID, err
	}

	_, err = updateCount.Exec(count)
	if err != nil {
		return newID, err
	}

	return newID, nil
}

// returns the game representation of the dbgame
func (g *dbgame) game() *Game {

	var newGame Game

	newGame.GameID = g.gameid
	newGame.Players = [2]uint{g.player0, g.player1}

	comprBoard := [9]uint{g.box0, g.box1, g.box2, g.box3, g.box4, g.box5, g.box6, g.box7, g.box8}
	newGame.Board.Decompress(comprBoard)

	newGame.Turn = g.turn
	newGame.MoveHistory.Decompress(g.movehistory0, g.movehistory1)
	newGame.Started = g.started
	newGame.Modified = g.modified

	// checks whether the game is won
	newGame.CheckVictor()

	return &newGame
}

// updates database version to be equal the dbgame
func (g *dbgame) update() (sql.Result, error) {
	err := db.Db.Ping()
	if err != nil {
		return nil, err
	}

	updateGame, err := db.Db.Prepare("UPDATE games SET turn=?, box0=?, box1=?, box2=?, box3=?, box4=?, box5=?, box6=?, box7=?, box8=?, movehistory0=?, movehistory1=?, modified=? WHERE gameid=?")
	if err != nil {
		return nil, err
	}

	return updateGame.Exec(g.turn, g.box0, g.box1, g.box2, g.box3, g.box4, g.box5, g.box6, g.box7, g.box8, g.movehistory0, g.movehistory1, g.modified, g.gameid)
}

// adds the game into the database
func (g *dbgame) upload() (sql.Result, error) {
	err := db.Db.Ping()
	if err != nil {
		return nil, err
	}

	addGame, err := db.Db.Prepare("INSERT INTO games VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	return addGame.Exec(g.gameid, g.player0, g.player1, g.turn, g.box0, g.box1, g.box2, g.box3, g.box4, g.box5, g.box6, g.box7, g.box8, g.movehistory0, g.movehistory1, g.started, g.modified)
}
