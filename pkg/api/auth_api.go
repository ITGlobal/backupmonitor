package api

import (
	"github.com/gin-gonic/gin"
	"github.com/itglobal/backupmonitor/pkg/model"
	"github.com/itglobal/backupmonitor/pkg/service"
)

func (s *server) ConfigureAuthAPI() {
	repository := service.GetUserRepository(s.services)
	jwt := service.GetJwt(s.services)

	controller := &authController{repository, jwt}

	s.router.POST("/api/authorize", controller.Authorize)
	s.authorized.GET("/api/me", controller.GetMe)
	s.authorized.POST("/api/me/password", controller.ChangePassword)
}

type authController struct {
	repository service.UserRepository
	jwt        service.Jwt
}

// @Summary Get an access token
// @Router /api/authorize [post]
// @Accept json
// @Produce json
// @Param account body model.AuthRequest true "Request"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} model.Error
func (t *authController) Authorize(c *gin.Context) {
	var request model.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid request parameters"))
		return
	}

	user, err := t.repository.Get(request.Username)
	if err == nil {
		if user.CheckPassword(request.Password) {
			token, err := t.jwt.GenerateToken(user)
			if err != nil {
				processError(c, err)
				return
			}

			c.JSON(200, &model.AuthResponse{
				Token: token, 
				User: user})
			return
		}
	}

	c.JSON(400, model.NewError(model.EBadRequest, "invalid credentials"))
}

// @Summary Get current user
// @Router /api/me [get]
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
func (t *authController) GetMe(c *gin.Context) {
	u, _ := c.Get(gin.AuthUserKey)
	user := u.(*model.User)

	c.JSON(200, user)
}

// @Summary Change current user's password
// @Router /api/me/password [post]
// @Accept json
// @Produce json
// @Param account body model.UserChangePasswordRequest true "Request"
// @Success 200 {object} model.EmptyResponse
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 403 {object} model.Error
func (t *authController) ChangePassword(c *gin.Context) {
	var request model.UserChangePasswordRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid request parameters"))
		return
	}

	u, _ := c.Get(gin.AuthUserKey)
	user := u.(*model.User)

	user, err = t.repository.Get(user.UserName)
	if err != nil {
		panic(err)
	}

	if !user.CheckPassword(request.OldPassword) {
		c.JSON(400, model.NewError(model.EBadRequest, "invalid old password"))
		return
	}

	err = t.repository.SetPassword(user.ID, request.NewPassword)
	if err != nil {
		panic(err)
	}

	c.JSON(200, &model.EmptyResponse{})
}
