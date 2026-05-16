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

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
    collection := db.Collection("users")
    
    // Создаем индексы
    ctx := context.Background()
    indexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"deleted_at": nil}),
        },
        {
            Keys:    bson.D{{Key: "yandex_id", Value: 1}},
            Options: options.Index().SetUnique(true).SetSparse(true),
        },
        {
            Keys: bson.D{{Key: "deleted_at", Value: 1}},
        },
    }
    collection.Indexes().CreateMany(ctx, indexes)
    
    return &UserRepository{collection: collection}
}

func (r *UserRepository) Create(user *models.User) error {
    ctx := context.Background()
    user.ID = primitive.NewObjectID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    _, err := r.collection.InsertOne(ctx, user)
    return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    ctx := context.Background()
    filter := bson.M{
        "email": email,
        "deleted_at": nil,
    }
    
    var user models.User
    err := r.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, err
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) FindByID(id primitive.ObjectID) (*models.User, error) {
    ctx := context.Background()
    filter := bson.M{
        "_id": id,
        "deleted_at": nil,
    }
    
    var user models.User
    err := r.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) FindByYandexID(yandexID string) (*models.User, error) {
    ctx := context.Background()
    filter := bson.M{
        "yandex_id": yandexID,
        "deleted_at": nil,
    }
    
    var user models.User
    err := r.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
    ctx := context.Background()
    user.UpdatedAt = time.Now()
    
    filter := bson.M{"_id": user.ID}
    update := bson.M{"$set": user}
    
    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}

// SoftDelete мягко удаляет пользователя
func (r *UserRepository) SoftDelete(id primitive.ObjectID) error {
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