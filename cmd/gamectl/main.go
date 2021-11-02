package main

import (
	"log"
	"os"

	"github.com/indeedhat/gamectl/internal/config"
	"github.com/indeedhat/gamectl/internal/models"
	"github.com/indeedhat/gamectl/internal/performance"
	"github.com/indeedhat/gamectl/internal/router"
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

	if monitor := performance.GetMonitor(); monitor == nil {
		log.Fatal("Failed to initialize performance monitor")
	}

	router := router.BuildRoutes()

	if err := router.Run(os.Getenv("GIN_PORT")); err != nil {
		log.Fatalf("Run failed: %s", err)
	}
}
