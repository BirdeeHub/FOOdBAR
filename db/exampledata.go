package db

import (
	"fmt"
	"FOOdBAR/views/viewutils"

	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func ReadTabData(tt viewutils.TabType, userID uuid.UUID) (*viewutils.TabData, error) {
	db, err := sql.Open("sqlite3", "~/.local/share/FOOdBAR/db.db")
	if err != nil {
		return &viewutils.TabData{}, err
	}
	defer db.Close()

	tableName := fmt.Sprintf("%s_%s", userID, tt.String())

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		return &viewutils.TabData{}, err
	}

	rows, err := db.Query("SELECT name FROM " + tableName)
	if err != nil {
		return &viewutils.TabData{}, err
	}

	var items []*viewutils.TabItem
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return &viewutils.TabData{}, err
		}
		items = append(items, &viewutils.TabItem{Ttype: tt})
	}

	return &viewutils.TabData{Items: items, Ttype: tt}, nil
}
