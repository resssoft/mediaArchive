package main

import (
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/controller"
	"github.com/resssoft/mediaArchive/database"
	"github.com/rs/zerolog/log"
)

func main() {
	mongoDbApp, err := database.ProvideMongo()
	if err != nil {
		return
	}
	if err := controller.Routing(mongoDbApp, config.ApiUrl()); err != nil {
		log.Fatal().
			Err(err).
			Str("address", config.ApiUrl()).
			Msg("cannot start server")
	}
}
