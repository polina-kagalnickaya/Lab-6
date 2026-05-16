package service

import (
    "errors"
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

type AuthService struct {
    userRepo   *repository.UserRepository
    redisCache *cache.RedisCache
    jwtSecret  string
    jwtExp     time.Duration
    refreshExp time.Duration
    profileCacheTTL time.Duration
}

func NewAuthService(
    userRepo *repository.UserRepository,
    redisCache *cache.RedisCache,
    jwtSecret string,
    jwtExp time.Duration,
    refreshExp time.Duration,
    profileCacheTTL time.Duration,
) *AuthService {
    return &AuthService{
        userRepo:   userRepo,
        redisCache: redisCache,
        jwtSecret:  jwtSecret,
        jwtExp:     jwtExp,
        refreshExp: refreshExp,
        profileCacheTTL: profileCacheTTL,
    }
}

func (s *AuthService) GetJWTSecret() string     { return s.jwtSecret }
func (s *AuthService) GetJWTExp() time.Duration { return s.jwtExp }
func (s *AuthService) GetRefreshExp() time.Duration { return s.refreshExp }
func (s *AuthService) GetRedisCache() *cache.RedisCache { return s.redisCache }

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
    _, err := s.userRepo.FindByEmail(req.Email)
    if err == nil {
        return nil, errors.New("user already exists")
    }
    
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        Email:    req.Email,
        Password: hashedPassword,
        FullName: req.FullName,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return &dto.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        FullName:  user.FullName,
        CreatedAt: user.CreatedAt,
    }, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (accessToken, refreshToken string, user *dto.UserResponse, err error) {
    userModel, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return "", "", nil, errors.New("invalid credentials")
        }
        return "", "", nil, err
    }
    
    if !utils.CheckPassword(req.Password, userModel.Password) {
        return "", "", nil, errors.New("invalid credentials")
    }
    
    userIDHex := userModel.ID.Hex()
    accessToken, accessJTI, err := utils.GenerateAccessToken(userIDHex, s.jwtSecret, s.jwtExp)
    if err != nil {
        return "", "", nil, err
    }
    
    refreshToken, refreshJTI, err := utils.GenerateRefreshToken()
    if err != nil {
        return "", "", nil, err
    }
    
    // Конвертируем ObjectID в uint для Redis
    uintUserID := utils.ObjectIDToUint(userModel.ID)
    
    if err := s.redisCache.SaveRefreshToken(refreshJTI, uintUserID, s.refreshExp); err != nil {
        return "", "", nil, err
    }
    
    s.redisCache.GetClient().SAdd(s.redisCache.GetContext(), 
        fmt.Sprintf("user:refresh:%s", userIDHex), refreshJTI)
    
    _ = accessJTI
    
    user = &dto.UserResponse{
        ID:        userModel.ID,
        Email:     userModel.Email,
        FullName:  userModel.FullName,
        CreatedAt: userModel.CreatedAt,
    }
    
    return accessToken, refreshToken, user, nil
}

func (s *AuthService) RefreshToken(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {
    tokenHash := utils.HashToken(oldRefreshToken)
    
    key := fmt.Sprintf("refresh:hash:%s", tokenHash)
    userIDHex, err := s.redisCache.GetClient().Get(s.redisCache.GetContext(), key).Result()
    if err != nil {
        return "", "", errors.New("invalid refresh token")
    }
    
    s.redisCache.GetClient().Del(s.redisCache.GetContext(), key)
    
    newAccessToken, accessJTI, err := utils.GenerateAccessToken(userIDHex, s.jwtSecret, s.jwtExp)
    if err != nil {
        return "", "", err
    }
    _ = accessJTI
    
    newRefreshToken, refreshJTI, err := utils.GenerateRefreshToken()
    if err != nil {
        return "", "", err
    }
    
    newHash := utils.HashToken(newRefreshToken)
    newKey := fmt.Sprintf("refresh:hash:%s", newHash)
    
    if err := s.redisCache.GetClient().Set(s.redisCache.GetContext(), newKey, userIDHex, s.refreshExp).Err(); err != nil {
        return "", "", err
    }
    
    s.redisCache.GetClient().SAdd(s.redisCache.GetContext(), 
        fmt.Sprintf("user:refresh:%s", userIDHex), refreshJTI)
    
    return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
    tokenHash := utils.HashToken(refreshToken)
    return s.redisCache.GetClient().Del(s.redisCache.GetContext(), "refresh:hash:"+tokenHash).Err()
}

func (s *AuthService) LogoutAll(userID primitive.ObjectID) error {
    uintUserID := utils.ObjectIDToUint(userID)
    return s.redisCache.DeleteAllUserRefreshTokens(uintUserID)
}

func (s *AuthService) ValidateAccessToken(tokenString string) (primitive.ObjectID, string, error) {
    claims, err := utils.ValidateAccessToken(tokenString, s.jwtSecret)
    if err != nil {
        return primitive.NilObjectID, "", err
    }
    
    revoked, err := s.redisCache.IsAccessTokenRevoked(claims.JTI)
    if err != nil {
        return primitive.NilObjectID, "", err
    }
    if revoked {
        return primitive.NilObjectID, "", errors.New("token revoked")
    }
    
    objID, err := primitive.ObjectIDFromHex(claims.UserID)
    if err != nil {
        return primitive.NilObjectID, "", err
    }
    
    return objID, claims.JTI, nil
}

func (s *AuthService) RevokeAccessToken(jti string, ttl time.Duration) error {
    return s.redisCache.BlacklistAccessToken(jti, ttl)
}

func (s *AuthService) GetUserByID(userID primitive.ObjectID) (*dto.UserResponse, error) {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }
    
    response := &dto.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        FullName:  user.FullName,
        CreatedAt: user.CreatedAt,
    }
    
    return response, nil
}

func (s *AuthService) ForgotPassword(email string) error {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return nil
    }
    
    resetToken, _, err := utils.GenerateRefreshToken()
    if err != nil {
        return err
    }
    
    key := fmt.Sprintf("password_reset:%s", resetToken)
    if err := s.redisCache.SetString(key, user.ID.Hex(), 15*time.Minute); err != nil {
        return err
    }

    fmt.Printf("Password reset token for %s: %s\n", email, resetToken)
    
    return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
    key := fmt.Sprintf("password_reset:%s", token)
    userIDHex, err := s.redisCache.GetString(key)
    if err != nil {
        return errors.New("invalid or expired reset token")
    }
    
    s.redisCache.Delete(key)
    
    objID, err := primitive.ObjectIDFromHex(userIDHex)
    if err != nil {
        return errors.New("invalid user ID")
    }
    
    user, err := s.userRepo.FindByID(objID)
    if err != nil {
        return errors.New("user not found")
    }
    
    hashedPassword, err := utils.HashPassword(newPassword)
    if err != nil {
        return err
    }
    
    user.Password = hashedPassword
    if err := s.userRepo.Update(user); err != nil {
        return err
    }
    
    return s.LogoutAll(objID)
}