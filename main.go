package main

import (
	"access-control-system/config"
	"access-control-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	// Подключаем CORS
	config.SetupCORS(router)

	// Регистрируем маршруты
	routes.UserRoutes(router)
	routes.RegisterRoutes(router)
	routes.RoomRoutes(router)
	routes.ScheduleRoutes(router)

	router.Run(":8080")
}
