package main

import (
	"log"

	"github.com/indeedhat/command-center/app"
	"github.com/indeedhat/command-center/app/config"
	"github.com/indeedhat/command-center/app/models"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to setup config: %s", err)
	}

	if err := models.Connect(); err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	if err := models.Migrate(); err != nil {
		log.Fatalf("Failed to migrate schema: %s", err)
	}

	if err := config.ReloadAppConfig(); err != nil {
		log.Fatalf("Failed to load app config: %s", err)
	}

	router := app.BuildRoutes()

	if err := router.Run(); err != nil {
		log.Fatalf("Run failed: %s", err)
	}
}
