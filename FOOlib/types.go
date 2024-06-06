package foodlib

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

const PagePrefix = "/FOOdBAR"
const JWTcookiename = "foodie"

var HtmxPath = "https://unpkg.com/htmx.org/dist/htmx.min.js"
var HyperscriptPath = "https://unpkg.com/hyperscript.org@0.9.12"

// TODO: make sorting rules
type SortMethod string

const (
	Inactive          SortMethod = ""
	CreatedDescending            = "created_at DESC;"
	CreatedAscending             = "created_at ASC;"
	RecencyDescending            = "last_modified DESC;"
	RecencyAscending             = "last_modified ASC;"
	NameDescending               = "name DESC;"
	NameAscending                = "name ASC;"
	NameRandom                   = "name RANDOM();"
	// These are incorrect, the values are stored as json lists, not just single value strings
	// DietaryDescending                = "dietary DESC;"
	// DietaryAscending                 = "dietary ASC;"
	// DietaryRandom                    = "dietary RANDOM();"
	// CategoryDescending               = "category DESC;"
	// CategoryAscending                = "category ASC;"
	// CategoryRandom                   = "category RANDOM();"
	// IngredientsDescending            = "ingredients DESC;"
	// IngredientsAscending             = "ingredients ASC;"
	// IngredientsRandom                = "ingredients RANDOM();"

	// NameCustom                = "CASE WHEN condition THEN value ELSE value END;"
)

func GetSortMethods() [8]SortMethod {
	return [...]SortMethod{
		Inactive,
		CreatedDescending,
		CreatedAscending,
		RecencyDescending,
		RecencyAscending,
		NameDescending,
		NameAscending,
		NameRandom,
		// DietaryDescending,
		// DietaryAscending,
		// DietaryRandom,
		// CategoryDescending,
		// CategoryAscending,
		// CategoryRandom,
		// IngredientsDescending,
		// IngredientsAscending,
		// IngredientsRandom,
	}
}

func GetSortMethodByNumber(k int) (SortMethod, error) {
	list := GetSortMethods()
	if k >= len(list) || k < 0 {
		return Inactive, errors.New("sort method does not exist")
	}
	return list[k], nil
}

type TabType string

const (
	Invalid  TabType = ""
	Recipe           = "Recipe"
	Pantry           = "Pantry"
	Menu             = "Menu"
	Customer         = "Customer"
	Events           = "Events"
	Shopping         = "Shopping"
	Preplist         = "Preplist"
	Earnings         = "Earnings"
)

func GetTabTypes() [8]TabType {
	return [...]TabType{Recipe, Pantry, Menu, Shopping, Preplist, Earnings, Customer, Events}
}

func (t *TabType) String() string {
	return string(*t)
}

// will return viewutils.Invalid if no match
func String2TabType(str string) TabType {
	for _, t := range GetTabTypes() {
		if t.String() == str {
			return t
		}
	}
	return Invalid
}

// PageData and its children

type TabButtonData struct {
	Ttype  TabType `json:"tab_type"`
	Active bool    `json:"active"`
}

type PageData struct {
	UserID    uuid.UUID  `json:"user_id"`
	SessionID uuid.UUID  `json:"session_id"`
	TabID     string     `json:"tab_id"`
	TabDatas  []*TabData `json:"tab_datas"`
}

type TabData struct {
	Ttype   TabType    `json:"tab_type"`
	Items   []*TabItem `json:"items"`
	OrderBy []int      `json:"order_by"`
	Flipped uuid.UUID  `json:"flipped"`
}

type TabItem struct {
	ItemID uuid.UUID `json:"item_id"`
	Ttype  TabType   `json:"tab_type"`

	Selected bool `json:"selected"`
	Expanded bool `json:"expanded"`
}

func (tbd *TabData) AddTabItem(ti *TabItem) *TabItem {
	ti.Ttype = tbd.Ttype

	if ti.ItemID == uuid.Nil {
		ti.ItemID = uuid.New()
	}
	if tbd.Items == nil {
		tbd.Items = []*TabItem{}
	}
	tbd.Items = append(tbd.Items, ti)
	return ti
}

func (tbd *TabData) GetTabItem(itemID uuid.UUID) (*TabItem, error) {
	if tbd.Items == nil {
		tbd.Items = []*TabItem{}
	} else {
		for _, v := range tbd.Items {
			if v.ItemID == itemID {
				return v, nil
			}
		}
	}
	return nil, errors.New("no item found")
}

func (pd *PageData) IsActive(tt TabType) bool {
	for _, t := range pd.TabDatas {
		if t.Ttype == tt {
			return true
		}
	}
	return false
}

func (pd *PageData) SetActive(td *TabData, v bool) {
	if td == nil {
		return
	}
	for i, t := range pd.TabDatas {
		if t.Ttype == td.Ttype {
			if v {
				return
			} else {
				pd.TabDatas = append(pd.TabDatas[:i], pd.TabDatas[i+1:]...)
			}
		}
	}
	if v {
		pd.TabDatas = append(pd.TabDatas, td)
	}
}

