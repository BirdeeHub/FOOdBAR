package main

import (
	"FOOdBAR/srvapi"
	"os"
)
func main() {
	// TODO: get a much better key from a file
	signingKey := []byte("secret-passphrase-willitwork")
	dbpath := os.Getenv("FOOdBAR_STATE")
	listenOn := ":42069"
	srvapi.Init(dbpath, signingKey, listenOn)
}
