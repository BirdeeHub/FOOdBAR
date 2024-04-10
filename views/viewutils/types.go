package viewutils

import (
	"github.com/google/uuid"
)

type PageData struct {
	TabDatas []TabData
}

type TabType string

const (
	Recipe   TabType = "Recipe"
	Pantry           = "Pantry"
	Menu             = "Menu"
	Shopping         = "Shopping"
	Preplist         = "Preplist"
	Earnings         = "Earnings"
)

type TabItem struct {
	ItemID   uuid.UUID
	ItemName string
	Ttype    TabType
}

type TabData struct {
	Ttype TabType
	Items []TabItem
}
