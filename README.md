# Tax Priority API

REST API для управления пользователями и продуктами, разработанное с использованием Go, Gin и GORM.

## Особенности

- CRUD операции для пользователей и продуктов
- Фильтрация и пагинация
- Поиск по полям
- Batch операции
- Поддержка CORS
- PostgreSQL база данных
- UUID для идентификаторов

## Структура проекта

```
api/
├── cmd/
│   └── main.go              # Точка входа приложения
├── src/
│   ├── database/
│   │   └── database.go      # Конфигурация базы данных
│   ├── handlers/
│   │   ├── user_handler.go  # Обработчики для пользователей
│   │   └── product_handler.go # Обработчики для продуктов
│   ├── models/
│   │   ├── user.go          # Модель пользователя
│   │   └── product.go       # Модель продукта
│   └── routes/
│       └── routes.go        # Настройка маршрутов
├── go.mod                   # Зависимости Go
└── README.md
```

## Установка и запуск

### Предварительные требования

- Go 1.21 или выше
- PostgreSQL

### Установка зависимостей

```bash
go mod download
```

### Настройка переменных окружения

Создайте файл `.env` или установите переменные окружения:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tax_priority
DB_SSLMODE=disable

# Server Configuration
PORT=8080
GIN_MODE=release
```

### Создание базы данных

```sql
CREATE DATABASE tax_priority;
```

### Запуск приложения

```bash
go run cmd/main.go
```

API будет доступно по адресу: `http://localhost:8080`

## API Endpoints

### Пользователи

- `GET /api/users` - Получить список пользователей
- `GET /api/users/count` - Получить количество пользователей
- `GET /api/users/:id` - Получить пользователя по ID
- `POST /api/users` - Создать нового пользователя
- `PUT /api/users/:id` - Обновить пользователя
- `DELETE /api/users/:id` - Удалить пользователя
- `POST /api/users/batch` - Получить пользователей по списку ID
- `POST /api/users/bulk-delete` - Массовое удаление пользователей
- `PATCH /api/users/:id/status` - Изменить статус пользователя
- `GET /api/users/by-role/:role` - Получить пользователей по роли
- `GET /api/users/search` - Поиск пользователей

### Продукты

- `GET /api/products` - Получить список продуктов
- `GET /api/products/count` - Получить количество продуктов
- `GET /api/products/:id` - Получить продукт по ID
- `POST /api/products` - Создать новый продукт
- `PUT /api/products/:id` - Обновить продукт
- `DELETE /api/products/:id` - Удалить продукт
- `POST /api/products/batch` - Получить продукты по списку ID
- `POST /api/products/bulk-delete` - Массовое удаление продуктов
- `GET /api/products/category/:category` - Получить продукты по категории
- `PATCH /api/products/:id/stock` - Обновить количество продукта
- `GET /api/products/search` - Поиск продуктов

### Системные

- `GET /health` - Проверка состояния API

## Параметры запросов

### Пагинация

- `_limit` - Количество записей на странице (по умолчанию: 10)
- `_offset` - Смещение (по умолчанию: 0)

### Сортировка

- `_sort` - Поле для сортировки (по умолчанию: createdAt)
- `_order` - Направление сортировки (asc/desc, по умолчанию: desc)

### Фильтрация

#### Пользователи

- `email` - Поиск по email (частичное совпадение)
- `name` - Поиск по имени (частичное совпадение)
- `role` - Фильтр по роли (admin, user, moderator)
- `status` - Фильтр по статусу (active, inactive, banned)
- `createdAt_gte` - Дата создания от
- `createdAt_lte` - Дата создания до

#### Продукты

- `name` - Поиск по названию (частичное совпадение)
- `category` - Фильтр по категории
- `currency` - Фильтр по валюте (USD, EUR, GBP)
- `inStock` - Фильтр по наличию (true/false)
- `price_gte` - Цена от
- `price_lte` - Цена до
- `quantity_gte` - Количество от
- `quantity_lte` - Количество до
- `createdAt_gte` - Дата создания от
- `createdAt_lte` - Дата создания до

## Примеры запросов

### Создание пользователя

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user",
    "status": "active"
  }'
```

### Создание продукта

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sample Product",
    "description": "This is a sample product",
    "price": 29.99,
    "currency": "USD",
    "category": "electronics",
    "quantity": 100,
    "tags": ["sample", "test"]
  }'
```

### Получение пользователей с фильтрацией

```bash
curl "http://localhost:8080/api/users?role=admin&_limit=5&_sort=name&_order=asc"
```

### Поиск продуктов

```bash
curl "http://localhost:8080/api/products/search?q=laptop"
```

## Модели данных

### User

```json
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "role": "admin|user|moderator",
  "status": "active|inactive|banned",
  "avatar": "string",
  "createdAt": "datetime",
  "updatedAt": "datetime"
}
```

### Product

```json
{
  "id": "uuid",
  "name": "string",
  "description": "string",
  "price": "number",
  "currency": "USD|EUR|GBP",
  "category": "string",
  "tags": ["string"],
  "inStock": "boolean",
  "quantity": "number",
  "images": ["string"],
  "createdAt": "datetime",
  "updatedAt": "datetime"
}
```

## Соответствие с AdminJS клиентом

Это API полностью совместимо с вашим AdminJS клиентом:

- Поддерживает все необходимые эндпоинты
- Использует те же параметры фильтрации и пагинации
- Возвращает данные в ожидаемом формате
- Поддерживает batch операции
- Имеет специальные эндпоинты для кастомных операций

## Разработка

### Запуск в режиме разработки

```bash
go run cmd/main.go
```

### Сборка для продакшена

```bash
go build -o bin/api cmd/main.go
```

### Запуск собранного бинарника

```bash
./bin/api
``` 