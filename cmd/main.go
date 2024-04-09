package main

import (
	"net/http"

	views "foodbar/views"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HTML(c echo.Context, code int, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Status = code
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func mkRecipeItem(name string) views.TabItem {
	return views.TabItem{
		ItemName: name,
		Ttype:    views.Recipe,
		ItemID:   uuid.New(),
	}
}

func newExampleRecipeTabData() views.TabData {
	return views.TabData{
		Items: []views.TabItem{mkRecipeItem("Chicken"), mkRecipeItem("turd sandwich"), mkRecipeItem("chicken masala"), mkRecipeItem("tacos caliente")},
		Ttype: views.Recipe,
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		return templ.NewCSSMiddleware(hndl, views.StaticStyles...)
	}))

	data := views.PageData{
		TabDatas: []views.TabData{newExampleRecipeTabData()},
	}

	e.Static("/images", "images")

	e.GET("/", func(c echo.Context) error {
		e.Logger.Print(c)
		return HTML(c, http.StatusOK, views.Homepage(data))
	})

	e.Logger.Fatal(e.Start(":42069"))
}
