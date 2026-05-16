package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "newyear-app/internal/cache"
    "newyear-app/internal/handler"
    "newyear-app/internal/middleware"
    "newyear-app/internal/oauth"
    "newyear-app/internal/repository"
    "newyear-app/internal/service"
    "newyear-app/internal/templates"
)

func main() {
    // Загрузка .env
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not loaded, using environment variables")
    }

    cfg := &Config{
        ServerPort: getEnv("SERVER_PORT", "4200"),
        MongoURI:   getEnv("MONGO_URI", "mongodb://admin:admin123@localhost:27017"),
        MongoDB:    getEnv("MONGO_DB", "newyear_app"),
        
        RedisHost:     getEnv("REDIS_HOST", "localhost"),
        RedisPort:     getEnv("REDIS_PORT", "6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        RedisDB:       getEnvAsInt("REDIS_DB", 0),

        JWTSecret:       getEnv("JWT_ACCESS_SECRET", "change-me-in-production"),
        AccessTokenTTL:  getEnvAsDuration("JWT_ACCESS_EXPIRATION", 15*time.Minute),
        RefreshTokenTTL: getEnvAsDuration("JWT_REFRESH_EXPIRATION", 168*time.Hour),

        WishesCacheTTL:  getEnvAsDuration("CACHE_WISHES_TTL", 5*time.Minute),
        ProfileCacheTTL: getEnvAsDuration("CACHE_PROFILE_TTL", 10*time.Minute),

        YandexClientID:     getEnv("YANDEX_CLIENT_ID", ""),
        YandexClientSecret: getEnv("YANDEX_CLIENT_SECRET", ""),
        YandexCallbackURL:  getEnv("YANDEX_CALLBACK_URL", "http://localhost:4200/auth/oauth/yandex/callback"),
    }

    
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    defer client.Disconnect(ctx)
    
    
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }
    log.Println("Connected to MongoDB")
    
    db := client.Database(cfg.MongoDB)

    
    redisCache, err := cache.NewRedisCache(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    defer redisCache.Close()
    log.Println("Connected to Redis")

    
    userRepo := repository.NewUserRepository(db)
    wishRepo := repository.NewWishRepository(db)

    
    authService := service.NewAuthService(
        userRepo,
        redisCache,
        cfg.JWTSecret,
        cfg.AccessTokenTTL,
        cfg.RefreshTokenTTL,
        cfg.ProfileCacheTTL,
    )

    wishService := service.NewWishService(wishRepo, redisCache, cfg.WishesCacheTTL)

    yandexProvider := oauth.NewYandexProvider(
        cfg.YandexClientID,
        cfg.YandexClientSecret,
        cfg.YandexCallbackURL,
    )

    oauthService := service.NewOAuthService(
        userRepo,
        redisCache,
        authService,
        yandexProvider,
    )

    
    authHandler := handler.NewAuthHandler(authService, oauthService)
    wishHandler := handler.NewWishHandler(wishService)
    authMiddleware := middleware.NewAuthMiddleware(authService)

    
    pageHandlers := templates.NewPageHandlers(authService, wishService)
    templates.SetupPageRoutes(http.DefaultServeMux, pageHandlers)

    
    http.HandleFunc("/info", handler.InfoHandler)
    http.HandleFunc("/auth/register", authHandler.Register)
    http.HandleFunc("/auth/login", authHandler.Login)
    http.HandleFunc("/auth/refresh", authHandler.Refresh)
    http.HandleFunc("/auth/forgot-password", authHandler.ForgotPassword)
    http.HandleFunc("/auth/reset-password", authHandler.ResetPassword)
    http.HandleFunc("/auth/oauth/yandex", authHandler.YandexAuth)
    http.HandleFunc("/auth/oauth/yandex/callback", authHandler.YandexCallback)

    
    http.HandleFunc("/auth/logout", authMiddleware.Authenticate(authHandler.Logout))
    http.HandleFunc("/auth/logout-all", authMiddleware.Authenticate(authHandler.LogoutAll))
    http.HandleFunc("/auth/whoami", authMiddleware.Authenticate(authHandler.Whoami))
    http.HandleFunc("/wishes", authMiddleware.Authenticate(wishHandler.WishesHandler))

    
    serverAddr := ":" + cfg.ServerPort
    log.Printf("Server starting on http://localhost%s", serverAddr)

    if err := http.ListenAndServe(serverAddr, nil); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}

type Config struct {
    ServerPort       string
    MongoURI         string
    MongoDB          string
    RedisHost        string
    RedisPort        string
    RedisPassword    string
    RedisDB          int
    JWTSecret        string
    AccessTokenTTL   time.Duration
    RefreshTokenTTL  time.Duration
    WishesCacheTTL   time.Duration
    ProfileCacheTTL  time.Duration
    YandexClientID     string
    YandexClientSecret string
    YandexCallbackURL  string
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}