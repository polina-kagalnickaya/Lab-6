package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisCache(host, port, password string, db int) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", host, port),
        Password: password,
        DB:       db,
    })

    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("redis connection failed: %w", err)
    }

    return &RedisCache{
        client: client,
        ctx:    ctx,
    }, nil
}


func (r *RedisCache) Close() error {
    return r.client.Close()
}


func (r *RedisCache) SaveRefreshToken(jti string, userID uint, ttl time.Duration) error {
    key := fmt.Sprintf("refresh:jti:%s", jti)
    data := map[string]interface{}{
        "user_id":    userID,
        "created_at": time.Now().Unix(),
    }
    if err := r.client.HSet(r.ctx, key, data).Err(); err != nil {
        return err
    }
    return r.client.Expire(r.ctx, key, ttl).Err()
}


func (r *RedisCache) IsRefreshTokenValid(jti string) (bool, error) {
    key := fmt.Sprintf("refresh:jti:%s", jti)
    exists, err := r.client.Exists(r.ctx, key).Result()
    return exists == 1, err
}


func (r *RedisCache) DeleteRefreshToken(jti string) error {
    key := fmt.Sprintf("refresh:jti:%s", jti)
    return r.client.Del(r.ctx, key).Err()
}


func (r *RedisCache) DeleteAllUserRefreshTokens(userID uint) error {
    key := fmt.Sprintf("user:refresh:%d", userID)
    jtis, err := r.client.SMembers(r.ctx, key).Result()
    if err != nil {
        return err
    }
    for _, jti := range jtis {
        r.client.Del(r.ctx, fmt.Sprintf("refresh:jti:%s", jti))
    }
    return r.client.Del(r.ctx, key).Err()
}


func (r *RedisCache) BlacklistAccessToken(jti string, ttl time.Duration) error {
    key := fmt.Sprintf("access:blacklist:%s", jti)
    return r.client.Set(r.ctx, key, "revoked", ttl).Err()
}


func (r *RedisCache) IsAccessTokenRevoked(jti string) (bool, error) {
    key := fmt.Sprintf("access:blacklist:%s", jti)
    exists, err := r.client.Exists(r.ctx, key).Result()
    return exists == 1, err
}


func (r *RedisCache) GetWishesList(userID uint, page, limit int) (interface{}, error) {
    key := fmt.Sprintf("wishes:list:user:%d:page:%d:limit:%d", userID, page, limit)
    data, err := r.client.Get(r.ctx, key).Bytes()
    if err == redis.Nil {
        return nil, nil 
    }
    if err != nil {
        return nil, err
    }
    var result interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    return result, nil
}


func (r *RedisCache) SetWishesList(userID uint, page, limit int, data interface{}, ttl time.Duration) error {
    key := fmt.Sprintf("wishes:list:user:%d:page:%d:limit:%d", userID, page, limit)
    bytes, err := json.Marshal(data)
    if err != nil {
        return err
    }
    return r.client.Set(r.ctx, key, bytes, ttl).Err()
}


func (r *RedisCache) InvalidateWishesCache(userID uint) error {
    pattern := fmt.Sprintf("wishes:list:user:%d:*", userID)
    iter := r.client.Scan(r.ctx, 0, pattern, 0).Iterator()
    for iter.Next(r.ctx) {
        r.client.Del(r.ctx, iter.Val())
    }
    return iter.Err()
}


func (r *RedisCache) GetUserProfile(userID uint) (interface{}, error) {
    key := fmt.Sprintf("user:profile:%d", userID)
    data, err := r.client.Get(r.ctx, key).Bytes()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    var result interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    return result, nil
}

func (r *RedisCache) SetUserProfile(userID uint, data interface{}, ttl time.Duration) error {
    key := fmt.Sprintf("user:profile:%d", userID)
    bytes, err := json.Marshal(data)
    if err != nil {
        return err
    }
    return r.client.Set(r.ctx, key, bytes, ttl).Err()
}

func (r *RedisCache) InvalidateUserProfile(userID uint) error {
    key := fmt.Sprintf("user:profile:%d", userID)
    return r.client.Del(r.ctx, key).Err()
}


func (r *RedisCache) GetClient() *redis.Client {
    return r.client
}


func (r *RedisCache) GetContext() context.Context {
    return r.ctx
}


func (r *RedisCache) SetString(key string, value string, ttl time.Duration) error {
    return r.client.Set(r.ctx, key, value, ttl).Err()
}


func (r *RedisCache) GetString(key string) (string, error) {
    return r.client.Get(r.ctx, key).Result()
}


func (r *RedisCache) Delete(key string) error {
    return r.client.Del(r.ctx, key).Err()
}


func (r *RedisCache) SetAdd(key string, member interface{}) error {
    return r.client.SAdd(r.ctx, key, member).Err()
}


func (r *RedisCache) SetRemove(key string, member interface{}) error {
    return r.client.SRem(r.ctx, key, member).Err()
}


func (r *RedisCache) GetSetMembers(key string) ([]string, error) {
    return r.client.SMembers(r.ctx, key).Result()
}