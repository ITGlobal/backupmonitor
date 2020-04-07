package service

import (
	"log"

	"github.com/itglobal/backupmonitor/pkg/database"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/jinzhu/gorm"
	"github.com/m1/go-generate-password/generator"
	"github.com/sarulabs/di"
)

// UserRepository contains methods to manage user accounts
type UserRepository interface {
	// Get a user by its username
	Get(username string) (*model.User, error)

	// Set user's password
	SetPassword(id int, password string) error
}

const userRepositoryKey = "UserRepository"

// GetUserRepository returns an implementation of UserRepository from DI container
func GetUserRepository(c di.Container) UserRepository {
	return c.Get(userRepositoryKey).(UserRepository)
}

// An implementation of UserRepository
type userRepository struct {
	logger   *log.Logger
	provider database.Provider
}

func (s *userRepository) Initialize() error {
	db, err := s.provider.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	count := 0
	err = tx.Model(&database.User{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	// Generate new admin user
	g, err := generator.New(&generator.Config{
		Length:                     16,
		IncludeSymbols:             false,
		IncludeNumbers:             true,
		IncludeLowercaseLetters:    true,
		IncludeUppercaseLetters:    true,
		ExcludeSimilarCharacters:   false,
		ExcludeAmbiguousCharacters: false,
	})
	if err != nil {
		return err
	}

	password, err := g.Generate()
	if err != nil {
		return err
	}

	mUser := &model.User{
		UserName: "admin",
	}
	mUser.SetPassword(*password)

	eUser := &database.User{}
	eUser.CopyFromModel(mUser)

	tx.Create(eUser)
	tx.Commit()

	s.logger.Printf("new admin user #%d has been generated", eUser.ID)
	s.logger.Printf("use the following credentials to log in:")
	s.logger.Printf("  username: \"%s\"", mUser.UserName)
	s.logger.Printf("  password: \"%s\"", *password)

	return nil
}

// Get a user by its username
func (s *userRepository) Get(username string) (*model.User, error) {
	db, err := s.provider.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	eUser := &database.User{}
	err = db.Where("username = ?", username).First(eUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewError(model.ENotFound, "user \"%s\" doesn't exist", username)
		}

		return nil, err
	}

	mUser := &model.User{}
	eUser.CopyToModel(mUser)

	return mUser, nil
}

// Set user's password
func (s *userRepository) SetPassword(id int, password string) error {
	db, err := s.provider.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	eUser := &database.User{}
	err = db.Where("id = ?", id).First(eUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.NewError(model.ENotFound, "user #%d doesn't exist", id)
		}

		return err
	}

	mUser := &model.User{}
	eUser.CopyToModel(mUser)
	mUser.SetPassword(password)
	eUser.CopyFromModel(mUser)

	err = tx.Save(eUser).Error
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	s.logger.Printf("password for user #%d has been changed", eUser.ID)

	return nil
}
