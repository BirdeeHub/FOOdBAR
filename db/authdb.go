package db

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func AuthUser(username string, password string) (uuid.UUID, error) {
	db, err := sql.Open("sqlite3", "/home/birdee/.local/share/FOOdBAR/auth.db")
	if err != nil {
		return uuid.Nil, err
	}
	defer db.Close()

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
	db, err := sql.Open("sqlite3", "/home/birdee/.local/share/FOOdBAR/auth.db")
	if err != nil {
		return uuid.Nil, err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_auth_table (
						username TEXT PRIMARY KEY,
						password BLOB,
						userID TEXT
					)`)
	if err != nil {
		return uuid.Nil, err
	}
	hashedPassword := sha256.Sum256([]byte(password))

	userID := uuid.New()

	_, err = db.Exec("INSERT INTO user_auth_table (username, password, userID) VALUES (?, ?, ?)", username, hashedPassword[:], userID.String())
	if err != nil {
		// Check if the error is due to a constraint violation (duplicate username)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return uuid.Nil, errors.New("username already exists")
		}
		return uuid.Nil, err
	}

	return userID, nil
}
