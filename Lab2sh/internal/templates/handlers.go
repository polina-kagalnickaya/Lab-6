package templates

import (
    "fmt"
    "html/template"
    "net/http"
    "time"

    "newyear-app/internal/service"
)

// PageHandlers содержит обработчики для HTML страниц
type PageHandlers struct {
    authService *service.AuthService
    wishService *service.WishService
}

// NewPageHandlers создает новый экземпляр PageHandlers
func NewPageHandlers(authService *service.AuthService, wishService *service.WishService) *PageHandlers {
    return &PageHandlers{
        authService: authService,
        wishService: wishService,
    }
}

// getBasePageData получает базовые данные для страницы
func (h *PageHandlers) getBasePageData(r *http.Request, title string) BasePageData {
    data := BasePageData{
        Title: title,
    }
    
    // Получаем дни до Нового года
    now := time.Now()
    nextYear := now.Year() + 1
    newYear := time.Date(nextYear, time.January, 1, 0, 0, 0, 0, now.Location())
    data.DaysLeft = int(newYear.Sub(now).Hours() / 24)
    
    // Проверяем авторизацию
    if cookie, err := r.Cookie("access_token"); err == nil {
        if userID, _, err := h.authService.ValidateAccessToken(cookie.Value); err == nil {
            data.IsAuthenticated = true
            if user, err := h.authService.GetUserByID(userID); err == nil {
                data.UserName = user.FullName
            }
        }
    }
    
    return data
}

// LoginPageHandler обрабатывает страницу логина
func (h *PageHandlers) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
    // Если пользователь уже авторизован, перенаправляем на главную
    if cookie, err := r.Cookie("access_token"); err == nil {
        if _, _, err := h.authService.ValidateAccessToken(cookie.Value); err == nil {
            http.Redirect(w, r, "/", http.StatusFound)
            return
        }
    }
    
    data := LoginPageData{
        BasePageData: h.getBasePageData(r, "Вход"),
    }
    
    // Проверяем ошибки из query параметров
    if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
        data.Error = errorMsg
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, RenderLoginPage(data))
}

// RegisterPageHandler обрабатывает страницу регистрации
func (h *PageHandlers) RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
    // Если пользователь уже авторизован, перенаправляем на главную
    if cookie, err := r.Cookie("access_token"); err == nil {
        if _, _, err := h.authService.ValidateAccessToken(cookie.Value); err == nil {
            http.Redirect(w, r, "/", http.StatusFound)
            return
        }
    }
    
    data := RegisterPageData{
        BasePageData: h.getBasePageData(r, "Регистрация"),
    }
    
    if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
        data.Error = errorMsg
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, RenderRegisterPage(data))
}

// HomePageHandler обрабатывает главную страницу
func (h *PageHandlers) HomePageHandler(w http.ResponseWriter, r *http.Request) {
    data := HomePageData{
        BasePageData: h.getBasePageData(r, "Главная"),
    }
    
    // Получаем желания
    var wishesHTML template.HTML
    
    if data.IsAuthenticated {
        // Получаем желания пользователя
        if cookie, _ := r.Cookie("access_token"); cookie != nil {
            if userID, _, err := h.authService.ValidateAccessToken(cookie.Value); err == nil {
                if wishes, err := h.wishService.GetAll(1, 50, userID); err == nil {
                    wishCards := make([]WishCardData, len(wishes.Data))
                    for i, w := range wishes.Data {
                        wishCards[i] = WishCardData{
                            ID:       w.ID.Hex(), // Конвертируем ObjectID в строку
                            Text:     w.Text,
                            Author:   w.Author,
                            Priority: w.Priority,
                            Date:     w.CreatedAt.Format("02.01.2006 15:04"),
                        }
                    }
                    wishesHTML = GenerateWishesHTML(wishCards, true)
                }
            }
        }
    } else {
        // Для неавторизованных показываем публичные желания
        if wishes, err := h.wishService.GetAllPublic(1, 10); err == nil {
            wishCards := make([]WishCardData, len(wishes))
            for i, w := range wishes {
                wishCards[i] = WishCardData{
                    ID:       w.ID.Hex(), // Конвертируем ObjectID в строку
                    Text:     w.Text,
                    Author:   w.Author,
                    Priority: w.Priority,
                    Date:     w.CreatedAt.Format("02.01.2006 15:04"),
                }
            }
            wishesHTML = GenerateWishesHTML(wishCards, false)
        }
    }
    
    data.WishesHTML = wishesHTML
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, RenderHomePage(data))
}

// PublicWishesPageHandler обрабатывает страницу публичных желаний
func (h *PageHandlers) PublicWishesPageHandler(w http.ResponseWriter, r *http.Request) {
    data := PublicWishesPageData{
        BasePageData: h.getBasePageData(r, "Публичные желания"),
    }
    
    if wishes, err := h.wishService.GetAllPublic(1, 50); err == nil {
        wishCards := make([]WishCardData, len(wishes))
        for i, w := range wishes {
            wishCards[i] = WishCardData{
                ID:       w.ID.Hex(), // Конвертируем ObjectID в строку
                Text:     w.Text,
                Author:   w.Author,
                Priority: w.Priority,
                Date:     w.CreatedAt.Format("02.01.2006 15:04"),
            }
        }
        data.WishesHTML = GenerateWishesHTML(wishCards, false)
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, RenderPublicWishesPage(data))
}

// ProfilePageHandler обрабатывает страницу профиля
func (h *PageHandlers) ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
    // Проверяем авторизацию
    cookie, err := r.Cookie("access_token")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }
    
    userID, _, err := h.authService.ValidateAccessToken(cookie.Value)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }
    
    data := ProfilePageData{
        BasePageData: h.getBasePageData(r, "Профиль"),
    }
    
    // Получаем данные пользователя
    if user, err := h.authService.GetUserByID(userID); err == nil {
        data.UserEmail = user.Email
        data.UserFullName = user.FullName
        data.UserCreatedAt = user.CreatedAt.Format("02.01.2006")
    }
    
    // Получаем количество желаний
    if wishes, err := h.wishService.GetAll(1, 1, userID); err == nil {
        data.WishesCount = int(wishes.Meta.Total)
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, RenderProfilePage(data))
}

// ForgotPasswordPageHandler обрабатывает страницу восстановления пароля
func (h *PageHandlers) ForgotPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
    data := ForgotPasswordPageData{
        BasePageData: h.getBasePageData(r, "Восстановление пароля"),
    }
    
    if message := r.URL.Query().Get("message"); message != "" {
        data.Message = message
    }
    if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
        data.Error = errorMsg
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, RenderForgotPasswordPage(data))
}

// ResetPasswordPageHandler обрабатывает страницу сброса пароля
func (h *PageHandlers) ResetPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
    data := ResetPasswordPageData{
        BasePageData: h.getBasePageData(r, "Сброс пароля"),
    }
    
    data.Token = r.URL.Query().Get("token")
    
    if message := r.URL.Query().Get("message"); message != "" {
        data.Message = message
    }
    if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
        data.Error = errorMsg
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, RenderResetPasswordPage(data))
}