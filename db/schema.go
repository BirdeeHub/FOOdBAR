package db

import (
	foodlib "FOOdBAR/FOOlib"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// returns open db handle, tablename for that user and tab type, and an error
func CreateTabTableIfNotExists(userID uuid.UUID, tt foodlib.TabType) (*sql.DB, string, error) {
	var err error
	fooDB := filepath.Join(dbpath, "FOOdBAR", "FOOdb.db")
	fooDB, err = foodlib.CreateEmptyFileIfNotExists(fooDB)
	if err != nil {
		return nil, "", err
	}
	db, err := sql.Open("sqlite3", fooDB)
	if err != nil {
		return nil, "", err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_meta_table (
		id TEXT PRIMARY KEY,
		postfix TEXT UNIQUE,
		recipe_table TEXT DEFAULT "",
		menu_table TEXT DEFAULT "",
		pantry_table TEXT DEFAULT "",
		customer_table TEXT DEFAULT "",
		events_table TEXT DEFAULT "",
		preplist_table TEXT DEFAULT "",
		shopping_table TEXT DEFAULT "",
		earnings_table TEXT DEFAULT ""
	)`)
	if err != nil {
		return nil, "", err
	}

	var fieldname string
	switch tt {
	case foodlib.Recipe:
		fieldname = "recipe_table"
	case foodlib.Menu:
		fieldname = "menu_table"
	case foodlib.Pantry:
		fieldname = "pantry_table"
	case foodlib.Customer:
		fieldname = "customer_table"
	case foodlib.Events:
		fieldname = "events_table"
	case foodlib.Preplist:
		fieldname = "preplist_table"
	case foodlib.Shopping:
		fieldname = "shopping_table"
	case foodlib.Earnings:
		fieldname = "earnings_table"
	}
	var tableName string
	var postfix string
	err = db.QueryRow(fmt.Sprintf("SELECT postfix, %s FROM user_meta_table WHERE id = ?", fieldname), userID).Scan(&postfix, &tableName)
	if err != nil && err == sql.ErrNoRows {
		for err != nil {
			postfix = strings.ReplaceAll(uuid.New().String(), "-", "_")
			tableName = fmt.Sprintf("%s_%s", tt.String(), postfix)
			_, err = db.Exec(fmt.Sprintf(`INSERT INTO user_meta_table (id, postfix, %s) VALUES (?, ?, ?)`, fieldname), userID, postfix, tableName)
			if err != nil {
				// if not unique we make a new postfix and also tableName and try again
				// otherwise we just return an error
				if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
					return nil, "", err
				}
			}
		}
	} else if err != nil {
		return nil, "", err
	}

	var createTable string
	switch tt {
	case foodlib.Recipe:
		createTable = `CREATE TABLE IF NOT EXISTS %s (
				id TEXT PRIMARY KEY,
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
				created_at TEXT,
				last_modified TEXT,
				last_author TEXT,
				event_id TEXT,
				menu_id TEXT
				)`
	default:
		return nil, "", errors.New("invalid tab type")
	}
	cmd := fmt.Sprintf(createTable, tableName)
	_, err = db.Exec(cmd)
	if err != nil {
		return nil, "", err
	}
	err = makeAuditTriggers(db, tableName)
	if err != nil {
		return nil, "", err
	}
	return db, tableName, nil
}

func makeAuditTriggers(db *sql.DB, tableName string) error {

	createTrigger := func(name string, method string, field string, old string) string {
		trigger := `CREATE TRIGGER IF NOT EXISTS %s_%s_audit
			%s ON %s
			FOR EACH ROW  
			BEGIN
				UPDATE %s
				SET %s = CURRENT_TIMESTAMP
				WHERE id = %s.id;
			END;`
		return fmt.Sprintf(trigger, name, tableName, method, tableName, tableName, field, old)
	}

	var err error = nil
	_, err = db.Exec(createTrigger("add_created", "AFTER INSERT", "created_at", "NEW"))
	if err != nil {
		return err
	}
	_, err = db.Exec(createTrigger("update_modified_insert", "AFTER INSERT", "last_modified", "NEW"))
	if err != nil {
		return err
	}
	_, err = db.Exec(createTrigger("update_modified_update", "AFTER UPDATE", "last_modified", "OLD"))
	if err != nil {
		return err
	}
	_, err = db.Exec(createTrigger("update_modified_delete", "AFTER DELETE", "last_modified", "OLD"))
	return err
}
