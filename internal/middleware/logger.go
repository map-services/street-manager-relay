package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger returns a Gin middleware that logs HTTP requests in structured JSON.
func RequestLogger(logger *slog.Logger, ignoredPaths ...string) gin.HandlerFunc {
	ignored := make(map[string]struct{})
	for _, path := range ignoredPaths {
		ignored[path] = struct{}{}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		if _, ok := ignored[path]; ok {
			return
		}

		end := time.Now()
		latency := end.Sub(start)

		logger.Info("HTTP request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
		)
	}
}
