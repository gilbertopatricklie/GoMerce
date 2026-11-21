package main

import (
	"gomerce/internal/logger"
	"gomerce/internal/router"
	"log/slog"
)

func main()  {
	//New() method from internal/logger/logger.go
	appLogger := logger.New()

	//SetupRouter func from internal/router/router.go
	app := router.SetupRouter(appLogger) //SetupRouter() returns an engine, so we could use engine method i.e Run()

	appLogger.Info("Starting GoMerce server...")

	if err := app.Run(":8080"); err != nil {
		appLogger.Error("Failed to start server", slog.String("error", err.Error()))
	}
}

