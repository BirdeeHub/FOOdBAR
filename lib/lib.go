package lib

import (
	"os"
	"path/filepath"
)

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

func MapSlice[T *any, V *any](f func(T) V, list []T) ([]V) {
    var ret []V
    for _, item := range list {
        ret = append(ret, f(item))
    }
    return ret
}

func FilterSlice[T *any](f func(T) bool, list []T) ([]T) {
    var ret []T
    for _, item := range list {
        if f(item) {
            ret = append(ret, item)
        }
    }
    return ret
}

func MapMap[T *any, V *any, R *any](f func(T, V) R, m map[T]V) map[T]R {
    ret := make(map[T]R)
    for k, v := range m {
        ret[k] = f(k, v)
    }
    return ret
}

func FilterMap[T *any, V *any](f func(T, V) bool, m map[T]V) map[T]V {
    ret := make(map[T]V)
    for k, v := range m {
        if f(k, v) {
            ret[k] = v
        }
    }
    return ret
}

func FilterPointerSlice[T any](f func(*T) bool, list []*T) ([]*T) {
    var ret []*T
    for _, item := range list {
        if f(item) {
            ret = append(ret, item)
        }
    }
    return ret
}

func MapPointerSlice[T any, V any](f func(*T) *V, list []*T) ([]*V) {
    var ret []*V
    for _, item := range list {
        ret = append(ret, f(item))
    }
    return ret
}

func MapPointerMap[T any, V any, R any](f func(*T, *V) *R, m map[*T]*V) map[*T]*R {
    ret := make(map[*T]*R)
    for k, v := range m {
        ret[k] = f(k, v)
    }
    return ret
}

func FilterPointerMap[T any, V any](f func(*T, *V) bool, m map[*T]*V) map[*T]*V {
    ret := make(map[*T]*V)
    for k, v := range m {
        if f(k, v) {
            ret[k] = v
        }
    }
    return ret
}
