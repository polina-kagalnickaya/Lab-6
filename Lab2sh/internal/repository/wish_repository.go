package repository

import (
    "context"
    "time"
    
    "newyear-app/internal/models"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type WishRepository struct {
    collection *mongo.Collection
}

func NewWishRepository(db *mongo.Database) *WishRepository {
    collection := db.Collection("wishes")
    
    // Создаем индексы
    ctx := context.Background()
    indexes := []mongo.IndexModel{
        {
            Keys: bson.D{{Key: "user_id", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "deleted_at", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "priority", Value: -1}},
        },
        {
            Keys: bson.D{{Key: "created_at", Value: -1}},
        },
    }
    collection.Indexes().CreateMany(ctx, indexes)
    
    return &WishRepository{collection: collection}
}

func (r *WishRepository) Create(wish *models.Wish) error {
    ctx := context.Background()
    wish.ID = primitive.NewObjectID()
    wish.CreatedAt = time.Now()
    wish.UpdatedAt = time.Now()
    
    if wish.Author == "" {
        wish.Author = "Anonymous"
    }
    if wish.Priority == 0 {
        wish.Priority = 1
    }
    
    _, err := r.collection.InsertOne(ctx, wish)
    return err
}

func (r *WishRepository) FindByID(id primitive.ObjectID) (*models.Wish, error) {
    ctx := context.Background()
    filter := bson.M{
        "_id": id,
        "deleted_at": nil,
    }
    
    var wish models.Wish
    err := r.collection.FindOne(ctx, filter).Decode(&wish)
    if err != nil {
        return nil, err
    }
    
    return &wish, nil
}

func (r *WishRepository) FindAll(page, limit int, userID primitive.ObjectID) ([]models.Wish, int64, error) {
    ctx := context.Background()
    
    filter := bson.M{
        "user_id": userID,
        "deleted_at": nil,
    }
    
    // Получаем общее количество
    total, err := r.collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }
    
    // Настройки пагинации и сортировки
    opts := options.Find().
        SetSkip(int64((page - 1) * limit)).
        SetLimit(int64(limit)).
        SetSort(bson.D{
            {Key: "priority", Value: -1},
            {Key: "created_at", Value: -1},
        })
    
    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)
    
    var wishes []models.Wish
    if err := cursor.All(ctx, &wishes); err != nil {
        return nil, 0, err
    }
    
    return wishes, total, nil
}

func (r *WishRepository) FindAllPublic(page, limit int) ([]models.Wish, int64, error) {
    ctx := context.Background()
    
    filter := bson.M{
        "deleted_at": nil,
    }
    
    // Получаем общее количество
    total, err := r.collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }
    
    // Настройки пагинации и сортировки
    opts := options.Find().
        SetSkip(int64((page - 1) * limit)).
        SetLimit(int64(limit)).
        SetSort(bson.D{
            {Key: "priority", Value: -1},
            {Key: "created_at", Value: -1},
        })
    
    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)
    
    var wishes []models.Wish
    if err := cursor.All(ctx, &wishes); err != nil {
        return nil, 0, err
    }
    
    return wishes, total, nil
}

func (r *WishRepository) Update(wish *models.Wish) error {
    ctx := context.Background()
    wish.UpdatedAt = time.Now()
    
    filter := bson.M{"_id": wish.ID}
    update := bson.M{"$set": wish}
    
    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}

func (r *WishRepository) Delete(id primitive.ObjectID) error {
    ctx := context.Background()
    now := time.Now()
    
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
            "deleted_at": now,
            "updated_at": now,
        },
    }
    
    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}