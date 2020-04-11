package api

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/service"
)

func (s *server) ConfigureProjectsAPI() {
	projectRepository := service.GetProjectRepository(s.services)
	accessKeyRepository := service.GetAccessKeyRepository(s.services)
	backupRepository := service.GetBackupRepository(s.services)

	controller := &projectController{
		projectRepository:   projectRepository,
		accessKeyRepository: accessKeyRepository,
		backupRepository:    backupRepository,
	}

	s.authorized.GET("/api/projects", controller.List)
	s.authorized.GET("/api/projects/:id", controller.Get)
	s.authorized.POST("/api/projects", controller.Post)
	s.authorized.PUT("/api/projects/:id", controller.Put)
	s.authorized.DELETE("/api/projects/:id", controller.Delete)
}

type projectController struct {
	projectRepository   service.ProjectRepository
	accessKeyRepository service.AccessKeyRepository
	backupRepository    service.BackupRepository
}

// @Summary List projects
// @Router /api/projects [get]
// @Accept json
// @Produce json
// @Success 200 {object} model.Projects
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
func (controller *projectController) List(c *gin.Context) {
	list, err := controller.projectRepository.List()
	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, list)
}

// @Summary Get a project
// @Router /api/projects/:id [get]
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} model.Project
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *projectController) Get(c *gin.Context) {
	id := c.Param("id")

	p, err := controller.projectRepository.Get(id)
	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, p)
}

// @Summary Create new project
// @Router /api/projects [post]
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param body body model.ProjectCreateParams true "Body"
// @Success 201 {object} model.Project
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 409 {object} model.Error
func (controller *projectController) Post(c *gin.Context) {
	var req model.ProjectCreateParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid request parameters"))
		return
	}

	p, err := controller.projectRepository.Create(&req)
	if err != nil {
		processError(c, err)
		return
	}

	c.Header("Location", fmt.Sprintf("/api/projects/%s", url.QueryEscape(p.ID)))
	c.JSON(201, p)
}

// @Summary Update existing project
// @Router /api/projects/:id [post]
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param body body model.ProjectUpdateParams true "Body"
// @Success 201 {object} model.Project
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *projectController) Put(c *gin.Context) {
	id := c.Param("id")

	var req model.ProjectUpdateParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid request parameters"))
		return
	}

	p, err := controller.projectRepository.Update(id, &req)
	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, p)
}

// @Summary Delete existing project
// @Router /api/projects/:id [delete]
// @Accept json
// @Param id path string true "ID"
// @Success 204
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *projectController) Delete(c *gin.Context) {
	id := c.Param("id")

	// Mark projects as disabled
	isActive := false
	req := &model.ProjectUpdateParams{
		IsActive: &isActive,
	}
	_, err := controller.projectRepository.Update(id, req)
	if err != nil {
		processError(c, err)
		return
	}

	// Drop all access keys
	accessKeys, err := controller.accessKeyRepository.List(id)
	if err != nil {
		processError(c, err)
		return
	}
	for _, accessKey := range accessKeys {
		err = controller.accessKeyRepository.Delete(id, accessKey.ID)
		if err != nil {
			processError(c, err)
			return
		}
	}

	// Drop all backups
	backups, err := controller.backupRepository.List(id)
	if err != nil {
		processError(c, err)
		return
	}
	for _, backup := range backups {
		err = controller.backupRepository.Delete(backup.ID, "project deletion")
		if err != nil {
			processError(c, err)
			return
		}
	}

	// Drop project
	err = controller.projectRepository.Delete(id)
	if err != nil {
		processError(c, err)
		return
	}

	c.Status(204)
}
