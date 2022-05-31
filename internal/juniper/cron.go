package juniper

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/1set/cronrange"
	"github.com/go-playground/validator"
	"gopkg.in/yaml.v3"
)

type CronTask struct {
	Command  string   ` yaml:"command" validate:"required"`
	Schedule string   ` yaml:"schedule" validate:"required"`
	Args     []string `yaml:"args"`
}

type cronTasks struct {
	Tasks []CronTask `validate:"dive"`
}

// shouldRun checks if the given time matches the schedule set on the task
func (task CronTask) ShouldRun(now time.Time) bool {
	schedule, err := cronrange.ParseString(task.Schedule)
	if err != nil {
		return false
	}

	return schedule.IsWithin(now)
}

// ParseCronSchedule file by its file path
func ParseCronSchedule(configPath string) ([]CronTask, error) {
	var tasks []CronTask

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	if err := validateCronSchedule(tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func validateCronSchedule(tasks []CronTask) error {
	v := validator.New()
	return v.Struct(cronTasks{Tasks: tasks})
}

var _ error = (*CronErrors)(nil)

type CronErrors map[string]error

// Error implementation keeps CronErrors adheering to the error interface
func (ce CronErrors) Error() string {
	errString := strings.Builder{}

	for key, err := range ce {
		errString.WriteString(fmt.Sprintf("%s: %s", key, err))
	}

	return errString.String()
}

// RunCranTasks that have hit their trigger
func RunCranTasks(tasks []CronTask, register CliCommandEntries) CronErrors {
	var (
		wg     sync.WaitGroup
		errors = make(CronErrors)
		now    = time.Now()
	)

	for _, task := range tasks {
		wg.Add(1)

		go func(task CronTask) {
			defer wg.Done()

			if !task.ShouldRun(now) {
				return
			}

			entry := register.Find(task.Command)
			if entry == nil {
				log.Printf("Command '%s' not found", task.Command)
				return
			}

			if err := entry.Run(task.Args); err != nil {
				errors[task.Command] = err
			}
		}(task)
	}

	wg.Wait()

	if len(errors) > 0 {
		return errors
	}

	return nil
}
