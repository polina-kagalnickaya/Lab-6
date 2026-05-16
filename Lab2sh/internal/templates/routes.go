package templates

import (
    "net/http"
)

// SetupPageRoutes настраивает маршруты для HTML страниц
func SetupPageRoutes(mux *http.ServeMux, handlers *PageHandlers) {
    // Публичные страницы
    mux.HandleFunc("/", handlers.HomePageHandler)
    mux.HandleFunc("/login", handlers.LoginPageHandler)
    mux.HandleFunc("/register", handlers.RegisterPageHandler)
    mux.HandleFunc("/public", handlers.PublicWishesPageHandler)
    mux.HandleFunc("/forgot-password", handlers.ForgotPasswordPageHandler)
    mux.HandleFunc("/reset-password", handlers.ResetPasswordPageHandler)
    
    // Защищенные страницы
    mux.HandleFunc("/profile", handlers.ProfilePageHandler)
}