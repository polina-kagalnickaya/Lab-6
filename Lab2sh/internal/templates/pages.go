package templates

import (
	"fmt"
    "html/template"
    "strings"
)

// BasePageData базовая структура для всех страниц
type BasePageData struct {
    Title        string
    IsAuthenticated bool
    UserName     string
    DaysLeft     int
}

// LoginPageData данные для страницы логина
type LoginPageData struct {
    BasePageData
    Error string
}

// RegisterPageData данные для страницы регистрации
type RegisterPageData struct {
    BasePageData
    Error string
}

// HomePageData данные для главной страницы
type HomePageData struct {
    BasePageData
    WishesHTML   template.HTML
}

// PublicWishesPageData данные для страницы публичных желаний
type PublicWishesPageData struct {
    BasePageData
    WishesHTML   template.HTML
}

// ProfilePageData данные для страницы профиля
type ProfilePageData struct {
    BasePageData
    UserEmail    string
    UserFullName string
    UserCreatedAt string
    WishesCount  int
}

// ForgotPasswordPageData данные для страницы восстановления пароля
type ForgotPasswordPageData struct {
    BasePageData
    Message string
    Error   string
}

// ResetPasswordPageData данные для страницы сброса пароля
type ResetPasswordPageData struct {
    BasePageData
    Token   string
    Message string
    Error   string
}

// RenderLoginPage генерирует HTML страницу логина
func RenderLoginPage(data LoginPageData) string {
    tmpl := getBaseTemplate() + `
    <div class="auth-container">
        <div class="auth-card">
            <div class="auth-header">
                <h1>🎄 Вход в аккаунт</h1>
                <p>Войдите, чтобы управлять своими желаниями</p>
            </div>
            {{if .Error}}
            <div class="alert alert-error">
                <span class="alert-icon">⚠️</span>
                <span>{{.Error}}</span>
            </div>
            {{end}}
            <form action="/auth/login" method="POST" class="auth-form" id="loginForm">
                <div class="form-group">
                    <label for="email">
                        <span class="label-icon">📧</span>
                        Email
                    </label>
                    <input type="email" id="email" name="email" placeholder="your@email.com" required>
                </div>
                <div class="form-group">
                    <label for="password">
                        <span class="label-icon">🔒</span>
                        Пароль
                    </label>
                    <div class="password-input-wrapper">
                        <input type="password" id="password" name="password" placeholder="Введите пароль" required>
                        <button type="button" class="toggle-password" onclick="togglePassword('password')">👁️</button>
                    </div>
                </div>
                <button type="submit" class="btn btn-primary btn-block">Войти</button>
            </form>
            <div class="auth-footer">
                <a href="/auth/oauth/yandex" class="btn btn-yandex btn-block">
                    <span class="btn-icon">🔑</span>
                    Войти через Яндекс
                </a>
                <div class="auth-links">
                    <a href="/register">Создать аккаунт</a>
                    <a href="/forgot-password">Забыли пароль?</a>
                </div>
                <a href="/" class="back-link">← Вернуться на главную</a>
            </div>
        </div>
    </div>
    ` + getFooterTemplate()
    
    t := template.Must(template.New("login").Parse(tmpl))
    var result strings.Builder
    t.Execute(&result, data)
    return result.String()
}

