package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    ServerPort string

    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string

    RedisHost     string
    RedisPort     string
    RedisPassword string
    RedisDB       int

    JWTSecret       string
    AccessTokenTTL  time.Duration
    RefreshTokenTTL time.Duration

    WishesCacheTTL  time.Duration
    ProfileCacheTTL time.Duration

    YandexClientID     string
    YandexClientSecret string
    YandexCallbackURL  string
    VKClientID         string
    VKClientSecret     string
    VKCallbackURL      string
}

func Load() (*Config, error) {
    return &Config{
        ServerPort: getEnv("SERVER_PORT", "4200"),

        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "newyear_app"),

        RedisHost:     getEnv("REDIS_HOST", "localhost"),
        RedisPort:     getEnv("REDIS_PORT", "6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        RedisDB:       getEnvAsInt("REDIS_DB", 0),

        JWTSecret:       getEnv("JWT_ACCESS_SECRET", "change-me"),
        AccessTokenTTL:  getEnvAsDuration("JWT_ACCESS_EXPIRATION", 15*time.Minute),
        RefreshTokenTTL: getEnvAsDuration("JWT_REFRESH_EXPIRATION", 168*time.Hour),

        WishesCacheTTL:  getEnvAsDuration("CACHE_WISHES_TTL", 5*time.Minute),
        ProfileCacheTTL: getEnvAsDuration("CACHE_PROFILE_TTL", 10*time.Minute),

        YandexClientID:     getEnv("YANDEX_CLIENT_ID", ""),
        YandexClientSecret: getEnv("YANDEX_CLIENT_SECRET", ""),
        YandexCallbackURL:  getEnv("YANDEX_CALLBACK_URL", "http://localhost:4200/auth/oauth/yandex/callback"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if v := os.Getenv(key); v != "" {
        if i, err := strconv.Atoi(v); err == nil {
            return i
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if v := os.Getenv(key); v != "" {
        if d, err := time.ParseDuration(v); err == nil {
            return d
        }
    }
    return defaultValue
}