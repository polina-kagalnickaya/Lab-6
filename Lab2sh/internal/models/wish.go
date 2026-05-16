package models

import (
    "time"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Wish представляет желание в MongoDB
type Wish struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
    Text      string             `bson:"text" json:"text"`
    Author    string             `bson:"author" json:"author"`
    Priority  int                `bson:"priority" json:"priority"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
    DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"-"`
}