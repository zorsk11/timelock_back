package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

func LogRoutes(router *gin.Engine) {
	router.POST("/logs", controllers.CreateLog)
	router.GET("/logs", controllers.GetLogs)
}
