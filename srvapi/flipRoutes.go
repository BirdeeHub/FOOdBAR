package srvapi

import (
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/db"
	"FOOdBAR/views/tabviews"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetupFlipAPIroutes(e *echo.Group) error {

	e.GET("/api/itemEditFlip/open/:type/:itemID", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		itemID, err := uuid.Parse(c.Param("itemID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("invalid tab type"))
		}
		td := pageData.GetTabDataByType(tt)
		if td.Items == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("error: No tab open"))
		}
		_, err = td.GetTabItem(itemID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("item with that ID not found in this tab"))
		}
		td.Flipped = itemID
		pageData.SetActive(td, true)
		err = db.SavePageData(c, pageData)
		if err != nil {
			echo.NewHTTPError(http.StatusTeapot, "Cannot marshal page data")
		}
		return HTML(c, http.StatusOK, tabviews.OOBflipTab(pageData, td))
	})

	e.GET("/api/itemEditFlip/close/:type", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("invalid tab type"))
		}
		td := pageData.GetTabDataByType(tt)
		if td.Items == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("error: No tab open"))
		}
		td.Flipped = uuid.Nil
		pageData.SetActive(td, true)
		err = db.SavePageData(c, pageData)
		if err != nil {
			echo.NewHTTPError(http.StatusTeapot, "Cannot marshal page data")
		}
		return HTML(c, http.StatusOK, tabviews.OOBflipTab(pageData, td))
	})

	e.GET("/api/itemCreateFlip/open/:type", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("invalid tab type"))
		}
		td := pageData.GetTabDataByType(tt)
		item := td.AddTabItem(&foodlib.TabItem{})
		td.Flipped = item.ItemID
		pageData.SetActive(td, true)
		err = db.SavePageData(c, pageData)
		if err != nil {
			echo.NewHTTPError(http.StatusTeapot, "Cannot marshal page data")
		}
		return HTML(c, http.StatusOK, tabviews.OOBflipTab(pageData, td))
	})

	return nil
}
