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

		args := []any{
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"response_body_size", c.Writer.Size(),
		}

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			args = append(args, "request_payload_size", c.Request.ContentLength)
		}

		logger.Info("request", args...)
	}
}
