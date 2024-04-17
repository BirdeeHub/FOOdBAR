package srvapi

import (
	"errors"
	"fmt"
	"foodbar/db"
	"foodbar/views"
	"foodbar/views/viewutils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func initPageData(userID uuid.UUID) *viewutils.PageData {
	return &viewutils.PageData{
		UserID: userID,
		TabDatas: []*viewutils.TabData{
			{
				Active: true,
				Ttype: viewutils.Recipe,
				Items: nil,
			},
			{
				Active: true,
				Ttype: viewutils.Pantry,
				Items: nil,
			},
			{
				Active: true,
				Ttype: viewutils.Menu,
				Items: nil,
			},
			{
				Active: true,
				Ttype: viewutils.Preplist,
				Items: nil,
			},
			{
				Active: true,
				Ttype: viewutils.Shopping,
				Items: nil,
			},
			{
				Active: true,
				Ttype: viewutils.Earnings,
				Items: nil,
			},
		},
	}
}

func SetupAPIroutes(e *echo.Echo) error {
	// TODO: prepopulate with empty tabDatas
	pageData := initPageData(uuid.New())

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

		tabdata, err := db.ReadTabData(*tt, pageData.UserID)
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

		tabdata, err := db.ReadTabData(*tt, pageData.UserID)
		return RenderTab(TabMaximizeRenderer, c, pageData, tabdata)
	})
	return nil
}
