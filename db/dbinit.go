package db

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var dbpath string

func SetDBpath(newdbpath string) {
	if newdbpath == "" {
		if runtime.GOOS == "windows" {
			dbpath = os.Getenv("TEMP") 
		} else {
			dbpath = "/tmp"
		}
	} else {
		dbpath = newdbpath
	}
}

func init() {
	SetDBpath("")
	go func() {
		for range time.Tick(time.Hour) {
			if err := CleanSessionBlacklist(); err != nil {
				fmt.Print(err)
			}
		}
	}()
	go func() {
		for range time.Tick(time.Minute * 5) {
			if err := CleanPageDataDB(); err != nil {
				fmt.Print(err)
			}
		}
	}()
}
