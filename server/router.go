package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kneed/iot-device-simulator/controller/httpv1"
	"github.com/kneed/iot-device-simulator/pkg/logging"
	"github.com/kneed/iot-device-simulator/settings"
	log "github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func NewRouter() *gin.Engine {
	env := settings.AppSetting.LogLevel
	if env != "Debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	// 日志中间件
	router.Use(logging.LoggerMiddleware())
	// panic处理中间件
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Error(err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// swagger文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Route := router.Group("/simulator/api/v1")
	{
		deviceGroup := v1Route.Group("devices")
		{
			deviceGroup.POST("", httpv1.CreateDevice)
			deviceGroup.GET("", httpv1.GetDevices)
		}
		protocolGroup := v1Route.Group("protocols")
		{

			protocolGroup.GET("", httpv1.GetProtocols)
			protocolGroup.POST("", httpv1.CreateProtocol)
		}
	}
	return router

}
