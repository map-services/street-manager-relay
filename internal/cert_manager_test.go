package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyMessageSignatureURL(t *testing.T) {
	cm := &CachedCertManager{}

	tests := []struct {
		name    string
		url     string
		wantErr bool
		msg     string
	}{
		{
			name:    "valid AWS URL",
			url:     "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-0123456789.pem",
			wantErr: false,
		},
		{
			name:    "valid another AWS URL",
			url:     "https://sns.eu-west-1.amazonaws.com/cert.pem",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			url:     "http://sns.us-east-1.amazonaws.com/cert.pem",
			wantErr: true,
			msg:     "SigningCertURL was not using HTTPS",
		},
		{
			name:    "invalid domain",
			url:     "https://malicious-site.com/cert.pem",
			wantErr: true,
			msg:     "SigningCertURL host is not a trusted amazonaws.com domain",
		},
		{
			name:    "subdomain of amazonaws.com but not AWS",
			url:     "https://fake.amazonaws.com.malicious.com/cert.pem",
			wantErr: true,
			msg:     "SigningCertURL host is not a trusted amazonaws.com domain",
		},
		{
			name:    "invalid URL",
			url:     "://invalid",
			wantErr: true,
			msg:     "invalid URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cm.verifyMessageSignatureURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.msg != "" {
					assert.Contains(t, err.Error(), tt.msg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
