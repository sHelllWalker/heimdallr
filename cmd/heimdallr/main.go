package main

import (
	"log"

	"github.com/sHelllWalker/heimdallr/internal/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Listen()

	if err != nil {
		log.Fatal(err)
	}
}
