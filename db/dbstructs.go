package db

import (
	// foodlib "FOOdBAR/FOOlib"
	// "errors"
	"time"

	"github.com/google/uuid"
	// "database/sql"
	// _ "github.com/mattn/go-sqlite3"
)

type RecipeResult struct {
	ID           uuid.UUID   `json:"id"`
	UserID       uuid.UUID   `json:"userid"`
	CreatedAt    time.Time   `json:"createdat"`
	LastModified time.Time   `json:"lastmodified"`
	LastAuthor   string      `json:"lastauthor"`
	Name         string      `json:"name"`
	Category     []string    `json:"category"`
	Dietary      []string    `json:"dietary"`
	Ingredients  []uuid.UUID `json:"ingredients"`
	Instructions string      `json:"instructions"`
}

type MenuResult struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID   `json:"userid"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	MenuID       uuid.UUID `json:"menuid"`
	Name         string    `json:"name"`
	Number       int       `json:"number"`
}

type PantryResult struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID   `json:"userid"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	Name         string    `json:"name"`
	Dietary      []string  `json:"dietary"`
	Amount       float32   `json:"amount"`
	Units        string    `json:"units"`
}

type CustomerResult struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID   `json:"userid"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Name         string    `json:"name"`
	Dietary      []string  `json:"dietary"`
}

type EventResult struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID   `json:"userid"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	Name         string    `json:"name"`
	MenuID       uuid.UUID `json:"menuid"`
	Date         string    `json:"date"`
	Location     string    `json:"location"`
	Customer     string    `json:"customer"`
}

type PreplistResult struct {
	ID           uuid.UUID             `json:"id"`
	UserID       uuid.UUID   `json:"userid"`
	CreatedAt    time.Time             `json:"createdat"`
	LastModified time.Time             `json:"lastmodified"`
	LastAuthor   string                `json:"lastauthor"`
	EventID      uuid.UUID             `json:"eventid"`
	MenuID       uuid.UUID             `json:"menuid"`
	Ingredients  map[uuid.UUID]float32 `json:"ingredients"`
	Recipes      []uuid.UUID           `json:"recipes"`
}

type ShoppingResult struct {
	ID           uuid.UUID   `json:"id"`
	UserID       uuid.UUID   `json:"userid"`
	CreatedAt    time.Time   `json:"createdat"`
	LastModified time.Time   `json:"lastmodified"`
	LastAuthor   string      `json:"lastauthor"`
	EventID      uuid.UUID   `json:"eventid"`
	MenuID       uuid.UUID   `json:"menuid"`
	Amount       float32     `json:"amount"`
	Units        string      `json:"units"`
	Ingredients  []uuid.UUID `json:"ingredients"`
}

type EarningsResult struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID   `json:"userid"`
	CreatedAt    time.Time `json:"createdat"`
	LastModified time.Time `json:"lastmodified"`
	LastAuthor   string    `json:"lastauthor"`
	EventID      uuid.UUID `json:"eventid"`
	MenuID       uuid.UUID `json:"menuid"`
}
