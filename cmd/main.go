package main

import (
	"errors"
	"fmt"
	"net/http"

	"foodbar/db"
	"foodbar/views"
	"foodbar/views/viewutils"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HTML(c echo.Context, code int, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Status = code
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

type TabRenderer interface {
	func(*viewutils.TabType, echo.Context, *viewutils.PageData, *viewutils.TabData) error
}

func TabDeactivateRenderer(tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	if !tt.IsActive() {
		tt.ToggleActive()
		data.TabDatas = append(data.TabDatas, *td)
		HTML(c, http.StatusOK, views.TabButton(*tt))
		return HTML(c, http.StatusOK, views.OOBtabViewContainer(*td))
	} else {
		return HTML(c, http.StatusOK, views.TabButton(*tt))
	}
}

func TabActivateRenderer(tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	if !tt.IsActive() {
		tt.ToggleActive()
		data.TabDatas = append(data.TabDatas, *td)
		HTML(c, http.StatusOK, views.TabButton(*tt))
		return HTML(c, http.StatusOK, views.OOBtabViewContainer(*td))
	} else {
		return HTML(c, http.StatusOK, views.TabButton(*tt))
	}
}

func TabMaximizeRenderer(tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	if !tt.IsActive() {
		tt.ToggleActive()
		data.TabDatas = append(data.TabDatas, *td)
		HTML(c, http.StatusOK, views.OOBtabButtonToggle(*tt))
	}
	return HTML(c, http.StatusOK, views.TabContainer(*td))
}

func RenderTab[TR TabRenderer](tr TR, tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
    return tr(tt, c, data, td)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		cssmiddleware := templ.NewCSSMiddleware(hndl, views.StaticStyles...)
		cssmiddleware.Path = fmt.Sprintf("%sstyles/templ.css", viewutils.PagePrefix)
		return cssmiddleware
	}))
	e.Static(fmt.Sprintf("%simages",  viewutils.PagePrefix), fmt.Sprintf("%simages", viewutils.PagePrefixNoSlash))

	pageData := viewutils.PageData{TabDatas: []viewutils.TabData{}}

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
				errors.New("not a valid tab button type"),
			)
		}
		RenderTab(TabDeactivateRenderer, tt, c, &pageData, nil)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			errors.New("not a valid tab button type"),
		)
	})

	e.GET(fmt.Sprintf("%sapi/tabButton/activate/:type", viewutils.PagePrefix), func(c echo.Context) error {
		e.Logger.Print(c)
		tt, err := viewutils.String2TabType(c.Param("type"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// TODO: fetch these tabDatas from database
		var tabdata viewutils.TabData
		switch *tt {
		case viewutils.Recipe:
			tabdata = db.NewExampleRecipeTabData()
		case viewutils.Pantry:
			tabdata = db.NewExamplePantryTabData()
		case viewutils.Menu:
			tabdata = db.NewExampleMenuTabData()
		case viewutils.Shopping:
			tabdata = db.NewExampleShoppingTabData()
		case viewutils.Preplist:
			tabdata = db.NewExamplePreplistTabData()
		case viewutils.Earnings:
			tabdata = db.NewExampleEarningsTabData()
		}
		return RenderTab(TabActivateRenderer, tt, c, &pageData, &tabdata)
	})

	e.POST(fmt.Sprintf("%sapi/tabButton/maximize/:type", viewutils.PagePrefix), func(c echo.Context) error {
		e.Logger.Print(c)
		tt, err := viewutils.String2TabType(c.Param("type"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		for _, s := range pageData.TabDatas {
			if s.Ttype.IsActive() {
				s.Ttype.ToggleActive()
				HTML(c, http.StatusOK, views.OOBtabButtonToggle(s.Ttype))
			}
		}
		pageData.TabDatas = []viewutils.TabData{}

		// TODO: fetch these tabDatas from database
		var tabdata viewutils.TabData
		switch *tt {
		case viewutils.Recipe:
			tabdata = db.NewExampleRecipeTabData()
		case viewutils.Pantry:
			tabdata = db.NewExamplePantryTabData()
		case viewutils.Menu:
			tabdata = db.NewExampleMenuTabData()
		case viewutils.Shopping:
			tabdata = db.NewExampleShoppingTabData()
		case viewutils.Preplist:
			tabdata = db.NewExamplePreplistTabData()
		case viewutils.Earnings:
			tabdata = db.NewExampleEarningsTabData()
		}
		return RenderTab(TabMaximizeRenderer, tt, c, &pageData, &tabdata)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
