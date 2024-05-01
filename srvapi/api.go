package srvapi

import (
	// "FOOdBAR/db"
	"FOOdBAR/views"
	foodlib "FOOdBAR/FOOlib"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)


func SetupAPIroutes(e *echo.Group, dbpath string) error {

	var mainPage func(echo.Context) error = func(c echo.Context) error {
		pd, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		pd.SavePageData(c)
		return HTML(c, http.StatusOK, views.Homepage(pd))
	}

	e.GET("", func(c echo.Context) error {
		return mainPage(c)
	})

	e.POST("", func(c echo.Context) error {
		return mainPage(c)
	})

	e.GET("/", func(c echo.Context) error {
		return mainPage(c)
	})

	e.POST("/", func(c echo.Context) error {
		return mainPage(c)
	})

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
		return RenderTab(TabActivateRenderer, c, pageData, pageData.GetTabDataByType(foodlib.String2TabType(c.Param("type"))))
	})

	e.POST("/api/tabButton/maximize/:type", func(c echo.Context) error {
		// TODO: fetch appropriate TabData.Items from database
		// based on sort. Implement infinite scroll for them.
		pageData, err := foodlib.GetPageData(c)
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

	err := SetupModalAPIroutes(e, dbpath)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}
	
	return nil
}
