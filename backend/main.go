package main

import (
	"github.com/bercivarga/website-builder/cmd/server"
)

func main() {
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
