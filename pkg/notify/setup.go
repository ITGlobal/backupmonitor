package notify

import (
	"log"

	"github.com/itglobal/backupmonitor/pkg/component"
	"github.com/sarulabs/di"
)

const serviceKey = "NotificationService"

// Setup configures package services
func Setup(builder component.Builder) {
	builder.AddService(di.Def{
		Name: serviceKey,
		Build: func(c di.Container) (interface{}, error) {
			logger := log.New(log.Writer(), "[notify] ", log.Flags())

			s := &serviceImpl{
				slack:    createSlackNotifyer(logger),
				telegram: createTelegramNotifyer(logger),
				webhook:  createWebhookNotifyer(logger),
			}
			return s, nil
		},
	})
}

// SlackMessage is a content for Slack notification
type SlackMessage struct {
	To    []string
	Title string
	Text  string
	Emoji string
}

// TelegramMessage is a content for Telegram notification
type TelegramMessage struct {
	To    []string
	Title string
	Text  string
	Emoji string
}

// WebhookMessage is a content for webhook notification
type WebhookMessage struct {
	To          []string
	PayloadJSON interface{}
}

// Service provides methods to send notifications
type Service interface {
	// Send a notification via Slack
	NotifySlack(msg *SlackMessage) error

	// Send a notification via Telegram
	NotifyTelegram(msg *TelegramMessage) error

	// Send a notification via webhook
	NotifyWebhook(msg *WebhookMessage) error
}

// GetService returns an implementation Service from DI container
func GetService(c di.Container) Service {
	return c.Get(serviceKey).(Service)
}

type serviceImpl struct {
	slack    slackNotifyer
	telegram telegramNotifyer
	webhook  webhookNotifyer
}

// Send a notification via Slack
func (s *serviceImpl) NotifySlack(msg *SlackMessage) error {
	err := s.slack.Notify(msg)
	return err
}

// Send a notification via Telegram
func (s *serviceImpl) NotifyTelegram(msg *TelegramMessage) error {
	err := s.telegram.Notify(msg)
	return err
}

// Send a notification via webhook
func (s *serviceImpl) NotifyWebhook(msg *WebhookMessage) error {
	err := s.webhook.Notify(msg)
	return err
}
