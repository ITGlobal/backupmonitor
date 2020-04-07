package api

import (
	_ "github.com/itglobal/backupmonitor/doc" // auto generated swagger docs

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title BackupManager API
// @version 1.0
// @BasePath /

func (s *server) ConfigureSwagger() {
	s.router.StaticFile("/api/swagger.json", "doc/swagger.json")
	s.router.StaticFile("/api/swagger.yaml", "doc/swagger.yaml")

	url := ginSwagger.URL("/api/swagger.yaml")
	s.router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
