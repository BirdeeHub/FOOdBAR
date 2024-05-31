package main

import (
	foodlib "FOOdBAR/FOOlib"
	"FOOdBAR/db"
	"FOOdBAR/srvapi"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
)

//go:embed static/* FOOstatic/*
var staticFiles embed.FS

func main() {
	embeddedHTMX, err := isFilePresent(staticFiles, "htmx.min.js")
	if embeddedHTMX && err == nil {
		foodlib.HtmxPath = "/static/htmx.min.js"
	}
	embeddedHyperscript, err := isFilePresent(staticFiles, "_hyperscript.min.js")
	if embeddedHyperscript && err == nil {
		foodlib.HyperscriptPath = "/static/_hyperscript.min.js"
	}

	var signingKeyPath string
	flag.StringVar(&signingKeyPath, "keypath", os.Getenv("FOOdBAR_SIGNING_KEY"), "key file to use for signed cookies (overrides FOOdBAR_SIGNING_KEY env var)")
	var dbpath string
	flag.StringVar(&dbpath, "dbpath", os.Getenv("FOOdBAR_STATE"), "path to database directory (overrides FOOdBAR_STATE env var), defaults to /tmp or windows temp folder if not set")
	var port int
	flag.IntVar(&port, "port", 42069, "port to listen on")
	var ip string
	flag.StringVar(&ip, "ip", "localhost", "IP address to bind to")
	flag.Parse()

	listenOn := fmt.Sprintf("%s:%d", ip, port)
	db.SetDBpath(dbpath)

	var signingKey []byte
	if signingKeyPath != "" {
		signingKey, err = os.ReadFile(signingKeyPath)
	} else {
		err = fmt.Errorf("keypath not set")
	}
	if err != nil {
		signingKey = []byte("secret-passphrase-willitwork")
		fmt.Println("Danger: using default signing key due to reason: " + err.Error() + "\n Use -keypath or FOOdBAR_SIGNING_KEY environment var to set it.")
	}

	srvapi.InitServer(signingKey, listenOn, staticFiles)
}

func isFilePresent(filesystem fs.FS, filename string) (bool, error) {
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
