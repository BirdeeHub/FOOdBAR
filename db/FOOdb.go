package db

import (
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/views/viewutils"
	"database/sql"
	"errors"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func CreateTabTableIfNotExists(userID uuid.UUID, dbpath string, tt viewutils.TabType) (*sql.DB, error) {
	var fooDB string = dbpath
	if fooDB == "" {
		fooDB = "/tmp"
	}
	var err error
	fooDB = filepath.Join(dbpath, "FOOdBAR", "FOOdb.db")
	fooDB, err = foodlib.CreateEmptyFileIfNotExists(fooDB)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", fooDB)

	var createTable string
	switch tt {
	case viewutils.Recipe:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				name TEXT UNIQUE,
				category TEXT,
				dietary TEXT,
				ingredients TEXT,
				instructions TEXT,
				)`
	case viewutils.Menu:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				menu_id TEXT,
				name TEXT,
				number INTEGER,
				)`
	case viewutils.Pantry:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				name TEXT UNIQUE,
				dietary TEXT,
				amount TEXT,
				units TEXT,
				)`
	case viewutils.Customer:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				email TEXT UNIQUE,
				phone TEXT UNIQUE,
				name TEXT,
				dietary TEXT,
				)`
	case viewutils.Events:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				name TEXT,
				menu_id TEXT,
				date TEXT,
				location TEXT,
				customer TEXT,
				)`
	case viewutils.Preplist:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT,
				)`
	case viewutils.Shopping:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT,
				ingredient TEXT,
				amount TEXT,
				units TEXT,
				)`
	case viewutils.Earnings:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT,
				)`
	default:
		return nil, errors.New("invalid tab type")
	}
	createStmt, err := db.Prepare(createTable)
	if err != nil {
		return nil, err
	}
	defer createStmt.Close()
	_, err = createStmt.Exec(userID, tt.String())
	if err != nil {
		return nil, err
	}
	err = makeAuditTriggers(db, userID, tt)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func makeAuditTriggers(db *sql.DB, userID uuid.UUID, tt viewutils.TabType) error {

	updateTrigger := `CREATE TRIGGER IF NOT EXISTS update_?_?_audit 
		AFTER INSERT OR UPDATE OR DELETE ON ?_?
		FOR EACH ROW
		BEGIN
			UPDATE ?_?
			SET last_modified = CURRENT_TIMESTAMP
			WHERE id = OLD.id;
		END;`

	insertTrigger := `CREATE TRIGGER IF NOT EXISTS add_created_?_?_audit
		AFTER INSERT ON ?_?
		FOR EACH ROW  
		BEGIN
			UPDATE ?_?
			SET created_at = CURRENT_TIMESTAMP
			WHERE id = OLD.id;
		END;`

	updateStmt, err := db.Prepare(updateTrigger)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	_, err = updateStmt.Exec(userID, tt.String(), userID, tt.String(), userID, tt.String())
	if err != nil {
		return err
	}

	insertStmt, err := db.Prepare(insertTrigger)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec(userID, tt.String(), userID, tt.String(), userID, tt.String())
	return err
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
	// where the key string is a column name

	return nil
}
