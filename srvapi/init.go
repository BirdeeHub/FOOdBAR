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

func GenerateJWTfromIDandKey(userID uuid.UUID, key []byte) (*http.Cookie, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     "user",
		Value:    t,
		Path:     fmt.Sprintf("%s", viewutils.PagePrefix),
		SameSite: http.SameSiteStrictMode,
	}, nil
}

func Init() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	e.GET(fmt.Sprintf("%s/login", viewutils.PagePrefix), func(c echo.Context) error {
		c.Logger().Print(c)
		userID := uuid.New()
		key := []byte("secret")
		cookie, err := GenerateJWTfromIDandKey(userID, key)
		if err != nil {
			c.SetCookie(cookie)
		}
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	r := e.Group(fmt.Sprintf("%s", viewutils.PagePrefix))

	config := echojwt.Config{
		ContextKey: "user",
		TokenLookup: "cookie:user",
		ErrorHandler: func(c echo.Context, err error) error {
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login", viewutils.PagePrefix))
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

	err := SetupAPIroutes(r)
	if err != nil {
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}

	e.Logger.Fatal(e.Start(":42069"))
}
