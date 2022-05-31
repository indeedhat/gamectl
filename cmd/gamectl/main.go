package main

import (
	"flag"
	"log"

	"github.com/indeedhat/gamectl/internal/cli"
	"github.com/indeedhat/gamectl/internal/cli/server"
	"github.com/indeedhat/gamectl/internal/config"
	"github.com/indeedhat/gamectl/internal/juniper"
	"github.com/indeedhat/gamectl/internal/models"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to setup config: %s", err)
	}

	if err := models.Connect(); err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	if err := models.Migrate(); err != nil {
		log.Fatalf("Failed to migrate schema: %s", err)
	}

	if err := config.ReloadAppConfig(); err != nil {
		log.Fatalf("Failed to load app config: %s", err)
	}

	cli.GenerateRegister(models.DB)
	commandKey := flag.String("cmd", server.ServerKey, "Command to run")
	flag.Usage = juniper.CliUsage(
		"Dungeon Crawler Server",
		"Server env for the dungeon crawler game",
		"server",
		cli.CommandRegister,
	)
	flag.Parse()

	switch *commandKey {
	case "":
		tmp := server.ServerKey
		commandKey = &tmp
		fallthrough

	default:
		command := cli.CommandRegister.Find(*commandKey)
		if command == nil {
			panic("Command not found")
		}

		if err := command.Run(flag.Args()); err != nil {
			log.Fatalf("%s", err)
		}
	}
}
