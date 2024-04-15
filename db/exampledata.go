package db

import (
	"fmt"
	"foodbar/views/viewutils"

	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func readTab(tt viewutils.TabType, userID uuid.UUID) (viewutils.TabData, error) {
	db, err := sql.Open("sqlite3", "~/.local/share/FOOdBAR/db.db")
	if err != nil {
		return viewutils.TabData{}, err
	}
	defer db.Close()

	tableName := fmt.Sprintf("%s_%s", userID, tt)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		return viewutils.TabData{}, err
	}

	rows, err := db.Query("SELECT name FROM " + tableName)
	if err != nil {
		return viewutils.TabData{}, err
	}

	var items []viewutils.TabItem
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return viewutils.TabData{}, err
		}
		items = append(items, viewutils.TabItem{Ttype: tt})
	}

	return viewutils.TabData{Items: items, Ttype: tt}, nil
}

func MkRecipeItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Recipe,
		ItemID:   uuid.New(),
	}
}

func MkPantryItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Pantry,
		ItemID:   uuid.New(),
	}
}

func MkMenuItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Menu,
		ItemID:   uuid.New(),
	}
}

func MkShoppingItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Shopping,
		ItemID:   uuid.New(),
	}
}

func MkPreplistItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Preplist,
		ItemID:   uuid.New(),
	}
}

func MkEarningsItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Earnings,
		ItemID:   uuid.New(),
	}
}

func NewExampleRecipeTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			MkRecipeItem("Chicken"),
			MkRecipeItem("turd sandwich"),
			MkRecipeItem("chicken masala"),
			MkRecipeItem("tacos caliente"),
		},
		Ttype: viewutils.Recipe,
	}
}

func NewExamplePantryTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			MkPantryItem("Chicken"),
			MkPantryItem("turd sandwich"),
			MkPantryItem("chicken masala"),
			MkPantryItem("tacos caliente"),
		},
		Ttype: viewutils.Pantry,
	}
}

func NewExampleMenuTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			MkMenuItem("Chicken"),
			MkMenuItem("turd sandwich"),
			MkMenuItem("chicken masala"),
			MkMenuItem("tacos caliente"),
		},
		Ttype: viewutils.Menu,
	}
}

func NewExampleShoppingTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			MkShoppingItem("Chicken"),
			MkShoppingItem("turd sandwich"),
			MkShoppingItem("chicken masala"),
			MkShoppingItem("tacos caliente"),
		},
		Ttype: viewutils.Shopping,
	}
}

func NewExamplePreplistTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			MkPreplistItem("Chicken"),
			MkPreplistItem("turd sandwich"),
			MkPreplistItem("chicken masala"),
			MkPreplistItem("tacos caliente"),
		},
		Ttype: viewutils.Preplist,
	}
}

func NewExampleEarningsTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			MkEarningsItem("Chicken"),
			MkEarningsItem("turd sandwich"),
			MkEarningsItem("chicken masala"),
			MkEarningsItem("tacos caliente"),
		},
		Ttype: viewutils.Earnings,
	}
}
