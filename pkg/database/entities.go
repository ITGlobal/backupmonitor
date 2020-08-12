package database

import (
	"strings"
	"time"

	"github.com/itglobal/backupmonitor/pkg/model"
)

// User contains information about application user
type User struct {
	ID           int    `gorm:"column:id;auto_increment;primary_key"`
	UserName     string `gorm:"column:username;unique_index"`
	PasswordHash string `gorm:"column:password"`
}

// TableName returns database table name
func (User) TableName() string {
	return "users"
}

// CopyToModel copies entity data to model
func (p *User) CopyToModel(m *model.User) {
	m.ID = p.ID
	m.UserName = p.UserName
	m.PasswordHash = p.PasswordHash
}

// CopyFromModel copies model data to entity
func (p *User) CopyFromModel(m *model.User) {
	p.ID = m.ID
	p.UserName = m.UserName
	p.PasswordHash = m.PasswordHash
}

// Project contains information about project
type Project struct {
	ID                  string             `gorm:"column:id;type:varchar(64);primary_key"`
	Name                string             `gorm:"column:name;type:varchar(256)"`
	BackupRetention     int                `gorm:"column:backup_retention"`
	BackupFrequency     int                `gorm:"column:backup_frequency"`
	IsActive            bool               `gorm:"column:is_active"`
	EnableNotifications bool               `gorm:"column:enable_notifications"`
	LastNotification    *time.Time         `gorm:"column:last_notified"`
	SlackUsers          string             `gorm:"column:notify_slack;type:varchar(256)"`
	TelegramUsers       string             `gorm:"column:notify_telegram;type:varchar(256)"`
	Webhooks            string             `gorm:"column:notify_webhook;type:varchar(1024)"`
	BackupStatus        model.BackupStatus `gorm:"column:backup_status"`
	Backups             []*Backup          `gorm:"foreignkey:project_id"`
	AccessKeys          []*AccessKey       `gorm:"foreignkey:project_id"`
}

// TableName returns database table name
func (Project) TableName() string {
	return "projects"
}

// ToModel creates new model and copies entity data to it
func (p *Project) ToModel() *model.Project {
	m := &model.Project{}
	p.CopyToModel(m)
	return m
}

// CopyToModel copies entity data to model
func (p *Project) CopyToModel(m *model.Project) {
	m.ID = p.ID
	m.Name = p.Name
	m.BackupRetention = p.BackupRetention
	m.BackupFrequency = p.BackupFrequency
	m.IsActive = p.IsActive
	m.BackupStatus = p.BackupStatus
	m.LastNotification = p.LastNotification

	if m.Notifications == nil {
		m.Notifications = &model.NotificationParams{}
	}

	m.Notifications.Enabled = p.EnableNotifications
	m.Notifications.SlackUsers = commaSeparatedToStringArray(p.SlackUsers)
	m.Notifications.TelegramUsers = commaSeparatedToStringArray(p.TelegramUsers)
	m.Notifications.Webhooks = commaSeparatedToStringArray(p.Webhooks)
}

// CopyFromModel copies model data to entity
func (p *Project) CopyFromModel(m *model.Project) {
	p.ID = m.ID
	p.Name = m.Name
	p.BackupRetention = m.BackupRetention
	p.BackupFrequency = m.BackupFrequency
	p.IsActive = m.IsActive
	p.BackupStatus = m.BackupStatus
	p.LastNotification = m.LastNotification

	if m.Notifications != nil {
		p.EnableNotifications = m.Notifications.Enabled
		p.SlackUsers = stringArrayToCommaSeparated(m.Notifications.SlackUsers)
		p.TelegramUsers = stringArrayToCommaSeparated(m.Notifications.TelegramUsers)
		p.Webhooks = stringArrayToCommaSeparated(m.Notifications.Webhooks)
	} else {
		p.EnableNotifications = false
		p.SlackUsers = ""
		p.TelegramUsers = ""
		p.Webhooks = ""
	}

}

func stringArrayToCommaSeparated(array []string) string {
	str := strings.Join(array, ";")
	return str
}

func commaSeparatedToStringArray(str string) []string {
	if str == "" {
		return make([]string, 0)
	}

	array := strings.Split(str, ";")
	return array
}

// Backup contains information about project's backup
type Backup struct {
	ID              string           `gorm:"column:id;type:varchar(128);primary_key"`
	ProjectID       string           `gorm:"column:project_id;type:varchar(128);foreignkey"`
	FileName        string           `gorm:"column:filename;type:varchar(256)"`
	StorageFilePath string           `gorm:"column:storage_path;type:varchar(256);unique_index"`
	Time            time.Time        `gorm:"column:time"`
	Type            model.BackupType `gorm:"column:type"`
	Length          int64            `gorm:"column:length";default:-1`
}

// TableName returns database table name
func (Backup) TableName() string {
	return "backups"
}

// ToModel creates new model and copies entity data to it
func (p *Backup) ToModel() *model.Backup {
	m := &model.Backup{}
	p.CopyToModel(m)
	return m
}

// CopyToModel copies entity data to model
func (p *Backup) CopyToModel(m *model.Backup) {
	m.ID = p.ID
	m.ProjectID = p.ProjectID
	m.FileName = p.FileName
	m.StorageFilePath = p.StorageFilePath
	m.Time = p.Time
	m.Type = p.Type
	m.Length = p.Length
}

// CopyFromModel copies model data to entity
func (p *Backup) CopyFromModel(m *model.Backup) {
	p.ID = m.ID
	p.ProjectID = m.ProjectID
	p.FileName = m.FileName
	p.StorageFilePath = m.StorageFilePath
	p.Time = m.Time
	p.Type = m.Type
	p.Length = m.Length
}

// AccessKey contains information about project's access key
type AccessKey struct {
	ID        int    `gorm:"column:id;auto_increment;primary_key"`
	Label     string `gorm:"column:label;type:varchar(256)"`
	ProjectID string `gorm:"column:project_id;type:varchar(128);foreignkey"`
	Key       string `gorm:"column:key;type:varchar(1024);unique_index"`
}

// TableName returns database table name
func (AccessKey) TableName() string {
	return "access_keys"
}

// CopyToModel copies entity data to model
func (p *AccessKey) CopyToModel(m *model.AccessKey) {
	m.ID = p.ID
	m.Label = p.Label
	m.ProjectID = p.ProjectID
	m.Key = p.Key
}

// CopyFromModel copies model data to entity
func (p *AccessKey) CopyFromModel(m *model.AccessKey) {
	p.ID = m.ID
	p.Label = m.Label
	p.ProjectID = m.ProjectID
	p.Key = m.Key
}
