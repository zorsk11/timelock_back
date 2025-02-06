package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes регистрирует все маршруты для проверки доступа
func RegisterRoutes(router *gin.Engine) {
	// Проверка доступа пользователя к комнате
	// Здесь :user_id – идентификатор пользователя (в виде hex строки ObjectID)
	// :room_number – номер комнаты, к которой осуществляется доступ
	router.GET("/access/:user_id/room/:room_number", controllers.CheckAccess)
}
