package controllers

import (
	"access-control-system/config"
	"access-control-system/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckAccess проверяет, имеет ли пользователь доступ к комнате на основе расписания.
// Если день недели или текущее время не совпадают с данными из расписания, доступ будет запрещён.
func CheckAccess(c *gin.Context) {
	// Получаем идентификаторы пользователя и комнаты из URL-параметров
	userIDParam := c.Param("user_id")
	roomNumber := c.Param("room_number")

	// Преобразуем userID в primitive.ObjectID, если он хранится в таком виде
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат user_id"})
		return
	}

	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки временной зоны"})
		return
	}

	now := time.Now().In(loc)
	currentDay := now.Weekday().String()

	var user models.User
	err = config.DB.Database("ENU").Collection("users").
		FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Проверяем, есть ли запрошенная комната в списке доступных для пользователя
	hasRoomAccess := false
	for _, room := range user.AccessRooms {
		if room == roomNumber {
			hasRoomAccess = true
			break
		}
	}
	if !hasRoomAccess {
		c.JSON(http.StatusForbidden, gin.H{"message": "У пользователя нет доступа к этой комнате"})
		return
	}

	var schedule models.Schedule
	err = config.DB.Database("ENU").Collection("schedule").
		FindOne(context.TODO(), bson.M{
			"user_id":     userID,
			"room_number": roomNumber,
			"day":         currentDay,
		}).Decode(&schedule)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Нет расписания для этой комнаты в данный день"})
		return
	}

	if schedule.Day != currentDay {
		c.JSON(http.StatusForbidden, gin.H{"message": "День доступа не соответствует расписанию"})
		return
	}

	// Если время расписания задано без секунд (например, "17:10"), добавляем ":00"
	if len(schedule.StartTime) == 5 {
		schedule.StartTime = schedule.StartTime + ":00"
	}
	if len(schedule.EndTime) == 5 {
		schedule.EndTime = schedule.EndTime + ":00"
	}

	// Формируем полное время начала и окончания расписания, привязывая их к сегодняшней дате
	dateStr := now.Format("2006-01-02")
	startDateTimeStr := fmt.Sprintf("%s %s", dateStr, schedule.StartTime)
	endDateTimeStr := fmt.Sprintf("%s %s", dateStr, schedule.EndTime)

	// Парсинг времени с учётом временной зоны для Астаны
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startDateTimeStr, loc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка в формате времени начала расписания"})
		return
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", endDateTimeStr, loc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка в формате времени окончания расписания"})
		return
	}

	// Проверяем, находится ли текущее время в пределах расписания.
	// Если время входа не попадает в указанный интервал, вход запрещается.
	if now.After(startTime) && now.Before(endTime) {
		c.JSON(http.StatusOK, gin.H{"message": "Доступ разрешен"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "Время доступа не соответствует расписанию"})
	}
}
