package db

import (
	foodlib "FOOdBAR/FOOlib"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

//TODO: This function (should be applicable for both submit AND update, retaining old values for empty fields)
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
	insertStmt, err := db.Prepare(fmt.Sprintf(`INSERT INTO %s (id, userid, last_author, name, dietary, amount, units) VALUES (?, ?, ?, ?, ?, ?, ?)`, item.Ttype.String()))
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec(item.ItemID.String(), pd.UserID, pd.UserID, name, dietary, amount, units)
	return err 
}

func CreateTabTableIfNotExists(userID uuid.UUID, tt foodlib.TabType) (*sql.DB, error) {
	var err error
	fooDB := filepath.Join(dbpath, "FOOdBAR", "FOOdb.db")
	fooDB, err = foodlib.CreateEmptyFileIfNotExists(fooDB)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", fooDB)
	if err != nil {
		return nil, err
	}

	var createTable string
	switch tt {
	case foodlib.Recipe:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				name TEXT UNIQUE,
				category TEXT,
				dietary TEXT,
				ingredients TEXT,
				instructions TEXT
				)`
	case foodlib.Menu:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				menu_id TEXT,
				name TEXT,
				number INTEGER
				)`
	case foodlib.Pantry:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				name TEXT UNIQUE,
				dietary TEXT,
				amount TEXT,
				units TEXT
				)`
	case foodlib.Customer:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				email TEXT UNIQUE,
				phone TEXT UNIQUE,
				name TEXT,
				dietary TEXT
				)`
	case foodlib.Events:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				name TEXT,
				menu_id TEXT,
				date TEXT,
				location TEXT,
				customer TEXT
				)`
	case foodlib.Preplist:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT,
				ingredients TEXT,
				recipes TEXT
				)`
	case foodlib.Shopping:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT,
				ingredients TEXT,
				amount TEXT,
				units TEXT
				)`
	case foodlib.Earnings:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT
				)`
	default:
		return nil, errors.New("invalid tab type")
	}
	_, err = db.Exec(fmt.Sprintf(createTable, tt.String()))
	if err != nil {
		return nil, err
	}
	err = makeAuditTriggers(db, tt)
	if err != nil {
		fmt.Printf("error: %s path: %s command: %s\n", err.Error(), fooDB, fmt.Sprintf(createTable, tt.String()))
		return nil, err
	}
	return db, nil
}

func makeAuditTriggers(db *sql.DB, tt foodlib.TabType) error {

	updateTrigger := `CREATE TRIGGER IF NOT EXISTS update_%s_audit 
		DELETE ON %s
		FOR EACH ROW
		BEGIN
			UPDATE %s
			SET last_modified = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
		CREATE TRIGGER IF NOT EXISTS update_%s_audit 
		UPDATE ON %s
		FOR EACH ROW
		BEGIN
			UPDATE %s
			SET last_modified = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
		CREATE TRIGGER IF NOT EXISTS update_%s_audit 
		INSERT ON %s
		FOR EACH ROW
		BEGIN
			UPDATE %s
			SET last_modified = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;`

	insertTrigger := `CREATE TRIGGER IF NOT EXISTS add_created_%s_audit
		INSERT ON %s
		FOR EACH ROW  
		BEGIN
			UPDATE %s
			SET created_at = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;`

	var err error = nil
	_, err = db.Exec(fmt.Sprintf(
		updateTrigger,
		tt.String(),
		tt.String(),
		tt.String(),
		tt.String(),
		tt.String(),
		tt.String(),
		tt.String(),
		tt.String(),
		tt.String(),
	))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf(
		insertTrigger,
		tt.String(),
		tt.String(),
		tt.String(),
	))
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
