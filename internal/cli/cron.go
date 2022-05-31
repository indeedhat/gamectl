package cli

import "github.com/indeedhat/gamectl/internal/juniper"

const (
	CronTriggerKey   = "cron:trigger"
	CronTriggerUsage = "Trigger the appropriate command based on the cron definition in the cron.yml file"
)

const configPath = "./configs/cron.yml"

// TriggerCronTasks that are due
//
// All due tasks will be ran concurrently on their own coroutine
func TriggerCronTasks(register *juniper.CliCommandEntries) juniper.CliCommandFunc {
	return func([]string) error {
		schedule, err := juniper.ParseCronSchedule(configPath)
		if err != nil {
			return err
		}

		return juniper.RunCranTasks(schedule, *register)
	}
}
