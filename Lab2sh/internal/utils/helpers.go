package utils

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDToUint конвертирует ObjectID в uint для использования в Redis
func ObjectIDToUint(id primitive.ObjectID) uint {
    if id.IsZero() {
        return 0
    }
    // Используем timestamp из ObjectID
    return uint(id.Timestamp().Unix())
}