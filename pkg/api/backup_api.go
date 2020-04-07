package api

import (
	"fmt"
	"io"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/service"
)

func (s *server) ConfigureBackupAPI() {
	controller := &backupController{
		accessRepo: service.GetAccessKeyRepository(s.services),
		backupRepo: service.GetBackupRepository(s.services),
	}

	s.router.GET("/api/backup/:id", controller.Download)
	s.router.POST("/api/backup", controller.Upload)

	s.authorized.GET("/api/projects/:id/backup", controller.List)
	s.authorized.DELETE("/api/backup/:id", controller.Delete)
}

type backupController struct {
	accessRepo service.AccessKeyRepository
	backupRepo service.BackupRepository
}

// @Summary Download backup file
// @Router /api/backup/:id [get]
// @Accept json
// @Produce application/octet-stream
// @Success 200
// @Failure 404 {object} model.Error
func (controller *backupController) Download(c *gin.Context) {
	id := c.Param("id")

	result, err := controller.backupRepo.Download(id)
	if err != nil {
		processError(c, err)
		return
	}

	defer result.File.Close()

	header := c.Writer.Header()
	header["Content-type"] = []string{"application/octet-stream"}
	header["Content-Disposition"] = []string{"attachment; filename= " + result.Backup.FileName}

	c.Status(200)
	io.Copy(c.Writer, result.File)
}

// @Summary Upload backup file
// @Router /api/backup [post]
// @Accept json
// @Produce json
// @Param key query string true "Access key"
// @Success 200 {object} model.Backup
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *backupController) Upload(c *gin.Context) {
	key := c.Query("key")

	if key == "" {
		key = c.GetHeader("Authorization")
	}

	accessKey, err := controller.accessRepo.Get(key)
	if accessKey == nil {
		c.JSON(403, model.NewError(model.EBadRequest, "access denied"))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(403, model.NewError(model.EBadRequest, "not a multipart form"))
		return
	}

	for _, files := range form.File {
		for _, file := range files {
			filename := filepath.Base(file.Filename)

			f, err := file.Open()
			if err != nil {
				processError(c, err)
				return
			}

			defer f.Close()

			backup, err := controller.backupRepo.Upload(accessKey.ProjectID, filename, f)
			if err != nil {
				processError(c, err)
				return
			}

			c.Header("Location", fmt.Sprintf("/api/backups/%s", url.QueryEscape(backup.ID)))
			c.JSON(201, backup)
			return
		}
	}

	c.JSON(400, model.NewError(model.EBadRequest, "no files uploaded"))
}

// @Summary List project's backups
// @Router /api/projects/:id/backup [get]
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 201 {object} model.Backups
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *backupController) List(c *gin.Context) {
	projectID := c.Param("id")

	list, err := controller.backupRepo.List(projectID)
	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, list)
}

// @Summary Delete a backup
// @Router /api/backup/:id [delete]
// @Accept json
// @Param id path string true "ID"
// @Success 204
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 404 {object} model.Error
func (controller *backupController) Delete(c *gin.Context) {
	id := c.Param("id")

	err := controller.backupRepo.Delete(id)
	if err != nil {
		processError(c, err)
		return
	}

	c.Status(204)
}
