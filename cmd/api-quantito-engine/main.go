package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	conf := getConfig()
	setupLogger(conf.Log)

	log.Info().Msgf("Starting quantito-engine service on %s", conf.API.Listen)

	router := initRouter(conf)
	err := router.Start(conf.API.Listen)
	log.Fatal().Err(err)
}
