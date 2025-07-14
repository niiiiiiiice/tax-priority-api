# Пример правильного вызова API

## До исправления (неправильно)
```bash
curl -X 'POST' \
  'http://localhost:8081/api/api/faq' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "answer": "Для подачи налоговой декларации необходимо...",
    "category": "налоги",
    "priority": 50,
    "question": "Как подать налоговую декларацию?"
  }'
```

## После исправления (правильно)
```bash
curl -X 'POST' \
  'http://localhost:8081/api/faq' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "answer": "Для подачи налоговой декларации необходимо...",
    "category": "налоги",
    "priority": 50,
    "question": "Как подать налоговую декларацию?"
  }'
```

## Доступ к Swagger UI

### Теперь доступно:
- `http://localhost:8081/swagger` - автоматический редирект
- `http://localhost:8081/swagger/index.html` - прямой доступ

### Исправленные проблемы:
1. ✅ Убрано дублирование `/api` в URL
2. ✅ Добавлен редирект с `/swagger` на `/swagger/index.html`
3. ✅ Обновлена документация

### Другие примеры правильных вызовов:

```bash
# Получить все FAQ
curl -X 'GET' 'http://localhost:8081/api/faq' -H 'accept: application/json'

# Получить FAQ по ID
curl -X 'GET' 'http://localhost:8081/api/faq/{id}' -H 'accept: application/json'

# Поиск FAQ
curl -X 'GET' 'http://localhost:8081/api/faq/search?q=налоги' -H 'accept: application/json'
``` 