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

	listenOn := fmt.Sprintf("%s:%d", ip, port)
	db.SetDBpath(dbpath)

	var signingKey []byte
	var err error
	if signingKeyPath != "" {
		signingKey, err = os.ReadFile(signingKeyPath)
	} else {
		err = fmt.Errorf("keypath not set")
	}
	if err != nil {
		keystring := os.Getenv("FOOdBAR_SIGNING_KEY")
		if keystring != "" {
			signingKey = []byte(os.Getenv(keystring))
		} else {
			fmt.Println("Error: ", err)
			signingKey = []byte("secret-passphrase-willitwork")
			fmt.Println("Danger: using default signing key. Use -keypath or FOOdBAR_SIGNING_KEY environment var to set it.")
		}
	}
	srvapi.InitServer(signingKey, listenOn, staticFiles)
}
