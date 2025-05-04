package main

import (
	"github.com/bercivarga/website-builder/cmd/server"
)

func main() {
	app, err := server.Start()
	if err != nil {
		panic(err)
	}

	defer app.DB.Close()
}
