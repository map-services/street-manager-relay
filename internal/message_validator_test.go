package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyMessageSignatureVersion(t *testing.T) {
	tests := []struct {
		version string
		wantErr bool
	}{
		{"1", false},
		{"2", false},
		{"3", true},
		{"", true},
		{"0", true},
	}

	for _, tt := range tests {
		t.Run("version "+tt.version, func(t *testing.T) {
			err := verifyMessageSignatureVersion(tt.version)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMessageToSign(t *testing.T) {
	tests := []struct {
		name     string
		message  *SNSMessage
		expected string
	}{
		{
			name: "Notification",
			message: &SNSMessage{
				Type:      "Notification",
				Message:   "hello",
				MessageId: "123",
				Subject:   "test",
				Timestamp: "2023-01-01T00:00:00Z",
				TopicArn:  "arn:aws:sns:us-east-1:123456789012:MyTopic",
			},
			expected: "Message\nhello\nMessageId\n123\nSubject\ntest\nTimestamp\n2023-01-01T00:00:00Z\nTopicArn\narn:aws:sns:us-east-1:123456789012:MyTopic\nType\nNotification\n",
		},
		{
			name: "SubscriptionConfirmation",
			message: &SNSMessage{
				Type:         "SubscriptionConfirmation",
				Message:      "confirm",
				MessageId:    "456",
				SubscribeURL: "https://sns.url",
				Timestamp:    "2023-01-01T00:00:00Z",
				Token:        "tok",
				TopicArn:     "arn:aws:sns:us-east-1:123456789012:MyTopic",
			},
			expected: "Message\nconfirm\nMessageId\n456\nSubscribeURL\nhttps://sns.url\nTimestamp\n2023-01-01T00:00:00Z\nToken\ntok\nTopicArn\narn:aws:sns:us-east-1:123456789012:MyTopic\nType\nSubscriptionConfirmation\n",
		},
		{
			name: "Unknown Type",
			message: &SNSMessage{
				Type: "Unknown",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getMessageToSign(tt.message)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
