package middleware

import (
	"log/slog"
	"time"
	"github.com/gin-gonic/gin"
)

//This is a custom/self-made middleware logger
func AccessLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID := c.GetString("request_id")
		path := c.FullPath()

		if path == "" {
			path = c.Request.URL.Path
		}

		attrs := []any{
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("client_ip", c.ClientIP()),
		}

		switch {
		case status >= 500:
			logger.Error("HTTP Request completed", attrs...)
		case status >= 400:
			logger.Warn("HTTP Request completed", attrs...)
		default:
			logger.Info("HTTP Request completed", attrs...)
		}
	}
}