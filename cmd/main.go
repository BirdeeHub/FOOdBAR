package main

import (
	"errors"
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

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		return templ.NewCSSMiddleware(hndl, views.StaticStyles...)
	}))

	pageData := viewutils.PageData{TabDatas: []viewutils.TabData{}}

	e.Static("/images", "images")

	e.GET("/", func(c echo.Context) error {
		e.Logger.Print(c)
		return HTML(c, http.StatusOK, views.Homepage(pageData))
	})

	e.DELETE("/api/tabButton/deactivate/:type", func(c echo.Context) error {
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
			return TabToggleRenderer(true, tt, c, &pageData, nil)
		case viewutils.Pantry:
			return TabToggleRenderer(true, tt, c, &pageData, nil)
		case viewutils.Menu:
			return TabToggleRenderer(true, tt, c, &pageData, nil)
		case viewutils.Shopping:
			return TabToggleRenderer(true, tt, c, &pageData, nil)
		case viewutils.Preplist:
			return TabToggleRenderer(true, tt, c, &pageData, nil)
		case viewutils.Earnings:
			return TabToggleRenderer(true, tt, c, &pageData, nil)
		}
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			errors.New("not a valid tab button type"),
		)
	})

	e.POST("/api/tabButton/activate/:type", func(c echo.Context) error {
		e.Logger.Print(c)
		tt, err := viewutils.String2TabType(c.Param("type"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// TODO: fetch these tabDatas from database
		// TODO: implement buffered scrolling for infinite scrolling capabilities
		switch *tt {
		case viewutils.Recipe:
			tabdata := db.NewExampleRecipeTabData(true)
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Pantry:
			tabdata := db.NewExamplePantryTabData(true)
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Menu:
			tabdata := db.NewExampleMenuTabData(true)
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Shopping:
			tabdata := db.NewExampleShoppingTabData(true)
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Preplist:
			tabdata := db.NewExamplePreplistTabData(true)
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		case viewutils.Earnings:
			tabdata := db.NewExampleEarningsTabData(true)
			return TabToggleRenderer(true, tt, c, &pageData, &tabdata)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("not a valid tab button type"))
	})

	e.Logger.Fatal(e.Start(":42069"))
}
