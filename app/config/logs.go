package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	LogTypeFile    = "file"
	LogTypeCommand = "command"
)

type ServerLog struct {
	Path        string
	Description string
	Type        string
}

// Watch will create a channel with file updates
func (log ServerLog) Watch(done chan bool) (chan string, error) {
	switch log.Type {
	case LogTypeFile:
		return log.watchLogFile(done)

	case LogTypeCommand:
		return log.watchLogCommand(done)

	default:
		return nil, errors.New("Invalid type")
	}
}

func (log ServerLog) watchLogFile(done chan bool) (chan string, error) {
	var err error

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if watcher.Add(log.Path) != nil {
		return nil, err
	}

	output := make(chan string)
	go func() {
		var pos int64

		data, _ := readToEndOfFile(log.Path, &pos)
		output <- data

		for {
			select {
			case <-done:
				close(output)
				return

			case event := <-watcher.Events:
				data, size := readToEndOfFile(log.Path, &pos)
				if size > 0 {
					output <- data
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename {

					output <- "Log file closed"
					close(output)
					return
				}
			}
		}
	}()

	return output, nil
}

// watchLogCommand
//
// This fucking method has taken me longer to get working than anything else on this project
// by a LOOOONG way
//
// i have no idea what i actually did to make it work
func (log ServerLog) watchLogCommand(done chan bool) (chan string, error) {
	cmd := exec.Command(log.Path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	output := make(chan string)
	go func() {
		for range time.NewTicker(100 * time.Millisecond).C {
			var buf bytes.Buffer

			for {
				bucket := make([]byte, 1024)

				size, _ := stdout.Read(bucket)
				buf.Write(bucket)

				if size < 1024 {
					break
				}
			}

			output <- buf.String()
		}
	}()

	time.Sleep(1 * time.Second)

	return output, nil
}

func readToEndOfFile(filePath string, pos *int64) (string, int64) {
	fh, err := os.Open(filePath)
	if err != nil {
		return "", 0
	}
	defer fh.Close()

	_, err = fh.Seek(*pos, 0)
	if err != nil {
		return "", 0
	}

	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return "", 0
	}

	*pos += int64(len(data))

	return string(data), int64(len(data))
}
