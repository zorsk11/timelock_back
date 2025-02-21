package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

// LogRoutes регистрирует маршруты, связанные с логированием.
func LogRoutes(router *gin.Engine) {
	// POST /logs - создание записи лога
	router.POST("/logs", controllers.CreateLog)
	// GET /logs - получение всех логов
	router.GET("/logs", controllers.GetLogs)
}
