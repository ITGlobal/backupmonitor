package model

import (
	"github.com/itglobal/backupmonitor/pkg/notify"
)

// TestSlackNotificationRequest contains parameters to send test Slack notification
type TestSlackNotificationRequest struct {
	Target string `json:"target"`
}

// String converts an object to string
func (p *TestSlackNotificationRequest) String() string {
	return toJSON(&p)
}

// ToMessage converts request values to a SlackMessage
func (p *TestSlackNotificationRequest) ToMessage() *notify.SlackMessage {
	msg := &notify.SlackMessage{
		To:    []string{p.Target},
		Title: "Test Slack notification",
	}
	return msg
}

// TestTelegramNotificationRequest contains parameters to send test Telegram notification
type TestTelegramNotificationRequest struct {
	Target string `json:"target"`
}

// String converts an object to string
func (p *TestTelegramNotificationRequest) String() string {
	return toJSON(&p)
}

// ToMessage converts request values to a TelegramMessage
func (p *TestTelegramNotificationRequest) ToMessage() *notify.TelegramMessage {
	msg := &notify.TelegramMessage{
		To:    []string{p.Target},
		Title: "Test Telegram notification",
	}
	return msg
}

// TestWebhookNotificationPayload is a JSON payload for Webhook notification tests
type TestWebhookNotificationPayload struct {
	Test    bool   `json:"test"`
	Message string `json:"message"`
}

// TestWebhookNotificationRequest contains parameters to send test Webhook notification
type TestWebhookNotificationRequest struct {
	Target string `json:"target"`
}

// String converts an object to string
func (p *TestWebhookNotificationRequest) String() string {
	return toJSON(&p)
}

// ToMessage converts request values to a WebhookMessage
func (p *TestWebhookNotificationRequest) ToMessage() *notify.WebhookMessage {
	msg := &notify.WebhookMessage{
		To: []string{p.Target},
		PayloadJSON: &TestWebhookNotificationPayload{
			Test:    true,
			Message: "Test Webhook notification",
		},
	}
	return msg
}
