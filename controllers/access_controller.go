package controllers

import (
	"access-control-system/config"
	"access-control-system/models"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ScheduleWithUser объединяет данные расписания и данные пользователя.
type ScheduleWithUser struct {
	models.Schedule `bson:",inline"`
	UserInfo        struct {
		FirstName  string `bson:"first_name"`
		SecondName string `bson:"second_name"`
	} `bson:"user_info"`
}

// CheckAccess проверяет, имеет ли пользователь доступ к комнате на основе расписания.
func CheckAccess(c *gin.Context) {
	// Получаем идентификаторы пользователя и комнаты из URL-параметров
	userIDParam := c.Param("user_id")
	roomNumber := c.Param("room_number")

	// Преобразуем userID в primitive.ObjectID
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

	// Разбиваем каждый элемент AccessRooms по запятой, чтобы получить список отдельных номеров
	var rooms []string
	for _, roomStr := range user.AccessRooms {
		splitted := strings.Split(roomStr, ",")
		rooms = append(rooms, splitted...)
	}

	// Проверяем, есть ли запрошенная комната в списке доступных для пользователя
	hasRoomAccess := false
	for _, room := range rooms {
		if room == roomNumber {
			hasRoomAccess = true
			break
		}
	}
	if !hasRoomAccess {
		LogEvent("unauthorized_door_access",
			"Пользователь "+user.FirstName+" "+user.SecondName+" пытался получить доступ к комнате "+roomNumber+" без разрешения",
			nil)
		c.JSON(http.StatusForbidden, gin.H{"message": "У пользователя нет доступа к этой комнате"})
		return
	}

	// Агрегация с left join: получаем расписание с данными пользователя (FirstName, SecondName)
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"user_id", userID},
			{"room_number", roomNumber},
			{"day", currentDay},
		}}},
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user_info"},
		}}},
		{{"$unwind", "$user_info"}},
	}

	cursor, err := config.DB.Database("ENU").Collection("schedule").Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка агрегирования расписания"})
		return
	}
	var results []ScheduleWithUser
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка декодирования результатов"})
		return
	}
	if len(results) == 0 {
		LogEvent("unauthorized_schedule_access",
			"Пользователь "+user.FirstName+" "+user.SecondName+" пытался войти в комнату "+roomNumber+" без действующего расписания на "+currentDay,
			&user.ID)
		c.JSON(http.StatusForbidden, gin.H{"message": "Нет расписания для этой комнаты в данный день"})
		return
	}
	scheduleWithUser := results[0]
	schedule := scheduleWithUser.Schedule

	// Если время расписания задано без секунд (например, "17:10"), добавляем ":00"
	if len(schedule.StartTime) == 5 {
		schedule.StartTime = schedule.StartTime + ":00"
	}
	if len(schedule.EndTime) == 5 {
		schedule.EndTime = schedule.EndTime + ":00"
	}

	// Формируем полное время начала и окончания расписания
	dateStr := now.Format("2006-01-02")
	startDateTimeStr := fmt.Sprintf("%s %s", dateStr, schedule.StartTime)
	endDateTimeStr := fmt.Sprintf("%s %s", dateStr, schedule.EndTime)

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
	if now.After(startTime) && now.Before(endTime) {
		// В ответе можно вернуть имя пользователя вместо user_id
		c.JSON(http.StatusOK, gin.H{
			"message":    "Доступ разрешен",
			"firstName":  scheduleWithUser.UserInfo.FirstName,
			"secondName": scheduleWithUser.UserInfo.SecondName,
		})
	} else {
		LogEvent("unauthorized_time_access",
			"Пользователь "+scheduleWithUser.UserInfo.FirstName+" "+scheduleWithUser.UserInfo.SecondName+" пытался войти в комнату "+roomNumber+" в неразрешённое время ("+now.Format("15:04:05")+
				"). Допустимое время: "+schedule.StartTime+" - "+schedule.EndTime,
			&user.ID)
		c.JSON(http.StatusForbidden, gin.H{"message": "Время доступа не соответствует расписанию"})
	}
}
