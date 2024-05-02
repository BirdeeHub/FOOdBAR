package db

import (
	// foodlib "FOOdBAR/FOOlib"
	// "errors"
	"time"

	"github.com/google/uuid"
	// "database/sql"
	// _ "github.com/mattn/go-sqlite3"
)

type Recipe struct {
	ID           uuid.UUID   `json:"id"`
	CreatedAt    time.Time   `json:"createdat"`
	LastModified time.Time   `json:"lastmodified"`
	LastAuthor   string      `json:"lastauthor"`
	Name         string      `json:"name"`
	Category     []string    `json:"category"`
	Dietary      []string    `json:"dietary"`
	Ingredients  []uuid.UUID `json:"ingredients"`
	Instructions string      `json:"instructions"`
}

type Menu struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	MenuID       uuid.UUID `json:"menuid"`
	Name         string    `json:"name"`
	Number       int       `json:"number"`
}

type Pantry struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	Name         string    `json:"name"`
	Dietary      []string  `json:"dietary"`
	Amount       float32   `json:"amount"`
	Units        string    `json:"units"`
}

type Customer struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Name         string    `json:"name"`
	Dietary      []string  `json:"dietary"`
}

type Event struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	Name         string    `json:"name"`
	MenuID       uuid.UUID `json:"menuid"`
	Date         string    `json:"date"`
	Location     string    `json:"location"`
	Customer     string    `json:"customer"`
}

type Preplist struct {
	ID           uuid.UUID             `json:"id"`
	CreatedAt    time.Time             `json:"createdat"`
	LastModified time.Time             `json:"lastmodified"`
	LastAuthor   string                `json:"lastauthor"`
	EventID      uuid.UUID             `json:"eventid"`
	MenuID       uuid.UUID             `json:"menuid"`
	Ingredients  map[uuid.UUID]float32 `json:"ingredients"`
	Recipes      []uuid.UUID           `json:"recipes"`
}

type Shopping struct {
	ID           uuid.UUID   `json:"id"`
	CreatedAt    time.Time   `json:"createdat"`
	LastModified time.Time   `json:"lastmodified"`
	LastAuthor   string      `json:"lastauthor"`
	EventID      uuid.UUID   `json:"eventid"`
	MenuID       uuid.UUID   `json:"menuid"`
	Amount       float32     `json:"amount"`
	Units        string      `json:"units"`
	Ingredients  []uuid.UUID `json:"ingredients"`
}

type Earnings struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	EventID      uuid.UUID `json:"eventid"`
	MenuID       uuid.UUID `json:"menuid"`
}
