package db

import (
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/views/viewutils"
	"database/sql"
	"fmt"

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

	switch tbd.Ttype {
		case viewutils.Recipe:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tbd.Ttype.String()))
		case viewutils.Menu:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tbd.Ttype.String()))
		case viewutils.Pantry:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tbd.Ttype.String()))
		case viewutils.Preplist:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tbd.Ttype.String()))
		case viewutils.Shopping:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tbd.Ttype.String()))
		case viewutils.Earnings:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tbd.Ttype.String()))
	}

	return nil
}
