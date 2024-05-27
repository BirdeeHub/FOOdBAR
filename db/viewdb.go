package db

import (
	foodlib "FOOdBAR/FOOlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// TODO: Get this from db instead of cookie (cookie was a bad idea, it doesnt hold the actual data but it's still too much)
// Luckily, all you need to change is this function, everything gets its pagaData via this function.
// When you do so, make it able to have multiple sessions per browser
// So that all tabs dont have the exact same view

// TODO: make db table creating function
// db schema should be tab_id TEXT, user_id TEXT, session_id TEXT, page_data BLOB, last_modified DATETIME
// A db trigger will be used to keep track of which is most recent,
// any create, update, or delete will update last_modified
func GetPageData(c echo.Context) (*foodlib.PageData, error) {
	userID, err := foodlib.GetUserFromClaims(foodlib.GetClaimsFromContext(c))
	if err != nil {
		return nil, err
	}
	SID, err := foodlib.GetSessionIDFromClaims(foodlib.GetClaimsFromContext(c))
	if err != nil {
		return nil, err
	}
	tabID := c.Request().Header.Get("tab_id")
	c.Logger().Printf("tabID: %s", tabID)
	if tabID == "" {
		tabID = uuid.New().String()
		c.Logger().Printf("tabID: %s", tabID)
	}
	// TODO: replace with db query for pageData
	// If tabID is not in db, search for the most recent
	// pageData for that SessionID. If still not found, create a new one
	pdcookie, err := c.Cookie(tabID)
	if err != nil {
		pd := foodlib.InitPageData(userID, SID, tabID)
		return pd, nil
	}
	pd := &foodlib.PageData{}
	pdmarshalled, err := base64.StdEncoding.DecodeString(pdcookie.Value)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(pdmarshalled, pd)
	if err != nil {
		return nil, err
	}
	if pd.UserID != userID || pd.SessionID != SID || pd.TabID != tabID {
		pd := foodlib.InitPageData(userID, SID, tabID)
		return pd, nil
	}
	return pd, nil
}

// TODO: save to db insted of cookie
// When you do so, make it able to have multiple sessions per browser
// So that all tabs dont have the exact same view
func SavePageData(c echo.Context, pd *foodlib.PageData) error {
	pdmarshalled, err := json.Marshal(pd)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     pd.TabID,
		Value:    base64.StdEncoding.EncodeToString(pdmarshalled),
		Path:     fmt.Sprintf("%s", foodlib.PagePrefix),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}
	c.Response().Header().Add("tab_id", pd.TabID)
	c.SetCookie(cookie)
	return nil
}
