package config

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const ConfigDirectoryPattern = "./.config/*.app.yml"

var appCache map[string]App

var (
	ErrBadExtension         = errors.New("Extension must be .zip")
	ErrWorldDirectoryNotSet = errors.New("World Directory not set")
)

// AppStatus
type AppStatus struct {
	Online      bool `json:"online"`
	MaxPlayers  int8 `json:"max_players"`
	PlayerCount int8 `json:"players"`
	Uptime      uint `json:"uptime"`
	Extra       []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"extra"`
}

// App configuration
type App struct {
	Title       string
	Description string
	Icon        string

	WorldDirectory string `yaml:"worldDirectory"`

	Commands struct {
		Status string
		Start  string
		Stop   string
	}

	Files map[string]struct {
		Path        string
		Description string
		Mode        string
	}

	Logs map[string]ServerLog
}

// Start the application
func (app App) Start() (string, error) {
	output, err := runCommand(app.Commands.Start)

	return string(output), err
}

// Stop the application
func (app App) Stop() error {
	_, err := runCommand(app.Commands.Stop)

	return err
}

// Status gives the current status of the application
func (app App) Status() (*AppStatus, error) {
	var status AppStatus
	data, err := runCommand(app.Commands.Status)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// FileKeys will return the keys for any config files defined on the application
func (app App) ConfigFiles() map[string]map[string]string {
	fileList := make(map[string]map[string]string)

	for key, info := range app.Files {
		fileList[key] = map[string]string{
			"description": info.Description,
			"mode":        info.Mode,
		}
	}

	return fileList
}

// LogFiles will return the keys for any log file/script defined on the application
func (app App) LogFiles() map[string]string {
	logList := make(map[string]string)

	for key, info := range app.Logs {
		logList[key] = info.Description
	}

	return logList
}

// BackupWorldDirectory will generate a new zip file in the .temp directory from the sourceDirecotry
//
// on success it will return the relative path to the archive (including the .temp/ prefix
//
// I know that this would run a lot faster if i wrote the archive to memory given that im just going to be
// reading it back out again for the download but as this tool is designed to be run on low end vps's
// memory is likely to be at a premium and thus i prefer this solution
func (app App) BackupWorldDirectory(name string) (string, error) {
	if app.WorldDirectory == "" {
		return "", ErrWorldDirectoryNotSet
	}

	destinationFile := fmt.Sprintf("%s.zip", name)
	extension := path.Ext(destinationFile)

	if extension != ".zip" {
		return "", ErrBadExtension
	}

	archivePath := path.Join(os.Getenv("TEMP_DIR"), destinationFile)
	outputFile, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	writer := zip.NewWriter(outputFile)
	defer writer.Close()

	if err = filepath.Walk(app.WorldDirectory, buildRecursiveFileWalker(app.WorldDirectory, writer)); err != nil {
		return "", err
	}

	return archivePath, nil
}

// buildRecursiveFileWalker creates the file walker specifically for use by the BackupWorldDirectory function
func buildRecursiveFileWalker(sourceDirectory string, writer *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		destFile, err := writer.Create(strings.TrimPrefix(path, sourceDirectory))
		if err != nil {
			return err
		}

		_, err = io.Copy(destFile, sourceFile)

		return err
	}
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

// GetApp will get an app by its key
func GetApp(key string) *App {
	apps := Apps()
	app, ok := (*apps)[key]
	if !ok {
		return nil
	}

	return &app
}

// ReloadAppConfig from the yaml files in the config directory
func ReloadAppConfig() error {
	appConfig := make(map[string]App)

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

		appConfig[appKey(file)] = app
	}

	appCache = appConfig
	log.Print(appCache)
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
func runCommand(cmdString string) ([]byte, error) {
	var err error

	cmd := exec.Command(cmdString)

	if cmd.Dir, err = os.Getwd(); err != nil {
		return []byte{}, err
	}

	return cmd.Output()
}
