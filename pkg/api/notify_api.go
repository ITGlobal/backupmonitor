package api

import (
	"github.com/gin-gonic/gin"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/notify"
)

func (s *server) ConfigureNotifyAPI() {
	svc := notify.GetService(s.services)
	controller := &notifyController{svc}

	s.authorized.POST("/api/notify/slack", controller.NotifySlack)
	s.authorized.POST("/api/notify/telegram", controller.NotifyTelegram)
	s.authorized.POST("/api/notify/webhook", controller.NotifyWebhook)
}

type notifyController struct {
	service notify.Service
}

// @Summary Send a test Slack notification
// @Router /api/notify/slack [post]
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param body body model.TestSlackNotificationRequest true "Body"
// @Success 200 {object} model.Empty
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 500 {object} model.Error
func (controller *notifyController) NotifySlack(c *gin.Context) {
	var req model.TestSlackNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid request parameters"))
		return
	}

	msg := req.ToMessage()

	err := controller.service.NotifySlack(msg)
	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, model.Empty{})
}

// @Summary Send a test Telegram notification
// @Router /api/notify/telegram [post]
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param body body model.TestTelegramNotificationRequest true "Body"
// @Success 200 {object} model.Empty
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 500 {object} model.Error
func (controller *notifyController) NotifyTelegram(c *gin.Context) {
	var req model.TestTelegramNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid request parameters"))
		return
	}

	msg := req.ToMessage()

	err := controller.service.NotifyTelegram(msg)
	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, model.Empty{})
}

// @Summary Send a test Webhook notification
// @Router /api/notify/webhook [post]
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param body body model.TestWebhookNotificationRequest true "Body"
// @Success 200 {object} model.Empty
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
// @Failure 500 {object} model.Error
func (controller *notifyController) NotifyWebhook(c *gin.Context) {
	var req model.TestWebhookNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid request parameters"))
		return
	}

	msg := req.ToMessage()

	err := controller.service.NotifyWebhook(msg)
	if err != nil {
		processError(c, err)
		return
	}

	c.JSON(200, model.Empty{})
}
