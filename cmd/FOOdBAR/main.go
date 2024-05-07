package main

import (
	"FOOdBAR/srvapi"
	"FOOdBAR/db"
	"os"
)
func main() {
	// TODO: get this stuff from arguments eventually
	// (and get a better key which isnt stored here)
	signingKey := []byte("secret-passphrase-willitwork")
	dbpath := os.Getenv("FOOdBAR_STATE")
	listenOn := ":42069"
	db.SetDBpath(dbpath)
	srvapi.InitServer(signingKey, listenOn)
}
