package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchURL(t *testing.T) {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantErr  bool
		expected string
		errMsg   string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprint(w, "success body")
			},
			wantErr:  false,
			expected: "success body",
		},
		{
			name: "not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: true,
			errMsg:  "status code: 404",
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
			errMsg:  "status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			actual, err := FetchURL(server.URL)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}

	t.Run("unreachable", func(t *testing.T) {
		_, err := FetchURL("http://localhost:12345/unreachable")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error fetching URL")
	})

	t.Run("timeout", func(t *testing.T) {
		originalTimeout := httpClient.Timeout
		httpClient.Timeout = 100 * time.Millisecond
		defer func() { httpClient.Timeout = originalTimeout }()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		_, err := FetchURL(server.URL)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Client.Timeout exceeded")
	})
}
