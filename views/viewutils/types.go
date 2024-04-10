package viewutils

import (
	"errors"

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

func GetTabTypes() [6]TabType {
	return [...]TabType{Recipe, Pantry, Menu, Shopping, Preplist, Earnings}
}

func (t *TabType) String() string {
	return string(*t)
}

func String2TabType(str string) (*TabType, error) {
	for _, t := range GetTabTypes() {
		if t.String() == str {
			return &t, nil
		}
	}
	return nil, errors.New("invalid TabType")
}

var tabToggleButtons = map[TabType]bool{
	Recipe: false,
	Pantry: false,
	Menu: false,
	Shopping: false,
	Preplist: false,
	Earnings: false,
}

func (t *TabType) IsActive() bool {
	return tabToggleButtons[*t]
}

func (t *TabType) ToggleActive() {
	tabToggleButtons[*t] = !tabToggleButtons[*t]
}

type TabItem struct {
	ItemID   uuid.UUID
	ItemName string
	Ttype    TabType
}

type TabData struct {
	Ttype TabType
	Items []TabItem
}
