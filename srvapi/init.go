package srvapi

import (
	"errors"
	"fmt"
	"net/http"

	"FOOdBAR/db"
	"FOOdBAR/views/loginPage"
	"FOOdBAR/views/viewutils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

)

func Init(dbpath string, signingKey []byte, listenOn string) {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// TODO: figure out how to HTTPS
	// e.Pre(middleware.HTTPSRedirect())

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	e.GET(fmt.Sprintf("%s/login", viewutils.PagePrefix), func(c echo.Context) error {
		WipeAuth(c)
		return HTML(c, http.StatusOK, loginPage.LoginPage("login", nil))
	})

	e.GET(fmt.Sprintf("%s/loginform/:formtype", viewutils.PagePrefix), func(c echo.Context) error {
		formtype := c.Param("formtype")
		if formtype == "login" {
			return HTML(c, http.StatusOK, loginPage.LoginPageContents(loginPage.LoginType, nil))
		} else if formtype == "signup" {
			return HTML(c, http.StatusOK, loginPage.LoginPageContents(loginPage.SignupType, nil))
		} else {
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPageContents(loginPage.LoginType, errors.New("Invalid formtype")))
		}
	})

	e.POST(fmt.Sprintf("%s/submitlogin", viewutils.PagePrefix), func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		userID, err := db.AuthUser(username, password, dbpath)
		if err != nil {
			WipeAuth(c)
			c.Logger().Print(err)
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
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	e.POST(fmt.Sprintf("%s/submitsignup", viewutils.PagePrefix), func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		confirmpassword := c.FormValue("confirmpassword")
		if password != confirmpassword {
			WipeAuth(c)
			err := errors.New("Passwords don't match")
			c.Logger().Print(err)
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPage(loginPage.SignupType, err))
		}
		userID, err := db.CreateUser(username, password, dbpath)
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
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", viewutils.PagePrefix))
	})

	// NOTE: Authenticated routes below

	r := e.Group(fmt.Sprintf("%s", viewutils.PagePrefix))

	r.Use(GetJWTmiddlewareWithConfig(signingKey))

	r.Use(middleware.Logger())
	r.Static("/static", "FOOstatic")
	// r.Use(echo.WrapMiddleware(func(hndl http.Handler) http.Handler {
	// 	cssmiddleware := templ.NewCSSMiddleware(hndl, views.StaticStyles...)
	// 	cssmiddleware.Path = fmt.Sprintf("%s/styles/templ.css", viewutils.PagePrefix)
	// 	return cssmiddleware
	// }))

	err := SetupAPIroutes(r, dbpath)
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
