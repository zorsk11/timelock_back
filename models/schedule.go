package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Schedule struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Day        string             `json:"day" bson:"day"`
	StartTime  string             `json:"start_time" bson:"start_time"`
	EndTime    string             `json:"end_time" bson:"end_time"`
	RoomNumber string             `json:"room_number" bson:"room_number"`
	Subject    string             `json:"subject" bson:"subject"`
}
