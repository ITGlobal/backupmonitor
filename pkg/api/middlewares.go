package api

import (
	"fmt"
	"github.com/itglobal/backupmonitor/pkg/model"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itglobal/backupmonitor/pkg/service"
	"github.com/sarulabs/di"
)

func createLoggerMiddleware(log *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if c.Writer.Status() >= 400 {
			latency := time.Now().Sub(start)
			var request string
			if c.Request.URL.RawQuery != "" {
				request = fmt.Sprintf(
					"%s %s?%s",
					c.Request.Method,
					c.Request.URL.Path,
					c.Request.URL.RawQuery)
			} else {
				request = fmt.Sprintf(
					"%s %s",
					c.Request.Method,
					c.Request.URL.Path)
			}

			log.Printf(
				"%s -> %d (%s) %s",
				request,
				c.Writer.Status(),
				latency,
				c.Errors.ByType(gin.ErrorTypePrivate))
		}
	}
}

func createJwtMiddleware(c di.Container) gin.HandlerFunc {
	jwt := service.GetJwt(c)

	var processJwtError = func(c *gin.Context, e *service.JwtError) {
		if e.StatusCode == 403 {
			c.Header("WWW-Authenticate", "bearer")
		}
		c.JSON(e.StatusCode, model.NewError(model.EAccessDenied, e.Message))
		c.Abort()
	}

	var middleware gin.HandlerFunc = func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			processJwtError(c, service.NewJwtError(403, "missing access token"))
			return
		}

		splitted := strings.Split(header, " ")
		if len(splitted) != 2 {
			processJwtError(c, service.NewJwtError(403, "malformed access token"))
			return
		}

		if strings.ToLower(splitted[0]) != "bearer" {
			processJwtError(c, service.NewJwtError(403, "invalid realm"))
			return
		}

		user, e := jwt.ValidateToken(splitted[1])
		if e != nil {
			processJwtError(c, e)
			return
		}

		c.Set(gin.AuthUserKey, user)
		c.Next()
	}

	return middleware
}
