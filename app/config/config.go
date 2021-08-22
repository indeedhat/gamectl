package config

import "github.com/joho/godotenv"

const ConfigDirectory = "./config"

// Init will attempt to initialise all of the config gubbins
func Init() error {
	return godotenv.Load()
}
