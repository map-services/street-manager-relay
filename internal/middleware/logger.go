package middleware

import (
	"log/slog"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger returns a Gin middleware that logs HTTP requests in structured JSON.
func RequestLogger(logger *slog.Logger, excludedPaths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if slices.Contains(excludedPaths, c.Request.URL.Path) {
			return
		}

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"body_size", c.Request.ContentLength,
		)
	}
}
