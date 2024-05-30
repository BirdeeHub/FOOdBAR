package main

import (
	"FOOdBAR/db"
	"FOOdBAR/srvapi"
	"embed"
	"flag"
	"fmt"
	"os"
)

//go:embed static/* FOOstatic/*
var staticFiles embed.FS

func main() {
	var signingKeyPath string
	flag.StringVar(&signingKeyPath, "keypath", "", "key file to use for signed cookies (overrides FOOdBAR_SIGNING_KEY env var)")
	var dbpath string
	flag.StringVar(&dbpath, "dbpath", os.Getenv("FOOdBAR_STATE"), "path to database directory (overrides FOOdBAR_STATE env var)")
	var port int
	flag.IntVar(&port, "port", 42069, "port to listen on")
	var ip string
	flag.StringVar(&ip, "ip", "localhost", "IP address to bind to")
	flag.Parse()

	signingKey, err := os.ReadFile(signingKeyPath)
	if err != nil {
		keystring := os.Getenv("FOOdBAR_SIGNING_KEY")
		if keystring != "" {
			signingKey = []byte(os.Getenv(keystring))
		} else {
			signingKey = []byte("secret-passphrase-willitwork")
		}
	}
	listenOn := fmt.Sprintf("%s:%d", ip, port)
	db.SetDBpath(dbpath)
	srvapi.InitServer(signingKey, listenOn, staticFiles)
}
