package service

import (
	"log"

	"github.com/itglobal/backupmonitor/pkg/database"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/util"
	"github.com/jinzhu/gorm"
	"github.com/sarulabs/di"
)

// AccessKeyRepository contains methods to manage project access keys
type AccessKeyRepository interface {
	// List project's access keys
	List(projectID string) ([]*model.AccessKey, error)

	// Get an access key by its value
	Get(key string) (*model.AccessKey, error)

	// Get an access key by project ID and access key ID
	GetByID(projectID string, accessKeyID int) (*model.AccessKey, error)

	// Create new access key
	Create(projectID string, args *model.AccessKeyCreateParams) (*model.AccessKey, error)

	// Delete an access key
	Delete(projectID string, accessKeyID int) error
}

const accessKeyRepositoryKey = "AccessKeyRepository"

// GetAccessKeyRepository returns an implementation AccessKeyRepository from DI container
func GetAccessKeyRepository(c di.Container) AccessKeyRepository {
	return c.Get(accessKeyRepositoryKey).(AccessKeyRepository)
}

type accessKeyRepository struct {
	logger   *log.Logger
	provider database.Provider
}

// List project's access keys
func (s *accessKeyRepository) List(projectID string) ([]*model.AccessKey, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetchs access keys
	var eAccessKeys []*database.AccessKey
	err = db.Where("project_id = ?", projectID).Order("id asc").Find(&eAccessKeys).Error
	if err != nil {
		return nil, err
	}

	// Emit results
	mAccessKeys := make([]*model.AccessKey, len(eAccessKeys))
	for i, eAccessKey := range eAccessKeys {
		mAccessKey := &model.AccessKey{}
		eAccessKey.CopyToModel(mAccessKey)
		mAccessKeys[i] = mAccessKey
	}

	return mAccessKeys, nil
}

// Get an access key by its value
func (s *accessKeyRepository) Get(key string) (*model.AccessKey, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetchs access key
	eAccessKey := &database.AccessKey{}
	err = db.Where("key = ?", key).First(&eAccessKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewError(model.ENotFound, "access key \"%s\" doesn't exist", key)
		}

		return nil, err
	}

	// Emit result
	mAccessKey := &model.AccessKey{}
	eAccessKey.CopyToModel(mAccessKey)

	return mAccessKey, nil
}

// Get an access key by project ID and access key ID
func (s *accessKeyRepository) GetByID(projectID string, accessKeyID int) (*model.AccessKey, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Fetchs access key
	eAccessKey := &database.AccessKey{}
	err = db.Where("project_id = ? and id = ?", projectID, accessKeyID).First(eAccessKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewError(model.ENotFound, "access key #%d (project \"%s\") doesn't exist", accessKeyID, projectID)
		}

		return nil, err
	}

	// Emit result
	mAccessKey := &model.AccessKey{}
	eAccessKey.CopyToModel(mAccessKey)

	return mAccessKey, nil
}

// Create new access key
func (s *accessKeyRepository) Create(projectID string, args *model.AccessKeyCreateParams) (*model.AccessKey, error) {
	args.Normalize()

	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	// Fetch project
	eProject := &database.Project{}
	err = tx.Where("id = ?", projectID).First(eProject).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewError(model.ENotFound, "project \"%s\" doesn't exist", projectID)
		}

		return nil, err
	}

	// Create access key
	eAccessKey := &database.AccessKey{
		ProjectID: eProject.ID,
		Key:       util.GenerateToken(),
		Label:     args.Label,
	}
	err = tx.Create(eAccessKey).Error
	if err != nil {
		return nil, err
	}

	tx.Commit()
	s.logger.Printf("new access key #%d \"%s\" (project \"%s\") has been created", eAccessKey.ID, eAccessKey.Label, eAccessKey.ProjectID)

	// Emit result
	mAccessKey := &model.AccessKey{}
	eAccessKey.CopyToModel(mAccessKey)

	return mAccessKey, nil
}

// Delete an access key
func (s *accessKeyRepository) Delete(projectID string, accessKeyID int) error {
	db, err := s.provider.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	eAccessKey := &database.AccessKey{}
	err = db.Where("project_id = ? and id = ?", projectID, accessKeyID).First(eAccessKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.NewError(model.ENotFound, "access key #%d (project \"%s\") doesn't exist", accessKeyID, projectID)
		}

		return err
	}

	err = tx.Delete(eAccessKey).Error
	if err != nil {
		return err
	}

	tx.Commit()

	s.logger.Printf("access key #%d \"%s\" (project \"%s\") has been deleted", eAccessKey.ID, eAccessKey.Label, eAccessKey.ProjectID)
	return nil
}
