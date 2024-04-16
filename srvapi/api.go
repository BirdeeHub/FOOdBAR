package srvapi

import (
	"errors"
	"fmt"
	"foodbar/db"
	"foodbar/views"
	"foodbar/views/viewutils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupAPIroutes(e *echo.Echo) error {
	pageData := viewutils.PageData{TabDatas: []viewutils.TabData{}}

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
		return RenderTab(TabDeactivateRenderer, tt, c, &pageData, nil)
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

		// TODO: fetch these tabDatas from database
		var tabdata viewutils.TabData
		switch *tt {
		case viewutils.Recipe:
			tabdata = db.NewExampleRecipeTabData()
		case viewutils.Pantry:
			tabdata = db.NewExamplePantryTabData()
		case viewutils.Menu:
			tabdata = db.NewExampleMenuTabData()
		case viewutils.Shopping:
			tabdata = db.NewExampleShoppingTabData()
		case viewutils.Preplist:
			tabdata = db.NewExamplePreplistTabData()
		case viewutils.Earnings:
			tabdata = db.NewExampleEarningsTabData()
		}
		return RenderTab(TabActivateRenderer, tt, c, &pageData, &tabdata)
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

		// TODO: fetch these tabDatas from database
		var tabdata viewutils.TabData
		switch *tt {
		case viewutils.Recipe:
			tabdata = db.NewExampleRecipeTabData()
		case viewutils.Pantry:
			tabdata = db.NewExamplePantryTabData()
		case viewutils.Menu:
			tabdata = db.NewExampleMenuTabData()
		case viewutils.Shopping:
			tabdata = db.NewExampleShoppingTabData()
		case viewutils.Preplist:
			tabdata = db.NewExamplePreplistTabData()
		case viewutils.Earnings:
			tabdata = db.NewExampleEarningsTabData()
		}
		return RenderTab(TabMaximizeRenderer, tt, c, &pageData, &tabdata)
	})
	return nil
}
