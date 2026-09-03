package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	conf, configPath := getConfig()
	setupLogger(conf.Log)

	log.Info().
		Str("config_file", configPath).
		Str("listen", conf.API.Listen).
		Str("db_host", conf.DB.Host).
		Str("db_schema", conf.DB.Schema).
		Str("log_level", conf.Log.Level).
		Str("logger", conf.Log.Logger).
		Msg("Configuration loaded successfully")

	log.Info().Msgf("Starting quantito-engine service on %s", conf.API.Listen)

	router := initRouter(conf)
	err := router.Start(conf.API.Listen)
	log.Fatal().Err(err)
}

