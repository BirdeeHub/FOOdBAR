package srvapi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"FOOdBAR/views"
	"FOOdBAR/views/loginPage"
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

	// TODO: figure out how to HTTPS
	// e.Pre(middleware.HTTPSRedirect())

	// TODO: get a much better key from a file
	signingKey := []byte("secret-passphrasewillitwork")

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	e.GET(fmt.Sprintf("%s/login", viewutils.PagePrefix), func(c echo.Context) error {
		return HTML(c, http.StatusOK, loginPage.LoginPage("login"))
	})

	e.GET(fmt.Sprintf("%s/loginform/:formtype", viewutils.PagePrefix), func(c echo.Context) error {
		formtype := c.Param("formtype")
		if formtype == "login" {
			return HTML(c, http.StatusOK, loginPage.LoginPageContents(loginPage.LoginType))
		} else if formtype == "signup" {
			return HTML(c, http.StatusOK, loginPage.LoginPageContents(loginPage.SignupType))
		} else {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, errors.New("Invalid formtype"))
		}
	})

	e.POST(fmt.Sprintf("%s/submitlogin", viewutils.PagePrefix), func(c echo.Context) error {
		// TODO: check login and retrieve uuid
		// TODO: return visible error message if fail
		username := c.FormValue("username")
		password := c.FormValue("password")
		c.Logger().Print(username)
		c.Logger().Print(password)
		c.Logger().Print(c)
		userID := uuid.New()
		cookie, err := GenerateJWTfromIDandKey(userID, signingKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusTeapot, err)
		}
		c.SetCookie(cookie)
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	e.POST(fmt.Sprintf("%s/submitsignup", viewutils.PagePrefix), func(c echo.Context) error {
		// TODO: generate uuid and store with user and pass in db
		// TODO: return visible error message if fail
		username := c.FormValue("username")
		password := c.FormValue("password")
		c.Logger().Print(username)
		c.Logger().Print(password)
		c.Logger().Print(c)
		userID := uuid.New()
		cookie, err := GenerateJWTfromIDandKey(userID, signingKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusTeapot, err)
		}
		c.SetCookie(cookie)
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	// NOTE: Authenticated routes below

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
		SigningKey: signingKey,
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

	// TODO: figure out how to HTTPS
	// e.Logger.Fatal(e.StartTLS(":42069", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.Start(":42069"))
}
