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

// Поиск ID пользователя по имени и фамилии
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

// Создание расписания
func CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск пользователя по имени и фамилии для заполнения UserID
	userID, err := getUserIDByName(schedule.FirstName, schedule.SecondName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска пользователя"})
		return
	}

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

// Получение списка расписаний
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

// Обновление расписания (динамическое обновление только переданных полей)
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

	updateData := bson.M{}

	// Обновляем поля, если они переданы и не пустые
	if updatedSchedule.FirstName != "" {
		updateData["first_name"] = updatedSchedule.FirstName
	}
	if updatedSchedule.SecondName != "" {
		updateData["second_name"] = updatedSchedule.SecondName
	}
	// Если оба поля заданы, пробуем обновить UserID
	if updatedSchedule.FirstName != "" && updatedSchedule.SecondName != "" {
		if userID, err := getUserIDByName(updatedSchedule.FirstName, updatedSchedule.SecondName); err == nil {
			updateData["user_id"] = userID
		}
	}
	if updatedSchedule.Day != "" {
		updateData["day"] = updatedSchedule.Day
	}
	if updatedSchedule.StartTime != "" {
		updateData["start_time"] = updatedSchedule.StartTime
	}
	if updatedSchedule.EndTime != "" {
		updateData["end_time"] = updatedSchedule.EndTime
	}
	if updatedSchedule.RoomNumber != "" {
		updateData["room_number"] = updatedSchedule.RoomNumber
	}
	if updatedSchedule.Subject != "" {
		updateData["subject"] = updatedSchedule.Subject
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет данных для обновления"})
		return
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

// Удаление расписания
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
