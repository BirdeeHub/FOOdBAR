package srvapi

import (
	// "FOOdBAR/db"
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/db"
	"FOOdBAR/views"
	"net/http"

	"github.com/labstack/echo/v4"
)


func SetupTabCtlroutes(e *echo.Group) error {

	mainPage := func(c echo.Context) error {
		pd, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		pd.SavePageData(c)
		return HTML(c, http.StatusOK, views.Homepage(pd))
	}

	e.GET("", mainPage)

	e.POST("", mainPage)

	e.GET("/", mainPage)

	e.POST("/", mainPage)

	e.DELETE("/api/tabButton/deactivate/:type", func(c echo.Context) error {
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabDeactivateRenderer, c, pageData, pageData.GetTabDataByType(foodlib.String2TabType(c.Param("type"))))
	})

	e.GET("/api/tabButton/activate/:type", func(c echo.Context) error {
		// TODO: fetch appropriate TabData.Items from database
		// based on sort. Implement infinite scroll for them.
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		tabdata := pageData.GetTabDataByType(tt)
		err = db.FillXTabItems(pageData.UserID, tabdata, 50)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabActivateRenderer, c, pageData, tabdata)
	})

	e.POST("/api/tabButton/maximize/:type", func(c echo.Context) error {
		// TODO: fetch appropriate TabData.Items from database
		// based on sort. Implement infinite scroll for them.
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		tabdata := pageData.GetTabDataByType(tt)
		err = db.FillXTabItems(pageData.UserID, tabdata, 50)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabMaximizeRenderer, c, pageData, pageData.GetTabDataByType(foodlib.String2TabType(c.Param("type"))))
	})

	e.POST("/api/mediaQuery", func(c echo.Context) error {
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		if c.FormValue("query") == "(prefers-color-scheme: dark)" && c.FormValue("value") == "light" {
			pageData.Palette = foodlib.Light
			pageData.SavePageData(c)
		} else {
			pageData.Palette = foodlib.Dark
			pageData.SavePageData(c)
		}
		return c.NoContent(http.StatusOK)
	})
	
	return nil
}
