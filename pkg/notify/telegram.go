package notify

import (
	"fmt"
	"log"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hackebrot/turtle"
	"github.com/spf13/viper"
)

type telegramNotifyer interface {
	Notify(msg *TelegramMessage) error
}

func createTelegramNotifyer(logger *log.Logger) telegramNotifyer {

	apiToken := viper.GetString("TELEGRAM_TOKEN")
	if apiToken != "" {
		bot, err := telegram.NewBotAPI(apiToken)
		if err == nil {
			user, err := bot.GetMe()
			if err == nil {
				logger.Printf("connected to telegram as \"%s\"", user.UserName)
				return &enabledTelegramNotifyer{
					logger:   logger,
					telegram: bot,
				}
			}
		}

		logger.Printf("unable to connect to telegram: %v", err)
	} else {
		logger.Printf("telegram integration is disabled")
	}

	return &disabledTelegramNotifyer{logger: logger}
}

type disabledTelegramNotifyer struct {
	logger *log.Logger
}

func (s *disabledTelegramNotifyer) Notify(msg *TelegramMessage) error {
	s.logger.Printf("unable to deliver telegram message { title: \"%s\", text: \"%s\", emoji: \"%s\" } to [ %s ]: telegram integration is disabled",
		msg.Title,
		msg.Text,
		msg.Emoji,
		strings.Join(msg.To, ", "))

	return nil
}

type enabledTelegramNotifyer struct {
	logger   *log.Logger
	telegram *telegram.BotAPI
}

func (s *enabledTelegramNotifyer) Notify(msg *TelegramMessage) error {
	var text string

	text = msg.Title
	if msg.Text !="" {
		text = fmt.Sprintf("%s\n\n%s", msg.Title, text)
	}

	if msg.Emoji != "" {
		emoji, ok := turtle.Emojis[msg.Emoji]
		if ok {
			text = fmt.Sprintf("%s %s", emoji, text)
		}
	}

	for _, to := range msg.To {
		chat, err := s.telegram.GetChat(telegram.ChatConfig{
			SuperGroupUsername: to,
		})
		if err != nil {
			s.logger.Printf("unable to send telegram message to \"%s\": %v", to, err)
			return err
		}

		m := telegram.NewMessage(chat.ID, text)
		r, err := s.telegram.Send(m)
		if err != nil {
			s.logger.Printf("unable to send telegram message to \"%s\": %v", to, err)
			return err
		}

		s.logger.Printf("telegram message %d has been sent to %d", r.MessageID, r.Chat.ID)
	}

	return nil
}
