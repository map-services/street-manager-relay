package internal

import (
	"log/slog"
	"net/url"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/kofalt/go-memoize"
)

type CertManager interface {
	Download(certURL string) (string, error)
}

type CachedCertManager struct {
	cache *memoize.Memoizer
}

func NewCertManager(cache *memoize.Memoizer) CertManager {
	return &CachedCertManager{cache: cache}
}

func (cm *CachedCertManager) Download(certURL string) (string, error) {
	certificate, err, _ := memoize.Call(cm.cache, certURL, func() (string, error) {
		return cm.download(certURL)
	})
	return certificate, errors.Wrapf(err, "unable to download from: %s", certURL)
}

func (cm *CachedCertManager) verifyMessageSignatureURL(certURL string) error {
	parsedURL, err := url.Parse(certURL)
	if err != nil {
		return errors.Wrapf(err, "invalid URL: %s", certURL)
	}

	if parsedURL.Scheme != "https" {
		return errors.New("SigningCertURL was not using HTTPS")
	}

	if !strings.HasPrefix(parsedURL.Host, "sns.") || !strings.HasSuffix(parsedURL.Host, ".amazonaws.com") {
		return errors.New("SigningCertURL host is not a trusted AWS SNS domain")
	}

	return nil
}

func (cm *CachedCertManager) download(certURL string) (string, error) {
	slog.Info("Downloading certificate", "url", certURL)
	if err := cm.verifyMessageSignatureURL(certURL); err != nil {
		return "", errors.Wrap(err, "failed to verify signature URL")
	}

	return FetchURL(certURL)
}