// RenderRegisterPage генерирует HTML страницу регистрации
func RenderRegisterPage(data RegisterPageData) string {
    tmpl := getBaseTemplate() + `
    <div class="auth-container">
        <div class="auth-card">
            <div class="auth-header">
                <h1>✨ Регистрация</h1>
                <p>Создайте аккаунт и начните загадывать желания</p>
            </div>
            {{if .Error}}
            <div class="alert alert-error">
                <span class="alert-icon">⚠️</span>
                <span>{{.Error}}</span>
            </div>
            {{end}}
            <form action="/auth/register" method="POST" class="auth-form" id="registerForm">
                <div class="form-group">
                    <label for="full_name">
                        <span class="label-icon">👤</span>
                        Имя
                    </label>
                    <input type="text" id="full_name" name="full_name" placeholder="Ваше имя" required>
                </div>
                <div class="form-group">
                    <label for="email">
                        <span class="label-icon">📧</span>
                        Email
                    </label>
                    <input type="email" id="email" name="email" placeholder="your@email.com" required>
                </div>
                <div class="form-group">
                    <label for="password">
                        <span class="label-icon">🔒</span>
                        Пароль
                    </label>
                    <div class="password-input-wrapper">
                        <input type="password" id="password" name="password" placeholder="Минимум 8 символов" required minlength="8">
                        <button type="button" class="toggle-password" onclick="togglePassword('password')">👁️</button>
                    </div>
                </div>
                <div class="form-group">
                    <label for="confirm_password">
                        <span class="label-icon">🔒</span>
                        Подтвердите пароль
                    </label>
                    <div class="password-input-wrapper">
                        <input type="password" id="confirm_password" name="confirm_password" placeholder="Повторите пароль" required minlength="8">
                        <button type="button" class="toggle-password" onclick="togglePassword('confirm_password')">👁️</button>
                    </div>
                </div>
                <button type="submit" class="btn btn-primary btn-block">Создать аккаунт</button>
            </form>
            <div class="auth-footer">
                <a href="/auth/oauth/yandex" class="btn btn-yandex btn-block">
                    <span class="btn-icon">🔑</span>
                    Зарегистрироваться через Яндекс
                </a>
                <div class="auth-links">
                    <a href="/login">Уже есть аккаунт? Войти</a>
                </div>
                <a href="/" class="back-link">← Вернуться на главную</a>
            </div>
        </div>
    </div>
    ` + getFooterTemplate()
    
    t := template.Must(template.New("register").Parse(tmpl))
    var result strings.Builder
    t.Execute(&result, data)
    return result.String()
}

// RenderHomePage генерирует HTML главной страницы
func RenderHomePage(data HomePageData) string {
    tmpl := getBaseTemplate() + `
    <div class="container">
        <!-- Хедер -->
        <header class="main-header">
            <div class="header-content">
                <div class="logo">
                    <h1>🎄 Новогодние желания 2026</h1>
                    <p class="subtitle">Загадайте желание, и оно обязательно сбудется!</p>
                </div>
                <div class="new-year-counter">
                    <div class="counter-number">{{.DaysLeft}}</div>
                    <div class="counter-label">дней до Нового года</div>
                </div>
            </div>
            <nav class="main-nav">
                <a href="/" class="nav-link active">Главная</a>
                <a href="/public" class="nav-link">Публичные желания</a>
                {{if .IsAuthenticated}}
                <a href="/profile" class="nav-link">Мой профиль</a>
                <a href="/logout" class="nav-link nav-link-danger">Выйти</a>
                {{else}}
                <a href="/login" class="nav-link">Войти</a>
                <a href="/register" class="nav-link nav-link-accent">Регистрация</a>
                {{end}}
            </nav>
        </header>

        <!-- Основной контент -->
        <main>
            {{if .IsAuthenticated}}
            <div class="user-welcome">
                <h2>👋 Добро пожаловать, {{.UserName}}!</h2>
                <p>Здесь вы можете создавать и управлять своими желаниями</p>
            </div>
            
            <!-- Форма добавления желания -->
            <div class="card">
                <div class="card-header">
                    <h3>✨ Добавить новое желание</h3>
                </div>
                <div class="card-body">
                    <form id="addWishForm" class="wish-form">
                        <div class="form-group">
                            <label for="text">Текст желания *</label>
                            <textarea id="text" name="text" rows="3" placeholder="Опишите ваше желание..." required minlength="3" maxlength="500"></textarea>
                            <div class="char-count">
                                <span id="charCount">0</span>/500
                            </div>
                        </div>
                        <div class="form-row">
                            <div class="form-group">
                                <label for="author">Автор</label>
                                <input type="text" id="author" name="author" placeholder="Кто загадал?" maxlength="100">
                            </div>
                            <div class="form-group">
                                <label for="priority">Приоритет</label>
                                <select id="priority" name="priority">
                                    <option value="1">⭐ Низкий</option>
                                    <option value="2">⭐⭐ Средний</option>
                                    <option value="3" selected>⭐⭐⭐ Хороший</option>
                                    <option value="4">⭐⭐⭐⭐ Высокий</option>
                                    <option value="5">⭐⭐⭐⭐⭐ Очень важно!</option>
                                </select>
                            </div>
                        </div>
                        <button type="submit" class="btn btn-primary">
                            <span>➕</span> Добавить желание
                        </button>
                    </form>
                </div>
            </div>
            {{else}}
            <div class="hero">
                <div class="hero-content">
                    <h2>🎅 Добро пожаловать в мир новогодних желаний!</h2>
                    <p>Войдите или зарегистрируйтесь, чтобы создавать свои желания и следить за их исполнением</p>
                    <div class="hero-actions">
                        <a href="/login" class="btn btn-primary">Войти</a>
                        <a href="/register" class="btn btn-secondary">Зарегистрироваться</a>
                    </div>
                </div>
            </div>
            {{end}}

            <!-- Список желаний -->
            <div class="card">
                <div class="card-header">
                    <h3>📋 {{if .IsAuthenticated}}Мои желания{{else}}Популярные желания{{end}}</h3>
                </div>
                <div class="card-body">
                    <div id="wishesContainer">
                        {{.WishesHTML}}
                    </div>
                </div>
            </div>
        </main>
    </div>
    ` + getFooterTemplate()
    
    t := template.Must(template.New("home").Parse(tmpl))
    var result strings.Builder
    t.Execute(&result, data)
    return result.String()
}

