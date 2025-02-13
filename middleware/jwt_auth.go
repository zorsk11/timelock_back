package middleware

import (
	"access-control-system/controllers"
	"access-control-system/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTAuthMiddleware проверяет JWT-токен, переданный в заголовке Authorization,
// и, при успешной валидации, помещает данные пользователя в контекст.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует заголовок авторизации"})
			c.Abort()
			return
		}

		// Ожидается формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат заголовка авторизации"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &controllers.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Используем тот же секретный ключ, что и при генерации токена.
			return []byte("my_secret_key"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный или просроченный токен"})
			c.Abort()
			return
		}

		// Формируем объект пользователя на основе данных из claims.
		user := models.User{
			Role: claims.Role,
		}
		if strings.Contains(claims.Identifier, "@") {
			user.Email = claims.Identifier
		} else {
			user.Phone = claims.Identifier
		}

		// Сохраняем данные пользователя в контекст, чтобы их можно было использовать в дальнейшем.
		c.Set("user", user)
		c.Next()
	}
}
