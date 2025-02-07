package config

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS настраивает CORS middleware для Gin
func SetupCORS(router *gin.Engine) {
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Используем config.Debug
	if Debug {
		corsConfig.AllowOrigins = []string{"*"} // В режиме отладки разрешаем все источники
	}

	router.Use(cors.New(corsConfig))
}