// RenderPublicWishesPage генерирует HTML страницы публичных желаний
func RenderPublicWishesPage(data PublicWishesPageData) string {
    tmpl := getBaseTemplate() + `
    <div class="container">
        <header class="main-header">
            <div class="header-content">
                <div class="logo">
                    <h1>🎄 Публичные желания</h1>
                    <p class="subtitle">Что загадывают другие люди</p>
                </div>
                <div class="new-year-counter">
                    <div class="counter-number">{{.DaysLeft}}</div>
                    <div class="counter-label">дней до Нового года</div>
                </div>
            </div>
            <nav class="main-nav">
                <a href="/" class="nav-link">Главная</a>
                <a href="/public" class="nav-link active">Публичные желания</a>
                {{if .IsAuthenticated}}
                <a href="/profile" class="nav-link">Мой профиль</a>
                <a href="/logout" class="nav-link nav-link-danger">Выйти</a>
                {{else}}
                <a href="/login" class="nav-link">Войти</a>
                <a href="/register" class="nav-link nav-link-accent">Регистрация</a>
                {{end}}
            </nav>
        </header>

        <main>
            <div class="public-header">
                <h2>🌍 Желания со всего мира</h2>
                <p>Вдохновляйтесь желаниями других людей</p>
            </div>
            
            <div class="wishes-grid">
                {{.WishesHTML}}
            </div>
        </main>
    </div>
    ` + getFooterTemplate()
    
    t := template.Must(template.New("public").Parse(tmpl))
    var result strings.Builder
    t.Execute(&result, data)
    return result.String()
}

