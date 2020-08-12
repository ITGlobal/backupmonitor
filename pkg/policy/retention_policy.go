package policy

import (
	"log"
	"sync"
	"time"

	"github.com/itglobal/backupmonitor/pkg/component"
	"github.com/itglobal/backupmonitor/pkg/database"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/service"
	"github.com/sarulabs/di"
)

type retentionPolicy struct {
	logger            *log.Logger
	dbProvider        database.Provider
	projectRepository service.ProjectRepository
	backupRepository  service.BackupRepository
}

func createRetentionPolicy(c di.Container) (component.T, error) {
	logger := log.New(log.Writer(), "[policy] ", log.Flags())

	s := &retentionPolicy{
		logger:            logger,
		dbProvider:        database.GetProvider(c),
		projectRepository: service.GetProjectRepository(c),
		backupRepository:  service.GetBackupRepository(c),
	}
	return s, nil
}

func (s *retentionPolicy) Start(group *sync.WaitGroup, stop chan interface{}) {
	period := time.Minute
	t := time.NewTicker(period)

	group.Add(1)
	go func() {
		for range t.C {
			err := s.Execute()
			if err != nil {
				s.logger.Printf("unable to execute background task: %v", err)
			}
		}
	}()

	go func() {
		for range stop {
		}

		t.Stop()
		group.Done()
	}()
}

func (s *retentionPolicy) Execute() error {
	projects, err := s.projectRepository.List()
	if err != nil {
		return err
	}

	for _, project := range projects {
		err = s.Apply(project)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *retentionPolicy) Apply(project *model.Project) error {
	err := s.DropOldBackups(project)
	if err != nil {
		return err
	}

	err = s.UpdateBackupStatus(project)
	if err != nil {
		return err
	}

	return nil
}

func (s *retentionPolicy) DropOldBackups(project *model.Project) error {
	backups, err := s.backupRepository.List(project.ID)
	if err != nil {
		return err
	}

	// backups are ordered by time desc
	// remove all but first N backups

	if len(backups) <= project.BackupRetention {
		return nil
	}

	backups = backups[project.BackupRetention:]
	for _, backup := range backups {
		err = s.backupRepository.Delete(backup.ID, "by retention policy")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *retentionPolicy) UpdateBackupStatus(project *model.Project) error {
	db, err := s.dbProvider.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	err = s.projectRepository.UpdateBackupStatus(tx, project.ID)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
