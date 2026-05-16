package middleware

import (
    "context"
    "net/http"
    
    "newyear-app/internal/service"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuthMiddleware struct {
    authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
    return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("access_token")
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        userID, jti, err := m.authService.ValidateAccessToken(cookie.Value)
        if err != nil {
            http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
            return
        }
        
        _ = jti
        
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next(w, r.WithContext(ctx))
    }
}

func GetUserIDFromContext(ctx context.Context) (primitive.ObjectID, bool) {
    userID, ok := ctx.Value(UserIDKey).(primitive.ObjectID)
    return userID, ok
}