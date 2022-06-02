package apps

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/indeedhat/gamectl/internal/config"
	"github.com/indeedhat/gamectl/internal/juniper"
	"gorm.io/gorm"
)

const (
	StatusKey   = "apps:status"
	StatusUsage = "Get the app status"
)

func Status(*gorm.DB) juniper.CliCommandFunc {
	return func(args []string) error {
		if len(args) != 1 {
			return errors.New("expected 1 arg  (./gamectl -cmd apps:status [app_slug])")
		}

		app := config.GetApp(args[0])
		if app == nil {
			return errors.New("app not found")
		}

		status, err := app.Status()
		if err != nil {
			return errors.New("failed to get app status")
		}

		data, _ := json.MarshalIndent(status, "", "    ")
		fmt.Println(string(data))

		return nil
	}
}
