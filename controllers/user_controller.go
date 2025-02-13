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

// isValidRole проверяет, что переданная роль является допустимой (администратор, учитель или персонал).
func isValidRole(role models.Role) bool {
	switch role {
	case models.RoleAdmin, models.RoleTeacher, models.RoleStaff:
		return true
	default:
		return false
	}
}

// CreateUser создаёт нового пользователя.
// Если роль не указана, назначается роль по умолчанию ("учитель").
// Пароль хэшируется с помощью bcrypt, а keyID устанавливается равным ID в строковом представлении.
func CreateUser(c *gin.Context) {
	log.Println("Началась обработка запроса POST /users")
	var user models.User

	// Привязываем JSON к структуре пользователя.
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Ошибка при привязке JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка наличия пароля.
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль не может быть пустым"})
		return
	}

	// Хэширование пароля.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хэшировании пароля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя"})
		return
	}
	user.Password = string(hashedPassword)

	// Если роль не передана, назначаем роль по умолчанию.
	if user.Role == "" {
		user.Role = models.RoleTeacher // по умолчанию — "учитель"
	}

	// Проверка корректности переданной роли.
	if !isValidRole(user.Role) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Недопустимая роль пользователя. Допустимые роли: '%s', '%s', '%s'",
				models.RoleAdmin, models.RoleTeacher, models.RoleStaff),
		})
		return
	}

	// Если роль администратора, даём доступ ко всем комнатам.
	if user.Role == models.RoleAdmin {
		user.AccessRooms = []string{"*"}
	}

	// Генерация нового ObjectID для MongoDB.
	user.ID = primitive.NewObjectID()
	// Используем строковое представление ID как keyID.
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

// GetUsers возвращает список всех пользователей.
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

// UpdateUser обновляет данные пользователя.
// Если передан новый пароль, он хэшируется с помощью bcrypt.
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

	// Если передана роль, проверяем её корректность.
	if updatedUser.Role != "" && !isValidRole(updatedUser.Role) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Недопустимая роль пользователя. Допустимые роли: '%s', '%s', '%s'",
				models.RoleAdmin, models.RoleTeacher, models.RoleStaff),
		})
		return
	}

	// Если роль — администратор, даём доступ ко всем комнатам.
	if updatedUser.Role == models.RoleAdmin {
		updatedUser.AccessRooms = []string{"*"}
	}

	// Если обновляется пароль, хэшируем его.
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Ошибка при хэшировании пароля: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить пользователя"})
			return
		}
		updatedUser.Password = string(hashedPassword)
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updatedUser}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Printf("Ошибка при обновлении пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить пользователя"})
		return
	}

	log.Printf("Пользователь (ID: %s) успешно обновлен", id)
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно обновлен"})
}

// DeleteUser удаляет пользователя.
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
