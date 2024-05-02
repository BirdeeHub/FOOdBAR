package db

import (
	// foodlib "FOOdBAR/FOOlib"
	// "database/sql"
	// "errors"
	"time"

	"github.com/google/uuid"
	// _ "github.com/mattn/go-sqlite3"
)

type Recipe struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastModified time.Time
	LastAuthor   string
	Name         string
	Category     []string
	Dietary      []string
	Ingredients  []uuid.UUID
	Instructions string
}

type Menu struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastModified time.Time
	LastAuthor   string
	MenuID       uuid.UUID
	Name         string
	Number       int
}

type Pantry struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastModified time.Time
	LastAuthor   string
	Name         string
	Dietary      []string
	Amount       float32
	Units        string
}

type Customer struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastModified time.Time
	LastAuthor   string
	Email        string
	Phone        string
	Name         string
	Dietary      []string
}

type Event struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastModified time.Time
	LastAuthor   string
	Name         string
	MenuID       uuid.UUID
	Date         string
	Location     string
	Customer     string
}

type Preplist struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastModified time.Time
	LastAuthor   string
	EventID      uuid.UUID
	MenuID       uuid.UUID
	ingredients  map[uuid.UUID]float32
	recipes      []uuid.UUID
}

type Shopping struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastModified time.Time
	LastAuthor   string
	EventID      uuid.UUID
	MenuID       uuid.UUID
	Amount       float32
	Units        string
	ingredients  []uuid.UUID
}

type Earnings struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastModified time.Time
	LastAuthor   string
	EventID      uuid.UUID
	MenuID       uuid.UUID
}
