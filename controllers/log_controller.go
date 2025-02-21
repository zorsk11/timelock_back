package controllers

import (
	"access-control-system/config"
	"access-control-system/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LogRequest описывает структуру входящего запроса для создания лога.
type LogRequest struct {
	EventType string `json:"event_type" binding:"required"`
	Message   string `json:"message" binding:"required"`
	UserID    string `json:"user_id"` // может быть пустой строкой, если отсутствует
}

// CreateLog принимает JSON-запрос, валидирует данные и вызывает функцию LogEvent для записи лога.
func CreateLog(c *gin.Context) {
	var req LogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
		return
	}

	var userID *primitive.ObjectID
	if req.UserID != "" {
		oid, err := primitive.ObjectIDFromHex(req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат UserID"})
			return
		}
		userID = &oid
	}

	// Вызываем функцию логирования, которая записывает событие в MongoDB
	LogEvent(req.EventType, req.Message, userID)

	c.JSON(http.StatusOK, gin.H{"status": "Лог записан"})
}

// LogEvent записывает событие в коллекцию "logs" базы данных ENU.
// Если у события есть привязка к конкретному пользователю, можно передать его идентификатор.
func LogEvent(eventType, message string, userID *primitive.ObjectID) {
	logCollection := config.DB.Database("ENU").Collection("logs")
	logEntry := models.Log{
		ID:        primitive.NewObjectID(),
		EventType: eventType,
		Message:   message,
		Timestamp: time.Now(),
	}

	if userID != nil {
		logEntry.UserID = *userID
	}

	_, err := logCollection.InsertOne(context.TODO(), logEntry)
	if err != nil {
		log.Printf("Ошибка записи лога: %v", err)
	}
}

// GetLogs получает все записи логов из коллекции "logs".
func GetLogs(c *gin.Context) {
	logCollection := config.DB.Database("ENU").Collection("logs")
	cursor, err := logCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения логов: " + err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var logs []models.Log
	if err = cursor.All(context.TODO(), &logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения логов: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}