// RenderProfilePage генерирует HTML страницы профиля
func RenderProfilePage(data ProfilePageData) string {
    tmpl := getBaseTemplate() + `
    <div class="container">
        <header class="main-header">
            <div class="header-content">
                <div class="logo">
                    <h1>🎄 </h1>
                    <p class="subtitle">Управление аккаунтом</p>
                </div>
                <div class="new-year-counter">
                    <div class="counter-number">{{.DaysLeft}}</div>
                    <div class="counter-label">дней до Нового года</div>
                </div>
            </div>
            <nav class="main-nav">
                <a href="/" class="nav-link">Главная</a>
                <a href="/public" class="nav-link">Публичные желания</a>

                <a href="/logout" class="nav-link nav-link-danger">Выйти</a>
            </nav>
        </header>

        <main>
            <div class="profile-container">
                <!-- Информация о пользователе -->
                <div class="card">
                    <div class="card-header">
                        <h3>👤 Информация профиля</h3>
                    </div>
                    <div class="card-body">
                        <div class="profile-info">
                            <div class="profile-avatar">
                                {{.UserFullName | avatarFirstChar}}
                            </div>
                            <div class="profile-details">
                                <div class="detail-item">
                                    <span class="detail-label">Имя</span>
                                    <span class="detail-value">{{.UserFullName}}</span>
                                </div>
                                <div class="detail-item">
                                    <span class="detail-label">Email</span>
                                    <span class="detail-value">{{.UserEmail}}</span>
                                </div>
                                <div class="detail-item">
                                    <span class="detail-label">Дата регистрации</span>
                                    <span class="detail-value">{{.UserCreatedAt}}</span>
                                </div>
                                <div class="detail-item">
                                    <span class="detail-label">Всего желаний</span>
                                    <span class="detail-value">{{.WishesCount}}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Действия с аккаунтом -->
                <div class="card">
                    <div class="card-header">
                        <h3>⚙️ Управление аккаунтом</h3>
                    </div>
                    <div class="card-body">
                        <div class="account-actions">
                            <a href="/forgot-password" class="btn btn-secondary">
                                <span>🔑</span> Сменить пароль
                            </a>
                            <button onclick="logoutAll()" class="btn btn-danger">
                                <span>🚪</span> Выйти со всех устройств
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Статистика -->
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-icon"></div>
                        <div class="stat-number">{{.WishesCount}}</div>
                        <div class="stat-label">Желаний создано</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-icon"></div>
                        <div class="stat-number">{{.DaysLeft}}</div>
                        <div class="stat-label">Дней до Нового года</div>
                    </div>
                </div>
            </div>
        </main>
    </div>
    ` + getFooterTemplate()
    
    t := template.Must(template.New("profile").Parse(tmpl))
    var result strings.Builder
    t.Execute(&result, data)
    return result.String()
}

// RenderForgotPasswordPage генерирует HTML страницы восстановления пароля
func RenderForgotPasswordPage(data ForgotPasswordPageData) string {
    tmpl := getBaseTemplate() + `
    <div class="auth-container">
        <div class="auth-card">
            <div class="auth-header">
                <h1> Восстановление пароля</h1>
                <p>Введите email, и мы отправим вам инструкцию</p>
            </div>
            {{if .Message}}
            <div class="alert alert-success">
                <span class="alert-icon"></span>
                <span>{{.Message}}</span>
            </div>
            {{end}}
            {{if .Error}}
            <div class="alert alert-error">
                <span class="alert-icon"></span>
                <span>{{.Error}}</span>
            </div>
            {{end}}
            <form action="/auth/forgot-password" method="POST" class="auth-form">
                <div class="form-group">
                    <label for="email">
                        <span class="label-icon"></span>
                        Email
                    </label>
                    <input type="email" id="email" name="email" placeholder="your@email.com" required>
                </div>
                <button type="submit" class="btn btn-primary btn-block">Отправить инструкцию</button>
            </form>
            <div class="auth-footer">
                <a href="/login" class="back-link">← Вернуться ко входу</a>
            </div>
        </div>
    </div>
    ` + getFooterTemplate()
    
    t := template.Must(template.New("forgot").Parse(tmpl))
    var result strings.Builder
    t.Execute(&result, data)
    return result.String()
}

