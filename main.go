package main

import (
	"log"

	"github.com/indeedhat/command-center/app"
	"github.com/indeedhat/command-center/app/config"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to setup config: %s", err)
	}

	router := app.BuildRoutes()

	router.Run()
}
