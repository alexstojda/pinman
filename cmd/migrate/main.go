package main

import (
	"github.com/rs/zerolog/log"
	"pinman/internal/models"
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

	err = gormDb.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal().Err(err).Msg("migration failed")
	}

	log.Info().Msg("migration completed.")
}
