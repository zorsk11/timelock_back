package controllers

import (
	"access-control-system/config"
	"access-control-system/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)


var jwtKey = []byte("my_secret_key")

type Claims struct {
	ID         string      `json:"id"`
	Identifier string      `json:"identifier"` 
	Role       models.Role `json:"role"`
	jwt.RegisteredClaims
}


func Login(c *gin.Context) {
	log.Println("Началась обработка запроса POST /login")

	var credentials struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Printf("Ошибка при привязке JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	filter := bson.M{
		"$or": []bson.M{
			{"email": credentials.Identifier},
			{"phone": credentials.Identifier},
		},
	}

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Printf("Ошибка при поиске пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		log.Printf("Неверный пароль для пользователя %s: %v", credentials.Identifier, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	if user.Role != models.RoleAdmin {
		LogEvent("admin_access_denied",
			"Пользователь "+user.FirstName+" "+user.SecondName+" с ролью '"+string(user.Role)+"' пытался авторизоваться как администратор",
			nil)
		log.Printf("Попытка авторизации неадминистратора: %s", credentials.Identifier)
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ разрешен только администраторам"})
		return
	}

	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &Claims{
		ID:         user.ID.Hex(),
		Identifier: credentials.Identifier,
		Role:       user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Ошибка при генерации JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать токен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Вы успешно авторизовались",
		"token":   tokenString,
		"user": gin.H{
			"id":           user.ID.Hex(),
			"key_id":       user.KeyID,
			"first_name":   user.FirstName,
			"second_name":  user.SecondName,
			"email":        user.Email,
			"phone":        user.Phone,
			"role":         user.Role,
			"access_rooms": user.AccessRooms,
			"photos":       user.Photos,
			"address":      user.Address,
			"country":      user.Country,
			"city":         user.City,
		},
	})
}
