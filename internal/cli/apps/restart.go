package apps

import (
	"errors"

	"github.com/indeedhat/gamectl/internal/config"
	"github.com/indeedhat/gamectl/internal/juniper"
	"gorm.io/gorm"
)

const (
	RestartKey   = "apps:restart"
	RestartUsage = "Restart an app\nThis will only restart the app if is already running"
)

func Restart(*gorm.DB) juniper.CliCommandFunc {
	return func(args []string) error {
		if len(args) != 1 {
			return errors.New("expected 1 arg  (./server -cmd apps:restart [app_slug])")
		}

		app := config.GetApp(args[0])
		if app == nil {
			return errors.New("app not found")
		}

		if status, err := app.Status(); err != nil || !status.Online {
			return errors.New("app is not running")
		}

		if err := app.Stop(); err != nil {
			return errors.New("failed to stop app")
		}

		if _, err := app.Start(); err != nil {
			return errors.New("failed to restart app")
		}

		return nil
	}
}
