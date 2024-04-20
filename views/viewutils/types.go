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
	Invalid  TabType = ""
	Recipe           = "Recipe"
	Pantry           = "Pantry"
	Menu             = "Menu"
	Shopping         = "Shopping"
	Preplist         = "Preplist"
	Earnings         = "Earnings"
)

type ColorScheme string

const (
	None  ColorScheme = ""
	Dark              = "dark"
	Light             = "light"
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
	Palette    ColorScheme
	LastActive time.Time
}

type SortMethod string

const (
	Inactive   SortMethod = ""
	Descending            = "DESC;"
	Ascending             = "ASC;"
	Random                = "RANDOM();"
	Custom                = "END;"
	// Others can be made with CASE WHEN condition THEN value ELSE value END
	// when using this syntax, put the CASE WHEN... etc... into the OrderBy key
	// and then put Custom as the SortMethod
)

type TabData struct {
	Parent  *PageData
	Active  bool
	Ttype   TabType
	Items   map[uuid.UUID]*TabItem
	OrderBy map[string]SortMethod
}

type TabItem struct {
	Parent *TabData
	ItemID uuid.UUID
	Ttype  TabType

	Expanded bool
	ItemType string
}

func (tbd *TabData) AddTabItem(ti *TabItem) {
	tbd.Parent.LastActive = time.Now()
	ti.Parent = tbd
	ti.Ttype = tbd.Ttype

	if ti.ItemID == uuid.Nil {
		ti.ItemID = uuid.New()
	}
	tbd.Items[ti.ItemID] = ti
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

// TODO: Make it so that it clears old ones every so often
// TODO: Implement client side caching of pageData in case it times out a still-valid login?
var pageDatas = make(map[uuid.UUID]*PageData)

func GetPageData(userID uuid.UUID) *PageData {
	if pageDatas[userID] == nil {
		pageDatas[userID] = InitPageData(userID)
	}
	pageDatas[userID].LastActive = time.Now()
	return pageDatas[userID]
}
