package middleware

import (
	"access-control-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOnly проверяет, что пользователь аутентифицирован и имеет роль "администратор".
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Предполагается, что объект пользователя уже установлен в контекст (например, после проверки JWT).
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Вы не авторизованы"})
			c.Abort()
			return
		}

		// Приводим к типу models.User.
		currentUser, ok := user.(models.User)
		if !ok || currentUser.Role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ разрешен только администраторам"})
			c.Abort()
			return
		}

		// Если проверки прошли, передаем управление следующему обработчику.
		c.Next()
	}
}
