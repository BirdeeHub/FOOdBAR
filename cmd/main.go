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

func newExampleRecipeTabData() views.TabData {
	recipeItem := views.TabItem{
		ItemName: "Chicken",
		Ttype:    views.Recipe,
		ItemID:   uuid.New(),
	}
	recipeItem1 := views.TabItem{
		ItemName: "turd sandwich",
		Ttype:    views.Recipe,
		ItemID:   uuid.New(),
	}
	recipeItem2 := views.TabItem{
		ItemName: "chicken masala",
		Ttype:    views.Recipe,
		ItemID:   uuid.New(),
	}
	recipeItem3 := views.TabItem{
		ItemName: "tacos caliente",
		Ttype:    views.Recipe,
		ItemID:   uuid.New(),
	}

	recipeTabdata := views.TabData{
		Items: []views.TabItem{recipeItem, recipeItem1, recipeItem2, recipeItem3},
		Ttype: views.Recipe,
	}

	return recipeTabdata
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
