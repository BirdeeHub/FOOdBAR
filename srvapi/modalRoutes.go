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

func SetupModalAPIroutes(e *echo.Group) error {

	e.POST("/api/submitItemInfo/:type/:itemID", func(c echo.Context) error {
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("Invalid tab type"))
		}
		itemID, err := uuid.Parse(c.Param("itemID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("itemID is not a valid UUID"))
		}
		td := pageData.GetTabDataByType(tt)
		c.Logger().Print(td)
		ti := td.GetTabItem(itemID)
		c.Logger().Print(ti)
		switch tt {
		case foodlib.Recipe:
			err = db.SubmitPantryItem(c, pageData, td, ti)
		case foodlib.Pantry:
			err = db.SubmitPantryItem(c, pageData, td, ti)
		case foodlib.Menu:
			err = db.SubmitPantryItem(c, pageData, td, ti)
		case foodlib.Preplist:
			err = db.SubmitPantryItem(c, pageData, td, ti)
		case foodlib.Shopping:
			err = db.SubmitPantryItem(c, pageData, td, ti)
		case foodlib.Events:
			err = db.SubmitPantryItem(c, pageData, td, ti)
		case foodlib.Customer:
			err = db.SubmitPantryItem(c, pageData, td, ti)
		case foodlib.Earnings:
			err = db.SubmitPantryItem(c, pageData, td, ti)
		}
		if err != nil {
			c.Logger().Print(err)
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("failed to submit item"))
		}
		return nil
	})

	//TODO: this route
	e.GET("/api/modalGetNewField/:type/:itemID/:field", func(c echo.Context) error {
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("Invalid tab type"))
		}
		itemID, err := uuid.Parse(c.Param("itemID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("itemID is not a valid UUID"))
		}
		td := pageData.GetTabDataByType(tt)
		c.Logger().Print(td)
		ti := td.GetTabItem(itemID)
		c.Logger().Print(ti)
		return HTML(c, http.StatusOK, tabviews.OOBExtraField(c.Param("field")))
	})

	e.GET("/api/itemEditModal/open/:type/:itemID", func(c echo.Context) error {
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		itemID, err := uuid.Parse(c.Param("itemID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("Invalid tab type"))
		}
		td := pageData.GetTabDataByType(tt)
		if td.Items == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("Error: No tab open"))
		}
		item, ok := td.Items[itemID]
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("Item with that ID not found in this tab"))
		}
		return HTML(c, http.StatusOK, tabviews.ItemEditModal(tabviews.RenderModalContent(pageData, item, nil)))
	})

	e.GET("/api/itemCreateModal/open/:type", func(c echo.Context) error {
		pageData, err := foodlib.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("Invalid tab type"))
		}
		td := pageData.GetTabDataByType(tt)
		item := td.AddTabItem(&foodlib.TabItem{Expanded: false})
		return HTML(c, http.StatusOK, tabviews.ItemEditModal(tabviews.RenderModalContent(pageData, item, nil)))
	})

	return nil
}
