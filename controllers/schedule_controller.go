package controllers

import (
	"access-control-system/config"
	"access-control-system/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// getUserIDByName ищет пользователя по имени и фамилии и возвращает его ObjectID
func getUserIDByName(firstName, secondName string) (primitive.ObjectID, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Database("ENU").Collection("users")
	filter := bson.M{"first_name": firstName, "second_name": secondName}

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return user.ID, nil
}

// CreateSchedule создает новое расписание
func CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Найти пользователя по имени и фамилии
	userID, err := getUserIDByName(schedule.FirstName, schedule.SecondName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска пользователя"})
		return
	}

	// Назначить UserID найденного пользователя и создать новый ObjectID для расписания
	schedule.UserID = userID
	schedule.ID = primitive.NewObjectID()

	collection := config.DB.Database("ENU").Collection("schedule")
	_, err = collection.InsertOne(context.TODO(), schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetSchedules извлекает все расписания из базы данных
func GetSchedules(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нет соединения с базой данных"})
		return
	}

	cursor, err := config.DB.Database("ENU").Collection("schedule").Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных", "details": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var schedules []models.Schedule
	if err := cursor.All(context.TODO(), &schedules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при преобразовании данных", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Список расписаний", "data": schedules})
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

	// При необходимости обновляем userID по имени и фамилии
	if userID, err := getUserIDByName(updatedSchedule.FirstName, updatedSchedule.SecondName); err == nil {
		updatedSchedule.UserID = userID
	}

	// Формируем объект обновления вручную, чтобы не перезаписывать поля неверными или пустыми значениями
	updateData := bson.M{
		"userId":      updatedSchedule.UserID,
		"first_name":  updatedSchedule.FirstName,
		"second_name": updatedSchedule.SecondName,
		"day":         updatedSchedule.Day,
		"startTime":   updatedSchedule.StartTime,
		"endTime":     updatedSchedule.EndTime,
		"roomNumber":  updatedSchedule.RoomNumber,
		"subject":     updatedSchedule.Subject,
	}

	collection := config.DB.Database("ENU").Collection("schedule")
	update := bson.M{"$set": updateData}

	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить расписание"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Расписание не найдено"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Расписание успешно обновлено", "data": updateData})
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