func (pgd *PageData) GetTabDataByType(tt TabType) *TabData {
	if tt == Invalid {
		return nil
	}
	for _, t := range pgd.TabDatas {
		if t.Ttype == tt {
			return t
		}
	}
	td := &TabData{
		Ttype:   tt,
		Items:   []*TabItem{},
		OrderBy: []int{},
		Flipped: uuid.Nil,
	}
	pgd.TabDatas = append(pgd.TabDatas, td)
	return td
}

func InitPageData(userID uuid.UUID, sessionID uuid.UUID, tabID string) *PageData {
	pd := &PageData{
		SessionID: sessionID,
		UserID:    userID,
		TabID:     tabID,
		TabDatas:  []*TabData{},
	}
	return pd
}

func (pd *PageData) GetTabButtonData() []TabButtonData {
	var retval []TabButtonData
	for _, tt := range GetTabTypes() {
		active := false
		for _, tbd := range pd.TabDatas {
			if tbd.Ttype == tt {
				active = true
			}
		}
		retval = append(retval, TabButtonData{Ttype: tt, Active: active})
	}
	return retval
}

// Marshalling

func (ti *TabItem) MarshalJSON() ([]byte, error) {
	configpre := struct {
		ItemID   string `json:"item_id"`
		Ttype    string `json:"tab_type"`
		Selected bool   `json:"selected"`
		Expanded bool   `json:"expanded"`
	}{
		ItemID:   ti.ItemID.String(),
		Ttype:    ti.Ttype.String(),
		Expanded: ti.Expanded,
		Selected: ti.Selected,
	}
	marshalled, err := json.Marshal(configpre)
	return marshalled, err
}

func (ti *TabItem) UnmarshalJSON(data []byte) error {
	var irJson struct {
		ItemID   string `json:"item_id"`
		Ttype    string `json:"tab_type"`
		Selected bool   `json:"selected"`
		Expanded bool   `json:"expanded"`
	}
	err := json.Unmarshal(data, &irJson)
	if err != nil {
		return err
	}
	ti.Expanded = irJson.Expanded
	ti.Selected = irJson.Selected
	ti.Ttype = String2TabType(irJson.Ttype)
	ti.ItemID, err = uuid.Parse(irJson.ItemID)
	return err
}

func (tbd *TabData) MarshalJSON() ([]byte, error) {
	items := MapFilterSlice(tbd.Items,
		func(t *TabItem) TabItem { return *t },
		func(t *TabItem) bool { return t != nil },
	)
	configpre := struct {
		Ttype   string    `json:"tab_type"`
		Items   []TabItem `json:"items"`
		OrderBy []int     `json:"order_by"`
		Flipped string    `json:"flipped"`
	}{
		Ttype:   tbd.Ttype.String(),
		Items:   items,
		OrderBy: tbd.OrderBy,
		Flipped: tbd.Flipped.String(),
	}
	marshalled, err := json.Marshal(configpre)
	return marshalled, err
}

func (tbd *TabData) UnmarshalJSON(data []byte) error {
	var irJson struct {
		Ttype   string     `json:"tab_type"`
		Items   []*TabItem `json:"items"`
		OrderBy []int      `json:"order_by"`
		Flipped string     `json:"flipped"`
	}
	err := json.Unmarshal(data, &irJson)
	if err != nil {
		return err
	}
	tbd.OrderBy = irJson.OrderBy
	flippedID, err := uuid.Parse(irJson.Flipped)
	if err != nil {
		return err
	}
	tbd.Flipped = flippedID
	tbd.Ttype = String2TabType(irJson.Ttype)
	tbd.Items = irJson.Items
	return nil
}

func (pd *PageData) MarshalJSON() ([]byte, error) {
	configpre := struct {
		SessionID string     `json:"session_id"`
		TabID     string     `json:"tab_id"`
		UserID    string     `json:"user_id"`
		TabDatas  []*TabData `json:"tab_datas"`
	}{
		SessionID: pd.SessionID.String(),
		UserID:    pd.UserID.String(),
		TabID:     pd.TabID,
		TabDatas:  pd.TabDatas,
	}
	marshalled, err := json.Marshal(configpre)
	return marshalled, err
}

func (pd *PageData) UnmarshalJSON(data []byte) error {
	var irJson struct {
		SessionID string     `json:"session_id"`
		UserID    string     `json:"user_id"`
		TabID     string     `json:"tab_id"`
		TabDatas  []*TabData `json:"tab_datas"`
	}
	err := json.Unmarshal(data, &irJson)
	if err != nil {
		return err
	}
	pd.TabID = irJson.TabID
	pd.TabDatas = irJson.TabDatas
	pd.SessionID, err = uuid.Parse(irJson.SessionID)
	if err != nil {
		return err
	}
	pd.UserID, err = uuid.Parse(irJson.UserID)
	return err
}
