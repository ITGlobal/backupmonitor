package model

import (
	"time"
)

// BackupType is a type of backup - last or archive
type BackupType string

const (
	// BackupTypeLast means last backup for a project
	BackupTypeLast BackupType = "last"

	// BackupTypeArchive means non-last backup for a project
	BackupTypeArchive BackupType = "archive"
)

// Backup contains information about project's backup
type Backup struct {
	ID              string     `json:"id"`
	FileName        string     `json:"filename"`
	Time            time.Time  `json:"time"`
	Type            BackupType `json:"type"`
	StorageFilePath string     `json:"-"`
	ProjectID       string     `json:"-"`
}

// String convers an object to string
func (p Backup) String() string {
	return toJSON(&p)
}

// Backups is a list of Backup
type Backups []*Backup
