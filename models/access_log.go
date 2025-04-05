package models

import (
	"time"
)

type AccessLog struct {
	KeyID      string    `json:"key_id" bson:"key_id" binding:"required"`
	RoomID     string    `json:"room_id" bson:"room_id" binding:"required"`
	AccessTime time.Time `json:"access_time" bson:"access_time" binding:"required"`
	Status     string    `json:"status" bson:"status" binding:"required"` 
}
