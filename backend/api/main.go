package main

import (
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/controllers"
	"github.com/resssoft/mediaArchive/database"
	"github.com/resssoft/mediaArchive/services/translation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	mongoDbApp, err := database.ProvideMongo()
	if err != nil {
		return
	}
	translatorApp, _ := translation.ProvideTranslator()
	if err := controllers.Routing(mongoDbApp, translatorApp, config.ApiUrl()); err != nil {
		log.Fatal().
			Err(err).
			Str("address", config.ApiUrl()).
			Msg("cannot start server")
	}
}
