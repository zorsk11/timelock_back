package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes регистрирует публичные маршруты.
func RegisterPublicRoutes(router *gin.Engine) {
	// Регистрация маршрута для логина.
	router.POST("/login", controllers.Login)
}
