package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/service"
)

func (s *server) ConfigureAccessAPI() {
	controller := &accessController{
		repository: service.GetAccessKeyRepository(s.services),
	}

	s.authorized.GET("/api/projects/:id/keys", controller.List)
	s.authorized.GET("/api/projects/:id/keys/:key", controller.Get)
	s.authorized.POST("/api/projects/:id/keys", controller.Post)
	s.authorized.DELETE("/api/projects/:id/keys/:key", controller.Delete)
}

type accessController struct {
	repository service.AccessKeyRepository
}

// @Summary List project's access keys
// @Router /api/projects/:id/keys [get]
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} model.AccessKeys
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *accessController) List(c *gin.Context) {
	projectID := c.Param("id")
	list, err := controller.repository.List(projectID)

	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, list)
}

// @Summary Get a project's access key
// @Router /api/projects/:id/keys/:key [get]
// @Accept json
// @Param id path string true "Project ID"
// @Param key path string true "Access Key ID"
// @Success 200 {object} model.AccessKey
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *accessController) Get(c *gin.Context) {
	projectID := c.Param("id")
	accessKeyID, err := strconv.Atoi(c.Param("key"))
	if err != nil {
		processError(c, err)
		return
	}

	p, err := controller.repository.GetByID(projectID, accessKeyID)
	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, p)
}

// @Summary Create new project's access key
// @Router /api/projects/:id/keys [post]
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param body body model.AccessKeyCreateParams true "Body"
// @Success 201 {object} model.AccessKey
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *accessController) Post(c *gin.Context) {
	projectID := c.Param("id")

	var req model.AccessKeyCreateParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid request parameters"))
		return
	}

	p, err := controller.repository.Create(projectID, &req)
	if err != nil {
		processError(c, err)
		return
	}

	c.Header("Location", fmt.Sprintf("/api/projects/%s/keys/%d", url.QueryEscape(p.ProjectID), p.ID))
	c.JSON(201, p)
}

// @Summary Delete a project's access key
// @Router /api/projects/:id/keys/:key [delete]
// @Accept json
// @Param id path string true "Project ID"
// @Param key path string true "Access Key ID"
// @Success 204
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *accessController) Delete(c *gin.Context) {
	projectID := c.Param("id")
	accessKeyID, err := strconv.Atoi(c.Param("key"))
	if err != nil {
		processError(c, err)
		return
	}

	err = controller.repository.Delete(projectID, accessKeyID)
	if err != nil {
		processError(c, err)
		return
	}

	c.Status(204)
}
