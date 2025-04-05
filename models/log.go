package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EventType string             `json:"event_type" bson:"event_type"` 
	Message   string             `json:"message" bson:"message"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}
