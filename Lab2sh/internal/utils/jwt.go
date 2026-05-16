package utils

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type JWTClaims struct {
    UserID string `json:"user_id"` // ObjectID как строка
    JTI    string `json:"jti"`
    jwt.RegisteredClaims
}

func GenerateAccessToken(userID string, secret string, expiration time.Duration) (string, string, error) {
    jti := uuid.New().String()
    
    claims := JWTClaims{
        UserID: userID,
        JTI:    jti,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ID:        jti,
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(secret))
    return tokenString, jti, err
}

func ValidateAccessToken(tokenString, secret string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}

func GenerateRefreshToken() (string, string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", "", err
    }
    token := base64.URLEncoding.EncodeToString(bytes)
    jti := uuid.New().String()
    return token, jti, nil
}

func HashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return base64.URLEncoding.EncodeToString(hash[:])
}