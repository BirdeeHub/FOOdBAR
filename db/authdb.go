package db

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"time"

	foodlib "FOOdBAR/FOOlib"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func CleanSessionBlacklist() error {
	authDB := filepath.Join(dbpath, "FOOdBAR", "auth.db")
	authdbpath, err := foodlib.CreateEmptyFileIfNotExists(authDB)
	if err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", authdbpath)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS session_blacklist_table (
						sessionID TEXT PRIMARY KEY,
						expires_at DATETIME
					)`)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM session_blacklist_table WHERE expires_at < DATETIME('now')")
	return err
}

func AddToSessionBlacklist(sessionID uuid.UUID, expiration time.Time) error {
	authDB := filepath.Join(dbpath, "FOOdBAR", "auth.db")
	authdbpath, err := foodlib.CreateEmptyFileIfNotExists(authDB)
	if err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", authdbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS session_blacklist_table (
						sessionID TEXT PRIMARY KEY,
						expires_at DATETIME
					)`)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO session_blacklist_table (sessionID, expires_at) VALUES (?, ?)", sessionID.String(), expiration)
	if err != nil {
		// If the session is already in the blacklist, then I guess it worked...
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil
		}
	}
	return nil
}

func QuerySessionBlacklist(sessionID uuid.UUID) (bool, error) {
	authDB := filepath.Join(dbpath, "FOOdBAR", "auth.db")
	authdbpath, err := foodlib.CreateEmptyFileIfNotExists(authDB)
	if err != nil {
		return false, err
	}
	db, err := sql.Open("sqlite3", authdbpath)
	if err != nil {
		return false, err
	}
	defer db.Close()

	// prepare statement (for input sanitization)
	stmt, err := db.Prepare("SELECT expires_at FROM session_blacklist_table WHERE sessionID = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(sessionID.String())
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var expiration time.Time
	err = rows.Scan(&expiration)
	if err != nil {
		return false, nil
	}

	if expiration.Before(time.Now()) {
		_, err = db.Exec("DELETE FROM session_blacklist_table WHERE sessionID = ?", sessionID.String())
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil

}

func AuthUser(username string, password string) (uuid.UUID, error) {
	authDB := filepath.Join(dbpath, "FOOdBAR", "auth.db")
	authdbpath, err := foodlib.CreateEmptyFileIfNotExists(authDB)
	if err != nil {
		return uuid.Nil, err
	}
	db, err := sql.Open("sqlite3", authdbpath)
	defer db.Close()
	if err != nil {
		return uuid.Nil, err
	}

	// prepare statement (for input sanitization)
	stmt, err := db.Prepare("SELECT password, userID FROM user_auth_table WHERE username = ?")
	if err != nil {
		return uuid.Nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return uuid.Nil, err
	}
	defer rows.Close()
	var dbPassword []byte
	var dbpass [32]byte
	var userID uuid.UUID

	if rows.Next() {
		if err := rows.Scan(&dbPassword, &userID); err != nil {
			return uuid.Nil, err
		}
	} else {
		return uuid.Nil, errors.New("user not found")
	}
	if len(dbPassword) == 32 {
		copy(dbpass[:], dbPassword)
	} else {
		return uuid.Nil, errors.New("invalid hash length")
	}

	hashedPassword := sha256.Sum256([]byte(password))
	passwordMatch := true
	for i := 0; i < len(dbpass); i++ {
		if hashedPassword[i] != dbpass[i] {
			passwordMatch = false
		}
	}
	if !passwordMatch {
		return uuid.Nil, errors.New("invalid password")
	}

	return userID, nil
}

func CreateUser(username string, password string) (uuid.UUID, error) {
	authDB := filepath.Join(dbpath, "FOOdBAR", "auth.db")
	authdbpath, err := foodlib.CreateEmptyFileIfNotExists(authDB)
	if err != nil {
		return uuid.Nil, err
	}
	db, err := sql.Open("sqlite3", authdbpath)
	if err != nil {
		return uuid.Nil, err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_auth_table (
						username TEXT PRIMARY KEY,
						password BLOB,
						userID TEXT UNIQUE
					)`)
	if err != nil {
		return uuid.Nil, err
	}
	hashedPassword := sha256.Sum256([]byte(password))

	userID := uuid.New()

	_, err = db.Exec("INSERT INTO user_auth_table (username, password, userID) VALUES (?, ?, ?)", username, hashedPassword[:], userID.String())
	if err != nil {
		// make sure the uuid is unique (just in case)
		if strings.Contains(err.Error(), "userID") && strings.Contains(err.Error(), "UNIQUE constraint failed") {
			for err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") && strings.Contains(err.Error(), "userID") {
				userID = uuid.New()
				_, err = db.Exec("INSERT INTO user_auth_table (username, password, userID) VALUES (?, ?, ?)", username, hashedPassword[:], userID.String())
			}
			if err != nil {
				return uuid.Nil, err
			} else {
				return userID, nil
			}
		} else {
			return uuid.Nil, err
		}
	}

	return userID, nil
}