// RenderResetPasswordPage генерирует HTML страницы сброса пароля
func RenderResetPasswordPage(data ResetPasswordPageData) string {
    tmpl := getBaseTemplate() + `
    <div class="auth-container">
        <div class="auth-card">
            <div class="auth-header">
                <h1> Сброс пароля</h1>
                <p>Введите новый пароль</p>
            </div>
            {{if .Message}}
            <div class="alert alert-success">
                <span class="alert-icon"></span>
                <span>{{.Message}}</span>
            </div>
            {{end}}
            {{if .Error}}
            <div class="alert alert-error">
                <span class="alert-icon"></span>
                <span>{{.Error}}</span>
            </div>
            {{end}}
            <form action="/auth/reset-password" method="POST" class="auth-form">
                <input type="hidden" name="token" value="{{.Token}}">
                <div class="form-group">
                    <label for="password">
                        <span class="label-icon"></span>
                        Новый пароль
                    </label>
                    <div class="password-input-wrapper">
                        <input type="password" id="password" name="password" placeholder="Минимум 8 символов" required minlength="8">
                        <button type="button" class="toggle-password" onclick="togglePassword('password')">👁️</button>
                    </div>
                </div>
                <button type="submit" class="btn btn-primary btn-block">Сбросить пароль</button>
            </form>
            <div class="auth-footer">
                <a href="/login" class="back-link">← Вернуться ко входу</a>
            </div>
        </div>
    </div>
    ` + getFooterTemplate()
    
    t := template.Must(template.New("reset").Parse(tmpl))
    var result strings.Builder
    t.Execute(&result, data)
    return result.String()
}

