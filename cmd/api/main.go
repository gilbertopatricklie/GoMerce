package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gomerce/internal/config"
	"gomerce/internal/database"
	"gomerce/internal/logger"
	"gomerce/internal/repository"
	"gomerce/internal/router"
	"gomerce/internal/security"
	"gomerce/internal/service"

	"github.com/joho/godotenv"
)

func main() {
    if err := run(); err != nil {
        slog.Error("Application stopped with error",
            slog.String("error", err.Error()),
        )

        os.Exit(1)
    }
}

func run() error {
    if err := godotenv.Load(); err != nil {
        if !errors.Is(err, os.ErrNotExist) {
            return fmt.Errorf("Load environment variables: %w", err)
        }
    }

    cfg, err := config.Load()
    if err != nil {
        slog.Error("Failed to load configuration", slog.String("error", err.Error()),
        )

        return err
    }

    tokenManager, err := security.NewTokenManager(cfg.JWTAccessSecret,cfg.JWTAccessTTL)
    
    if err != nil {
        return fmt.Errorf("initialize token manager: %w", err)
    }

    //New() method from internal/logger/logger.go
    appLogger := logger.New(cfg.LogLevel)

    dbPool, err := database.NewPostgresPool(context.Background(), cfg.DatabaseURL)

    if err != nil {
        return err
    }

    defer dbPool.Close()

    appLogger.Info("PostgreSQL connection established")

    productRepo := repository.NewProductRepository(dbPool)
    productService := service.NewProductService(productRepo)

    categoryRepo := repository.NewCategoryRepository(dbPool)
    categoryService := service.NewCategoryService(categoryRepo)

    userRepo := repository.NewUserRepository(dbPool)
    refreshRepo := repository.NewRefreshTokenRepository(dbPool)
    userService := service.NewUserService(userRepo, refreshRepo)

    //SetupRouter func from internal/router/router.go
     //SetupRouter() returns an engine, so we could use engine method i.e Run()
    app := router.SetupRouter(appLogger, productService, categoryService, userService, tokenManager, cfg.JWTRefreshTTL)

    // Configure HTTP server.
    server := &http.Server{
        Addr:              ":" + cfg.AppPort,
        Handler:           app,
        ReadHeaderTimeout: 5 * time.Second,
    }

    // Start server in a separate goroutine.
    serverErrors := make(chan error, 1)

    go func() {
        appLogger.Info("Starting GoMerce server...",
        slog.String("environment", cfg.AppEnv),
        slog.String("port", cfg.AppPort),
        )

        serverErrors <- server.ListenAndServe()
    }()

    // Listen for operating system shutdown signals.
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    defer signal.Stop(stop)

    // Wait for either a shutdown signal or server error.
    select {
    case sig := <-stop:
        appLogger.Info("Shutdown signal received",
            slog.String("signal", sig.String()),
        )

    case err := <-serverErrors:
        if !errors.Is(err, http.ErrServerClosed) {
            return err
        }

        return nil
    }

    // Give active requests up to 10 seconds to finish.
    ctx, cancel := context.WithTimeout(
        context.Background(),
        10*time.Second,
    )
    defer cancel()

    appLogger.Info("Shutting down GoMerce server...")

    if err := server.Shutdown(ctx); err != nil {
        appLogger.Error("Graceful shutdown failed",
            slog.String("error", err.Error()),
        )

        // Force shutdown if graceful shutdown times out.
        if closeErr := server.Close(); closeErr != nil {
            return errors.Join(err, closeErr)
        }

        return err
    }

    appLogger.Info("GoMerce server stopped successfully")

    return nil
}
