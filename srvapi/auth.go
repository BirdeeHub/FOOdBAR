package srvapi

import (
	"FOOdBAR/views/viewutils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
)

func GetJWTmiddlewareWithConfig(signingKey []byte) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:  "user",
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
	})
}

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
		Name:     "user",
		Value:    "",
		Path:     fmt.Sprintf("%s", viewutils.PagePrefix),
		Expires:  expiration,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Response().Writer, &cookie)
}

func GetClaimFromToken(c echo.Context, claim string) interface{} {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims[claim]
}

func GetUserFromToken(c echo.Context) (uuid.UUID, error) {
	switch userID := GetClaimFromToken(c, "sub").(type) {
	case string:
		return uuid.Parse(userID)
	default:
		return uuid.Nil, errors.New("invalid userID")
	}
}

