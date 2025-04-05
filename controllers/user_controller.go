package controllers

import (
	"access-control-system/config"
	"access-control-system/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func isValidRole(role models.Role) bool {
	switch role {
	case models.RoleAdmin, models.RoleTeacher, models.RoleStaff:
		return true
	default:
		return false
	}
}

func CreateUser(c *gin.Context) {
	log.Println("Началась обработка запроса POST /users")
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Ошибка при привязке JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль не может быть пустым"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хэшировании пароля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}
	user.Password = string(hashedPassword)

	if user.Role == "" {
		user.Role = models.RoleTeacher
	}

	if !isValidRole(user.Role) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Недопустимая роль пользователя. Допустимые роли: '%s', '%s', '%s'",
				models.RoleAdmin, models.RoleTeacher, models.RoleStaff),
		})
		return
	}

	if user.Role == models.RoleAdmin {
		user.AccessRooms = []string{"*"}
	}

	user.ID = primitive.NewObjectID()
	user.KeyID = user.ID.Hex()

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}

	log.Printf("Пользователь успешно создан (ID: %s, KeyID: %s)", user.ID.Hex(), user.KeyID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Пользователь успешно создан",
		"user":    user,
	})
	log.Println("Запрос POST /users обработан успешно")
}

func GetUsers(c *gin.Context) {
	log.Println("Началась обработка запроса GET /users")

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Ошибка при получении пользователей: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Не удалось получить пользователей: %v", err)})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Printf("Ошибка при декодировании пользователей: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Не удалось декодировать пользователей: %v", err)})
		return
	}

	if len(users) == 0 {
		log.Println("Пользователи не найдены")
		c.JSON(http.StatusOK, gin.H{"message": "Пользователи не найдены"})
		return
	}

	log.Printf("Найдено пользователей: %d", len(users))
	c.JSON(http.StatusOK, users)
	log.Println("Запрос GET /users обработан успешно")
}

func UpdateUser(c *gin.Context) {
	log.Println("Началась обработка запроса PUT /users/:id")

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Ошибка при конвертации ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		log.Printf("Ошибка при привязке JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Если роль передана и она некорректная, возвращаем ошибку
	if updatedUser.Role != "" && !isValidRole(updatedUser.Role) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Недопустимая роль пользователя. Допустимые роли: '%s', '%s', '%s'",
				models.RoleAdmin, models.RoleTeacher, models.RoleStaff),
		})
		return
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateFields := bson.M{}

	// Обновляем только непустые поля
	if updatedUser.FirstName != "" {
		updateFields["first_name"] = updatedUser.FirstName
	}
	if updatedUser.SecondName != "" {
		updateFields["second_name"] = updatedUser.SecondName
	}
	if updatedUser.Email != "" {
		updateFields["email"] = updatedUser.Email
	}
	if updatedUser.Phone != "" {
		updateFields["phone"] = updatedUser.Phone
	}
	if updatedUser.Address != "" {
		updateFields["address"] = updatedUser.Address
	}
	if updatedUser.City != "" {
		updateFields["city"] = updatedUser.City
	}
	if updatedUser.Country != "" {
		updateFields["country"] = updatedUser.Country
	}
	if updatedUser.Role != "" {
		updateFields["role"] = updatedUser.Role
		// Если роль - администратор, устанавливаем доступ ко всем комнатам
		if updatedUser.Role == models.RoleAdmin {
			updateFields["access_rooms"] = []string{"*"}
		}
	}
	// Если передан непустой список комнат, обновляем поле access_rooms
	if len(updatedUser.AccessRooms) > 0 {
		updateFields["access_rooms"] = updatedUser.AccessRooms
	}
	// Обновляем пароль, если он передан
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Ошибка при хэшировании пароля: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить пользователя"})
			return
		}
		updateFields["password"] = string(hashedPassword)
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет данных для обновления"})
		return
	}

	update := bson.M{"$set": updateFields}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Printf("Ошибка при обновлении пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить пользователя"})
		return
	}

	log.Printf("Пользователь (ID: %s) успешно обновлен", id)
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно обновлен"})
}

func DeleteUser(c *gin.Context) {
	log.Println("Началась обработка запроса DELETE /users/:id")

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Ошибка при конвертации ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Printf("Ошибка при удалении пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить пользователя"})
		return
	}

	if result.DeletedCount == 0 {
		log.Printf("Пользователь с ID %s не найден", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	log.Printf("Пользователь (ID: %s) успешно удален", id)
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
}
