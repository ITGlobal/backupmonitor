package database

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/itglobal/backupmonitor/pkg/component"
	"github.com/jinzhu/gorm"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/sqlite" // SQLite dialect
)

// Provider provides a database context
type Provider interface {
	// Open database connection
	Open() (*gorm.DB, error)
}

// Setup configures package services
func Setup(builder component.Builder) {
	builder.AddService(di.Def{
		Name: providerKey,
		Build: func(c di.Container) (interface{}, error) {
			logger := log.New(log.Writer(), "[db] ", log.Flags())
			filepath := path.Join(viper.GetString("VAR"), "db/sqlite.db")
			filepath = path.Clean(filepath)

			dbLogger := &dbLogger{logger}
			p := &provider{logger, dbLogger, filepath}
			err := p.Initialize()
			if err != nil {
				return nil, err
			}

			return p, nil
		},
	})
}

// GetProvider returns an implementation of Provider from DI container
func GetProvider(c di.Container) Provider {
	return c.Get(providerKey).(Provider)
}

const providerKey = "DbProvider"

// An implementation of Provider
type provider struct {
	logger   *log.Logger
	dbLogger *dbLogger
	filepath string
}

// Initialize database
func (p *provider) Initialize() error {
	dir := path.Dir(p.filepath)
	err := os.MkdirAll(dir, 0)
	if err != nil {
		p.logger.Printf("unable to create directory \"%s\": %v", dir, err)
		return err
	}

	db, err := p.Open()
	if err != nil {
		return err
	}

	defer db.Close()

	err = db.AutoMigrate(&User{}, &Project{}, &Backup{}, &AccessKey{}).Error
	if err != nil {
		p.logger.Printf("unable to migrate database \"%s\": %v", p.filepath, err)
		return err
	}

	p.logger.Printf("using \"%s\"", p.filepath)

	return nil
}

// Open database connection
func (p *provider) Open() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", p.filepath)
	if err != nil {
		p.logger.Printf("unable to open database \"%s\": %v", p.filepath, err)
		return nil, err
	}

	db.SetLogger(p.dbLogger)

	return db, nil
}

type dbLogger struct {
	logger *log.Logger
}

func (s *dbLogger) Print(values ...interface{}) {
	msg := ""
	for _, v := range values {
		msg += fmt.Sprintf("%s ", v)
	}

	s.logger.Printf("%s", msg)
}
