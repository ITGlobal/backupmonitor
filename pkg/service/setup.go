package service

import (
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/itglobal/backupmonitor/pkg/component"
	"github.com/itglobal/backupmonitor/pkg/database"
	"github.com/itglobal/backupmonitor/pkg/storage"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
)

// Setup configure package services
func Setup(builder component.Builder) {
	// JWT service
	builder.AddService(di.Def{
		Name: jwtKey,
		Build: func(c di.Container) (interface{}, error) {
			password := viper.GetString("JWT_KEY")
			logger := log.New(log.Writer(), "[jwt] ", log.Flags())
			userRepository := GetUserRepository(c)
			return &jwtService{logger, password, jwt.SigningMethodHS256, userRepository}, nil
		},
	})

	// User repository
	builder.AddService(di.Def{
		Name: userRepositoryKey,
		Build: func(c di.Container) (interface{}, error) {
			logger := log.New(log.Writer(), "[users] ", log.Flags())
			provider := database.GetProvider(c)
			repository := &userRepository{logger, provider}
			err := repository.Initialize()
			if err != nil {
				return nil, err
			}
			return repository, nil
		},
	})

	// Project repository
	builder.AddService(di.Def{
		Name: projectRepositoryKey,
		Build: func(c di.Container) (interface{}, error) {
			logger := log.New(log.Writer(), "[projects] ", log.Flags())
			provider := database.GetProvider(c)
			return &projectRepository{logger, provider}, nil
		},
	})

	// Access key repository
	builder.AddService(di.Def{
		Name: accessKeyRepositoryKey,
		Build: func(c di.Container) (interface{}, error) {
			logger := log.New(log.Writer(), "[access] ", log.Flags())
			provider := database.GetProvider(c)
			return &accessKeyRepository{logger, provider}, nil
		},
	})

	// Backup repository
	builder.AddService(di.Def{
		Name: backupRepositoryKey,
		Build: func(c di.Container) (interface{}, error) {
			logger := log.New(log.Writer(), "[backup] ", log.Flags())
			provider := database.GetProvider(c)
			store := storage.GetService(c)
			projectRepository := GetProjectRepository(c)
			return &backupRepository{logger, provider, store, projectRepository}, nil
		},
	})
}
