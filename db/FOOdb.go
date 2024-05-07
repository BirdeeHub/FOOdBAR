package db

import (
	foodlib "FOOdBAR/FOOlib"
	"database/sql"
	"errors"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func SubmitPantryItem(c echo.Context, pd *foodlib.PageData, td *foodlib.TabData, item *foodlib.TabItem) error {
	db, err := CreateTabTableIfNotExists(pd.UserID, td.Ttype)
	defer db.Close()
	if err != nil {
		return err
	}
	name := c.FormValue("itemName")
	dietary := c.FormValue("itemDietary")
	amount := c.FormValue("itemAmount")
	units := c.FormValue("itemUnits")
	insertStmt, err := db.Prepare(`INSERT INTO ?_? (id, last_author, name, dietary, amount, units) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec(pd.UserID, td.Ttype.String(), item.ItemID.String(), pd.UserID, name, dietary, amount, units)
	return err 
}

func CreateTabTableIfNotExists(userID uuid.UUID, tt foodlib.TabType) (*sql.DB, error) {
	var err error
	fooDB := filepath.Join(dbpath, "FOOdBAR", "FOOdb.db")
	fooDB, err = foodlib.CreateEmptyFileIfNotExists(fooDB)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbpath)

	var createTable string
	switch tt {
	case foodlib.Recipe:
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
	case foodlib.Menu:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				menu_id TEXT,
				name TEXT,
				number INTEGER,
				)`
	case foodlib.Pantry:
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
	case foodlib.Customer:
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
	case foodlib.Events:
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
	case foodlib.Preplist:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT,
				ingredients TEXT,
				recipes TEXT,
				)`
	case foodlib.Shopping:
		createTable = `CREATE TABLE IF NOT EXISTS ?_? (
				id TEXT PRIMARY KEY,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT,
				ingredients TEXT,
				amount TEXT,
				units TEXT,
				)`
	case foodlib.Earnings:
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

func makeAuditTriggers(db *sql.DB, userID uuid.UUID, tt foodlib.TabType) error {

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
*/
func FillXTabItems(userID uuid.UUID, tbd *foodlib.TabData, number int) error {
	db, err := CreateTabTableIfNotExists(userID, tbd.Ttype)
	defer db.Close()
	if err != nil {
		return err
	}
	// TODO: fill tbd.Items with X number of items based on tbd.OrderBy: SortMethod
	// where the key string is a column name

	return nil
}
