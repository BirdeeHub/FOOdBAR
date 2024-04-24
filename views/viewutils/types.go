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

type TabButtonData struct {
	Ttype  TabType `json:"tab_type"`
	Active bool    `json:"active"`
}

type PageData struct {
	UserID   uuid.UUID   `json:"user_id"`
	TabDatas []*TabData  `json:"tab_datas"`
	Palette  ColorScheme `json:"palette"`
}

type TabData struct {
	Ttype   TabType                `json:"tab_type"`
	Items   map[uuid.UUID]*TabItem `json:"items"`
	OrderBy map[string]SortMethod  `json:"order_by"`
}

type TabItem struct {
	ItemID uuid.UUID `json:"item_id"`
	Ttype  TabType   `json:"tab_type"`

	Expanded bool   `json:"expanded"`
	ItemType string `json:"item_type"`
}

func (tbd *TabData) AddTabItem(ti *TabItem) *TabItem {
	ti.Ttype = tbd.Ttype

	if ti.ItemID == uuid.Nil {
		ti.ItemID = uuid.New()
	}
	tbd.Items[ti.ItemID] = ti
	return ti
}

func (tbd *TabData) GetTabItems() []*TabItem {
	var tis []*TabItem
	for _, ti := range tbd.Items {
		tis = append(tis, ti)
	}
	return tis
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
	for _, t := range pgd.TabDatas {
		if t.Ttype == tt {
			return t
		}
	}
	td := &TabData{
		Ttype:   tt,
		Items:   make(map[uuid.UUID]*TabItem),
		OrderBy: make(map[string]SortMethod),
	}
	pgd.TabDatas = append(pgd.TabDatas, td)
	return td
}

func InitPageData(userID uuid.UUID) *PageData {
	pd := &PageData{
		UserID:   userID,
		Palette:  None,
		TabDatas: []*TabData{},
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
			c.Logger().Print("New Page Data")
			return pd, nil
		}
		pd := &PageData{}
		pdmarshalled, err := base64.StdEncoding.DecodeString(pdcookie.Value)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(pdmarshalled, pd)
		if err != nil {
			return nil, err
		}
		return pd, nil
	default:
		return nil, errors.New("invalid user id")
	}
}

func (pd *PageData) SavePageData(c echo.Context) error {
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
		ItemID   string `json:"item_id"`
		Ttype    string `json:"tab_type"`
		ItemType string `json:"item_type"`
		Expanded bool   `json:"expanded"`
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
		ItemID   string `json:"item_id"`
		Ttype    string `json:"tab_type"`
		ItemType string `json:"item_type"`
		Expanded bool   `json:"expanded"`
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
		Ttype   string             `json:"tab_type"`
		Items   map[string]TabItem `json:"items"`
		OrderBy map[string]string  `json:"order_by"`
	}{
		Ttype:   tbd.Ttype.String(),
		Items:   itemsmap,
		OrderBy: orderby,
	}
	marshalled, err := json.Marshal(configpre)
	return marshalled, err
}

func (tbd *TabData) UnmarshalJSON(data []byte) error {
	var irJson struct {
		Ttype   string             `json:"tab_type"`
		Items   map[string]TabItem `json:"items"`
		OrderBy map[string]string  `json:"order_by"`
	}
	err := json.Unmarshal(data, &irJson)
	if err != nil {
		return err
	}
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
		UserID   string     `json:"user_id"`
		TabDatas []*TabData `json:"tab_datas"`
		Palette  string     `json:"palette"`
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
		UserID   string     `json:"user_id"`
		TabDatas []*TabData `json:"tab_datas"`
		Palette  string     `json:"palette"`
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
