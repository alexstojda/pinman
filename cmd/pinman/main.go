package main

import (
	"github.com/rs/zerolog/log"
	"pinman/internal/utils"

	"pinman/web"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("could not load environment variables")
	}

	gormDb, err := utils.ConnectDB(&config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to DB")
	}

	server := web.NewServer(config.SPAPath, gormDb)
	server.StartServer()
}
