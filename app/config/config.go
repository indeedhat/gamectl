package config

import "github.com/joho/godotenv"

// Init will attempt to initialise all of the config gubbins
func Init() error {
	return initDotEnv()
}

// initDotEnv will load in the environment variables from the .env file in the project root
func initDotEnv() error {
	return godotenv.Load()
}
