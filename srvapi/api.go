package srvapi

import (
	"errors"
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

func SetupAPIroutes(e *echo.Group) error {
	var pageData map[uuid.UUID]*viewutils.PageData = make(map[uuid.UUID]*viewutils.PageData)


	e.GET("", func(c echo.Context) error {
		userID, err := GetUserFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		if pageData[userID] == nil {
			pageData[userID] = viewutils.InitPageData(userID)
		}
		c.Logger().Print(c)
		return HTML(c, http.StatusOK, views.Homepage(pageData[userID]))
	})

	e.DELETE("/api/tabButton/deactivate/:type", func(c echo.Context) error {
		userID, err := GetUserFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		if pageData[userID] == nil {
			pageData[userID] = viewutils.InitPageData(userID)
		}
		c.Logger().Print(c)
		tt, err := viewutils.String2TabType(c.Param("type"))
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				errors.New("not a valid tab type"),
			)
		}
		tabdata, err := pageData[userID].GetTabDataByType(*tt)
		return RenderTab(TabDeactivateRenderer, c, pageData[userID], tabdata)
	})

	e.GET("/api/tabButton/activate/:type", func(c echo.Context) error {
		userID, err := GetUserFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		if pageData[userID] == nil {
			pageData[userID] = viewutils.InitPageData(userID)
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
		tabdata, err := pageData[userID].GetTabDataByType(*tt)
		return RenderTab(TabActivateRenderer, c, pageData[userID], tabdata)
	})

	e.POST("/api/tabButton/maximize/:type", func(c echo.Context) error {
		userID, err := GetUserFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		if pageData[userID] == nil {
			pageData[userID] = viewutils.InitPageData(userID)
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
		tabdata, err := pageData[userID].GetTabDataByType(*tt)
		return RenderTab(TabMaximizeRenderer, c, pageData[userID], tabdata)
	})
	return nil
}
