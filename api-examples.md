# API Examples

Примеры использования Tax Priority API

## Запуск API

### Локальный запуск
```bash
# Установить зависимости
make deps

# Запустить API (требуется PostgreSQL)
make run

# Или запустить в режиме разработки с автоперезагрузкой
make dev
```

### Запуск с Docker
```bash
# Запустить API и PostgreSQL в контейнерах
make docker-run

# Остановить контейнеры
make docker-stop
```

API будет доступно по адресу: http://localhost:8080

## Примеры запросов

### Проверка здоровья API
```bash
curl http://localhost:8080/health
```

### Пользователи

#### Создание пользователя
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "name": "Admin User",
    "role": "admin",
    "status": "active",
    "avatar": "https://example.com/avatar.jpg"
  }'
```

#### Получение списка пользователей
```bash
# Все пользователи
curl http://localhost:8080/api/users

# С пагинацией
curl "http://localhost:8080/api/users?_limit=5&_offset=0"

# С сортировкой
curl "http://localhost:8080/api/users?_sort=name&_order=asc"

# С фильтрацией
curl "http://localhost:8080/api/users?role=admin&status=active"

# Комбинированный запрос
curl "http://localhost:8080/api/users?role=admin&_limit=10&_sort=createdAt&_order=desc"
```

#### Получение пользователя по ID
```bash
curl http://localhost:8080/api/users/{user-id}
```

#### Обновление пользователя
```bash
curl -X PUT http://localhost:8080/api/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name",
    "status": "inactive"
  }'
```

#### Изменение статуса пользователя
```bash
curl -X PATCH http://localhost:8080/api/users/{user-id}/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "banned"
  }'
```

#### Получение пользователей по роли
```bash
curl http://localhost:8080/api/users/by-role/admin
```

#### Поиск пользователей
```bash
curl "http://localhost:8080/api/users/search?q=john"
```

#### Получение количества пользователей
```bash
curl http://localhost:8080/api/users/count
```

#### Batch операции
```bash
# Получить пользователей по списку ID
curl -X POST http://localhost:8080/api/users/batch \
  -H "Content-Type: application/json" \
  -d '{
    "ids": ["user-id-1", "user-id-2"]
  }'

# Массовое удаление
curl -X POST http://localhost:8080/api/users/bulk-delete \
  -H "Content-Type: application/json" \
  -d '{
    "ids": ["user-id-1", "user-id-2"]
  }'
```

#### Удаление пользователя
```bash
curl -X DELETE http://localhost:8080/api/users/{user-id}
```

### Продукты

#### Создание продукта
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "Apple MacBook Pro 16-inch",
    "price": 2499.99,
    "currency": "USD",
    "category": "laptops",
    "tags": ["apple", "laptop", "professional"],
    "quantity": 50,
    "images": ["https://example.com/macbook1.jpg", "https://example.com/macbook2.jpg"]
  }'
```

#### Получение списка продуктов
```bash
# Все продукты
curl http://localhost:8080/api/products

# С пагинацией
curl "http://localhost:8080/api/products?_limit=10&_offset=0"

# С сортировкой по цене
curl "http://localhost:8080/api/products?_sort=price&_order=desc"

# С фильтрацией
curl "http://localhost:8080/api/products?category=laptops&inStock=true"

# Фильтрация по цене
curl "http://localhost:8080/api/products?price_gte=1000&price_lte=3000"

# Комбинированный запрос
curl "http://localhost:8080/api/products?category=laptops&price_gte=1000&_sort=price&_order=asc&_limit=5"
```

#### Получение продукта по ID
```bash
curl http://localhost:8080/api/products/{product-id}
```

#### Обновление продукта
```bash
curl -X PUT http://localhost:8080/api/products/{product-id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product Name",
    "price": 1999.99,
    "quantity": 30
  }'
```

#### Обновление количества продукта
```bash
curl -X PATCH http://localhost:8080/api/products/{product-id}/stock \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100
  }'
```

#### Получение продуктов по категории
```bash
curl http://localhost:8080/api/products/category/laptops
```

#### Поиск продуктов
```bash
curl "http://localhost:8080/api/products/search?q=macbook"
```

#### Получение количества продуктов
```bash
curl http://localhost:8080/api/products/count
```

#### Batch операции
```bash
# Получить продукты по списку ID
curl -X POST http://localhost:8080/api/products/batch \
  -H "Content-Type: application/json" \
  -d '{
    "ids": ["product-id-1", "product-id-2"]
  }'

# Массовое удаление
curl -X POST http://localhost:8080/api/products/bulk-delete \
  -H "Content-Type: application/json" \
  -d '{
    "ids": ["product-id-1", "product-id-2"]
  }'
```

#### Удаление продукта
```bash
curl -X DELETE http://localhost:8080/api/products/{product-id}
```

## Фильтрация и пагинация

### Параметры пагинации
- `_limit` - количество записей на странице (по умолчанию: 10)
- `_offset` - смещение (по умолчанию: 0)

### Параметры сортировки
- `_sort` - поле для сортировки (по умолчанию: createdAt)
- `_order` - направление сортировки (asc/desc, по умолчанию: desc)

### Фильтры для пользователей
- `email` - поиск по email (частичное совпадение)
- `name` - поиск по имени (частичное совпадение)
- `role` - точное совпадение роли (admin, user, moderator)
- `status` - точное совпадение статуса (active, inactive, banned)
- `createdAt_gte` - дата создания от (ISO 8601)
- `createdAt_lte` - дата создания до (ISO 8601)

### Фильтры для продуктов
- `name` - поиск по названию (частичное совпадение)
- `category` - точное совпадение категории
- `currency` - точное совпадение валюты (USD, EUR, GBP)
- `inStock` - наличие на складе (true/false)
- `price_gte` - цена от
- `price_lte` - цена до
- `quantity_gte` - количество от
- `quantity_lte` - количество до
- `createdAt_gte` - дата создания от (ISO 8601)
- `createdAt_lte` - дата создания до (ISO 8601)

## Примеры ответов

### Успешный ответ со списком
```json
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user",
      "status": "active",
      "avatar": "",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### Ответ с ошибкой
```json
{
  "error": "User not found"
}
```

### Ответ с количеством
```json
{
  "count": 42
}
```

## Интеграция с AdminJS

Этот API полностью совместим с вашим AdminJS клиентом:

1. **Эндпоинты**: Все необходимые эндпоинты реализованы
2. **Параметры**: Поддерживаются все параметры фильтрации и пагинации
3. **Формат ответов**: Данные возвращаются в ожидаемом формате
4. **Batch операции**: Поддерживаются массовые операции
5. **Кастомные методы**: Реализованы специальные эндпоинты для смены статуса и обновления остатков

Просто обновите конфигурацию `apiClient` в вашем AdminJS клиенте, указав базовый URL: `http://localhost:8080/api` 