package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// NotificationParams contains list of targets to send notifications to
type NotificationParams struct {
	Enabled       bool     `json:"enabled"`
	SlackUsers    []string `json:"slack"`
	TelegramUsers []string `json:"telegram"`
	Webhooks      []string `json:"webhook"`
}

// String convers an object to string
func (p *NotificationParams) String() string {
	return toJSON(&p)
}

// ApplyTo applies request values to a NotificationParams
func (p *NotificationParams) ApplyTo(proj *NotificationParams) {
	proj.Enabled = p.Enabled

	if p.SlackUsers != nil {
		proj.SlackUsers = p.SlackUsers
	}

	if p.TelegramUsers != nil {
		proj.SlackUsers = p.TelegramUsers
	}

	if p.Webhooks != nil {
		proj.Webhooks = p.Webhooks
	}
}

// BackupStatus represent backup status for a project
type BackupStatus string

const (
	// BackupStatusOk means than project backup exists and is up to date
	BackupStatusOk BackupStatus = "ok"

	// BackupStatusNone means than project backup doesn't exist
	BackupStatusNone BackupStatus = "none"

	// BackupStatusOutdated means than project backup exists but is out of date
	BackupStatusOutdated BackupStatus = "outdated"
)

// Project contains information about project
type Project struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	IsActive         bool                `json:"isActive"`
	BackupRetention  int                 `json:"backupRetention"`
	BackupFrequency  int                 `json:"backupFrequency"`
	Notifications    *NotificationParams `json:"notifications"`
	BackupStatus     BackupStatus        `json:"backupStatus"`
	LastBackup       *Backup             `json:"lastBackup"`
	LastNotification *time.Time          `json:"-"`
}

const (
	// DefaultRetain is a default value for Project.Retain
	DefaultRetain = 10
	// DefaultPeriod is a default value for Project.Period
	DefaultPeriod = (24 + 8) * 3600 // 1d 8h
)

// String convers an object to string
func (p *Project) String() string {
	return toJSON(&p)
}

// CalcBackupStatus evaluates project's backup status
func (p *Project) CalcBackupStatus(lastBackup *Backup) BackupStatus {
	if lastBackup == nil {
		return BackupStatusNone
	}

	t := time.Now().UTC().Sub(lastBackup.Time)
	if t.Seconds() > float64(p.BackupFrequency) {
		return BackupStatusOutdated
	}

	return BackupStatusOk
}

// Projects is a list of Project
type Projects []*Project

// ProjectCreateParams contains parameters for project creation
type ProjectCreateParams struct {
	ID              string              `json:"id" binding:"required"`
	Name            string              `json:"name" binding:"required"`
	BackupRetention *int                `json:"backupRetention"`
	BackupFrequency *int                `json:"backupFrequency"`
	Enable          *bool               `json:"isActive"`
	Notifications   *NotificationParams `json:"notifications"`
	Webhooks        *[]string           `json:"webhook"`
}

// Normalize normalizes request's fields
func (p *ProjectCreateParams) Normalize() {
	r := regexp.MustCompile(`[^a-zA-Z0-9_-]`)

	p.ID = strings.TrimSpace(p.ID)
	p.ID = strings.ToLower(p.ID)
	p.ID = r.ReplaceAllLiteralString(p.ID, "")

	p.Name = strings.TrimSpace(p.Name)
}

// Validate validates request's fields
func (p *ProjectCreateParams) Validate() error {
	r := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !r.MatchString(p.ID) {
		return NewError(EBadRequest, fmt.Sprintf("\"%s\" is not a valid project ID", p.ID))
	}

	if p.BackupRetention != nil && *p.BackupRetention < 0 {
		return NewError(EBadRequest, fmt.Sprintf("\"%d\" is not a valid backup retention", *p.BackupRetention))
	}

	if p.BackupFrequency != nil && *p.BackupFrequency < 0 {
		return NewError(EBadRequest, fmt.Sprintf("\"%d\" is not a valid backup check period", *p.BackupFrequency))
	}

	return nil
}

// ApplyTo applies request values to a Project
func (p *ProjectCreateParams) ApplyTo(proj *Project) {
	proj.ID = p.ID
	proj.Name = p.Name

	if p.BackupRetention != nil {
		proj.BackupRetention = *p.BackupRetention
	} else {
		proj.BackupRetention = DefaultRetain
	}

	if p.BackupFrequency != nil {
		proj.BackupFrequency = *p.BackupFrequency
	} else {
		proj.BackupFrequency = DefaultPeriod
	}

	if p.Enable != nil {
		proj.IsActive = *p.Enable
	} else {
		proj.IsActive = false
	}

	if p.Notifications != nil {
		if proj.Notifications != nil {
			p.Notifications.ApplyTo(proj.Notifications)
		} else {
			proj.Notifications = p.Notifications
		}
	} else {
		if proj.Notifications == nil {
			proj.Notifications = &NotificationParams{
				Enabled:       false,
				SlackUsers:    make([]string, 0),
				TelegramUsers: make([]string, 0),
				Webhooks:      make([]string, 0),
			}
		}
	}
}

// String convers an object to string
func (p *ProjectCreateParams) String() string {
	return toJSON(&p)
}

// ProjectUpdateParams contains parameters for project modification
type ProjectUpdateParams struct {
	Name             *string             `json:"name"`
	BackupRetention  *int                `json:"backupRetention"`
	BackupFrequency  *int                `json:"backupFrequency"`
	IsActive         *bool               `json:"isActive"`
	Notifications    *NotificationParams `json:"notifications"`
	LastNotification *time.Time          `json:"-"`
}

// Normalize normalizes request's fields
func (p *ProjectUpdateParams) Normalize() {
	if p.Name != nil {
		*p.Name = strings.TrimSpace(*p.Name)
	}
}

// Validate validates request's fields
func (p *ProjectUpdateParams) Validate() error {
	if p.BackupRetention != nil && *p.BackupRetention < 0 {
		return NewError(EBadRequest, fmt.Sprintf("\"%d\" is not a valid backup retention", *p.BackupRetention))
	}

	if p.BackupFrequency != nil && *p.BackupFrequency < 0 {
		return NewError(EBadRequest, fmt.Sprintf("\"%d\" is not a valid backup check period", *p.BackupFrequency))
	}

	return nil
}

// String convers an object to string
func (p *ProjectUpdateParams) String() string {
	return toJSON(&p)
}

// ApplyTo applies request values to a Project
func (p *ProjectUpdateParams) ApplyTo(proj *Project) {
	if p.Name != nil && *p.Name != "" {
		proj.Name = *p.Name
	}

	if p.BackupRetention != nil {
		proj.BackupRetention = *p.BackupRetention
	}

	if p.BackupFrequency != nil {
		proj.BackupFrequency = *p.BackupFrequency
	}

	if p.IsActive != nil {
		proj.IsActive = *p.IsActive
	}

	if p.Notifications != nil {
		if proj.Notifications != nil {
			p.Notifications.ApplyTo(proj.Notifications)
		} else {
			proj.Notifications = p.Notifications
		}
	} else {
		if proj.Notifications == nil {
			proj.Notifications = &NotificationParams{
				Enabled:       false,
				SlackUsers:    make([]string, 0),
				TelegramUsers: make([]string, 0),
				Webhooks:      make([]string, 0),
			}
		}
	}
}
