package server

import (
	"errors"
	"fmt"
	"os"

	"github.com/indeedhat/gamectl/internal/juniper"
	"github.com/indeedhat/gamectl/internal/performance"
	"github.com/indeedhat/gamectl/internal/router"
	"gorm.io/gorm"
)

const (
	ServerKey   = "serve"
	ServerUsage = "Start the server"
)

func Serve(*gorm.DB) juniper.CliCommandFunc {
	return func([]string) error {
		if monitor := performance.GetMonitor(); monitor == nil {
			return errors.New("Failed to initialize performance monitor")
		}

		router := router.BuildRoutes()

		if err := router.Run(os.Getenv("GIN_PORT")); err != nil {
			return fmt.Errorf("Run failed: %s", err)
		}

		return nil
	}
}
