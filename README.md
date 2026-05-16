[README.md](https://github.com/user-attachments/files/27844108/README.md)
# Новогодние желания (New Year Wishes App)

## Описание проекта

Веб-приложение для управления списком новогодних желаний. Пользователи могут создавать, просматривать, редактировать и удалять свои желания. Приложение поддерживает аутентификацию через JWT токены и OAuth провайдеров (Яндекс, ВКонтакте).

## Основные возможности

### Аутентификация и авторизация
- Регистрация пользователей (email, пароль, имя)
- Вход по email и паролю
- JWT access и refresh токены с хранением в HttpOnly cookies
- Refresh token rotation
- Выход из аккаунта
- Выход со всех устройств одновременно
- OAuth авторизация через Яндекс
- Восстановление пароля через email с временным токеном в Redis

### Управление желаниями
- Создание желаний (текст, автор, приоритет)
- Просмотр списка своих желаний с пагинацией
- Просмотр публичных желаний других пользователей
- Редактирование желаний
- Удаление желаний (мягкое удаление)
- Приоритизация желаний (от 1 до 5)
- Кеширование списка желаний в Redis

### Пользовательский интерфейс
- Многостраничный адаптивный интерфейс
- Отображение счетчика дней до Нового года
- Страница профиля с информацией о пользователе и статистикой
- Публичная страница со всеми желаниями
- AJAX отправка форм без перезагрузки страниц
- Валидация форм на клиенте и сервере

### Технические особенности
- Кеширование запросов в Redis
- Черный список отозванных access токенов
- Мягкое удаление записей (soft delete) через GORM
- Валидация входных данных
- Docker контейнеризация
- Миграции базы данных через goose

## Технологический стек

- **Backend**: Go
- **База данных**: PostgreSQL 16
- **ORM**: GORM
- **Кеширование**: Redis 7
- **Аутентификация**: JWT (golang-jwt/jwt)
- **Хеширование паролей**: bcrypt (golang.org/x/crypto)
- **Валидация**: go-playground/validator
- **OAuth**: Яндекс ID
- **Миграции**: goose
- **Контейнеризация**: Docker, Docker Compose
- **Шаблоны**: html/template

## Технологии

- Go 1.21
- PostgreSQL 16
- GORM (ORM)
- JWT для аутентификации
- Docker и Docker Compose
- Goose для миграций

## Установка и запуск

### Предварительные требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)

### Настройка переменных окружения

Создайте файл `.env` в корне проекта на основе примера:

```env
# База данных
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=newyear_app
DB_PORT=5432

# JWT настройки
JWT_ACCESS_SECRET=your-super-secret-jwt-key-change-this
JWT_ACCESS_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=168h

# OAuth Яндекс
YANDEX_CLIENT_ID=your_yandex_client_id
YANDEX_CLIENT_SECRET=your_yandex_client_secret
YANDEX_CALLBACK_URL=http://localhost:4200/auth/oauth/yandex/callback


# Остановить и удалить старые контейнеры (если были)
docker-compose down -v

# Собрать образ приложения
docker-compose build

# Запустить все сервисы
docker-compose up






Регистрация пользователя

$body = @{email="test@example.com"; password="password123"; full_name="Test User"} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:4200/auth/register" -Method POST -ContentType "application/json" -Body $body




Вход (сохраняем сессию)

$body = @{email="test@example.com"; password="password123"} | ConvertTo-Json
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$response = Invoke-WebRequest -Uri "http://localhost:4200/auth/login" -Method POST -ContentType "application/json" -Body $body -WebSession $session




Проверка кто я

Invoke-RestMethod -Uri "http://localhost:4200/auth/whoami" -Method GET -WebSession $session




Создать желание

$body = @{text="hihihi"; author="Hihik"; priority=1} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:4200/wishes" -Method POST -ContentType "application/json" -Body $body -WebSession $session



# Подключение через mongosh
docker exec -it newyear-mongodb mongosh -u admin -p admin123

# Переключиться на вашу базу данных
use newyear_app

# Показать всех пользователей
db.users.find().pretty()

# Показать только активных (не удаленных) пользователей
db.users.find({ deleted_at: null }).pretty()

# Показать все желания
db.wishes.find().pretty()

# Показать только активные желания
db.wishes.find({ deleted_at: null }).pretty()

# Посчитать количество пользователей
db.users.countDocuments()

# Посчитать количество желаний
db.wishes.countDocuments()

# Показать коллекции в базе
show collections

# Выйти
exit

