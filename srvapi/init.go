package srvapi

import (
	"errors"
	"fmt"
	"net/http"

	"foodbar/views"
	"foodbar/views/viewutils"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		cssmiddleware := templ.NewCSSMiddleware(hndl, views.StaticStyles...)
		cssmiddleware.Path = fmt.Sprintf("%sstyles/templ.css", viewutils.PagePrefix)
		return cssmiddleware
	}))
	e.Static(fmt.Sprintf("%simages", viewutils.PagePrefix), fmt.Sprintf("%simages", viewutils.PagePrefixNoSlash))

	err := SetupAPIroutes(e)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}

	e.Logger.Fatal(e.Start(":42069"))
}
