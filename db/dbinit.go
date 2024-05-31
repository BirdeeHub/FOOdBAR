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
		for {
			time.Sleep(time.Second * 10)
			err := CleanSessionBlacklist()
			if err != nil {
				fmt.Print(err)
			}
			time.Sleep(time.Hour - (time.Second * 10))
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 10)
			err := CleanPageDataDB()
			if err != nil {
				fmt.Print(err)
			}
			time.Sleep((time.Minute * 5) - (time.Second * 10))
		}
	}()
}
