# Swagger Documentation Setup

Swagger документация для Tax Priority API настроена и готова к использованию.

## Доступ к документации

### Swagger UI
После запуска API, Swagger UI будет доступен по адресам:
```
http://localhost:8081/swagger
http://localhost:8081/swagger/index.html
```

> **Примечание**: `/swagger` автоматически перенаправляет на `/swagger/index.html`

### JSON документация
Swagger JSON доступен по адресу:
```
http://localhost:8081/swagger/doc.json
```

### YAML документация
Swagger YAML доступен по адресу:
```
http://localhost:8081/swagger/swagger.yaml
```

## Команды для работы с документацией

### Генерация документации
```bash
# Генерировать Swagger документацию
make swagger

# Или напрямую через swag
swag init -g cmd/main.go -o docs --parseDependency
```

### Установка инструментов
```bash
# Установить все инструменты разработки (включая swag)
make install-tools

# Или только swag
go install github.com/swaggo/swag/cmd/swag@latest
```

## Структура документации

### Основные модели
- `CreateFAQRequest` - Запрос на создание FAQ
- `UpdateFAQRequest` - Запрос на обновление FAQ
- `FAQResponse` - Ответ с данными FAQ
- `PaginatedFAQResponse` - Пагинированный ответ
- `CommandResult` - Результат выполнения команды
- `ErrorResponse` - Ошибка

### Доступные эндпоинты

#### CRUD операции
- `POST /api/faq` - Создать FAQ
- `GET /api/faq/{id}` - Получить FAQ по ID
- `PUT /api/faq/{id}` - Обновить FAQ
- `DELETE /api/faq/{id}` - Удалить FAQ

#### Получение списков
- `GET /api/faq` - Получить список FAQ с пагинацией
- `GET /api/faq/category/{category}` - Получить FAQ по категории
- `GET /api/faq/search` - Поиск FAQ
- `GET /api/faq/categories` - Получить категории
- `GET /api/faq/count` - Получить количество FAQ

#### Batch операции
- `POST /api/faq/batch` - Получить FAQ по списку ID
- `POST /api/faq/bulk-delete` - Массовое удаление FAQ

#### Управление состоянием
- `PATCH /api/faq/{id}/activate` - Активировать FAQ
- `PATCH /api/faq/{id}/deactivate` - Деактивировать FAQ
- `PATCH /api/faq/{id}/priority` - Обновить приоритет FAQ

## Параметры запросов

### Пагинация
- `_limit` - Количество записей (по умолчанию: 10)
- `_offset` - Смещение (по умолчанию: 0)

### Сортировка
- `_sort` - Поле сортировки (по умолчанию: createdAt)
- `_order` - Направление (asc/desc, по умолчанию: desc)

### Фильтрация
- `category` - Фильтр по категории
- `isActive` - Фильтр по активности
- `q` - Поисковый запрос

## Примеры использования

### Создание FAQ
```bash
curl -X POST http://localhost:8081/api/faq \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Как подать налоговую декларацию?",
    "answer": "Для подачи налоговой декларации необходимо...",
    "category": "налоги",
    "priority": 50
  }'
```

### Получение списка FAQ
```bash
curl "http://localhost:8081/api/faq?_limit=10&_offset=0&_sort=priority&_order=desc"
```

### Поиск FAQ
```bash
curl "http://localhost:8081/api/faq/search?q=налоги&category=налоги&activeOnly=true"
```

## Автоматическое обновление документации

Документация автоматически генерируется из аннотаций в коде. При добавлении новых эндпоинтов:

1. Добавьте Swagger аннотации к функции-хендлеру
2. Используйте модели из `src/presentation/models/swagger_models.go`
3. Запустите `make swagger` для регенерации документации

### Пример аннотации
```go
// CreateFAQ создает новую FAQ
// @Summary Создать FAQ
// @Description Создает новую запись FAQ
// @Tags FAQ
// @Accept json
// @Produce json
// @Param faq body models.CreateFAQRequest true "Данные FAQ"
// @Success 201 {object} models.CommandResult
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/faq [post]
func (h *FAQHTTPHandler) CreateFAQ(c *gin.Context) {
    // implementation
}
```

## Конфигурация

Основная конфигурация Swagger находится в `cmd/main.go`:

```go
// @title Tax Priority API
// @version 1.0
// @description REST API для управления FAQ в системе Tax Priority
// @host localhost:8081
// @BasePath /api
```

## Файлы документации

- `docs/docs.go` - Сгенерированная Go документация
- `docs/swagger.json` - JSON схема
- `docs/swagger.yaml` - YAML схема
- `docs/swagger-ui.html` - Standalone HTML страница

## Troubleshooting

### Ошибка генерации
Если возникают ошибки при генерации:
```bash
# Очистить кеш модулей
go clean -modcache

# Пересобрать проект
go mod tidy
make swagger
```

### Swagger UI не загружается
1. Убедитесь, что API запущен
2. Проверьте, что маршрут `/swagger/*any` зарегистрирован
3. Проверьте файлы в папке `docs/`

### Модели не отображаются
1. Проверьте импорты в хендлерах
2. Убедитесь, что используются модели из `models` пакета
3. Перегенерируйте документацию: `make swagger` 