// getBaseTemplate возвращает базовую структуру HTML
func getBaseTemplate() string {
    return `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Новогодние желания 2026</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        /* Хедер */
        .main-header {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            padding: 25px 30px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }
        
        .header-content {
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 20px;
            margin-bottom: 20px;
        }
        
        .logo h1 {
            font-size: 28px;
            color: #764ba2;
        }
        
        .subtitle {
            color: #666;
            font-size: 14px;
            margin-top: 5px;
        }
        
        .new-year-counter {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 15px 25px;
            border-radius: 15px;
            text-align: center;
        }
        
        .counter-number {
            font-size: 36px;
            font-weight: bold;
            line-height: 1;
        }
        
        .counter-label {
            font-size: 14px;
            margin-top: 5px;
        }
        
        /* Навигация */
        .main-nav {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }
        
        .nav-link {
            padding: 10px 20px;
            border-radius: 25px;
            text-decoration: none;
            color: #666;
            font-weight: 500;
            transition: all 0.3s;
        }
        
        .nav-link:hover {
            background: #f0f0f0;
        }
        
        .nav-link.active {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        
        .nav-link-accent {
            background: #ffd700;
            color: #333 !important;
            font-weight: 600;
        }
        
        .nav-link-danger {
            color: #ff4757 !important;
        }
        
        /* Карточки */
        .card {
            background: white;
            border-radius: 20px;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
            margin-bottom: 30px;
            overflow: hidden;
        }
        
        .card-header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px 25px;
        }
        
        .card-header h3 {
            font-size: 20px;
        }
        
        .card-body {
            padding: 25px;
        }
        
        /* Формы */
        .form-group {
            margin-bottom: 20px;
        }
        
        .form-group label {
            display: block;
            margin-bottom: 8px;
            color: #555;
            font-weight: 600;
        }
        
        .form-group input,
        .form-group textarea,
        .form-group select {
            width: 100%;
            padding: 12px 15px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 14px;
            transition: border-color 0.3s;
            font-family: inherit;
        }
        
        .form-group input:focus,
        .form-group textarea:focus,
        .form-group select:focus {
            outline: none;
            border-color: #764ba2;
        }
        
        .form-group textarea {
            resize: vertical;
            min-height: 80px;
        }
        
        .form-row {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
        }
        
        .password-input-wrapper {
            position: relative;
        }
        
        .toggle-password {
            position: absolute;
            right: 10px;
            top: 50%;
            transform: translateY(-50%);
            background: none;
            border: none;
            cursor: pointer;
            font-size: 20px;
        }
        
        .char-count {
            text-align: right;
            color: #999;
            font-size: 12px;
            margin-top: 5px;
        }
        
        /* Кнопки */
        .btn {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 12px 25px;
            border-radius: 10px;
            border: none;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
            text-decoration: none;
        }
        
        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        
        .btn-secondary {
            background: #f0f0f0;
            color: #333;
        }
        
        .btn-secondary:hover {
            background: #e0e0e0;
        }
        
        .btn-danger {
            background: #ff4757;
            color: white;
        }
        
        .btn-danger:hover {
            background: #ff3838;
        }
        
        .btn-yandex {
            background: #ffd700;
            color: #333;
        }
        
        .btn-yandex:hover {
            background: #ffc107;
        }
        
        .btn-block {
            display: flex;
            width: 100%;
            justify-content: center;
        }
        
        /* Алерты */
        .alert {
            padding: 15px;
            border-radius: 10px;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .alert-error {
            background: #ffe0e0;
            color: #cc0000;
            border: 1px solid #ffcccc;
        }
        
        .alert-success {
            background: #e0ffe0;
            color: #006600;
            border: 1px solid #ccffcc;
        }
        
        .alert-icon {
            font-size: 20px;
        }
        
        /* Желания */
        .wish-item {
            background: #f8f9fa;
            border-radius: 15px;
            padding: 20px;
            margin-bottom: 15px;
            border-left: 5px solid #667eea;
            transition: all 0.3s;
        }
        
        .wish-item:hover {
            transform: translateX(5px);
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
        }
        
        .wish-header {
            display: flex;
            justify-content: space-between;
            align-items: start;
            margin-bottom: 10px;
        }
        
        .wish-text {
            font-size: 16px;
            line-height: 1.5;
            color: #333;
            flex: 1;
        }
        
        .wish-actions {
            display: flex;
            gap: 10px;
            margin-left: 15px;
        }
        
        .btn-icon {
            background: none;
            border: none;
            cursor: pointer;
            font-size: 18px;
            padding: 5px;
            border-radius: 5px;
            transition: all 0.3s;
        }
        
        .btn-icon:hover {
            background: #e0e0e0;
        }
        
        .wish-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 10px;
        }
        
        .wish-author {
            color: #667eea;
            font-weight: 600;
            font-size: 14px;
        }
        
        .wish-priority {
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: bold;
        }
        
        .priority-1 { background: #70a1ff; color: white; }
        .priority-2 { background: #7bed9f; color: #333; }
        .priority-3 { background: #ffd32a; color: #333; }
        .priority-4 { background: #ffa502; color: white; }
        .priority-5 { background: #ff4757; color: white; }
        
        .wish-date {
            color: #999;
            font-size: 12px;
        }
        
        .wish-user {
            font-style: italic;
            color: #666;
        }
        
        /* Аутентификация */
        .auth-container {
            max-width: 450px;
            margin: 40px auto;
            padding: 20px;
        }
        
        .auth-card {
            background: white;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
        }
        
        .auth-header {
            text-align: center;
            margin-bottom: 30px;
        }
        
        .auth-header h1 {
            font-size: 28px;
            color: #764ba2;
            margin-bottom: 10px;
        }
        
        .auth-header p {
            color: #666;
        }
        
        .auth-footer {
            margin-top: 25px;
            text-align: center;
        }
        
        .auth-links {
            display: flex;
            justify-content: space-between;
            margin: 20px 0;
        }
        
        .auth-links a {
            color: #667eea;
            text-decoration: none;
            font-size: 14px;
        }
        
        .auth-links a:hover {
            text-decoration: underline;
        }
        
        .back-link {
            color: #666;
            text-decoration: none;
            font-size: 14px;
            display: inline-block;
            margin-top: 15px;
        }
        
        .back-link:hover {
            color: #333;
        }
        
        /* Профиль */
        .profile-container {
            display: grid;
            gap: 30px;
        }
        
        .profile-info {
            display: flex;
            gap: 30px;
            align-items: center;
        }
        
        .profile-avatar {
            width: 100px;
            height: 100px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 40px;
            font-weight: bold;
        }
        
        .profile-details {
            flex: 1;
        }
        
        .detail-item {
            display: flex;
            justify-content: space-between;
            padding: 12px 0;
            border-bottom: 1px solid #f0f0f0;
        }
        
        .detail-label {
            color: #999;
        }
        
        .detail-value {
            font-weight: 600;
        }
        
        .account-actions {
            display: flex;
            gap: 15px;
            flex-wrap: wrap;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
        }
        
        .stat-card {
            background: white;
            border-radius: 15px;
            padding: 25px;
            text-align: center;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
        }
        
        .stat-icon {
            font-size: 40px;
            margin-bottom: 10px;
        }
        
        .stat-number {
            font-size: 32px;
            font-weight: bold;
            color: #764ba2;
        }
        
        .stat-label {
            color: #666;
            margin-top: 5px;
        }
        
        /* Hero секция */
        .hero {
            background: white;
            border-radius: 20px;
            padding: 60px 40px;
            text-align: center;
            margin-bottom: 30px;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
        }
        
        .hero-content h2 {
            font-size: 32px;
            color: #764ba2;
            margin-bottom: 15px;
        }
        
        .hero-content p {
            color: #666;
            font-size: 18px;
            margin-bottom: 30px;
        }
        
        .hero-actions {
            display: flex;
            gap: 15px;
            justify-content: center;
        }
        
        .user-welcome {
            background: white;
            border-radius: 20px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
        }
        
        .user-welcome h2 {
            color: #764ba2;
            margin-bottom: 10px;
        }
        
        .public-header {
            text-align: center;
            margin-bottom: 30px;
            color: white;
        }
        
        .public-header h2 {
            font-size: 28px;
            margin-bottom: 10px;
        }
        
        .wishes-grid {
            display: grid;
            gap: 20px;
        }
        
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #999;
            font-size: 18px;
        }
        
        .empty-state .empty-icon {
            font-size: 60px;
            margin-bottom: 20px;
        }
        
        /* Адаптивность */
        @media (max-width: 768px) {
            .header-content {
                flex-direction: column;
                text-align: center;
            }
            
            .main-nav {
                justify-content: center;
            }
            
            .form-row {
                grid-template-columns: 1fr;
            }
            
            .profile-info {
                flex-direction: column;
                text-align: center;
            }
            
            .wish-header {
                flex-direction: column;
            }
            
            .wish-actions {
                margin-left: 0;
                margin-top: 10px;
            }
            
            .account-actions {
                flex-direction: column;
            }
            
            .hero-actions {
                flex-direction: column;
            }
        }
    </style>
</head>
<body>
`
}

