package srvapi

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	foodlib "FOOdBAR/FOOlib"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitServer(signingKey []byte, listenOn string, staticFilesystem fs.FS, staticFilesAuthed fs.FS) {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// TODO: figure out how to force HTTPS
	// e.Pre(middleware.HTTPSRedirect())
	// Custom handler to serve pre-compressed files if they exist
	e.Use(getUseGZmiddleware(staticFilesystem, ""))

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
	r.Use(getUseGZmiddleware(echo.MustSubFS(staticFilesAuthed, "FOOstatic"), "/FOOdBAR/static"))
	// Authed static directory at /FOOdBAR/static
	r.StaticFS("/static", echo.MustSubFS(staticFilesAuthed, "FOOstatic"))
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

	err = SetupEditAPIroutes(r)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}

	// TODO: figure out how to force HTTPS
	// e.Logger.Fatal(e.StartTLS(":42069", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.Start(listenOn))
}

func getUseGZmiddleware(static fs.FS, prefix string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestPath := strings.TrimPrefix(c.Request().URL.Path, prefix)
			if len(requestPath) > 0 && requestPath[0] == '/' {
				requestPath = requestPath[1:]
			}
			gzippedFilePath := requestPath + ".gz"
			if _, err := fs.Stat(static, gzippedFilePath); err == nil {
				filebytes, err := fs.ReadFile(static, gzippedFilePath)
				if err == nil {
					contentType, err := foodlib.ContentTypeFromExt(filepath.Ext(requestPath))
					if err != nil {
						c.Logger().Printf("static GZ middleware substitutor error, returning original\nStack trace:\n%s", err.Error())
						return next(c)
					}
					return GZIP(c, http.StatusOK, contentType, filebytes)
				}
			}
			return next(c)
		}
	}
}

