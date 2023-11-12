package web

import (
	"github.com/KatsayArtemDev/verification/src/result"
	"github.com/KatsayArtemDev/verification/src/web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app app) success(c *gin.Context, object any) {
	c.JSON(http.StatusOK, result.Success(object))
}

func (app app) fail(c *gin.Context, code int, err error) {
	err = c.Error(err).SetType(gin.ErrorTypePrivate)
	// TODO: think that not all logs this must write, skip for example NotAuthorized
	app.logger.Errorln(err.Error())
	c.JSON(code, result.HttpFail(c.FullPath(), code, err))
}

func serverRouter(app app) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Cors())

	routerGroup := router.Group("/api/v1")

	routerGroup.POST("/receiving-email", app.getEmail)
	routerGroup.POST("/receiving-pin", app.getPin)
	routerGroup.POST("/resending-pin", app.resendPin)

	return router
}
