package routes

import (
	"access-control-system/controllers"
	"access-control-system/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.JWTAuthMiddleware(), middleware.AdminOnly())
	{
		adminGroup.GET("/dashboard", controllers.AdminDashboard)
	}
}
