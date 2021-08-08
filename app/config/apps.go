package config

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
