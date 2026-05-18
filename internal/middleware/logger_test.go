package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	buf, logger := setupSlogBuffer()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestLogger(logger))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "\"msg\":\"request\"")
	assert.Contains(t, logOutput, "\"method\":\"GET\"")
	assert.Contains(t, logOutput, "\"path\":\"/test\"")
	assert.Contains(t, logOutput, "\"status\":200")
}

// setupSlogBuffer returns a buffer and a logger that writes JSON to it.
func setupSlogBuffer() (*bytes.Buffer, *slog.Logger) {
	buf := new(bytes.Buffer)
	handler := slog.NewJSONHandler(buf, nil)
	return buf, slog.New(handler)
}
