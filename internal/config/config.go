package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv		    string
	AppPort		    string
	LogLevel	    slog.Level
    DatabaseURL     string
    JWTAccessSecret string
    JWTAccessTTL    time.Duration
    JWTRefreshTTL   time.Duration
}
//This load function purpose is to load all configuration needed by the backend service
func Load() (*Config, error) {
    port := getEnv("APP_PORT", "8080")

	//Port validation, strconv.Atoi convert string to number 
    portNumber, err := strconv.Atoi(port)
    if err != nil || portNumber < 1 || portNumber > 65535 {
        return nil, fmt.Errorf("invalid APP_PORT: %s", port)
    }

    logLevel, err := parseLogLevel(
        getEnv("LOG_LEVEL", "info"),
    )
    if err != nil {
        return nil, err
    }

    jwtTTL, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
    if err != nil {
        return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
    }

    if jwtTTL <= 0 {
        return nil, fmt.Errorf("JWT_ACCESS_TTL must be positive")
    }

    refreshTTL, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TTL"))
    if err != nil {
        return nil, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
    }

    if refreshTTL <= 0 {
        return nil, fmt.Errorf("JWT_REFRESH_TTL must be positive")
    }

    //To get DATABASE_URL from env file
    databaseURL := getEnv("DATABASE_URL", "")

    if databaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }

    cfg := &Config{
        AppEnv:   getEnv("APP_ENV", "development"),
        AppPort:  port,
        LogLevel: logLevel,
        DatabaseURL: databaseURL,
        JWTAccessSecret: os.Getenv("JWT_ACCESS_SECRET"),
        JWTAccessTTL: jwtTTL,
        JWTRefreshTTL: refreshTTL,
    }

    return cfg, nil
}
//This function is to get env, if there's none return default value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

//This function purpose is to parse log level, making sure only debug, info, warn, error is used
func parseLogLevel(value string) (slog.Level, error) {
    switch strings.ToLower(value) {
    case "debug":
        return slog.LevelDebug, nil
    case "info":
        return slog.LevelInfo, nil
    case "warn":
        return slog.LevelWarn, nil
    case "error":
        return slog.LevelError, nil
    default:
        return 0, fmt.Errorf("invalid LOG_LEVEL: %s", value)
    }
}