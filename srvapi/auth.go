package srvapi

import (
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/db"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const tokenDuration = time.Hour * 72

func AddSessionToBlacklist(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	sess, err := GetSessionIDFromClaims(claims)
	if err != nil {
		return err
	}
	expiration, err := GetExpirationFromToken(user)
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
		ContextKey:  "user",
		TokenLookup: "cookie:user",
		SuccessHandler: func(c echo.Context) {
			userID, err := GetUserFromClaims(GetClaimsFromContext(c))
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
			if _, err := GetExpirationFromToken(token); err != nil {
				return nil, &echojwt.TokenError{Token: token, Err: errors.New("invalid token")}
			}
			sessionID, err := GetSessionIDFromClaims(token.Claims.(jwt.MapClaims))
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

func GenerateJWTfromIDandKey(userID uuid.UUID, key []byte) (*http.Cookie, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ID: uuid.New().String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     "user",
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
		Name:     "user",
		Value:    "",
		Path:     fmt.Sprintf("%s", foodlib.PagePrefix),
		Expires:  expiration,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}
	http.SetCookie(c.Response().Writer, &cookie)
}

func GetClaimsFromContext(c echo.Context) map[string]interface{} {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}

func GetUserFromClaims(claims map[string]interface{}) (uuid.UUID, error) {
	switch userID := claims["sub"].(type) {
	case string:
		return uuid.Parse(userID)
	default:
		return uuid.Nil, errors.New("invalid userID")
	}
}

func GetSessionIDFromClaims(claims map[string]interface{}) (uuid.UUID, error) {
	switch sessionID := claims["jti"].(type) {
	case string:
		return uuid.Parse(sessionID)
	default:
		return uuid.Nil, errors.New("invalid userID")
	}
}

func GetExpirationFromToken(token *jwt.Token) (*time.Time, error) {
	t, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	return &t.Time, nil
}