// getFooterTemplate возвращает закрывающие теги и скрипты
func getFooterTemplate() string {
    return `
    <script>
        // Переключение видимости пароля
        function togglePassword(fieldId) {
            const field = document.getElementById(fieldId);
            field.type = field.type === 'password' ? 'text' : 'password';
        }
        
        // Подсчет символов в textarea
        const textArea = document.getElementById('text');
        const charCount = document.getElementById('charCount');
        if (textArea && charCount) {
            textArea.addEventListener('input', () => {
                charCount.textContent = textArea.value.length;
            });
        }
        
        // Добавление желания через AJAX
        const addWishForm = document.getElementById('addWishForm');
        if (addWishForm) {
            addWishForm.addEventListener('submit', async (e) => {
                e.preventDefault();
                
                const wish = {
                    text: document.getElementById('text').value,
                    author: document.getElementById('author').value,
                    priority: parseInt(document.getElementById('priority').value)
                };
                
                try {
                    const response = await fetch('/wishes', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(wish)
                    });
                    
                    if (response.ok) {
                        window.location.reload();
                    } else {
                        const error = await response.text();
                        alert('Ошибка: ' + error);
                    }
                } catch (error) {
                    alert('Ошибка сети');
                }
            });
        }
        
        
        // Редактирование желания
        async function editWish(id, text, author, priority) {
            const newText = prompt('Измените текст желания:', text);
            if (!newText) return;
            
            const newAuthor = prompt('Автор:', author);
            const newPriority = prompt('Приоритет (1-5):', priority);
            
            try {
                const response = await fetch('/wishes?id=' + id, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        text: newText,
                        author: newAuthor,
                        priority: parseInt(newPriority) || priority
                    })
                });
                
                if (response.ok) {
                    window.location.reload();
                } else {
                    const error = await response.text();
                    alert('Ошибка: ' + error);
                }
            } catch (error) {
                alert('Ошибка сети');
            }
        }
        
        // Удаление желания
        // Удаление желания
        async function deleteWish(id) {
            if (!confirm('Вы уверены, что хотите удалить это желание?')) return;
            
            try {
                const response = await fetch('/wishes?id=' + id, {
                    method: 'DELETE'
                });
                
                if (response.ok) {
                    window.location.reload();
                } else {
                    const error = await response.text();
                    alert('Ошибка: ' + error);
                }
            } catch (error) {
                alert('Ошибка сети');
            }
        }
        
        // Выход со всех устройств
        async function logoutAll() {
            if (!confirm('Вы уверены? Вы выйдете со всех устройств.')) return;
            
            try {
                const response = await fetch('/auth/logout-all', {
                    method: 'POST'
                });
                
                if (response.ok) {
                    window.location.href = '/login';
                } else {
                    alert('Ошибка при выходе');
                }
            } catch (error) {
                alert('Ошибка сети');
            }
        }
    </script>
</body>
</html>
`
}

