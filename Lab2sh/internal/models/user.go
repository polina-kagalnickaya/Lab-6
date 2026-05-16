package models

import (
    "time"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// User представляет пользователя в MongoDB
type User struct {
    ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
    Email     string              `bson:"email" json:"email"`
    Password  string              `bson:"password,omitempty" json:"-"`
    FullName  string              `bson:"full_name" json:"full_name"`
    YandexID  *string             `bson:"yandex_id,omitempty" json:"-"`
    CreatedAt time.Time           `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time           `bson:"updated_at" json:"updated_at"`
    DeletedAt *time.Time          `bson:"deleted_at,omitempty" json:"-"`
}