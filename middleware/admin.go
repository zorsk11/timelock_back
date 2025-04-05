package middleware

import (
	"access-control-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Вы не авторизованы"})
			c.Abort()
			return
		}

		currentUser, ok := user.(models.User)
		if !ok || currentUser.Role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ разрешен только администраторам"})
			c.Abort()
			return
		}

		c.Next()
	}
}
