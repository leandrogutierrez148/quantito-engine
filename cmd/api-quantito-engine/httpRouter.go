package main

import (
	"os"

	"quantito-engine/internal/controllers"
	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/repositories"
	"quantito-engine/internal/services"
	"quantito-engine/internal/sinks"
	"quantito-engine/internal/strategies"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
)

func initRouter(conf *Config) *echo.Echo {
	router := echo.New()
	initRoutes(router, conf)
	router.Use(middleware.Gzip())

	globalLogger := log.Logger.With().Str("component", "echo").Logger()

	logger := lecho.From(globalLogger)
	router.Logger = logger

	router.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))

	return router
}

func initRoutes(router *echo.Echo, conf *Config) {
	dbHandler, err := infraestructures.NewMysqlDB(&conf.DB)
	if err != nil {
		log.Panic().Err(err).Msg("Error connecting to DB")
	}

	experimentsRepo := repositories.NewExperimentsRepository(dbHandler)
	executionsRepo := repositories.NewExecutionsRepository(dbHandler)
	signalsRepo := repositories.NewSignalsRepository(dbHandler)

	quantitoUrl := os.Getenv("QUANTITO_API_URL")
	if quantitoUrl == "" {
		quantitoUrl = "http://localhost:8081"
	}
	quantitoApiClient := infraestructures.NewAPIClient(quantitoUrl, nil)

	strategyFactory := strategies.NewStrategyFactory(quantitoApiClient)
	consoleSink := sinks.NewConsoleSink()
	signalDispatcher := sinks.NewSignalDispatcher(consoleSink)
	executionMutex := infraestructures.NewExecutionMutex()

	executionService := services.NewExecutionService(
		experimentsRepo,
		executionsRepo,
		signalsRepo,
		strategyFactory,
		signalDispatcher,
		executionMutex,
	)

	executionController := controllers.NewExecutionController(executionService)

	router.GET("/ping", controllers.NewPingController().Ping)
	router.POST("/execute/:id", executionController.Execute)
}