// GenerateWishesHTML генерирует HTML для списка желаний
func GenerateWishesHTML(wishes []WishCardData, isOwner bool) template.HTML {
    if len(wishes) == 0 {
        return template.HTML(`
            <div class="empty-state">
                <div class="empty-icon"></div>
                <p>Пока нет желаний. Будьте первым!</p>
            </div>
        `)
    }
    
    result := ""
    for _, wish := range wishes {
        priorityText := ""
        switch wish.Priority {
        case 1: priorityText = "⭐ Низкий"
        case 2: priorityText = "⭐⭐ Средний"
        case 3: priorityText = "⭐⭐⭐ Хороший"
        case 4: priorityText = "⭐⭐⭐⭐ Высокий"
        case 5: priorityText = "⭐⭐⭐⭐⭐ Очень важно!"
        }
        
        ownerActions := ""
        if isOwner {
            ownerActions = fmt.Sprintf(`
                <div class="wish-actions">
                    <button class="btn-icon" onclick="editWish('%s', '%s', '%s', %d)" title="Редактировать">✏️</button>
                    <button class="btn-icon" onclick="deleteWish('%s')" title="Удалить">🗑️</button>
                </div>
            `, wish.ID, escapeJS(wish.Text), escapeJS(wish.Author), wish.Priority, wish.ID)
        }
        
        userInfo := ""
        if wish.UserName != "" {
            userInfo = fmt.Sprintf(`<div class="wish-user">👤 %s</div>`, escapeHTML(wish.UserName))
        }
        
        result += fmt.Sprintf(`
            <div class="wish-item">
                <div class="wish-header">
                    <div class="wish-text">%s</div>
                    %s
                </div>
                <div class="wish-meta">
                    <div class="wish-author">👤 %s</div>
                    <div class="wish-priority priority-%d">%s</div>
                    <div class="wish-date"> %s</div>
                    %s
                </div>
            </div>
        `, escapeHTML(wish.Text), ownerActions, escapeHTML(wish.Author), wish.Priority, priorityText, wish.Date, userInfo)
    }
    
    return template.HTML(result)
}

// WishCardData данные для отображения желания
// WishCardData данные для отображения желания
type WishCardData struct {
    ID       string // Изменено с uint на string для ObjectID Hex
    Text     string
    Author   string
    Priority int
    Date     string
    UserName string
}

// escapeHTML экранирует HTML символы
func escapeHTML(s string) string {
    s = strings.ReplaceAll(s, "&", "&amp;")
    s = strings.ReplaceAll(s, "<", "&lt;")
    s = strings.ReplaceAll(s, ">", "&gt;")
    s = strings.ReplaceAll(s, "\"", "&quot;")
    s = strings.ReplaceAll(s, "'", "&#39;")
    return s
}

// escapeJS экранирует строку для JavaScript
func escapeJS(s string) string {
    s = strings.ReplaceAll(s, "\\", "\\\\")
    s = strings.ReplaceAll(s, "'", "\\'")
    s = strings.ReplaceAll(s, "\"", "\\\"")
    s = strings.ReplaceAll(s, "\n", "\\n")
    return s
}