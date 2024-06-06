package srvapi

import (
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/db"
	"FOOdBAR/views/loginPage"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const tokenDuration = time.Hour * 72

var lockouts = make(map[string]*lockoutEntry)

type lockoutEntry struct {
	Num int
	Last time.Time
	IP string
}

type jwtCustomClaims struct {
	IP string   `json:"ip"`
	jwt.RegisteredClaims
}

func SetupLoginRoutes(e *echo.Echo, signingKey []byte) error {
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

	e.POST(fmt.Sprintf("%s/submitlogin", foodlib.PagePrefix), func(c echo.Context) error {
		if lockouts[c.RealIP()] != nil {
			if time.Since(lockouts[c.RealIP()].Last) > 10*time.Minute {
				lockouts[c.RealIP()] = nil
			} else if lockouts[c.RealIP()].Num > 9 {
				lockouts[c.RealIP()].Last = time.Now()
				return HTML(c, http.StatusNotAcceptable, loginPage.LoginPage(loginPage.LoginType, errors.New("Too many failed login attempts, please try again later")))
			}
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
		cookie, err := GenerateJWTfromIDandKey(userID, signingKey, c.RealIP())
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
		cookie, err := GenerateJWTfromIDandKey(userID, signingKey, c.RealIP())
		if err != nil {
			WipeAuth(c)
			c.Logger().Print(err)
			echo.NewHTTPError(http.StatusTeapot, err)
			return HTML(c, http.StatusUnprocessableEntity, loginPage.LoginPage(loginPage.SignupType, err))
		}
		c.SetCookie(cookie)
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s", foodlib.PagePrefix))
	})

	return nil
}

//TODO: Currently this is unused. Unsure that I want to add it to WipeAuth
//But not having it would be a mistake. So, todo, use this when you need it.
func AddSessionToBlacklist(c echo.Context) error {
	claims := foodlib.GetClaimsFromContext(c)
	sess, err := foodlib.GetSessionIDFromClaims(claims)
	if err != nil {
		return err
	}
	expiration, err := foodlib.GetExpirationFromClaims(claims)
	if err != nil {
		return err
	}
	err = db.AddToSessionBlacklist(sess, *expiration)
	return err
}

func getKeyFunc(signingKey *[]byte, signingKeys map[string]interface{}, signingMethod string) (func(*jwt.Token) (interface{}, error)) {
	return func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != signingMethod {
			return nil, &echojwt.TokenError{Token: token, Err: fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])}
		}
		if len(signingKeys) == 0 {
			if signingKey == nil {
				return nil, &echojwt.TokenError{Token: token, Err: fmt.Errorf("no signing keys provided, can't verify jwt token")}
			}
			return *signingKey, nil
		}

		if kid, ok := token.Header["kid"].(string); ok {
			if key, ok := signingKeys[kid]; ok {
				return key, nil
			}
		}
		return nil, &echojwt.TokenError{Token: token, Err: fmt.Errorf("unexpected jwt key id=%v", token.Header["kid"])}
	}
}

func GetJWTmiddlewareWithConfig(signingKey []byte) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:  foodlib.JWTcookiename,
		TokenLookup: fmt.Sprintf("cookie:%s", foodlib.JWTcookiename),
		SuccessHandler: func(c echo.Context) {
			claims := foodlib.GetClaimsFromContext(c)
			userID, err := foodlib.GetUserFromClaims(claims)
			if err != nil {
				WipeAuth(c)
				c.Logger().Print(err)
				return
			}
			expirationtime, err := foodlib.GetExpirationFromClaims(claims)
			if err != nil {
				WipeAuth(c)
				c.Logger().Print(err)
				return
			}
			if time.Until(*expirationtime) < time.Hour {
				cookie, err := GenerateJWTfromIDandKey(userID, signingKey, c.RealIP())
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
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login", foodlib.PagePrefix))
		},
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			token, err := jwt.ParseWithClaims(auth, jwt.MapClaims{}, getKeyFunc(&signingKey, map[string]interface{}{}, "HS256"))
			if err != nil {
				return nil, &echojwt.TokenError{Token: token, Err: err}
			}
			if !token.Valid {
				return nil, &echojwt.TokenError{Token: token, Err: errors.New("invalid token")}
			}
			ip, err := foodlib.GetIPfromClaims(token.Claims)
			if err != nil {
				return nil, &echojwt.TokenError{Token: token, Err: errors.New("invalid ip field")}
			}
			if ip != c.RealIP() {
				return nil, &echojwt.TokenError{Token: token, Err: errors.New("You changed IP! Please log in again.")}
			}
			if _, err := foodlib.GetExpirationFromClaims(token.Claims); err != nil {
				return nil, &echojwt.TokenError{Token: token, Err: errors.New("invalid token")}
			}
			sessionID, err := foodlib.GetSessionIDFromClaims(token.Claims)
			if err != nil {
				return nil, &echojwt.TokenError{Token: token, Err: errors.New("invalid sessionID")}
			}
			status, err := db.QuerySessionBlacklist(sessionID)
			if err != nil || status == true {
				return nil, &echojwt.TokenError{Token: token, Err: errors.New("blacklisted sessionID")}
			}
			return token, nil
		},
	})
}

func GenerateJWTfromIDandKey(userID uuid.UUID, key []byte, ip string) (*http.Cookie, error) {
	claims := jwtCustomClaims{
		ip,
		jwt.RegisteredClaims{
			Subject:   userID.String(),
			ID: uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     foodlib.JWTcookiename,
		Value:    t,
		Path:     fmt.Sprintf("%s", foodlib.PagePrefix),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}, nil
}

func WipeAuth(c echo.Context) {
	// Set the cookie with the same name and an expiration time in the past
	expiration := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:     foodlib.JWTcookiename,
		Value:    "",
		Path:     fmt.Sprintf("%s", foodlib.PagePrefix),
		Expires:  expiration,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}
	http.SetCookie(c.Response().Writer, &cookie)
}
