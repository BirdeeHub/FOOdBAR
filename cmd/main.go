package main

import (
	"net/http"

	"foodbar/views"
	"foodbar/views/viewutils"

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

func mkRecipeItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Recipe,
		ItemID:   uuid.New(),
	}
}

func mkPantryItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Pantry,
		ItemID:   uuid.New(),
	}
}

func mkMenuItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Menu,
		ItemID:   uuid.New(),
	}
}

func mkShoppingItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Shopping,
		ItemID:   uuid.New(),
	}
}

func mkPreplistItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Preplist,
		ItemID:   uuid.New(),
	}
}

func mkEarningsItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Earnings,
		ItemID:   uuid.New(),
	}
}

func newExampleRecipeTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			mkRecipeItem("Chicken"),
			mkRecipeItem("turd sandwich"),
			mkRecipeItem("chicken masala"),
			mkRecipeItem("tacos caliente"),
		},
		Ttype: viewutils.Recipe,
	}
}

func newExamplePantryTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			mkPantryItem("Chicken"),
			mkPantryItem("turd sandwich"),
			mkPantryItem("chicken masala"),
			mkPantryItem("tacos caliente"),
		},
		Ttype: viewutils.Pantry,
	}
}

func newExampleMenuTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			mkMenuItem("Chicken"),
			mkMenuItem("turd sandwich"),
			mkMenuItem("chicken masala"),
			mkMenuItem("tacos caliente"),
		},
		Ttype: viewutils.Menu,
	}
}

func newExampleShoppingTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			mkShoppingItem("Chicken"),
			mkShoppingItem("turd sandwich"),
			mkShoppingItem("chicken masala"),
			mkShoppingItem("tacos caliente"),
		},
		Ttype: viewutils.Shopping,
	}
}

func newExamplePreplistTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			mkPreplistItem("Chicken"),
			mkPreplistItem("turd sandwich"),
			mkPreplistItem("chicken masala"),
			mkPreplistItem("tacos caliente"),
		},
		Ttype: viewutils.Preplist,
	}
}

func newExampleEarningsTabData() viewutils.TabData {
	return viewutils.TabData{
		Items: []viewutils.TabItem{
			mkEarningsItem("Chicken"),
			mkEarningsItem("turd sandwich"),
			mkEarningsItem("chicken masala"),
			mkEarningsItem("tacos caliente"),
		},
		Ttype: viewutils.Earnings,
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		return templ.NewCSSMiddleware(hndl, views.StaticStyles...)
	}))

	data := viewutils.PageData{
		TabDatas: []viewutils.TabData{
			newExampleRecipeTabData(),
			// newExamplePantryTabData(),
			// newExampleMenuTabData(),
			// newExampleShoppingTabData(),
			// newExamplePreplistTabData(),
			// newExampleEarningsTabData(),
		},
	}

	e.Static("/images", "images")

	e.GET("/", func(c echo.Context) error {
		e.Logger.Print(c)
		return HTML(c, http.StatusOK, views.Homepage(data))
	})

	e.Logger.Fatal(e.Start(":42069"))
}
