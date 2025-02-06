package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

func RoomRoutes(r *gin.Engine) {
	r.POST("/rooms", controllers.CreateRooms)

	r.GET("/rooms", controllers.GetRooms)
}
