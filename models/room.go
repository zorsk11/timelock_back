package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomNumber         string             `bson:"room_number" json:"room_number"`
	Floor              int                `bson:"floor" json:"floor"`
	AccessControllerID string             `bson:"access_controller_id" json:"access_controller_id"`
}
	