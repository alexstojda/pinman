package main

import (
	"github.com/rs/zerolog/log"
	"pinman/internal/app"
	"pinman/internal/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load environment variables")
	}

	gormDb, err := utils.ConnectDB(config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to DB")
	}

	server := app.NewServer(config, gormDb)
	err = server.StartServer()
	if err != nil {
		log.Fatal().Err(err).Msg("server could not be started")
	}
}
