package service

import (
	"fmt"
	"io"
	"log"
	"path"
	"regexp"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/itglobal/backupmonitor/pkg/database"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/storage"
	"github.com/itglobal/backupmonitor/pkg/util"
	"github.com/jinzhu/gorm"
	"github.com/sarulabs/di"
)

// BackupFile is a readable backup file
type BackupFile struct {
	Backup *model.Backup
	File   io.ReadCloser
}

// BackupRepository contains methods to manage project backups
type BackupRepository interface {
	// Create new backup
	Upload(projectID, filename string, source io.Reader) (*model.Backup, error)

	// List project's backups
	List(projectID string) ([]*model.Backup, error)

	// Download project's backup content
	Download(id string) (*BackupFile, error)

	// Delete a backup
	Delete(id, reason string) error
}

const backupRepositoryKey = "BackupRepository"

// GetBackupRepository returns an implementation BackupRepository from DI container
func GetBackupRepository(c di.Container) BackupRepository {
	return c.Get(backupRepositoryKey).(BackupRepository)
}

type backupRepository struct {
	logger            *log.Logger
	provider          database.Provider
	store             storage.Service
	projectRepository ProjectRepository
}

// Create new backup
func (s *backupRepository) Upload(projectID, filename string, source io.Reader) (*model.Backup, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	// Load project
	project, err := s.projectRepository.Get(projectID)
	if err != nil {
		return nil, err
	}

	if !project.IsActive {
		return nil, model.NewError(model.EAccessDenied, "access denied")
	}

	// Create backup
	mBackup := &model.Backup{}
	mBackup.ID = util.GenerateToken()
	mBackup.ProjectID = projectID
	mBackup.FileName = filename
	mBackup.Type = model.BackupTypeLast
	mBackup.Time = time.Now().UTC()

	// Upload backup file
	sourceWrapper := &readWrapper{reader: source}
	fileRef := s.GenerateBackupFileName(project, filename)
	fileRef, err = s.store.Upload(fileRef, sourceWrapper)
	if err != nil {
		return nil, err
	}
	mBackup.StorageFilePath = string(fileRef)
	mBackup.Length = sourceWrapper.length

	// Save backup to DB
	eBackup := &database.Backup{}
	eBackup.CopyFromModel(mBackup)
	err = tx.Create(eBackup).Error
	if err != nil {
		return nil, err
	}
	eBackup.CopyToModel(mBackup)

	// Update statuses of project's backups
	err = s.UpdateBackupStatuses(tx, eBackup.ProjectID)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	s.logger.Printf(
		"new backup \"%s\" (project \"%s\") has been uploaded (%s, see \"%s\")",
		eBackup.ID,
		eBackup.ProjectID,
		humanize.Bytes(uint64(eBackup.Length)),
		eBackup.StorageFilePath)
	return mBackup, nil
}

// Generate file name for a backup file
func (s *backupRepository) GenerateBackupFileName(project *model.Project, filename string) storage.FileRef {
	_, filename = path.Split(filename)
	ext := path.Ext(filename)
	name := filename[0 : len(filename)-len(ext)]
	r := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	name = r.ReplaceAllLiteralString(name, "")
	backupFileName := fmt.Sprintf(
		"%s.%s.%s%s",
		name,
		time.Now().UTC().Format("20060102.150405"),
		util.GenerateShortToken(),
		ext)

	backupFileName = path.Join(project.ID, backupFileName)

	return storage.FileRef(backupFileName)
}

// List project's backups
func (s *backupRepository) List(projectID string) ([]*model.Backup, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Load entities
	eBackups := make([]*database.Backup, 0)
	err = db.Where("project_id = ?", projectID).Order("time desc").Find(&eBackups).Error
	if err != nil {
		return nil, err
	}

	// Emit results
	mBackups := make([]*model.Backup, len(eBackups))
	for i, eBackup := range eBackups {
		mBackups[i] = eBackup.ToModel()
	}

	return mBackups, nil
}

// Download project's backup content
func (s *backupRepository) Download(id string) (*BackupFile, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Load backup
	eBackup := &database.Backup{}
	err = db.Where("id = ?", id).First(eBackup).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewError(model.ENotFound, "backup \"%s\" doesn't exist", id)
		}

		return nil, err
	}

	// Open backup file
	file, err := s.store.Download(storage.FileRef(eBackup.StorageFilePath))
	if err != nil {
		return nil, err
	}

	// Emit result
	result := &BackupFile{
		Backup: eBackup.ToModel(),
		File:   file,
	}

	return result, nil
}

// Delete a backup
func (s *backupRepository) Delete(id, reason string) error {
	db, err := s.provider.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	// Load backup
	eBackup := &database.Backup{}
	err = db.Where("id = ?", id).First(eBackup).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.NewError(model.ENotFound, "backup \"%s\" doesn't exist", id)
		}

		return err
	}

	// Delete backup
	err = tx.Delete(eBackup).Error
	if err != nil {
		return err
	}

	// Update statuses of project's backups
	err = s.UpdateBackupStatuses(tx, eBackup.ProjectID)
	if err != nil {
		return err
	}

	// Delete backup file
	err = s.store.Delete(storage.FileRef(eBackup.FileName))
	if err != nil {
		return err
	}

	tx.Commit()

	s.logger.Printf("backup \"%s\" (project \"%s\") has been deleted %s", eBackup.ID, eBackup.ProjectID, reason)
	return nil
}

// Update statuses of project's backups
func (s *backupRepository) UpdateBackupStatuses(tx *gorm.DB, projectID string) error {
	// Load all backups
	eBackups := make([]*database.Backup, 0)
	err := tx.Where("project_id = ?", projectID).Order("time desc").Find(&eBackups).Error
	if err != nil {
		return err
	}

	// Mark backups as "archive"
	for _, eBackup := range eBackups {
		eBackup.Type = model.BackupTypeArchive

		err = tx.Save(eBackup).Error
		if err != nil {
			return err
		}
	}

	// Load last backup
	eLastBackup := &database.Backup{}
	err = tx.Where("project_id = ?", projectID).Order("time desc").First(&eLastBackup).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// Mark last backup as "last"
	if eLastBackup != nil {
		eLastBackup.Type = model.BackupTypeLast

		err = tx.Save(eLastBackup).Error
		if err != nil {
			return nil
		}
	}

	// Update project's backup status
	err = s.projectRepository.UpdateBackupStatus(tx, projectID)
	if err != nil {
		return nil
	}

	return nil
}

type readWrapper struct {
	reader io.Reader
	length int64
}

func (r *readWrapper) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if err == nil {
		r.length += int64(n)
	}
	return
}
