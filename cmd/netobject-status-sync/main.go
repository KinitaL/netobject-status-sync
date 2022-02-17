package main

import (
	"log"
	"netobject-status-sync/internal/application"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := application.NewApp()

	if err := app.Configure(); err != nil {
		return err
	}

	app.Run()

	return nil
}
