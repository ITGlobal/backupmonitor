package storage

import (
	"io"
	"log"

	"github.com/itglobal/backupmonitor/pkg/component"
	"github.com/sarulabs/di"
)

// FileRef is a reference to a storage file
type FileRef string

const emptyFileRef = FileRef("")

// Service defines methods to read, write and manage storage files
type Service interface {
	// Upload new file
	Upload(filename FileRef, source io.Reader) (FileRef, error)

	// Download existing file
	Download(file FileRef) (io.ReadCloser, error)

	// List existing files
	List() ([]FileRef, error)

	// Delete existing file
	Delete(file FileRef) error
}

type serviceInternal interface {
	Service

	// Initialize service
	Initialize() error
}

const serviceKey = "StorageService"

// GetService returns an implementation Service from DI container
func GetService(c di.Container) Service {
	return c.Get(serviceKey).(Service)
}

// Setup configures package services
func Setup(builder component.Builder) {
	builder.AddService(di.Def{
		Name: serviceKey,
		Build: func(c di.Container) (interface{}, error) {
			logger := log.New(log.Writer(), "[storage] ", log.Flags())

			s := createS3Service(c, logger)

			if s == nil {
				s = createFileSystemService(c, logger)
			}

			err := s.Initialize()
			if err != nil {
				return nil, err
			}

			return s, nil
		},
	})
}
