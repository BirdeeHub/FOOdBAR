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

func TabToggleRenderer(activate bool, tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td *viewutils.TabData) error {
	if activate {
		if !tt.IsActive() {
			tt.ToggleActive()
			data.TabDatas = append(data.TabDatas, *td)
			HTML(c, http.StatusOK, views.TabButton(*tt))
			return HTML(c, http.StatusOK, views.OOBtabViewContainer(*td))
		} else {
			return HTML(c, http.StatusOK, views.TabButton(*tt))
		}
	} else {
		if tt.IsActive() {
			tt.ToggleActive()
			for i, td := range data.TabDatas {
				if td.Ttype == *tt {
					data.TabDatas = append(data.TabDatas[:i], data.TabDatas[i+1:]...)
				}
			}
			return HTML(c, http.StatusOK, views.OOBtabButtonToggle(*tt))
		} else {
			return nil
		}
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
		switch *tt {
		case viewutils.Recipe:
			return TabToggleRenderer(false, tt, c, &pageData, nil)
		case viewutils.Pantry:
			return TabToggleRenderer(false, tt, c, &pageData, nil)
		case viewutils.Menu:
			return TabToggleRenderer(false, tt, c, &pageData, nil)
		case viewutils.Shopping:
			return TabToggleRenderer(false, tt, c, &pageData, nil)
		case viewutils.Preplist:
			return TabToggleRenderer(false, tt, c, &pageData, nil)
		case viewutils.Earnings:
			return TabToggleRenderer(false, tt, c, &pageData, nil)
		}
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
		switch *tt {
		case viewutils.Recipe:
			tabdata := db.NewExampleRecipeTabData()
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Pantry:
			tabdata := db.NewExamplePantryTabData()
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Menu:
			tabdata := db.NewExampleMenuTabData()
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Shopping:
			tabdata := db.NewExampleShoppingTabData()
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Preplist:
			tabdata := db.NewExamplePreplistTabData()
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Earnings:
			tabdata := db.NewExampleEarningsTabData()
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("not a valid tab button type"))
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
		switch *tt {
		case viewutils.Recipe:
			tabdata := db.NewExampleRecipeTabData()
			return TabMaximizeRenderer(tt, c, &pageData, &tabdata)
		case viewutils.Pantry:
			tabdata := db.NewExamplePantryTabData()
			return TabMaximizeRenderer(tt, c, &pageData, &tabdata)
		case viewutils.Menu:
			tabdata := db.NewExampleMenuTabData()
			return TabMaximizeRenderer(tt, c, &pageData, &tabdata)
		case viewutils.Shopping:
			tabdata := db.NewExampleShoppingTabData()
			return TabMaximizeRenderer(tt, c, &pageData, &tabdata)
		case viewutils.Preplist:
			tabdata := db.NewExamplePreplistTabData()
			return TabMaximizeRenderer(tt, c, &pageData, &tabdata)
		case viewutils.Earnings:
			tabdata := db.NewExampleEarningsTabData()
			return TabMaximizeRenderer(tt, c, &pageData, &tabdata)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("not a valid tab button type"))
	})

	e.Logger.Fatal(e.Start(":42069"))
}
