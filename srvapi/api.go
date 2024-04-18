package srvapi

import (
	"errors"
	// "FOOdBAR/db"
	"FOOdBAR/views"
	"FOOdBAR/views/viewutils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetupAPIroutes(e *echo.Echo) error {
	pageData := viewutils.InitPageData(uuid.New())

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/FOOdBAR")
	})

	e.GET("/FOOdBAR", func(c echo.Context) error {
		e.Logger.Print(c)
		return HTML(c, http.StatusOK, views.Homepage(pageData))
	})

	e.DELETE("/FOOdBAR/api/tabButton/deactivate/:type", func(c echo.Context) error {
		e.Logger.Print(c)
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

	e.GET("/FOOdBAR/api/tabButton/activate/:type", func(c echo.Context) error {
		e.Logger.Print(c)
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

	e.POST("/FOOdBAR/api/tabButton/maximize/:type", func(c echo.Context) error {
		e.Logger.Print(c)
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
