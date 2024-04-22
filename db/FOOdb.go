package db

import (
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/views/viewutils"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func FillXTabItems(userID uuid.UUID, dbpath string, tbd *viewutils.TabData, number int) error {
	var fooDB string = dbpath
	if fooDB == "" {
		fooDB = "/tmp"
	}
	var err error
	fooDB = fmt.Sprintf("%s/FOOdBAR/FOOdb.db", dbpath)
	fooDB, err = foodlib.CreateEmptyFileIfNotExists(fooDB)
	if err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", fooDB)
	defer db.Close()
	if err != nil {
		return err
	}
	tbd.Parent.LastActive = time.Now()

	return nil
}
