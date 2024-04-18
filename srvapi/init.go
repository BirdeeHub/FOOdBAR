package srvapi

import (
	"errors"
	"net/http"

	"FOOdBAR/views"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		cssmiddleware := templ.NewCSSMiddleware(hndl, views.StaticStyles...)
		cssmiddleware.Path = "/FOOdBAR/styles/templ.css"
		return cssmiddleware
	}))
	e.Static("/FOOdBAR/images", "images")

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/FOOdBAR")
	})

	userID := uuid.New()
	err := SetupAPIroutes(e, userID)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}

	e.Logger.Fatal(e.Start(":42069"))
}
