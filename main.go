package main

import (
	"log"

	"access-control-system/config"
	"access-control-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	// Устанавливаем режим работы
	if config.Debug {
		gin.SetMode(gin.DebugMode)
		log.Println("Running in Debug mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Println("Running in Release mode")
	}

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
