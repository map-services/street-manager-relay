package internal

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// SNSMessage represents the structure of an SNS message
type SNSMessage struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Subject          string `json:"Subject,omitempty"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	SubscribeURL     string `json:"SubscribeURL,omitempty"`
	Token            string `json:"Token,omitempty"`
}

func UnmarshalSNSMessage(data []byte) (SNSMessage, error) {
	var r SNSMessage
	err := json.Unmarshal(data, &r)
	return r, err
}

func IsValidSignature(body *SNSMessage, certManager CertManager) (bool, error) {
	if err := verifyMessageSignatureVersion(body.SignatureVersion); err != nil {
		return false, err
	}

	certificate, err := certManager.Download(body.SigningCertURL)
	if err != nil {
		return false, err
	}

	return validateSignature(body, certificate)
}

func verifyMessageSignatureVersion(version string) error {
	if version != "1" && version != "2" {
		return errors.Newf("unsupported signature version: %s", version)
	}
	return nil
}

func validateSignature(message *SNSMessage, certificate string) (bool, error) {
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return false, errors.New("failed to parse PEM certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse certificate")
	}

	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("certificate does not contain RSA public key")
	}

	messageToSign := getMessageToSign(message)
	if messageToSign == "" {
		return false, errors.New("unable to build message to sign")
	}

	var hashAlgorithm crypto.Hash
	var digest []byte

	if message.SignatureVersion == "2" {
		h := sha256.Sum256([]byte(messageToSign))
		digest = h[:]
		hashAlgorithm = crypto.SHA256
	} else {
		h := sha1.Sum([]byte(messageToSign))
		digest = h[:]
		hashAlgorithm = crypto.SHA1
	}

	signature, err := base64.StdEncoding.DecodeString(message.Signature)
	if err != nil {
		return false, errors.Wrap(err, "failed to decode signature")
	}

	err = rsa.VerifyPKCS1v15(rsaPubKey, hashAlgorithm, digest, signature)
	return err == nil, nil
}

func getMessageToSign(body *SNSMessage) string {
	var keys []string
	switch body.Type {
	case "SubscriptionConfirmation", "UnsubscribeConfirmation":
		keys = []string{"Message", "MessageId", "SubscribeURL", "Timestamp", "Token", "TopicArn", "Type"}
	case "Notification":
		keys = []string{"Message", "MessageId", "Subject", "Timestamp", "TopicArn", "Type"}
	default:
		return ""
	}

	var builder strings.Builder
	for _, key := range keys {
		val := getFieldValue(body, key)
		if val != "" {
			builder.WriteString(key + "\n")
			builder.WriteString(val + "\n")
		}
	}
	return builder.String()
}

func getFieldValue(body *SNSMessage, key string) string {
	switch key {
	case "Message":
		return body.Message
	case "MessageId":
		return body.MessageId
	case "Subject":
		return body.Subject
	case "SubscribeURL":
		return body.SubscribeURL
	case "Timestamp":
		return body.Timestamp
	case "Token":
		return body.Token
	case "TopicArn":
		return body.TopicArn
	case "Type":
		return body.Type
	default:
		return ""
	}
}
