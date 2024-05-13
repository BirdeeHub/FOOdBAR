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
	if item.Ttype == foodlib.Invalid {
		return errors.New("Invalid Tab Type")
	}
	userID, err := foodlib.GetUserFromClaims(foodlib.GetClaimsFromContext(c))
	if err != nil {
		return err
	}
	db, err := CreateTabTableIfNotExists(userID, td.Ttype)
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

//TODO: This function
func GetTabItemDataByUUID(c echo.Context, item foodlib.TabItem) error {
	if item.Ttype == foodlib.Invalid {
		return errors.New("Invalid Tab Type")
	}
	userID, err := foodlib.GetUserFromClaims(foodlib.GetClaimsFromContext(c))
	if err != nil {
		return err
	}
	db, err := CreateTabTableIfNotExists(userID, item.Ttype)
	defer db.Close()
	if err != nil {
		return err
	}
	return errors.New("not yet implemented")
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

	// createUserTable := `CREATE TABLE IF NOT EXISTS user_meta_table (
	// 	id TEXT PRIMARY KEY,
	// 	recipe_table TEXT,
	// 	menu_table TEXT,
	// 	pantry_table TEXT,
	// 	customer_table TEXT,
	// 	events_table TEXT,
	// 	preplist_table TEXT,
	// 	shopping_table TEXT,
	// 	earnings_table TEXT,
	// )`
	//TODO: check if table field is empty, if so, create new table with unique name and populate field.

	//TODO: completely redesign tables
	//There should be a table of user tables, which stores the other table names and the pageData struct
	var createTable string
	switch tt {
	case foodlib.Recipe:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
				userid TEXT,
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				name TEXT,
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
				name TEXT,
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
				email TEXT,
				phone TEXT,
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
		return nil, err
	}
	return db, nil
}

func makeAuditTriggers(db *sql.DB, tt foodlib.TabType) error {

	createTrigger := func(name string, method string, field string, old string, tt foodlib.TabType) string {
		trigger := `CREATE TRIGGER IF NOT EXISTS %s_%s_audit
			%s ON %s
			FOR EACH ROW  
			BEGIN
				UPDATE %s
				SET %s = CURRENT_TIMESTAMP
				WHERE id = %s.id;
			END;`
		return fmt.Sprintf(trigger, name, tt.String(), method, tt.String(), tt.String(), field, old)
	}

	var err error = nil
	_, err = db.Exec(createTrigger("add_created", "AFTER INSERT", "created_at", "NEW", tt))
	if err != nil {
		return err
	}
	_, err = db.Exec(createTrigger("update_modified_insert", "AFTER INSERT", "last_modified", "NEW", tt))
	if err != nil {
		return err
	}
	_, err = db.Exec(createTrigger("update_modified_update", "AFTER UPDATE", "last_modified", "OLD", tt))
	if err != nil {
		return err
	}
	_, err = db.Exec(createTrigger("update_modified_delete", "AFTER DELETE", "last_modified", "OLD", tt))
	return err
}

//TODO: should not fetch data, but instead, which tabItems to fetch data from
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
