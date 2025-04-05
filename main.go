package main

import (
	"log"

	"access-control-system/config"
	"access-control-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()


	router := gin.Default()

	config.SetupCORS(router)

	routes.RegisterPublicRoutes(router)
	routes.UserRoutes(router)
	routes.RegisterRoutes(router)
	routes.RoomRoutes(router)
	routes.ScheduleRoutes(router)
	routes.RegisterAdminRoutes(router)
	routes.LogRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
