package main

import (
	"log"

	"access-control-system/config"
	"access-control-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключаемся к базе данных
	config.ConnectDB()

	// Устанавливаем режим работы (Debug или Release)
	if config.Debug {
		gin.SetMode(gin.DebugMode)
		log.Println("Запуск в режиме Debug")
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Println("Запуск в режиме Release")
	}

	// Создаем экземпляр маршрутизатора
	router := gin.Default()

	// Подключаем CORS (настройки CORS определены в config.SetupCORS)
	config.SetupCORS(router)

	// Регистрируем публичные маршруты, включая /login
	routes.RegisterPublicRoutes(router)
	routes.UserRoutes(router)
	routes.RegisterRoutes(router)
	routes.RoomRoutes(router)
	routes.ScheduleRoutes(router)

	// Регистрируем маршруты для администраторов
	routes.RegisterAdminRoutes(router)

	// Запускаем сервер на порту 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
