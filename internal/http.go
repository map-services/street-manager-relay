package internal

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// FetchURL makes a GET request to the provided URL and returns the body as a string.
func FetchURL(url string) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "error fetching URL: %s", url)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("error closing response body", "url", url, "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Newf("error fetching URL: %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // Limit to 1MB to prevent potential DoS
	if err != nil {
		return "", errors.Wrapf(err, "error reading response body from: %s", url)
	}

	return string(body), nil
}
