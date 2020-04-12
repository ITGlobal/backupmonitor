package policy

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/itglobal/backupmonitor/pkg/component"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/notify"
	"github.com/itglobal/backupmonitor/pkg/service"
	"github.com/sarulabs/di"
)

const (
	notificationFrequency = 8 * time.Hour
	notificationQuietPeriod = 8 * time.Hour
)

type notificationPolicy struct {
	logger              *log.Logger
	projectRepository   service.ProjectRepository
	notificationService notify.Service
}

func createNotificationPolicy(c di.Container) (component.T, error) {
	logger := log.New(log.Writer(), "[policy] ", log.Flags())

	s := &notificationPolicy{
		logger:              logger,
		projectRepository:   service.GetProjectRepository(c),
		notificationService: notify.GetService(c),
	}
	return s, nil
}

func (s *notificationPolicy) Start(group *sync.WaitGroup, stop chan interface{}) {
	period := time.Minute
	t := time.NewTicker(period)

	group.Add(1)
	go func() {
		for range t.C {
			err := s.Execute()
			if err != nil {
				s.logger.Printf("unable to execute background task: %v", err)
			}
		}
	}()

	go func() {
		for range stop {
		}

		t.Stop()
		group.Done()
	}()
}

func (s *notificationPolicy) Execute() error {
	projects, err := s.projectRepository.List()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, project := range projects {
		if s.ShouldSendNotification(project, now) {
			err = s.SendNotification(project)
			if err != nil {
				return err
			}

			err = s.MarkNotificationAsSent(project, now)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *notificationPolicy) ShouldSendNotification(project *model.Project, now time.Time) bool {
	if !project.IsActive || !project.Notifications.Enabled {
		return false
	}

	if project.BackupStatus != model.BackupStatusOutdated {
		return false
	}

	if project.LastNotification == nil {
		return true
	}

	t := now.Sub(*project.LastNotification)

	if t.Seconds() < notificationFrequency.Seconds() + notificationQuietPeriod.Seconds() {
		return false
	}

	return true
}

func (s *notificationPolicy) SendNotification(project *model.Project) error {
	title := fmt.Sprintf("%s backup warning", project.ID)
	text := fmt.Sprintf(
		"No backups of %s (%s) were taken in a while. Please review and take actions.",
		project.ID,
		project.Name)

	emoji := "warning"

	// Send to slack
	err := s.notificationService.NotifySlack(&notify.SlackMessage{
		To:    project.Notifications.SlackUsers,
		Title: title,
		Text:  text,
		Emoji: emoji,
	})
	if err != nil {
		return err
	}

	// Send to telegram
	err = s.notificationService.NotifyTelegram(&notify.TelegramMessage{
		To:    project.Notifications.TelegramUsers,
		Title: title,
		Text:  text,
		Emoji: emoji,
	})
	if err != nil {
		return err
	}

	// Send to webhook
	payloadJSON := map[string]interface{}{
		"project": project.ID,
	}

	if project.LastBackup != nil {
		payloadJSON["lastBackupTime"] = project.LastBackup.Time
	}

	err = s.notificationService.NotifyWebhook(&notify.WebhookMessage{
		To:          project.Notifications.Webhooks,
		PayloadJSON: payloadJSON,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *notificationPolicy) MarkNotificationAsSent(project *model.Project, now time.Time) error {
	args := &model.ProjectUpdateParams{
		LastNotification: &now,
	}

	_, err := s.projectRepository.Update(project.ID, args)
	if err != nil {
		return err
	}

	return nil
}
