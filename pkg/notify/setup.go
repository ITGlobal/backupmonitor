package notify

import (
	"log"
	"sync"

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
				slack: newAsyncInit(func() interface{} {
					return createSlackNotifier(logger)
				}),

				telegram: newAsyncInit(func() interface{} {
					return createTelegramNotifier(logger)
				}),

				webhook: newAsyncInit(func() interface{} {
					return createWebhookNotifier(logger)
				}),
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
	slack    *asyncInit
	telegram *asyncInit
	webhook  *asyncInit
}

// Send a notification via Slack
func (s *serviceImpl) NotifySlack(msg *SlackMessage) error {
	notifier := s.slack.GetValue().(slackNotifier)
	err := notifier.Notify(msg)
	return err
}

// Send a notification via Telegram
func (s *serviceImpl) NotifyTelegram(msg *TelegramMessage) error {
	notifier := s.telegram.GetValue().(telegramNotifier)
	err := notifier.Notify(msg)
	return err
}

// Send a notification via webhook
func (s *serviceImpl) NotifyWebhook(msg *WebhookMessage) error {
	notifier := s.webhook.GetValue().(webhookNotifier)
	err := notifier.Notify(msg)
	return err
}

type asyncInit struct {
	wait  *sync.WaitGroup
	value interface{}
}

func newAsyncInit(factory func() interface{}) *asyncInit {
	s := &asyncInit{
		wait:  &sync.WaitGroup{},
		value: nil,
	}
	s.wait.Add(1)

	go func() {
		value := factory()

		s.value = value
		s.wait.Done()
	}()

	return s
}

func (s *asyncInit) GetValue() interface{} {
	s.wait.Wait()
	return s.value
}
