package srvapi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"FOOdBAR/db"
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
func WipeAuth(c echo.Context) {
    // Set the cookie with the same name and an expiration time in the past
    expiration := time.Now().AddDate(0, 0, -1)
    cookie := http.Cookie{
        Name:    "user",
        Value:   "",
		Path:     fmt.Sprintf("%s", viewutils.PagePrefix),
        Expires: expiration,
		SameSite: http.SameSiteStrictMode,
    }
    http.SetCookie(c.Response().Writer, &cookie)
}

func Init() {
	e := echo.New()
	e.Use(middleware.Logger())

	// TODO: figure out how to HTTPS
	// e.Pre(middleware.HTTPSRedirect())

	// TODO: get a much better key from a file
	signingKey := []byte("secret-passphrase-willitwork")

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
		// TODO: return visible error message if fail
		username := c.FormValue("username")
		password := c.FormValue("password")
		userID, err := db.AuthUser(username, password)
		if err != nil {
			// TODO: return visible error message if fail as hx oob swap
			WipeAuth(c)
			c.Logger().Print(err)
			return echo.NewHTTPError(http.StatusNotAcceptable, err)
		}
		c.Logger().Print(userID)
		cookie, err := GenerateJWTfromIDandKey(userID, signingKey)
		if err != nil {
			WipeAuth(c)
			c.Logger().Print(err)
			return echo.NewHTTPError(http.StatusTeapot, err)
		}
		c.SetCookie(cookie)
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	e.POST(fmt.Sprintf("%s/submitsignup", viewutils.PagePrefix), func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		confirmpassword := c.FormValue("confirmpassword")
		if password != confirmpassword {
			WipeAuth(c)
			// TODO: return visible error message if fail as hx oob swap
			return echo.NewHTTPError(http.StatusNotAcceptable, errors.New("Passwords don't match"))
		}
		userID, err := db.CreateUser(username, password)
		if err != nil {
			// TODO: return visible error message if fail as hx oob swap
			WipeAuth(c)
			c.Logger().Print(err)
			return echo.NewHTTPError(http.StatusNotAcceptable, err)
		}
		c.Logger().Print(userID)
		cookie, err := GenerateJWTfromIDandKey(userID, signingKey)
		if err != nil {
			WipeAuth(c)
			c.Logger().Print(err)
			return echo.NewHTTPError(http.StatusTeapot, err)
		}
		c.SetCookie(cookie)
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	// NOTE: Authenticated routes below

	r := e.Group(fmt.Sprintf("%s", viewutils.PagePrefix))

	jwtConfig := echojwt.Config{
		ContextKey: "user",
		TokenLookup: "cookie:user",
		SuccessHandler: func(c echo.Context) {
			userID, err := GetUserFromToken(c)
			if err != nil {
				WipeAuth(c)
				c.Logger().Print(err)
				return
			}
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			expirationtime, err := claims.GetExpirationTime()
			if err != nil {
				WipeAuth(c)
				c.Logger().Print(err)
				return
			}
			if time.Until(expirationtime.Time) < time.Hour {
				cookie, err := GenerateJWTfromIDandKey(userID, signingKey)
				if err != nil {
					WipeAuth(c)
					c.Logger().Print(err)
					return
				}
				c.SetCookie(cookie)
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			WipeAuth(c)
			c.Logger().Print(err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login", viewutils.PagePrefix))
		},
		SigningKey: signingKey,
	}
	r.Use(echojwt.WithConfig(jwtConfig))

	r.Use(middleware.Logger())
	r.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
		cssmiddleware := templ.NewCSSMiddleware(hndl, views.StaticStyles...)
		cssmiddleware.Path = fmt.Sprintf("%s/styles/templ.css", viewutils.PagePrefix)
		return cssmiddleware
	}))
	r.Static("/images", "images")

	err := SetupAPIroutes(r)
	if err != nil {
		e.Logger.Print(err)
		echo.NewHTTPError(
			http.StatusTeapot,
			errors.New("server api setup failed: "+err.Error()),
		)
	}

	// TODO: figure out how to HTTPS
	// e.Logger.Fatal(e.StartTLS(":42069", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.Start(":42069"))
}
