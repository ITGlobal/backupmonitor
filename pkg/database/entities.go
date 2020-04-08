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
	ID               string             `gorm:"column:id;type:varchar(64);primary_key"`
	Name             string             `gorm:"column:name;type:varchar(256)"`
	Retain           int                `gorm:"column:retain"`
	Period           int                `gorm:"column:period"`
	Enable           bool               `gorm:"column:enable"`
	Notify           bool               `gorm:"column:notify"`
	LastNotification *time.Time         `gorm:"column:last_notified"`
	BackupStatus     model.BackupStatus `gorm:"column:backup_status"`
	SlackUsers       string             `gorm:"column:slack;type:varchar(256)"`
	TelegramUsers    string             `gorm:"column:telegram;type:varchar(256)"`
	Webhooks         string             `gorm:"column:webhook;type:varchar(1024)"`
	Backups          []*Backup          `gorm:"foreignkey:project_id"`
	AccessKeys       []*AccessKey       `gorm:"foreignkey:project_id"`
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
	m.Retain = p.Retain
	m.Period = p.Period
	m.Enable = p.Enable
	m.Notify = p.Notify
	m.BackupStatus = p.BackupStatus
	m.LastNotification = p.LastNotification

	m.SlackUsers = p.commaSeparatedToStringArray(p.SlackUsers)
	m.TelegramUsers = p.commaSeparatedToStringArray(p.TelegramUsers)
	m.Webhooks = p.commaSeparatedToStringArray(p.Webhooks)
}

// CopyFromModel copies model data to entity
func (p *Project) CopyFromModel(m *model.Project) {
	p.ID = m.ID
	p.Name = m.Name
	p.Retain = m.Retain
	p.Period = m.Period
	p.Enable = m.Enable
	p.Notify = m.Notify
	p.BackupStatus = m.BackupStatus
	p.LastNotification = m.LastNotification

	p.SlackUsers = p.stringArrayToCommaSeparated(m.SlackUsers)
	p.TelegramUsers = p.stringArrayToCommaSeparated(m.TelegramUsers)
	p.Webhooks = p.stringArrayToCommaSeparated(m.Webhooks)
}

func (p *Project) stringArrayToCommaSeparated(array []string) string {
	str := strings.Join(array, ";")
	return str
}

func (p *Project) commaSeparatedToStringArray(str string) []string {
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
}

// CopyFromModel copies model data to entity
func (p *Backup) CopyFromModel(m *model.Backup) {
	p.ID = m.ID
	p.ProjectID = m.ProjectID
	p.FileName = m.FileName
	p.StorageFilePath = m.StorageFilePath
	p.Time = m.Time
	p.Type = m.Type
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
