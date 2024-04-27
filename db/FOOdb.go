package db

import (
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/views/viewutils"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func CreateTabTableIfNotExists(userID uuid.UUID, dbpath string, tt viewutils.TabType) (*sql.DB, error) {
	var fooDB string = dbpath
	if fooDB == "" {
		fooDB = "/tmp"
	}
	var err error
	fooDB = fmt.Sprintf("%s/FOOdBAR/FOOdb.db", dbpath)
	fooDB, err = foodlib.CreateEmptyFileIfNotExists(fooDB)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", fooDB)

	switch tt {
		case viewutils.Recipe:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tt.String()))
		case viewutils.Menu:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tt.String()))
		case viewutils.Pantry:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tt.String()))
		case viewutils.Preplist:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tt.String()))
		case viewutils.Shopping:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tt.String()))
		case viewutils.Earnings:
			_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s (
				id TEXT PRIMARY KEY,
				)`, userID, tt.String()))
	}
	if err != nil {
		return nil, err
	}

	return db, err
}

/*
NOTE:
type TabItem struct {
	ItemID uuid.UUID `json:"item_id"`
	Ttype  TabType   `json:"tab_type"`
	Expanded bool   `json:"expanded"`
}
	NOTE:
type SortMethod string
const (
	Inactive   SortMethod = ""
	Descending            = "DESC;"
	Ascending             = "ASC;"
	Random                = "RANDOM();"
	Custom                = "END;"

	// Others can be made with CASE WHEN condition THEN value ELSE value END
	// when using this syntax, put the CASE WHEN... etc... into the OrderBy key
	// and then put Custom as the SortMethod
)
*/
func FillXTabItems(userID uuid.UUID, dbpath string, tbd *viewutils.TabData, number int) error {
	db, err := CreateTabTableIfNotExists(userID, dbpath, tbd.Ttype)
	defer db.Close()
	if err != nil {
		return err
	}
	// TODO: fill tbd.Items with X number of items based on tbd.OrderBy: map[string]SortMethod
	// where string is a column name or other sql identifier that can be sorted by


	return nil
}
