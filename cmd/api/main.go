package main

import (
	"log/slog"
	"os"

	"gomerce/internal/config"
	"gomerce/internal/logger"
	"gomerce/internal/router"
)

func main()  {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", slog.String("error", err.Error()))

		os.Exit(1)
	}

	//New() method from internal/logger/logger.go
	appLogger := logger.New(cfg.LogLevel)

	//SetupRouter func from internal/router/router.go
	app := router.SetupRouter(appLogger) //SetupRouter() returns an engine, so we could use engine method i.e Run()

	appLogger.Info("Starting GoMerce server...",
					slog.String("environment", cfg.AppEnv),
					slog.String("port", cfg.AppPort))

	if err := app.Run(":" + cfg.AppPort); err != nil {
		appLogger.Error("Failed to start server", slog.String("error", err.Error()))

		os.Exit(1)
	}
}

