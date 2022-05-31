package cli

import (
	"github.com/indeedhat/gamectl/internal/cli/apps"
	"github.com/indeedhat/gamectl/internal/cli/server"
	"github.com/indeedhat/gamectl/internal/juniper"
	"gorm.io/gorm"
)

var CommandRegister juniper.CliCommandEntries

func GenerateRegister(db *gorm.DB) {
	CommandRegister = juniper.CliCommandEntries{
		// server
		{
			Key:   server.ServerKey,
			Usage: server.ServerUsage,
			Run:   server.Serve(db),
		},

		// app things
		{
			Key:   apps.RestartKey,
			Usage: apps.RestartUsage,
			Run:   apps.Restart(db),
		},
		{
			Key:   apps.BackupKey,
			Usage: apps.BackupUsage,
			Run:   apps.Backup(db),
		},

		// cron things
		{
			Key:   CronTriggerKey,
			Usage: CronTriggerUsage,
			Run:   TriggerCronTasks(&CommandRegister),
		},
	}
}
