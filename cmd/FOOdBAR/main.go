package main

import (
	"FOOdBAR/srvapi"
	"os"
	"runtime"
)
func main() {
	// TODO: get this stuff from arguments eventually
	// (and get a better key which isnt stored here)
	signingKey := []byte("secret-passphrase-willitwork")
	dbpath := os.Getenv("FOOdBAR_STATE")
	if dbpath == "" {
		if runtime.GOOS == "windows" {
			dbpath = os.Getenv("TEMP") 
		} else {
			dbpath = "/tmp"
		}
	}
	listenOn := ":42069"
	srvapi.Init(dbpath, signingKey, listenOn)
}
