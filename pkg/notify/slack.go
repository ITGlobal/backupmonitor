package notify

import (
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type slackNotifier interface {
	Notify(msg *SlackMessage) error
}

func createSlackNotifier(logger *log.Logger) slackNotifier {

	apiToken := viper.GetString("SLACK_TOKEN")
	if apiToken != "" {
		api := slack.New(apiToken)
		logger.Printf("connected to slack")
		return &enabledSlackNotifier{
			logger:   logger,
			slack:    api,
			username: viper.GetString("SLACK_USERNAME"),
		}
	} else {
		logger.Printf("slack integration is disabled")
	}

	return &disabledSlackNotifier{logger: logger}
}

type disabledSlackNotifier struct {
	logger *log.Logger
}

func (s *disabledSlackNotifier) Notify(msg *SlackMessage) error {
	s.logger.Printf("unable to deliver slack message { title: \"%s\", text: \"%s\", emoji: \"%s\" } to [ %s ]: slack integration is disabled",
		msg.Title,
		msg.Text,
		msg.Emoji,
		strings.Join(msg.To, ", "))

	return nil
}

type enabledSlackNotifier struct {
	logger   *log.Logger
	slack    *slack.Client
	username string
}

func (s *enabledSlackNotifier) Notify(msg *SlackMessage) error {
	options := make([]slack.MsgOption, 0)
	options = append(options, slack.MsgOptionText(msg.Title, true))

	if msg.Text != "" {
		a := slack.Attachment{
			Text: msg.Text,
		}
		options = append(options, slack.MsgOptionAttachments(a))
	}

	if msg.Emoji != "" {
		options = append(options, slack.MsgOptionIconEmoji(msg.Emoji))
	}

	if s.username != "" {
		options = append(options, slack.MsgOptionUsername(s.username))
	}

	for _, to := range msg.To {
		_, ts, _, err := s.slack.SendMessage(to, options...)
		if err != nil {
			s.logger.Printf("unable to send slack message to \"%s\": %v", to, err)
			return err
		}

		s.logger.Printf("slack message \"%s\" has been sent to \"%s\"", ts, to)
	}

	return nil
}
