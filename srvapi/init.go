package srvapi

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/db"
	"FOOdBAR/views/loginPage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var lockouts = make(map[string]*lockoutEntry)

type lockoutEntry struct {
	Num int
	Last time.Time
	IP string
}

func InitServer(signingKey []byte, listenOn string) {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// TODO: figure out how to HTTPS
	// e.Pre(middleware.HTTPSRedirect())

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s", foodlib.PagePrefix))
	})

	e.GET(fmt.Sprintf("%s/login", foodlib.PagePrefix), func(c echo.Context) error {
		WipeAuth(c)
		return HTML(c, http.StatusOK, loginPage.LoginPage("login", nil))
	})

	e.GET(fmt.Sprintf("%s/loginform/:formtype", foodlib.PagePrefix), func(c echo.Context) error {
		formtype := c.Param("formtype")
		if formtype == "login" {
			return HTML(c, http.StatusOK, loginPage.LoginPageContents(loginPage.LoginType, nil))
		} else if formtype == "signup" {
			return HTML(c, http.StatusOK, loginPage.LoginPageContents(loginPage.SignupType, nil))
		} else {
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPageContents(loginPage.LoginType, errors.New("Invalid formtype")))
		}
	})

	// TODO: implement timeout
	e.POST(fmt.Sprintf("%s/submitlogin", foodlib.PagePrefix), func(c echo.Context) error {
		if time.Since(lockouts[c.RealIP()].Last) > 10*time.Minute {
			lockouts[c.RealIP()] = nil
		} else if lockouts[c.RealIP()].Num > 9 {
			return HTML(c, http.StatusNotAcceptable, loginPage.LoginPage(loginPage.LoginType, errors.New("Too many failed login attempts, please try again later")))
		}
		username := c.FormValue("username")
		password := c.FormValue("password")
		beepboop := c.FormValue("beepboop")
		if beepboop != "" {
			WipeAuth(c)
			err := errors.New("Scraper no scraping!")
			c.Logger().Print(err)
			if lockouts[c.RealIP()] == nil {
				lockouts[c.RealIP()] = &lockoutEntry{IP: c.RealIP()}
			}
			lockouts[c.RealIP()].Num = 10
			lockouts[c.RealIP()].Last = time.Now()
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPage(loginPage.LoginType, err))
		}
		userID, err := db.AuthUser(username, password)
		if err != nil {
			WipeAuth(c)
			c.Logger().Print(err)
			if lockouts[c.RealIP()] == nil {
				lockouts[c.RealIP()] = &lockoutEntry{IP: c.RealIP()}
			}
			if strings.Contains(err.Error(), "invalid password") {
				lockouts[c.RealIP()].Num++
				lockouts[c.RealIP()].Last = time.Now()
			}
			return HTML(c, http.StatusNotAcceptable, loginPage.LoginPage(loginPage.LoginType, err))
		}
		cookie, err := GenerateJWTfromIDandKey(userID, signingKey)
		if err != nil {
			WipeAuth(c)
			c.Logger().Print(err)
			echo.NewHTTPError(http.StatusTeapot, err)
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPage(loginPage.LoginType, err))
		}
		c.SetCookie(cookie)
		lockouts[c.RealIP()] = nil
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", foodlib.PagePrefix))
	})

	e.POST(fmt.Sprintf("%s/submitsignup", foodlib.PagePrefix), func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		confirmpassword := c.FormValue("confirmpassword")
		beepboop := c.FormValue("beepboop")
		if beepboop != "" {
			WipeAuth(c)
			err := errors.New("Scraper no scraping!")
			c.Logger().Print(err)
			if lockouts[c.RealIP()] == nil {
				lockouts[c.RealIP()] = &lockoutEntry{IP: c.RealIP()}
			}
			lockouts[c.RealIP()].Num = 10
			lockouts[c.RealIP()].Last = time.Now()
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPage(loginPage.SignupType, err))
		}
		if password != confirmpassword {
			WipeAuth(c)
			err := errors.New("Passwords don't match")
			c.Logger().Print(err)
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPage(loginPage.SignupType, err))
		}
		userID, err := db.CreateUser(username, password)
		if err != nil {
			WipeAuth(c)
			c.Logger().Print(err)
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPage(loginPage.SignupType, err))
		}
		cookie, err := GenerateJWTfromIDandKey(userID, signingKey)
		if err != nil {
			WipeAuth(c)
			c.Logger().Print(err)
			echo.NewHTTPError(http.StatusTeapot, err)
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPage(loginPage.SignupType, err))
		}
		c.SetCookie(cookie)
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", foodlib.PagePrefix))
	})

	// NOTE: Authenticated routes below

	r := e.Group(fmt.Sprintf("%s", foodlib.PagePrefix))

	r.Use(GetJWTmiddlewareWithConfig(signingKey))

	r.Use(middleware.Logger())
	r.Static("/static", "FOOstatic")
	// r.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
	// 	cssmiddleware := templ.NewCSSMiddleware(hndl, views.StaticStyles...)
	// 	cssmiddleware.Path = fmt.Sprintf("%s/styles/templ.css", viewutils.PagePrefix)
	// 	return cssmiddleware
	// }))

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
	e.Logger.Fatal(e.Start(listenOn))
}
