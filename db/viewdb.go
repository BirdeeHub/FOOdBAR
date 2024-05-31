package db

import (
	foodlib "FOOdBAR/FOOlib"
	"database/sql"
	"encoding/json"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// NOTE: This should probably honestly be something like redis,
// but the database will do fine for this app most likely.

func GetPageDataDB() (*sql.DB, error) {
	viewDB := filepath.Join(dbpath, "FOOdBAR", "views.db")
	viewdbpath, err := foodlib.CreateEmptyFileIfNotExists(viewDB)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", viewdbpath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS viewstates (
						tab_id TEXT,
						user_id TEXT,
						session_id TEXT,
						page_data BLOB,
						last_modified DATETIME
					)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TRIGGER IF NOT EXISTS last_insert
			AFTER INSERT ON viewstates
			FOR EACH ROW  
			BEGIN
				UPDATE viewstates
				SET last_modified = CURRENT_TIMESTAMP
				WHERE tab_id = NEW.tab_id;
			END;`)
	if err != nil {
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TRIGGER IF NOT EXISTS last_update
			AFTER UPDATE ON viewstates
			FOR EACH ROW  
			BEGIN
				UPDATE viewstates
				SET last_modified = CURRENT_TIMESTAMP
				WHERE tab_id = NEW.tab_id;
			END;`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func CleanPageDataDB() error {
	db, err := GetPageDataDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// TODO: make timeout value configurable when running from command line
	_, err = db.Exec("DELETE FROM viewstates WHERE last_modified < DATETIME('now', '-1 hours')")
	return err
}

func GetPageData(c echo.Context) (*foodlib.PageData, error) {
	userID, err := foodlib.GetUserFromClaims(foodlib.GetClaimsFromContext(c))
	if err != nil {
		return nil, err
	}
	SID, err := foodlib.GetSessionIDFromClaims(foodlib.GetClaimsFromContext(c))
	if err != nil {
		return nil, err
	}
	tabID := c.Request().Header.Get("tab_id")
	c.Logger().Printf("tabID: %s", tabID)
	if tabID == "" {
		tabID = uuid.New().String()
		c.Logger().Printf("tabID: %s", tabID)
	}
	db, err := GetPageDataDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var pageDataBlob []byte
	err = db.QueryRow("SELECT page_data FROM viewstates WHERE tab_id = ?", tabID).Scan(&pageDataBlob)
	if err != nil {
		// If tabID is not in db, db.QueryRow for SessionID sorted by last_modified
		err = db.QueryRow("SELECT page_data FROM viewstates WHERE session_id = ? ORDER BY last_modified DESC LIMIT 1", SID).Scan(&pageDataBlob)
		if err != nil {
			// If still not found, create a new one and add it to the db
			pd := foodlib.InitPageData(userID, SID, tabID)
			pageDataBlob, err = json.Marshal(pd)
			if err != nil {
				return nil, err
			}
			_, err = db.Exec("INSERT INTO viewstates (tab_id, session_id, user_id, page_data) VALUES (?, ?, ?, ?)", tabID, SID, userID, pageDataBlob)
			if err != nil {
				return nil, err
			}
			return pd, nil
		}
	}

	pd := &foodlib.PageData{}
	err = json.Unmarshal(pageDataBlob, pd)
	if err != nil {
		return nil, err
	}
	if pd.UserID != userID || pd.SessionID != SID || pd.TabID != tabID {
		pd := foodlib.InitPageData(userID, SID, tabID)
		return pd, nil
	}
	return pd, nil
}

func SavePageData(c echo.Context, pd *foodlib.PageData) error {
	db, err := GetPageDataDB()
	if err != nil {
		return err
	}
	defer db.Close()

	pdmarshalled, err := json.Marshal(pd)
	if err != nil {
		return err
	}

	// Check if the row exists, if so, update it; otherwise, insert a new row
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM viewstates WHERE tab_id = ?)", pd.TabID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = db.Exec("UPDATE viewstates SET page_data = ? WHERE tab_id = ?", pdmarshalled, pd.TabID)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("INSERT INTO viewstates (tab_id, user_id, session_id, page_data) VALUES (?, ?, ?, ?)", pd.TabID, pd.UserID, pd.SessionID, pdmarshalled)
		if err != nil {
			return err
		}
	}

	c.Response().Header().Add("tab_id", pd.TabID)
	return nil
}
