package srvapi

import (
	"errors"
	"time"
	// "FOOdBAR/db"
	"FOOdBAR/views"
	"FOOdBAR/views/viewutils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetClaimFromToken(c echo.Context, claim string) interface{} {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims[claim]
}

func GetUserFromToken(c echo.Context) (uuid.UUID, error) {
	switch userID := GetClaimFromToken(c, "sub").(type) {
	case string:
		return uuid.Parse(userID)
	default:
		return uuid.Nil, errors.New("invalid userID")
	}
}

// TODO: Make it so that it clears old ones
// TODO: Implement client side caching of pageData in case it times out a still-valid login
var pageDatas map[uuid.UUID]*viewutils.PageData = make(map[uuid.UUID]*viewutils.PageData)

func GetPageData(userID uuid.UUID) *viewutils.PageData {
	if pageDatas[userID] == nil {
		pageDatas[userID] = viewutils.InitPageData(userID)
	}
	pageDatas[userID].LastActive = time.Now()
	return pageDatas[userID]
}

func SetupAPIroutes(e *echo.Group) error {

	e.GET("", func(c echo.Context) error {
		userID, err := GetUserFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		c.Logger().Print(c)
		return HTML(c, http.StatusOK, views.Homepage(GetPageData(userID)))
	})

	e.DELETE("/api/tabButton/deactivate/:type", func(c echo.Context) error {
		userID, err := GetUserFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		c.Logger().Print(c)
		tt, err := viewutils.String2TabType(c.Param("type"))
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				errors.New("not a valid tab type"),
			)
		}
		pageData := GetPageData(userID)
		tabdata, err := pageData.GetTabDataByType(*tt)
		return RenderTab(TabDeactivateRenderer, c, pageData, tabdata)
	})

	e.GET("/api/tabButton/activate/:type", func(c echo.Context) error {
		userID, err := GetUserFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		c.Logger().Print(c)
		tt, err := viewutils.String2TabType(c.Param("type"))
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				errors.New("not a valid tab type"),
			)
		}

		// TODO: fetch appropriate TabData.Items from database
		// based on sort. Implement infinite scroll for them.
		pageData := GetPageData(userID)
		tabdata, err := pageData.GetTabDataByType(*tt)
		return RenderTab(TabActivateRenderer, c, pageData, tabdata)
	})

	e.POST("/api/tabButton/maximize/:type", func(c echo.Context) error {
		userID, err := GetUserFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		c.Logger().Print(c)
		tt, err := viewutils.String2TabType(c.Param("type"))
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				errors.New("not a valid tab type"),
			)
		}

		// TODO: fetch appropriate TabData.Items from database
		// based on sort. Implement infinite scroll for them.
		pageData := GetPageData(userID)
		tabdata, err := pageData.GetTabDataByType(*tt)
		return RenderTab(TabMaximizeRenderer, c, pageData, tabdata)
	})
	return nil
}
