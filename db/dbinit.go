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
	go cleanSessBlacklist()
}

func cleanSessBlacklist() {
	for {
		err := CleanSessionBlacklist()
		if err != nil {
			fmt.Print(err)
		}
		time.Sleep(time.Hour)
	}
}
