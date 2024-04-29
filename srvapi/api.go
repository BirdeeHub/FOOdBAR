package srvapi

import (
	// "FOOdBAR/db"
	"FOOdBAR/views"
	"FOOdBAR/views/viewutils"
	"net/http"
	"github.com/labstack/echo/v4"
)

func SetupAPIroutes(e *echo.Group, dbpath string) error {

	e.GET("", func(c echo.Context) error {
		pd, err := viewutils.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		pd.SavePageData(c)
		return HTML(c, http.StatusOK, views.Homepage(pd))
	})

	e.POST("", func(c echo.Context) error {
		pd, err := viewutils.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		pd.SavePageData(c)
		return HTML(c, http.StatusOK, views.Homepage(pd))
	})

	e.POST("/api/mediaQuery", func(c echo.Context) error {
		pageData, err := viewutils.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		if c.FormValue("query") == "(prefers-color-scheme: dark)" && c.FormValue("value") != "dark" {
			pageData.Palette = viewutils.Light
			pageData.SavePageData(c)
		} else {
			pageData.Palette = viewutils.Dark
			pageData.SavePageData(c)
		}
		return c.NoContent(http.StatusOK)
	})

	e.DELETE("/api/tabButton/deactivate/:type", func(c echo.Context) error {
		pageData, err := viewutils.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabDeactivateRenderer, c, pageData, pageData.GetTabDataByType(viewutils.String2TabType(c.Param("type"))))
	})

	e.GET("/api/tabButton/activate/:type", func(c echo.Context) error {
		// TODO: fetch appropriate TabData.Items from database
		// based on sort. Implement infinite scroll for them.
		pageData, err := viewutils.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabActivateRenderer, c, pageData, pageData.GetTabDataByType(viewutils.String2TabType(c.Param("type"))))
	})

	e.POST("/api/tabButton/maximize/:type", func(c echo.Context) error {
		// TODO: fetch appropriate TabData.Items from database
		// based on sort. Implement infinite scroll for them.
		pageData, err := viewutils.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabMaximizeRenderer, c, pageData, pageData.GetTabDataByType(viewutils.String2TabType(c.Param("type"))))
	})
	return nil
}
