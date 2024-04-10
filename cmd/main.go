package main

import (
	"errors"
	"net/http"

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

func TabToggleRenderer(activate bool, tt *viewutils.TabType, c echo.Context, data *viewutils.PageData, td viewutils.TabData) error {
	if activate {
		if !tt.IsActive() {
			tt.ToggleActive()
			data.TabDatas = append(data.TabDatas, td)
			HTML(c, http.StatusOK, views.TabButton(!activate, *tt))
			return HTML(c, http.StatusOK, views.OOBtabViewContainer(td))
		} else {
			return HTML(c, http.StatusOK, views.TabButton(!activate, *tt))
		}
	} else {
		if tt.IsActive() {
			tt.ToggleActive()
			// TODO:: actually remove the tab from the page and page data
			return HTML(c, http.StatusOK, views.TabButton(!activate, *tt))
		} else {
			return HTML(c, http.StatusOK, views.TabButton(!activate, *tt))
		}
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		return templ.NewCSSMiddleware(hndl, views.StaticStyles...)
	}))

	data := viewutils.PageData{ TabDatas: []viewutils.TabData{} }

	e.Static("/images", "images")

	e.GET("/", func(c echo.Context) error {
		e.Logger.Print(c)
		return HTML(c, http.StatusOK, views.Homepage(data))
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
		if *tt == viewutils.Recipe {
			return TabToggleRenderer(false, tt, c, &data, newExampleRecipeTabData(true))
		} else if *tt == viewutils.Pantry {
			return TabToggleRenderer(false, tt, c, &data, newExamplePantryTabData(true))
		} else if *tt == viewutils.Menu {
			return TabToggleRenderer(false, tt, c, &data, newExampleMenuTabData(true))
		} else if *tt == viewutils.Shopping {
			return TabToggleRenderer(false, tt, c, &data, newExampleShoppingTabData(true))
		} else if *tt == viewutils.Preplist {
			return TabToggleRenderer(false, tt, c, &data, newExamplePreplistTabData(true))
		} else if *tt == viewutils.Earnings {
			return TabToggleRenderer(false, tt, c, &data, newExampleEarningsTabData(true))
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
		if *tt == viewutils.Recipe {
			return TabToggleRenderer(true, tt, c, &data, newExampleRecipeTabData(true))
		} else if *tt == viewutils.Pantry {
			return TabToggleRenderer(true, tt, c, &data, newExamplePantryTabData(true))
		} else if *tt == viewutils.Menu {
			return TabToggleRenderer(true, tt, c, &data, newExampleMenuTabData(true))
		} else if *tt == viewutils.Shopping {
			return TabToggleRenderer(true, tt, c, &data, newExampleShoppingTabData(true))
		} else if *tt == viewutils.Preplist {
			return TabToggleRenderer(true, tt, c, &data, newExamplePreplistTabData(true))
		} else if *tt == viewutils.Earnings {
			return TabToggleRenderer(true, tt, c, &data, newExampleEarningsTabData(true))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("not a valid tab button type"))
	})

	e.Logger.Fatal(e.Start(":42069"))
}
