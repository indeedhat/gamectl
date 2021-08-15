package config

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const ConfigDirectoryPattern = "./config/*.app.yml"

var appCache map[string]App

// AppStatus
type AppStatus struct {
	Running   bool
	UserCount uint8
	StartTime time.Time
}

// App configuration
type App struct {
	Title       string
	Description string
	Icon        string

	Commands struct {
		Status string
		Start  string
		Stop   string
	}

	Files []string
}

// Start the application
func (app App) Start() error {
	_, err := runCommand(app.Commands.Start)

	return err
}

// Stop the application
func (app App) Stop() error {
	_, err := runCommand(app.Commands.Stop)

	return err
}

// Status gives the current status of the application
func (app App) Status() (*AppStatus, error) {
	var status AppStatus
	data, err := runCommand(app.Commands.Start)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// Apps will get apps from cache
//
// populating the cache when necesarry
func Apps() *map[string]App {
	if appCache == nil {
		ReloadAppConfig()
	}

	return &appCache
}

// ReloadAppConfig from the yaml files in the config directory
func ReloadAppConfig() error {
	appCache = make(map[string]App)

	files, err := filepath.Glob(ConfigDirectoryPattern)
	if err != nil {
		return err
	}

	for _, file := range files {
		var app App

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal(data, &app); err != nil {
			return err
		}

		appCache[appKey(file)] = app
	}

	return nil
}

func appKey(file string) string {
	fname := path.Base(file)
	if !strings.HasSuffix(fname, ".app.yml") {
		return ""
	}

	return fname[:len(fname)-8]
}

// runCommand runs a command
// ...
// Who knew
func runCommand(cmd string) ([]byte, error) {
	return []byte{}, nil
}
