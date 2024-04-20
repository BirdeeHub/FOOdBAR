package viewutils

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const PagePrefix = "/FOOdBAR"

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

type PageData struct {
	UserID     uuid.UUID
	TabDatas   []*TabData
	LastActive time.Time
}

type TabData struct {
	Parent   *PageData
	Active bool
	Ttype  TabType
	Items  map[uuid.UUID]*TabItem
}

type TabItem struct {
	Parent   *TabData
	ItemID   uuid.UUID
	Ttype    TabType

	Expanded bool
}

func (tbd *TabData) AddTabItem(ti *TabItem) {
	tbd.Parent.LastActive = time.Now()
	ti.Parent = tbd
	itemID := uuid.New()
	ti.ItemID = itemID
	ti.Ttype = tbd.Ttype
	tbd.Items[itemID] = ti
}

func (tbd *TabData) GetTabItems() []*TabItem {
	tbd.Parent.LastActive = time.Now()
	var tis []*TabItem
	for _, ti := range tbd.Items {
		tis = append(tis, ti)
	}
	return tis
}

func (tbd *TabData) String() string {
	tbd.Parent.LastActive = time.Now()
	return tbd.Ttype.String()
}

func (tbd *TabData) IsActive() bool {
	tbd.Parent.LastActive = time.Now()
	return (*tbd).Active
}

func (tbd *TabData) ToggleActive() {
	tbd.Parent.LastActive = time.Now()
	(*tbd).Active = !(*tbd).Active
}

func (tbd *TabData) SetActive(v bool) {
	tbd.Parent.LastActive = time.Now()
	(*tbd).Active = v
}

func (pgd *PageData) GetTabDataByType(tt TabType) (*TabData, error) {
	pgd.LastActive = time.Now()
	for _, t := range pgd.TabDatas {
		if t.Ttype == tt {
			return t, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no %s tab", tt))
}

func InitPageData(userID uuid.UUID) *PageData {
	pd := &PageData{
		UserID:     userID,
		LastActive: time.Now(),
		TabDatas: []*TabData{
			{
				Active: false,
				Ttype:  Recipe,
				Items:  make(map[uuid.UUID]*TabItem),
			},
			{
				Active: false,
				Ttype:  Pantry,
				Items:  make(map[uuid.UUID]*TabItem),
			},
			{
				Active: false,
				Ttype:  Menu,
				Items:  make(map[uuid.UUID]*TabItem),
			},
			{
				Active: false,
				Ttype:  Preplist,
				Items:  make(map[uuid.UUID]*TabItem),
			},
			{
				Active: false,
				Ttype:  Shopping,
				Items:  make(map[uuid.UUID]*TabItem),
			},
			{
				Active: false,
				Ttype:  Earnings,
				Items:  make(map[uuid.UUID]*TabItem),
			},
		},
	}
	for _, td := range pd.TabDatas {
		td.Parent = pd
	}
	return pd
}
