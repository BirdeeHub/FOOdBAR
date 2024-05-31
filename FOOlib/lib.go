package foodlib

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func IsFilePresent(filesystem fs.FS, filename string) (bool, error) {
	found := false
	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == filename && !d.IsDir() {
			found = true
			return fs.SkipAll
			// return fs.SkipDir
		}
		return nil
	}
	err := fs.WalkDir(filesystem, ".", walkFn)
	if err != nil {
		return false, err
	}
	return found, nil
}

func CreateEmptyFileIfNotExists(filename string) (string, error) {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	if err != nil {
		return "", err
	}

	// Try to open the file in read-only mode
	filetry, err := os.Open(filename)
	filetry.Close()
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		defer file.Close()
		if err != nil {
			return "", err
		}
		err = file.Chmod(0600)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func MapSlice[T any, V any](f func(T) V, list []T) []V {
	var ret []V
	for _, item := range list {
		ret = append(ret, f(item))
	}
	return ret
}

func FilterSlice[T any](f func(T) bool, list []T) []T {
	var ret []T
	for _, item := range list {
		if f(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

func MapMap[T comparable, V any, R any](f func(T, V) R, m map[T]V) map[T]R {
	ret := make(map[T]R)
	for k, v := range m {
		ret[k] = f(k, v)
	}
	return ret
}

func FilterMap[T comparable, V any](f func(T, V) bool, m map[T]V) map[T]V {
	ret := make(map[T]V)
	for k, v := range m {
		if f(k, v) {
			ret[k] = v
		}
	}
	return ret
}

func GetClaimsFromContext(c echo.Context) map[string]interface{} {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}

func GetIPfromClaims(claims map[string]interface{}) (string, error) {
	switch ip := claims["ip"].(type) {
	case string:
		return ip, nil
	default:
		return "", errors.New("invalid userID")
	}
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
		return uuid.Nil, errors.New("invalid sessionID")
	}
}

func GetExpirationFromToken(token *jwt.Token) (*time.Time, error) {
	t, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	return &t.Time, nil
}
