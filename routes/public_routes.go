package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(router *gin.Engine) {
	router.POST("/login", controllers.Login)
}
