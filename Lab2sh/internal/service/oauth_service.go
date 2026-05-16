package service

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "log"
    
    "newyear-app/internal/cache"
    "newyear-app/internal/dto"
    "newyear-app/internal/models"
    "newyear-app/internal/oauth"
    "newyear-app/internal/repository"
    "newyear-app/internal/utils"
)

type OAuthService struct {
    userRepo       *repository.UserRepository
    redisCache     *cache.RedisCache
    authService    *AuthService
    yandexProvider oauth.Provider
}

func NewOAuthService(
    userRepo *repository.UserRepository,
    redisCache *cache.RedisCache,
    authService *AuthService,
    yandexProvider oauth.Provider,
) *OAuthService {
    return &OAuthService{
        userRepo:       userRepo,
        redisCache:     redisCache,
        authService:    authService,
        yandexProvider: yandexProvider,
    }
}

func (s *OAuthService) GenerateState() (string, error) {
    bytes := make([]byte, 32)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *OAuthService) GetYandexAuthURL(state string) string {
    url := s.yandexProvider.GetAuthURL(state)
    log.Printf("Yandex auth URL: %s", url)
    return url
}

func (s *OAuthService) HandleYandexCallback(code, state, storedState string) (accessToken, refreshToken string, user *dto.UserResponse, err error) {
    if state != storedState {
        return "", "", nil, errors.New("invalid state")
    }
    
    yandexToken, err := s.yandexProvider.ExchangeCode(code)
    if err != nil {
        return "", "", nil, err
    }
    
    userInfo, err := s.yandexProvider.GetUserInfo(yandexToken)
    if err != nil {
        return "", "", nil, err
    }
    
    return s.handleOAuthUser(userInfo)
}

func (s *OAuthService) handleOAuthUser(userInfo oauth.UserInfo) (accessToken, refreshToken string, user *dto.UserResponse, err error) {
    var userModel *models.User
    
    userModel, err = s.userRepo.FindByYandexID(userInfo.ID)
    
    if err != nil {
        userModel = &models.User{
            Email:    userInfo.Email,
            FullName: userInfo.FullName,
        }
        yandexID := userInfo.ID
        userModel.YandexID = &yandexID
        
        if err := s.userRepo.Create(userModel); err != nil {
            return "", "", nil, err
        }
    }
    
    // Используем Hex строку ObjectID для JWT
    userIDHex := userModel.ID.Hex()
    accessToken, accessJTI, err := utils.GenerateAccessToken(userIDHex, s.authService.GetJWTSecret(), s.authService.GetJWTExp())
    if err != nil {
        return "", "", nil, err
    }
    _ = accessJTI
    
    refreshToken, refreshJTI, err := utils.GenerateRefreshToken()
    if err != nil {
        return "", "", nil, err
    }
    
    // Конвертируем ObjectID в uint для Redis
    uintUserID := utils.ObjectIDToUint(userModel.ID)
    
    if err := s.redisCache.SaveRefreshToken(refreshJTI, uintUserID, s.authService.GetRefreshExp()); err != nil {
        return "", "", nil, err
    }
    
    user = &dto.UserResponse{
        ID:        userModel.ID,
        Email:     userModel.Email,
        FullName:  userModel.FullName,
        CreatedAt: userModel.CreatedAt,
    }
    
    return accessToken, refreshToken, user, nil
}