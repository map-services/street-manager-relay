package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "\"msg\":\"HTTP request\"") {
		t.Errorf("expected log message to contain msg field, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "\"method\":\"GET\"") {
		t.Errorf("expected log to contain method GET, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "\"path\":\"/test\"") {
		t.Errorf("expected log to contain path /test, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "\"status\":200") {
		t.Errorf("expected log to contain status 200, got: %s", logOutput)
	}
}

// SetupSlogBuffer returns a buffer and a logger that writes JSON to it.
func setupSlogBuffer() (*bytes.Buffer, *slog.Logger) {
	buf := new(bytes.Buffer)
	handler := slog.NewJSONHandler(buf, nil)
	return buf, slog.New(handler)
}
