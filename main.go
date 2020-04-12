package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"sync"

	"github.com/itglobal/backupmonitor/pkg/api"
	"github.com/itglobal/backupmonitor/pkg/component"
	"github.com/itglobal/backupmonitor/pkg/database"
	"github.com/itglobal/backupmonitor/pkg/notify"
	"github.com/itglobal/backupmonitor/pkg/policy"
	"github.com/itglobal/backupmonitor/pkg/service"
	"github.com/itglobal/backupmonitor/pkg/storage"
	"github.com/spf13/viper"
)

func configure() error {
	cwd, _ := os.Getwd()
	viper.SetDefault("VAR", path.Join(cwd, "var"))
	viper.SetDefault("JWT_KEY", "test")
	viper.SetDefault("LISTEN_ADDR", "0.0.0.0:8000")

	viper.AutomaticEnv()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	return nil
}

func setup() (component.Registry, error) {
	builder := component.NewBuilder()

	database.Setup(builder)
	storage.Setup(builder)
	service.Setup(builder)
	notify.Setup(builder)
	policy.Setup(builder)
	api.Setup(builder)

	registry, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return registry, nil
}

func main() {
	// Configure app
	err := configure()
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Lmsgprefix)

	// Create services and components
	registry, err := setup()
	if err != nil {
		panic(err)
	}

	// Run (with graceful shutdown)
	stop := make(chan interface{})
	group := &sync.WaitGroup{}

	for _, component := range registry.Components() {
		component.Start(group, stop)
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		log.Printf("SIGINT received, shutting down")
		close(stop)
	}()

	group.Wait()
	log.Printf("Good bye")
}
