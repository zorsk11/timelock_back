package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Проверка доступа пользователя к комнате

	router.GET("/access/:user_id/room/:room_number", controllers.CheckAccess)
}
