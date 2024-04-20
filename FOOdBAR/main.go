package main

import (
	"FOOdBAR/srvapi"
	"os"
)
func main() {
	dbpath := os.Getenv("FOOdBAR_STATE")
	srvapi.Init(dbpath)
}
