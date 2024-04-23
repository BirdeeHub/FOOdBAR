package viewutils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

func GetColorSchemes() [3]ColorScheme {
	return [...]ColorScheme{Dark, Light, None}
}
func GetSortMethods() [4]SortMethod {
	return [...]SortMethod{Descending, Ascending, Random, Custom}
}

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
		Palette: None,
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

// TODO: make this grab a PageData or something like that from echo context
func GetPageData(c echo.Context) (*PageData, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	switch stringID := claims["sub"].(type) {
	case string:
		userID, err := uuid.Parse(stringID)
		if err != nil {
			return nil, err
		}
		pdcookie, err := c.Cookie(userID.String())
		if err != nil {
			pd := InitPageData(userID)
			return pd, nil
		}
		pd := &PageData{}
		pdmarshalled, err := base64.StdEncoding.DecodeString(pdcookie.Value)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(pdmarshalled, pd)
		c.Logger().Print("GET DATA FUNC")
		c.Logger().Print(pd)
		for _, t := range pd.TabDatas {
			c.Logger().Print(t)
			c.Logger().Print(t.Active)
			c.Logger().Print(t.Ttype)
			c.Logger().Print(t.Items)
			c.Logger().Print(t.OrderBy)
		}
		return pd, nil
	default:
		return nil, errors.New("invalid user id")
	}
}

func (pd *PageData) SavePageData(c echo.Context) error {
	c.Logger().Print("SAVE DATA FUNC")
	c.Logger().Print(pd)
	for _, t := range pd.TabDatas {
		c.Logger().Print(t)
		c.Logger().Print(t.Active)
		c.Logger().Print(t.Ttype)
		c.Logger().Print(t.Items)
		c.Logger().Print(t.OrderBy)
	}
	pdmarshalled, err := json.Marshal(pd)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     pd.UserID.String(),
		Value:    base64.StdEncoding.EncodeToString(pdmarshalled),
		Path:     fmt.Sprintf("%s", PagePrefix),
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)
	return nil
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
		for _, x := range GetSortMethods() {
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
	for _, v := range GetColorSchemes() {
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
