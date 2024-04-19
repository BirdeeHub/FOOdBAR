package srvapi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"FOOdBAR/views"
	"FOOdBAR/views/viewutils"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/FOOdBAR")
	})

	userID := uuid.New()
	e.GET("/login", func(c echo.Context) error {
		claims := jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    "Bearer "+t,
			Path:     fmt.Sprintf("%s", viewutils.PagePrefix),
			SameSite: http.SameSiteNoneMode,
		})
		// return c.NoContent(http.StatusOK)
		return c.Redirect(http.StatusPermanentRedirect, "/FOOdBAR")
	})

	r := e.Group(fmt.Sprintf("%s", viewutils.PagePrefix))

	config := echojwt.Config{
		ContextKey: "token",
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			return err
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))

	r.Use(middleware.Logger())
	r.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		cssmiddleware := templ.NewCSSMiddleware(hndl, views.StaticStyles...)
		cssmiddleware.Path = fmt.Sprintf("%s/styles/templ.css", viewutils.PagePrefix)
		return cssmiddleware
	}))
	r.Static("/images", "images")

	err := SetupAPIroutes(r, userID)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}

	e.Logger.Fatal(e.Start(":42069"))
}
