package routes

import (
	"access-control-system/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/users", controllers.CreateUser)
	router.GET("/users", controllers.GetUsers)
	router.PUT("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)

}
