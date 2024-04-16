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

func SetupAPIroutes(e *echo.Echo) error {
	// TODO: prepopulate with empty tabDatas
	pageData := viewutils.PageData{TabDatas: []*viewutils.TabData{}}
	//TODO: get userID from session / figure out auth
	userID := uuid.New()

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
		return RenderTab(TabDeactivateRenderer, c, &pageData, tabdata)
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

		tabdata, err := db.ReadTabData(*tt, userID)
		return RenderTab(TabActivateRenderer, c, &pageData, tabdata)
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

		tabdata, err := db.ReadTabData(*tt, userID)
		return RenderTab(TabMaximizeRenderer, c, &pageData, tabdata)
	})
	return nil
}
