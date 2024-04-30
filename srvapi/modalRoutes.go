package srvapi

import (
	"FOOdBAR/views/tabviews"
	"FOOdBAR/views/viewutils"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetupModalAPIroutes(e *echo.Group, dbpath string) error {

	e.POST("/api/itemEditModal/open/:type/:itemID", func(c echo.Context) error {
		pageData, err := viewutils.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		itemID, err := uuid.Parse(c.Param("itemID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := viewutils.String2TabType(c.Param("type"))
		if tt == viewutils.Invalid {
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
		return HTML(c, http.StatusOK, tabviews.ItemEditModal(tabviews.RenderModalContent(item)))
	})

	e.POST("/api/itemCreateModal/open/:type", func(c echo.Context) error {
		pageData, err := viewutils.GetPageData(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		tt := viewutils.String2TabType(c.Param("type"))
		if tt == viewutils.Invalid {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("Invalid tab type"))
		}
		td := pageData.GetTabDataByType(tt)
		item := td.AddTabItem(&viewutils.TabItem{Expanded: false})
		return HTML(c, http.StatusOK, tabviews.ItemEditModal(tabviews.RenderModalContent(item)))
	})

	return nil
}
