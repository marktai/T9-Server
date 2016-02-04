package game_test

import (
	"db"
	"game"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestMakeGame(t *testing.T) {
	newDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexepected error using sqlmock.New(): %s", err)
	}
	db.Db = newDb

	rows := sqlmock.NewRows([]string{"count", "scale", "addConst"}).AddRow(1, 2, 3)
	mock.ExpectQuery("^SELECT count, scale, addConst FROM count WHERE type='games'$").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"EXISTS(SELECT 1 FROM games WHERE gameid=7)"}).AddRow(0)
	mock.ExpectQuery("^SELECT EXISTS\\(SELECT 1 FROM games WHERE gameid=\\?\\)").WillReturnRows(rows)

	mock.ExpectPrepare("^UPDATE count SET count=[0-9]* WHERE type='games'$")
	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec("UPDATE").WithArgs(2).WillReturnResult(result)

	newGame, err := game.MakeGame(0, 1)
	if err != nil {
		t.Errorf("Unexpected error when creating game: %s", err)
	}

	_ = newGame

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	db.Db.Close()
}
