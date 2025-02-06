package controllers

import (
	"access-control-system/config"
	"access-control-system/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateSchedule создает новое расписание
func CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule.ID = primitive.NewObjectID()

	// Исправленный доступ к коллекции schedule в базе данных ENU
	collection := config.DB.Database("ENU").Collection("schedule")
	_, err := collection.InsertOne(context.TODO(), schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetSchedules извлекает все расписания из базы данных
func GetSchedules(c *gin.Context) {
	// Проверка соединения с базой данных
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нет соединения с базой данных"})
		return
	}

	// Получаем все расписания
	cursor, err := config.DB.Database("ENU").Collection("schedule").Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных", "details": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	// Преобразуем все документы в слайс
	var schedules []models.Schedule
	if err := cursor.All(context.TODO(), &schedules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при преобразовании данных", "details": err.Error()})
		return
	}

	// Ответ клиенту
	c.JSON(http.StatusOK, gin.H{
		"message": "Список расписаний",
		"data":    schedules,
	})
}

// UpdateSchedule обновляет данные расписания
func UpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID расписания"})
		return
	}

	var updatedSchedule models.Schedule
	if err := c.ShouldBindJSON(&updatedSchedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := config.DB.Database("ENU").Collection("schedule")
	update := bson.M{"$set": updatedSchedule}
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить расписание"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Расписание не найдено"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Расписание успешно обновлено"})
}

// DeleteSchedule удаляет расписание
func DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID расписания"})
		return
	}

	collection := config.DB.Database("ENU").Collection("schedule")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить расписание"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Расписание не найдено"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Расписание успешно удалено"})
}
