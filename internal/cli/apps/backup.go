package apps

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/indeedhat/gamectl/internal/config"
	"github.com/indeedhat/gamectl/internal/juniper"
	"gorm.io/gorm"
)

const (
	BackupKey   = "apps:backup"
	BackupUsage = "Backup app server\nThis will restart the server if it is running"
)

func Backup(*gorm.DB) juniper.CliCommandFunc {
	return func(args []string) error {
		if len(args) != 2 {
			return errors.New("expected 2 arg  (./server -cmd apps:restart [app_slug] [save_dir])")
		}

		if _, err := os.Stat(args[1]); err != nil {
			return fmt.Errorf("path error: %s", err)
		}

		app := config.GetApp(args[0])
		if app == nil || app.WorldDirectory == "" {
			return errors.New("Cannot find app world to backup")
		}

		status, err := app.Status()
		if err != nil {
			return errors.New("failed to get app status")
		}

		if status.Online {
			if err := app.Stop(); err != nil {
				return errors.New("failed to stop server")
			}

			defer app.Start()
		}

		archivePath, err := app.BackupWorldDirectory(fmt.Sprintf("%s_%s",
			args[0],
			time.Now().Format("2007-01-02T15:04"),
		))
		if err != nil {
			return fmt.Errorf("archive failed: %s", err)
		}
		defer os.Remove(archivePath)

		archiveName := path.Base(archivePath)
		if err := os.Rename(archivePath, path.Join(args[1], archiveName)); err != nil {
			return fmt.Errorf("failed to move backup to destination: %s", err)
		}

		return nil
	}
}
