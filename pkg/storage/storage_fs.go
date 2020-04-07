package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/util"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
)

type filesystemServiceImpl struct {
	logger    *log.Logger
	directory string
}

func createFileSystemService(c di.Container, logger *log.Logger) serviceInternal {
	directory := path.Join(viper.GetString("VAR"), "blob")
	directory = path.Clean(directory)

	s := &filesystemServiceImpl{logger, directory}
	return s
}

// Initialize service
func (s *filesystemServiceImpl) Initialize() error {
	err := os.MkdirAll(s.directory, 0)
	if err != nil {
		s.logger.Printf("unable to create directory \"%s\": %v", s.directory, err)
		return err
	}

	s.logger.Printf("using file system as storage (see \"%s\")", s.directory)
	return nil
}

// Upload new file
func (s *filesystemServiceImpl) Upload(project *model.Project, filename string, source io.Reader) (FileRef, error) {
	// Create directory for project
	directory := path.Join(s.directory, project.ID)
	err := os.MkdirAll(directory, 0)
	if err != nil {
		s.logger.Printf("unable to create directory \"%s\": %v", directory, err)
		return emptyFileRef, err
	}

	// Generate file name
	_, filename = path.Split(filename)
	ext := path.Ext(filename)
	name := filename[0 : len(filename)-len(ext)]
	r := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	name = r.ReplaceAllLiteralString(name, "")
	filename = fmt.Sprintf(
		"%s.%s.%s%s",
		name,
		time.Now().UTC().Format("20060102.150405"),
		util.GenerateShortToken(),
		ext)
	fullFileName := path.Join(directory, filename)

	var n int64
	{
		file, err := os.Create(fullFileName)
		if err != nil {
			s.logger.Printf("unable to create file \"%s\": %v", fullFileName, err)
			return emptyFileRef, err
		}

		defer file.Close()

		n, err = io.Copy(file, source)
		if err != nil {
			s.logger.Printf("unable to write file \"%s\": %v", fullFileName, err)
			return emptyFileRef, err
		}
	}

	// Return result
	relativeFileName := path.Join(project.ID, filename)
	s.logger.Printf("new file has been written: \"%s\" (%d bytes)", fullFileName, n)

	return FileRef(relativeFileName), nil
}

// Download existing file
func (s *filesystemServiceImpl) Download(file FileRef) (io.ReadCloser, error) {
	fullFileName := path.Join(s.directory, string(file))
	f, err := os.Open(fullFileName)
	if err != nil {
		s.logger.Printf("unable to open file \"%s\": %v", fullFileName, err)
		return nil, err
	}

	return f, nil
}

// List existing files
func (s *filesystemServiceImpl) List() ([]FileRef, error) {
	items := make([]FileRef, 0)

	entries, err := ioutil.ReadDir(s.directory)
	if err != nil {
		if os.IsNotExist(err) {
			return items, nil
		}

		s.logger.Printf("unable to read directory \"%s\": %v", s.directory, err)
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			items, err = s.listDirectory(entry.Name(), items)
			if err != nil {
				return nil, err
			}
		}
	}

	return items, nil
}

func (s *filesystemServiceImpl) listDirectory(directory string, items []FileRef) ([]FileRef, error) {
	fillDirectoryName := path.Join(s.directory, directory)
	entries, err := ioutil.ReadDir(fillDirectoryName)
	if err != nil {
		if os.IsNotExist(err) {
			return items, nil
		}

		s.logger.Printf("unable to read directory \"%s\": %v", fillDirectoryName, err)
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := path.Join(directory, entry.Name())
		items = append(items, FileRef(filename))
	}

	return items, nil
}

// Delete existing file
func (s *filesystemServiceImpl) Delete(file FileRef) error {
	fullFileName := path.Join(s.directory, string(file))
	err := os.Remove(fullFileName)
	if err != nil {
		if os.IsNotExist(err) {
			s.logger.Printf("won't remove file \"%s\" since it doesn't exist", fullFileName)
			return nil
		}
		s.logger.Printf("unable to remove file \"%s\": %v", fullFileName, err)
		return err
	}

	s.logger.Printf("file \"%s\" has been removed", fullFileName)
	return nil
}
