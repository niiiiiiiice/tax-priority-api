# Настройка Air на Windows

Air - это инструмент для автоматической перезагрузки Go приложений при изменении кода.

## Установка

Air уже установлен в проекте. Если нужно переустановить:

```bash
go install github.com/air-verse/air@latest
```

## Конфигурация

Конфигурация Air находится в файле `.air.toml`. Основные настройки для Windows:

- `bin = "./tmp/main.exe"` - путь к исполняемому файлу с расширением .exe
- `cmd = "go build -o ./tmp/main.exe ./cmd/main.go"` - команда сборки для Windows

## Запуск

### Способ 1: Batch файл (рекомендуется)
```bash
start-dev.bat
```

### Способ 2: PowerShell скрипт
```powershell
.\start-dev.ps1
```

### Способ 3: Makefile
```bash
make dev-win
```

### Способ 4: Напрямую через Air
```bash
# Сначала установите переменные окружения
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=postgres
set DB_NAME=tax_priority
set DB_SSLMODE=disable
set PORT=8080

# Затем запустите Air
air
```

## Переменные окружения

Для локальной разработки используются следующие переменные:

- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=tax_priority`
- `DB_SSLMODE=disable`
- `PORT=8080`

## Требования

1. PostgreSQL должен быть запущен на localhost:5432
2. База данных `tax_priority` должна существовать
3. Пользователь `postgres` с паролем `postgres`

## Возможные проблемы

### Air не найден
Убедитесь, что `%GOPATH%\bin` добавлен в PATH:
```bash
go env GOPATH
```

### Ошибка сборки
Проверьте, что Go установлен и проект компилируется:
```bash
go build -o ./tmp/main.exe ./cmd/main.go
```

### Ошибка подключения к базе данных
Убедитесь, что PostgreSQL запущен и доступен:
```bash
psql -h localhost -p 5432 -U postgres -d tax_priority
```

## Полезные команды

```bash
# Показать доступные команды
make help

# Запустить без автоперезагрузки
make run-local

# Запустить с Docker
make docker-run

# Остановить Air
Ctrl+C
``` 