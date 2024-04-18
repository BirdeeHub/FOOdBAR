package srvapi

import (
	"errors"
	"fmt"
	// "FOOdBAR/db"
	"FOOdBAR/views"
	"FOOdBAR/views/viewutils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetupAPIroutes(e *echo.Echo) error {
	pageData := viewutils.InitPageData(uuid.New())

	e.GET(fmt.Sprintf("%s", viewutils.PagePrefix), func(c echo.Context) error {
		e.Logger.Print(c)
		return HTML(c, http.StatusOK, views.Homepage(pageData))
	})

	e.DELETE(fmt.Sprintf("%sapi/tabButton/deactivate/:type", viewutils.PagePrefix), func(c echo.Context) error {
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

	e.GET(fmt.Sprintf("%sapi/tabButton/activate/:type", viewutils.PagePrefix), func(c echo.Context) error {
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

	e.POST(fmt.Sprintf("%sapi/tabButton/maximize/:type", viewutils.PagePrefix), func(c echo.Context) error {
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
