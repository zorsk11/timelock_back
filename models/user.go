package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleAdmin   Role = "администратор"
	RoleTeacher Role = "учитель"
	RoleStaff   Role = "персонал"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	KeyID       string             `json:"key_id" bson:"key_id"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	SecondName  string             `json:"second_name" bson:"second_name"`
	Email       string             `json:"email" bson:"email"`
	AccessRooms []string           `json:"access_rooms" bson:"access_rooms"`
	Photos      []string           `json:"photos,omitempty" bson:"photos,omitempty"`
	Address     string             `json:"address,omitempty" bson:"address,omitempty"` // Адрес (необязательно)
	Phone       string             `json:"phone,omitempty" bson:"phone,omitempty"`     // Телефон (необязательно)
	Country     string             `json:"country,omitempty" bson:"country,omitempty"` // Страна (необязательно)
	City        string             `json:"city,omitempty" bson:"city,omitempty"`       // Город (необязательно)
	Role        Role               `json:"role" bson:"role"`                           // Роль пользователя
	Password    string             `json:"password,omitempty" bson:"password"`         // Пароль (необязательно)
}
