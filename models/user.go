package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User представляет пользователя в системе
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Role        string             `json:"role" bson:"role"`
	KeyID       string             `json:"key_id" bson:"key_id"`
	AccessRooms []string           `json:"access_rooms" bson:"access_rooms"` // Список комнат, доступных для пользователя
}
