package api

import (
	"context"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itglobal/backupmonitor/pkg/component"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
)

type server struct {
	services   di.Container
	logger     *log.Logger
	router     *gin.Engine
	authorized gin.IRoutes
}

const serverKey = "ApiServer"

// Setup configures package services
func Setup(builder component.Builder) {

	builder.AddService(di.Def{
		Name: serverKey,
		Build: func(c di.Container) (interface{}, error) {
			return newServer(c), nil
		},
	})

	builder.AddComponent(createComponent)
}

func newServer(c di.Container) *server {
	logger := log.New(log.Writer(), "[api] ", log.Flags())

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(createLoggerMiddleware(logger))
	router.Use(gin.Recovery())

	authorized := router.Group("/")
	authorized.Use(createJwtMiddleware(c))

	s := &server{c, logger, router, authorized}
	return s
}

func createComponent(c di.Container) (component.T, error) {
	server := c.Get(serverKey).(*server)

	server.ConfigureSwagger()
	server.ConfigureAuthAPI()
	server.ConfigureProjectsAPI()
	server.ConfigureBackupAPI()
	server.ConfigureAccessAPI()
	server.ConfigureNotifyAPI()
	server.ConfigureStaticFiles()

	http.Handle("/", server.router)

	return server, nil
}

func (s *server) ConfigureStaticFiles() {

	s.router.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./www/index.html")
		} else {
			c.File("./www" + path.Join(dir, file))
		}
	})
}

func (s *server) Start(group *sync.WaitGroup, stop chan interface{}) {
	group.Add(1)

	addr := viper.GetString("LISTEN_ADDR")
	server := &http.Server{Addr: addr}

	go func() {
		s.logger.Printf("listening on \"%s\"\n", server.Addr)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("could not listen on \"%s\": %v\n", server.Addr, err)
		}

		group.Done()
	}()

	go func() {
		for range stop {
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		err := server.Shutdown(ctx)
		if err != nil {
			s.logger.Fatalf("could not gracefully shutdown the server: %v\n", err)
		}
	}()
}

func processError(c *gin.Context, err error) {
	e, ok := err.(*model.Error)
	if ok {
		status := 400
		switch e.Code {
		case model.EAccessDenied:
			status = 403
			break
		case model.EBadRequest:
			status = 400
			break
		case model.ENotFound:
			status = 404
			break
		case model.EConflict:
			status = 409
			break
		}

		c.JSON(status, e)
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(500, model.NewError(model.EInternalError, "internal server error"))
	c.Error(err)
	c.Abort()
	return
}
