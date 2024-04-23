package viewutils

import (
	"encoding/json"
	"errors"
	"fmt"

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

var colorschemes = [3]ColorScheme{Dark, Light, None}

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
	UserID   uuid.UUID
	TabDatas []*TabData
	Palette  ColorScheme
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

var sortmethods = [4]SortMethod{Descending, Ascending, Random, Custom}

type TabData struct {
	Active  bool
	Ttype   TabType
	Items   map[uuid.UUID]*TabItem
	OrderBy map[string]SortMethod
}

type TabItem struct {
	ItemID uuid.UUID
	Ttype  TabType

	Expanded bool
	ItemType string
}

func (tbd *TabData) AddTabItem(ti *TabItem) {
	ti.Ttype = tbd.Ttype

	if ti.ItemID == uuid.Nil {
		ti.ItemID = uuid.New()
	}
	tbd.Items[ti.ItemID] = ti
}

func (tbd *TabData) GetTabItems() []*TabItem {
	var tis []*TabItem
	for _, ti := range tbd.Items {
		tis = append(tis, ti)
	}
	return tis
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
	pd := &PageData{
		UserID: userID,
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
	return pd
}

// TODO: Make it so that it clears old ones every so often
// TODO: Implement client side caching of pageData in case it times out a still-valid login?
var pageDatas = make(map[uuid.UUID]*PageData)

func GetPageData(userID uuid.UUID) *PageData {
	if pageDatas[userID] == nil {
		pageDatas[userID] = InitPageData(userID)
	}
	return pageDatas[userID]
}

func (ti *TabItem) MarshalJSON() ([]byte, error) {
	configpre := struct {
		ItemID   string
		Ttype    string
		ItemType string
		Expanded bool
	}{
		ItemID:   ti.ItemID.String(),
		Ttype:    ti.Ttype.String(),
		ItemType: ti.ItemType,
		Expanded: ti.Expanded,
	}
	marshalled, err := json.Marshal(configpre)
	return marshalled, err
}

func (ti *TabItem) UnmarshalJSON(data []byte) error {
	var irJson struct {
		ItemID   string
		Ttype    string
		ItemType string
		Expanded bool
	}
	err := json.Unmarshal(data, &irJson)
	if err != nil {
		return err
	}
	ti.Expanded = irJson.Expanded
	ttype, err := String2TabType(irJson.Ttype)
	if err != nil {
		return err
	}
	ti.Ttype = *ttype
	ti.ItemType = irJson.ItemType
	ti.ItemID, err = uuid.Parse(irJson.ItemID)
	if err != nil {
		return err
	}
	return nil
}

func (tbd *TabData) MarshalJSON() ([]byte, error) {
	itemsmap := make(map[string]TabItem)
	for k, v := range tbd.Items {
		itemsmap[k.String()] = *v
	}
	orderby := make(map[string]string)
	for k, v := range tbd.OrderBy {
		orderby[k] = string(v)
	}
	configpre := struct {
		Active  bool
		Ttype   string
		Items   map[string]TabItem
		OrderBy map[string]string
	}{
		Active:  tbd.Active,
		Ttype:   tbd.Ttype.String(),
		Items:   itemsmap,
		OrderBy: orderby,
	}
	marshalled, err := json.Marshal(configpre)
	return marshalled, err
}

func (tbd *TabData) UnmarshalJSON(data []byte) error {
	var irJson struct {
		Active  bool
		Ttype   string
		Items   map[string]TabItem
		OrderBy map[string]string
	}
	err := json.Unmarshal(data, &irJson)
	if err != nil {
		return err
	}
	tbd.Active = irJson.Active
	ttype, err := String2TabType(irJson.Ttype)
	if err != nil {
		return err
	}
	tbd.Ttype = *ttype
	for k, v := range irJson.Items {
		id, err := uuid.Parse(k)
		if err != nil {
			return err
		}
		tbd.Items[id] = &v
	}
	for k, v := range irJson.OrderBy {
		for _, x := range sortmethods {
			if v == string(x) {
				tbd.OrderBy[k] = x
			} else {
				return errors.New("invalid sort method")
			}
		}
	}
	return nil
}

func (pd *PageData) MarshalJSON() ([]byte, error) {
	configpre := struct {
		UserID   string
		TabDatas []*TabData
		Palette  string
	}{
		UserID:   pd.UserID.String(),
		TabDatas: pd.TabDatas,
		Palette:  string(pd.Palette),
	}
	marshalled, err := json.Marshal(configpre)
	return marshalled, err
}

func (pd *PageData) UnmarshalJSON(data []byte) error {
	var irJson struct {
		UserID   string
		TabDatas []*TabData
		Palette  string
	}
	err := json.Unmarshal(data, &irJson)
	if err != nil {
		return err
	}
	pd.TabDatas = irJson.TabDatas
	pd.UserID, err = uuid.Parse(irJson.UserID)
	if err != nil {
		return err
	}
	err = errors.New("invalid color scheme")
	for _, v := range colorschemes {
		if irJson.Palette == string(v) {
			pd.Palette = v
			err = nil
		}
	}
	if err != nil {
		return err
	}
	return nil
}
