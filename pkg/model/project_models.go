package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

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
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Retain       int          `json:"retain"`
	Enable       bool         `json:"enable"`
	Notify       bool         `json:"notify"`
	Period       int          `json:"period"`
	BackupStatus BackupStatus `json:"backup_status"`
	LastBackup   *Backup      `json:"last_backup"`
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
	if t.Seconds() > float64(p.Period) {
		return BackupStatusOutdated
	}

	return BackupStatusOk
}

// Projects is a list of Project
type Projects []*Project

// ProjectCreateParams contains parameters for project creation
type ProjectCreateParams struct {
	ID     string `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Retain *int   `json:"retain"`
	Period *int   `json:"period"`
	Enable *bool  `json:"enable"`
	Notify *bool  `json:"notify"`
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

	if p.Retain != nil && *p.Retain < 0 {
		return NewError(EBadRequest, fmt.Sprintf("\"%d\" is not a valid backup retention", *p.Retain))
	}

	if p.Period != nil && *p.Period < 0 {
		return NewError(EBadRequest, fmt.Sprintf("\"%d\" is not a valid backup check period", *p.Period))
	}

	return nil
}

// ApplyTo applies request values to a Project
func (p *ProjectCreateParams) ApplyTo(proj *Project) {
	proj.ID = p.ID
	proj.Name = p.Name

	if p.Retain != nil {
		proj.Retain = *p.Retain
	} else {
		proj.Retain = DefaultRetain
	}

	if p.Period != nil {
		proj.Period = *p.Period
	} else {
		proj.Period = DefaultPeriod
	}

	if p.Enable != nil {
		proj.Enable = *p.Enable
	} else {
		proj.Enable = false
	}

	if p.Notify != nil {
		proj.Notify = *p.Notify
	} else {
		proj.Notify = false
	}
}

// String convers an object to string
func (p *ProjectCreateParams) String() string {
	return toJSON(&p)
}

// ProjectUpdateParams contains parameters for project modification
type ProjectUpdateParams struct {
	Name   *string `json:"name"`
	Retain *int    `json:"retain"`
	Period *int    `json:"period"`
	Enable *bool   `json:"enable"`
	Notify *bool   `json:"notify"`
}

// Normalize normalizes request's fields
func (p *ProjectUpdateParams) Normalize() {
	if p.Name != nil {
		*p.Name = strings.TrimSpace(*p.Name)
	}
}

// Validate validates request's fields
func (p *ProjectUpdateParams) Validate() error {
	if p.Retain != nil && *p.Retain < 0 {
		return NewError(EBadRequest, fmt.Sprintf("\"%d\" is not a valid backup retention", *p.Retain))
	}

	if p.Period != nil && *p.Period < 0 {
		return NewError(EBadRequest, fmt.Sprintf("\"%d\" is not a valid backup check period", *p.Period))
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

	if p.Retain != nil {
		proj.Retain = *p.Retain
	}

	if p.Period != nil {
		proj.Period = *p.Period
	}

	if p.Enable != nil {
		proj.Enable = *p.Enable
	}

	if p.Notify != nil {
		proj.Notify = *p.Notify
	}
}
