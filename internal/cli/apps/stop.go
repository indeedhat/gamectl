package apps

import (
	"errors"

	"github.com/indeedhat/gamectl/internal/config"
	"github.com/indeedhat/gamectl/internal/juniper"
	"gorm.io/gorm"
)

const (
	StopKey   = "apps:stop"
	StopUsage = "Stop an app"
)

func Stop(*gorm.DB) juniper.CliCommandFunc {
	return func(args []string) error {
		if len(args) != 1 {
			return errors.New("expected 1 arg  (./gamectl -cmd apps:stop [app_slug])")
		}

		app := config.GetApp(args[0])
		if app == nil {
			return errors.New("app not found")
		}

		if status, err := app.Status(); err != nil {
			return errors.New("failed to get app status")
		} else if !status.Online {
			return errors.New("app is not running")
		}

		if err := app.Stop(); err != nil {
			return errors.New("failed to stop app")
		}

		return nil
	}
}
