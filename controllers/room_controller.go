package controllers

import (
	"access-control-system/config"
	"access-control-system/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRooms создает несколько комнат
func CreateRooms(c *gin.Context) {
	var rooms []models.Room // Слайс для нескольких комнат

	// Пробуем распарсить JSON в слайс
	if err := c.ShouldBindJSON(&rooms); err != nil {
		// Логируем ошибку и тело запроса для отладки
		fmt.Println("Ошибка парсинга JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка парсинга JSON",
			"details": err.Error(),
		})
		return
	}

	// Генерируем новый ObjectID для каждой комнаты
	for i := range rooms {
		rooms[i].ID = primitive.NewObjectID()
	}

	// Проверяем, инициализирована ли БД
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нет соединения с базой данных"})
		return
	}

	// Преобразуем слайс комнат в слайс интерфейсов
	var interfaceRooms []interface{}
	for _, room := range rooms {
		interfaceRooms = append(interfaceRooms, room)
	}

	// Вставляем несколько комнат в MongoDB
	_, err := config.DB.Database("ENU").Collection("rooms").InsertMany(context.TODO(), interfaceRooms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка вставки в MongoDB",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Комнаты созданы",
		"data":    rooms,
	})
}

// GetRooms извлекает все комнаты из базы данных
func GetRooms(c *gin.Context) {
	// Проверяем, инициализирована ли БД
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нет соединения с базой данных"})
		return
	}

	// Получаем все комнаты
	cursor, err := config.DB.Database("ENU").Collection("rooms").Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка получения данных из MongoDB",
			"details": err.Error(),
		})
		return
	}
	defer cursor.Close(context.TODO())

	// Преобразуем все документы в слайс
	var rooms []models.Room
	if err := cursor.All(context.TODO(), &rooms); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при преобразовании данных",
			"details": err.Error(),
		})
		return
	}

	// Ответ клиенту
	c.JSON(http.StatusOK, gin.H{
		"message": "Список комнат",
		"data":    rooms,
	})
}
