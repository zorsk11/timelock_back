package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminDashboard — обработчик защищённого маршрута для администраторов.
func AdminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Добро пожаловать в админскую панель",
	})
}
