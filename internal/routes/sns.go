package routes

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/map-services/street-manager-relay/generated"
	"github.com/map-services/street-manager-relay/internal"
	"github.com/map-services/street-manager-relay/models"
)

func HandleSNSMessage(repo *internal.DbRepository, certManager internal.CertManager, excludeActivityTypes []string) gin.HandlerFunc {
	excluded := make(map[string]bool)
	for _, t := range excludeActivityTypes {
		excluded[strings.ToLower(t)] = true
	}

	return func(c *gin.Context) {
		messageType := c.GetHeader("x-amz-sns-message-type")
		if messageType == "" {
			abortWithError(c, http.StatusBadRequest, "Missing x-amz-sns-message-type header", errors.New("missing x-amz-sns-message-type header"))
			return
		}

		bodyBytes, err := c.GetRawData()
		if err != nil {
			abortWithError(c, http.StatusBadRequest, "Failed to read request body", errors.Wrap(err, "error reading request body"))
			return
		}

		var body internal.SNSMessage
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			abortWithError(c, http.StatusBadRequest, "Invalid JSON", errors.Wrap(err, "error parsing JSON"))
			return
		}

		valid, err := internal.IsValidSignature(&body, certManager)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, "Signature validation failed", errors.Wrap(err, "signature validation failed"))
			return
		}

		if !valid {
			abortWithError(c, http.StatusUnauthorized, "Message signature is not valid", errors.New("message signature is not valid"))
			return
		}

		if err := handleMessage(repo, &body, excluded); err != nil {
			abortWithError(c, http.StatusInternalServerError, "Failed to handle message", errors.Wrap(err, "failed to handle message"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func handleMessage(repo *internal.DbRepository, body *internal.SNSMessage, excluded map[string]bool) error {
	switch body.Type {
	case "SubscriptionConfirmation":
		return confirmSubscription(body.SubscribeURL)
	case "Notification":
		return handleNotification(repo, body, excluded)
	default:
		slog.Warn("Unknown message type", "type", body.Type)
		return nil
	}
}

// confirmSubscription confirms the SNS subscription by making GET request to subscribe URL
func confirmSubscription(subscriptionURL string) error {
	_, err := internal.FetchURL(subscriptionURL)
	if err != nil {
		return errors.Wrap(err, "failed to confirm subscription")
	}

	slog.Info("Subscription confirmed")
	return nil
}

func handleNotification(repo *internal.DbRepository, body *internal.SNSMessage, excluded map[string]bool) error {
	event, err := generated.UnmarshalEventNotifierMessage([]byte(body.Message))
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal event")
	}

	activityType := ""
	if event.ObjectData.ActivityType != nil {
		activityType = *event.ObjectData.ActivityType
	}

	if excluded[strings.ToLower(activityType)] {
		slog.Info("Skipping excluded activity type", "type", activityType, "objectReference", event.ObjectReference)
		return nil
	}

	_, err = repo.UpsertSingle(models.NewEventFrom(event))
	return err
}
