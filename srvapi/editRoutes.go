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

func SetupEditAPIroutes(e *echo.Group) error {

	e.POST("/api/submitItemInfo/:type/:itemID", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("invalid tab type"))
		}
		itemID, err := uuid.Parse(c.Param("itemID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("itemID is not a valid UUID"))
		}
		isActive := false
		if pageData.IsActive(tt) {
			isActive = true
		}
		td := pageData.GetTabDataByType(tt)
		flipped := false
		if td.Flipped != uuid.Nil {
			flipped = true
		}
		c.Logger().Print(td)
		ti, err := td.GetTabItem(itemID)
		present := true
		if err != nil {
			present = false
			ti = &foodlib.TabItem{ItemID: itemID}
			td.AddTabItem(ti)
		}
		c.Logger().Print(ti)
		// TODO: submit method for other tables
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
			return HTML(c, http.StatusUnprocessableEntity, tabviews.OOBsendBackSubmitStatus(itemID, "", err))
		}
		if isActive && !flipped {
			// TODO: Test this better
			if present {
				// Re render the item oob if it was present
				HTML(c, http.StatusOK, tabviews.OOBRenderItemContainer(pageData, td, ti))
			} else {
				// check if it should be present now
				tis, err := db.FillXTabItems(pageData.UserID, td, 0, len(td.Items))
				if err != nil {
					c.Logger().Print(err)
					return HTML(c, http.StatusInternalServerError, tabviews.OOBsendBackSubmitStatus(itemID, "", err))
				}
				found := false;
				for i, val := range tis {
					// 3: if found is true, it needs to be before this one. So put it there.
					if found {
						HTML(c, http.StatusOK, tabviews.OOBRenderItemBeforeThis(pageData, td, ti, val.ItemID))
						break
					}
					// 1: check if we found it
					if val.ItemID == ti.ItemID {
						found = true
					}
					// 2: if its the last one, deal with it now just adding to end
					if found && i >= len(tis) - 1 {
						justMore := &foodlib.TabData{
							Items: []*foodlib.TabItem{val},
							Ttype: td.Ttype,
							OrderBy: td.OrderBy,
							Flipped: td.Flipped,
						}
						tabviews.OOBmoreTabItems(pageData, justMore, true)
					}
				}
			}
		}
		db.SavePageData(c, pageData)
		return HTML(c, http.StatusOK, tabviews.OOBsendBackSubmitStatus(itemID, "Item Saved Successfully!", nil))
	})

	e.GET("/api/submitGetNewField/:type/:itemID/:field", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("invalid tab type"))
		}
		itemID, err := uuid.Parse(c.Param("itemID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("itemID is not a valid UUID"))
		}
		td := pageData.GetTabDataByType(tt)
		c.Logger().Print(td)
		ti, err := td.GetTabItem(itemID)
		if err != nil {
			td.AddTabItem(&foodlib.TabItem{ItemID: itemID})
		}
		c.Logger().Print(ti)
		return HTML(c, http.StatusOK, tabviews.OOBExtraField(c.Param("field"), itemID))
	})

	e.GET("/api/itemEditModal/open/:type/:itemID", func(c echo.Context) error {
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
		item, err := td.GetTabItem(itemID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("item with that ID not found in this tab"))
		}
		return HTML(c, http.StatusOK, tabviews.ItemEditModal(tabviews.RenderSubmissionContent(pageData, item, "", nil)))
	})

	e.GET("/api/itemCreateModal/open/:type", func(c echo.Context) error {
		pageData, err := db.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := foodlib.String2TabType(c.Param("type"))
		if tt == foodlib.Invalid {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("invalid tab type"))
		}
		td := pageData.GetTabDataByType(tt)
		item := td.AddTabItem(&foodlib.TabItem{Expanded: false})
		return HTML(c, http.StatusOK, tabviews.ItemEditModal(tabviews.RenderSubmissionContent(pageData, item, "", nil)))
	})

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
