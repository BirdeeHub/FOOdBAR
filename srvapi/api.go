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

func SetupAPIroutes(e *echo.Group, userID uuid.UUID) error {
	var pageData *viewutils.PageData = nil

	e.GET("", func(c echo.Context) error {
		user := c.Get("token").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		c.Logger().Print(claims)
		c.Logger().Print(claims["exp"])
		if pageData == nil {
			pageData = viewutils.InitPageData(userID)
		}
		c.Logger().Print(c)
		return HTML(c, http.StatusOK, views.Homepage(pageData))
	})

	e.DELETE("/api/tabButton/deactivate/:type", func(c echo.Context) error {
		if pageData == nil {
			pageData = viewutils.InitPageData(userID)
		}
		c.Logger().Print(c)
		tt, err := viewutils.String2TabType(c.Param("type"))
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				errors.New("not a valid tab type"),
			)
		}
		tabdata, err := pageData.GetTabDataByType(*tt)
		return RenderTab(TabDeactivateRenderer, c, pageData, tabdata)
	})

	e.GET("/api/tabButton/activate/:type", func(c echo.Context) error {
		if pageData == nil {
			pageData = viewutils.InitPageData(userID)
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
		tabdata, err := pageData.GetTabDataByType(*tt)
		return RenderTab(TabActivateRenderer, c, pageData, tabdata)
	})

	e.POST("/api/tabButton/maximize/:type", func(c echo.Context) error {
		if pageData == nil {
			pageData = viewutils.InitPageData(userID)
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
		tabdata, err := pageData.GetTabDataByType(*tt)
		return RenderTab(TabMaximizeRenderer, c, pageData, tabdata)
	})
	return nil
}
