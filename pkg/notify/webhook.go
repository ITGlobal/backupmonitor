package notify

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type webhookNotifier interface {
	Notify(msg *WebhookMessage) error
}

func createWebhookNotifier(logger *log.Logger) webhookNotifier {
	return &webhookNotifierImpl{
		logger: logger,
	}
}

type webhookNotifierImpl struct {
	logger *log.Logger
}

func (s *webhookNotifierImpl) Notify(msg *WebhookMessage) error {
	json, err := json.Marshal(msg.PayloadJSON)
	if err != nil {
		s.logger.Printf("unable to trigger webhooks: %v", err)
		return err
	}

	contentType := "application/json"

	for _, to := range msg.To {
		body := bytes.NewBuffer(json)

		r, err := http.Post(to, contentType, body)
		if err != nil {
			s.logger.Printf("unable to trigger webhook \"%s\": %v", to, err)
			return err
		}

		s.logger.Printf("triggered webhook: POST %s -> %d", to, r.StatusCode)
	}

	return nil
}
