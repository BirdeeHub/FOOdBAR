package viewutils

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const PagePrefix = "/FOOdBAR"

type PageData struct {
	UserID   uuid.UUID
	TabDatas []*TabData
}

type TabItem struct {
	ItemID   uuid.UUID
	ItemName string
	Ttype    TabType
}

type TabData struct {
	UserID uuid.UUID
	Active bool
	Ttype TabType
	Items []*TabItem
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

func (tbd *TabData) String() string {
	return tbd.Ttype.String()
}

func (tbd *TabData) IsActive() bool {
	return (*tbd).Active
}

func (tbd *TabData) ToggleActive() {
	(*tbd).Active = !(*tbd).Active
}

func (tbd *TabData) SetActive(v bool) {
	(*tbd).Active = v
}

func (pgd *PageData) GetTabDataByType(tt TabType) (*TabData, error) {
	for _, t := range pgd.TabDatas {
		if t.Ttype == tt {
			return t, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no %s tab", tt))
}

func InitPageData(userID uuid.UUID) *PageData {
	return &PageData{
		UserID: userID,
		TabDatas: []*TabData{
			{
				UserID: userID,
				Active: false,
				Ttype: Recipe,
				Items: nil,
			},
			{
				UserID: userID,
				Active: false,
				Ttype: Pantry,
				Items: nil,
			},
			{
				UserID: userID,
				Active: false,
				Ttype: Menu,
				Items: nil,
			},
			{
				UserID: userID,
				Active: false,
				Ttype: Preplist,
				Items: nil,
			},
			{
				UserID: userID,
				Active: false,
				Ttype: Shopping,
				Items: nil,
			},
			{
				UserID: userID,
				Active: false,
				Ttype: Earnings,
				Items: nil,
			},
		},
	}
}

