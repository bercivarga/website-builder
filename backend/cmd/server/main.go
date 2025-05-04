package main

import (
	"github.com/bercivarga/website-builder/internal/app"
)

func main() {
	_, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
}
