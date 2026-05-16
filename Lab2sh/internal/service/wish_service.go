package service

import (
    "fmt"
    "time"
    
    "newyear-app/internal/cache"
    "newyear-app/internal/dto"
    "newyear-app/internal/models"
    "newyear-app/internal/repository"
    "newyear-app/internal/utils"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type WishService struct {
    repo       *repository.WishRepository
    redisCache *cache.RedisCache
    cacheTTL   time.Duration
}

func NewWishService(repo *repository.WishRepository, redisCache *cache.RedisCache, cacheTTL time.Duration) *WishService {
    return &WishService{
        repo:       repo,
        redisCache: redisCache,
        cacheTTL:   cacheTTL,
    }
}

func (s *WishService) Create(req dto.CreateWishRequest, userID primitive.ObjectID) (*dto.WishResponse, error) {
    wish := &models.Wish{
        UserID:   userID,
        Text:     req.Text,
        Priority: req.Priority,
    }

    if req.Author != "" {
        wish.Author = req.Author
    } else {
        wish.Author = "Anonymous"
    }

    if wish.Priority == 0 {
        wish.Priority = 1
    }

    if err := s.repo.Create(wish); err != nil {
        return nil, err
    }

    // Инвалидация кеша
    s.redisCache.InvalidateWishesCache(utils.ObjectIDToUint(userID))

    return s.toResponse(wish), nil
}

func (s *WishService) GetByID(id primitive.ObjectID, userID primitive.ObjectID) (*dto.WishResponse, error) {
    wish, err := s.repo.FindByID(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("wish not found")
        }
        return nil, err
    }

    if wish.UserID != userID {
        return nil, fmt.Errorf("permission denied")
    }

    return s.toResponse(wish), nil
}

func (s *WishService) GetAll(page, limit int, userID primitive.ObjectID) (*dto.PaginatedResponse, error) {
    wishes, total, err := s.repo.FindAll(page, limit, userID)
    if err != nil {
        return nil, err
    }
    
    responses := make([]dto.WishResponse, len(wishes))
    for i, wish := range wishes {
        responses[i] = *s.toResponse(&wish)
    }
    
    totalPages := (total + int64(limit) - 1) / int64(limit)
    result := &dto.PaginatedResponse{
        Data: responses,
        Meta: dto.Meta{
            Total:      total,
            Page:       page,
            Limit:      limit,
            TotalPages: totalPages,
        },
    }
    
    return result, nil
}

func (s *WishService) GetAllPublic(page, limit int) ([]models.Wish, error) {
    wishes, _, err := s.repo.FindAllPublic(page, limit)
    return wishes, err
}

func (s *WishService) Update(id primitive.ObjectID, req dto.UpdateWishRequest, userID primitive.ObjectID) (*dto.WishResponse, error) {
    wish, err := s.repo.FindByID(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("wish not found")
        }
        return nil, err
    }

    if wish.UserID != userID {
        return nil, fmt.Errorf("permission denied")
    }

    if req.Text != "" {
        wish.Text = req.Text
    }
    if req.Author != "" {
        wish.Author = req.Author
    }
    if req.Priority != 0 {
        wish.Priority = req.Priority
    }

    if err := s.repo.Update(wish); err != nil {
        return nil, err
    }

    s.redisCache.InvalidateWishesCache(utils.ObjectIDToUint(userID))

    return s.toResponse(wish), nil
}

func (s *WishService) Delete(id primitive.ObjectID, userID primitive.ObjectID) error {
    wish, err := s.repo.FindByID(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fmt.Errorf("wish not found")
        }
        return err
    }

    if wish.UserID != userID {
        return fmt.Errorf("permission denied")
    }

    if err := s.repo.Delete(id); err != nil {
        return err
    }

    s.redisCache.InvalidateWishesCache(utils.ObjectIDToUint(userID))

    return nil
}

func (s *WishService) toResponse(wish *models.Wish) *dto.WishResponse {
    return &dto.WishResponse{
        ID:        wish.ID,
        Text:      wish.Text,
        Author:    wish.Author,
        Priority:  wish.Priority,
        CreatedAt: wish.CreatedAt,
    }
}