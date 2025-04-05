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

func CreateRooms(c *gin.Context) {
	var rooms []models.Room // Слайс для нескольких комнат

	if err := c.ShouldBindJSON(&rooms); err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка парсинга JSON",
			"details": err.Error(),
		})
		return
	}

	for i := range rooms {
		rooms[i].ID = primitive.NewObjectID()
	}

	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нет соединения с базой данных"})
		return
	}

	var interfaceRooms []interface{}
	for _, room := range rooms {
		interfaceRooms = append(interfaceRooms, room)
	}

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

func GetRooms(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нет соединения с базой данных"})
		return
	}

	cursor, err := config.DB.Database("ENU").Collection("rooms").Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка получения данных из MongoDB",
			"details": err.Error(),
		})
		return
	}
	defer cursor.Close(context.TODO())

	var rooms []models.Room
	if err := cursor.All(context.TODO(), &rooms); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при преобразовании данных",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Список комнат",
		"data":    rooms,
	})
}
