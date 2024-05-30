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
		return HTML(c, http.StatusOK, views.Homepage())
	}

	e.GET("", mainPage)

	e.POST("", mainPage)

	e.GET("/", mainPage)

	e.POST("/", mainPage)

	e.GET("/bodycontents", func(c echo.Context) error {
		pd, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		db.SavePageData(c, pd)
		return HTML(c, http.StatusOK, views.BodyContents(pd))
	})

	e.DELETE("/api/tabButton/deactivate/:type", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabDeactivateRenderer, c, pageData, pageData.GetTabDataByType(foodlib.String2TabType(c.Param("type"))))
	})

	e.GET("/api/tabButton/activate/:type", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		tabdata := pageData.GetTabDataByType(tt)
		// TODO: Implement infinite scroll for these.
		// after the following function actually only returns X number of items
		err = db.FillXTabItems(pageData.UserID, tabdata, 50)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabActivateRenderer, c, pageData, tabdata)
	})

	e.POST("/api/tabButton/maximize/:type", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		tabdata := pageData.GetTabDataByType(tt)
		// TODO: Implement infinite scroll for these.
		// after the following function actually only returns X number of items
		err = db.FillXTabItems(pageData.UserID, tabdata, 50)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return RenderTab(TabMaximizeRenderer, c, pageData, pageData.GetTabDataByType(foodlib.String2TabType(c.Param("type"))))
	})
	
	return nil
}
