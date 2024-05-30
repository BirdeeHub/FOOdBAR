package srvapi

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	foodlib "FOOdBAR/FOOlib"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitServer(signingKey []byte, listenOn string, staticFilesystem fs.FS) {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// TODO: figure out how to HTTPS
	// e.Pre(middleware.HTTPSRedirect())

	e.StaticFS("/static", echo.MustSubFS(staticFilesystem, "static"))

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s", foodlib.PagePrefix))
	})

	err := SetupLoginRoutes(e, signingKey)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("login route setup failed: "+err.Error()),
		)
	}

	// NOTE: Authenticated routes below

	r := e.Group(fmt.Sprintf("%s", foodlib.PagePrefix))

	r.Use(GetJWTmiddlewareWithConfig(signingKey))

	r.Use(middleware.Logger())
	r.StaticFS("/static", echo.MustSubFS(staticFilesystem, "FOOstatic"))
	// r.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
	// 	cssmiddleware := templ.NewCSSMiddleware(hndl, views.StaticStyles...)
	// 	cssmiddleware.Path = fmt.Sprintf("%s/styles/templ.css", viewutils.PagePrefix)
	// 	return cssmiddleware
	// }))

	err = SetupTabCtlroutes(r)
	if err != nil {
		e.Logger.Print(err)
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}

	err = SetupModalAPIroutes(r)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}

	err = SetupFlipAPIroutes(r)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}



	// TODO: figure out how to HTTPS
	// e.Logger.Fatal(e.StartTLS(":42069", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.Start(listenOn))
}
