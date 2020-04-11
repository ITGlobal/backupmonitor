package service

import (
	"log"

	"github.com/itglobal/backupmonitor/pkg/database"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/jinzhu/gorm"
	"github.com/sarulabs/di"
)

// ProjectRepository contains methods to manage projects registry
type ProjectRepository interface {
	// List all registered projects
	List() ([]*model.Project, error)

	// Get a project by its ID
	Get(id string) (*model.Project, error)

	// Create new project
	Create(args *model.ProjectCreateParams) (*model.Project, error)

	// Update an existing project
	Update(id string, args *model.ProjectUpdateParams) (*model.Project, error)

	// Delete an existing project
	Delete(id string) error

	// Update status of project's backups
	UpdateBackupStatus(tx *gorm.DB, projectID string) error
}

const projectRepositoryKey = "ProjectRepository"

// GetProjectRepository returns an implementation of ProjectRepository from DI container
func GetProjectRepository(c di.Container) ProjectRepository {
	return c.Get(projectRepositoryKey).(ProjectRepository)
}

// An implementation of ProjectRepository
type projectRepository struct {
	logger   *log.Logger
	provider database.Provider
}

// List all registered projects
func (s *projectRepository) List() ([]*model.Project, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetchs projects
	var eProjects []*database.Project
	err = db.Order("id asc").Find(&eProjects).Error
	if err != nil {
		return nil, err
	}

	// Fetchs last backups
	var eBackups []*database.Backup
	err = db.Where("type = ?", model.BackupTypeLast).Find(&eBackups).Error
	if err != nil {
		return nil, err
	}

	// Emit results
	mProjects := make([]*model.Project, len(eProjects))
	mProjectsbyID := make(map[string]*model.Project)
	for i, eProject := range eProjects {
		mProject := eProject.ToModel()
		mProjects[i] = mProject
		mProjectsbyID[mProject.ID] = mProject
	}

	for _, eBackup := range eBackups {
		mProject, exists := mProjectsbyID[eBackup.ProjectID]
		if exists {
			mBackup := eBackup.ToModel()
			mProject.LastBackup = mBackup
		}
	}

	return mProjects, nil
}

// Get a project by its ID
func (s *projectRepository) Get(id string) (*model.Project, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetch project
	eProject := &database.Project{}
	err = db.Where("id = ?", id).First(eProject).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewError(model.ENotFound, "project \"%s\" doesn't exist", id)
		}

		return nil, err
	}
	mProject := eProject.ToModel()

	// Fetch last backup
	eBackup := &database.Backup{}
	err = db.Where("project_id = ? and type = ?", mProject.ID, model.BackupTypeLast).First(&eBackup).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err != gorm.ErrRecordNotFound {
		mProject.LastBackup = eBackup.ToModel()
	}

	return mProject, nil
}

// Create new project
func (s *projectRepository) Create(args *model.ProjectCreateParams) (*model.Project, error) {
	err := args.Validate()
	if err != nil {
		return nil, err
	}

	args.Normalize()

	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	eProject := &database.Project{}
	mProject := &model.Project{}

	// Check for conflicts
	err = tx.Where("id = ?", args.ID).First(eProject).Error
	if err != gorm.ErrRecordNotFound {
		if err != nil {
			return nil, err
		}

		return nil, model.NewError(model.EConflict, "project \"%s\" already exists", args.ID)
	}

	// Create new project
	args.ApplyTo(mProject)
	mProject.BackupStatus = model.BackupStatusNone
	eProject.CopyFromModel(mProject)
	err = tx.Create(eProject).Error
	if err != nil {
		return nil, err
	}

	eProject.CopyToModel(mProject)

	tx.Commit()

	s.logger.Printf("new project \"%s\" has been created: %s", mProject.ID, mProject)

	return mProject, nil
}

// Update an existing project
func (s *projectRepository) Update(id string, args *model.ProjectUpdateParams) (*model.Project, error) {
	err := args.Validate()
	if err != nil {
		return nil, err
	}

	args.Normalize()

	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	eProject := &database.Project{}
	mProject := &model.Project{}

	// Fetch project
	err = tx.Where("id = ?", id).First(eProject).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewError(model.ENotFound, "project \"%s\" doesn't exist", id)
		}

		return nil, err
	}
	eProject.CopyToModel(mProject)

	// Update project
	args.ApplyTo(mProject)
	eProject.CopyFromModel(mProject)
	err = tx.Save(eProject).Error
	if err != nil {
		return nil, err
	}

	tx.Commit()

	s.logger.Printf("project \"%s\" has been updated: %s", mProject.ID, mProject)

	// Return result (need this to fetch last backup)
	return s.Get(id)
}

// Delete an existing project
func (s *projectRepository) Delete(id string) error {
	db, err := s.provider.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	eProject := &database.Project{}
	err = db.Where("id = ?", id).First(eProject).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.NewError(model.ENotFound, "project \"%s\" doesn't exist", id)
		}

		return err
	}

	err = tx.Delete(eProject).Error
	if err != nil {
		return err
	}

	tx.Commit()

	s.logger.Printf("project \"%s\" has been deleted", id)
	return nil
}

// Update status of project's backups
func (s *projectRepository) UpdateBackupStatus(tx *gorm.DB, projectID string) error {
	// Load project
	eProject := &database.Project{}
	err := tx.Where("id = ?", projectID).First(eProject).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.NewError(model.ENotFound, "project \"%s\" doesn't exist", projectID)
		}

		return err
	}

	// Load last backup
	eLastBackup := &database.Backup{}
	err = tx.Where("project_id = ?", projectID).Order("time desc").First(&eLastBackup).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var mLastBackup *model.Backup
	if err == gorm.ErrRecordNotFound {
		mLastBackup = nil
	} else {
		mLastBackup = eLastBackup.ToModel()
	}

	// Evaluate project backup status
	mProject := eProject.ToModel()
	status := mProject.CalcBackupStatus(mLastBackup)
	if status == mProject.BackupStatus {
		return nil
	}

	mProject.BackupStatus = status
	mProject.LastNotification = nil
	eProject.CopyFromModel(mProject)

	// Update project
	err = tx.Save(eProject).Error
	if err != nil {
		return err
	}

	s.logger.Printf("backup status of project \"%s\" is now \"%s\"", projectID, status)
	return nil
}
