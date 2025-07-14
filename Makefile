.PHONY: build run test clean docker-build docker-run docker-stop deps

# Переменные
BINARY_NAME=tax-priority-api
DOCKER_IMAGE=tax-priority-api
GO_VERSION=1.24.5

# Сборка приложения
build:
	go build -o bin/$(BINARY_NAME) cmd/main.go

# Запуск приложения
run:
	go run cmd/main.go

# Запуск приложения с переменными окружения для локальной разработки
run-local:
	@echo "Starting API server locally..."
	@echo "Make sure PostgreSQL is running on localhost:5432"
	@echo "Database: tax_priority, User: postgres, Password: postgres"
	DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=tax_priority DB_SSLMODE=disable PORT=8080 go run cmd/main.go

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Очистка
clean:
	rm -rf bin/
	go clean

# Форматирование кода
fmt:
	go fmt ./...

# Проверка кода
vet:
	go vet ./...

# Тестирование
test:
	go test -v ./...

# Сборка Docker образа
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Запуск с Docker Compose
docker-run:
	docker-compose up --build

# Остановка Docker Compose
docker-stop:
	docker-compose down

# Полная очистка Docker (включая volumes)
docker-clean:
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE) || true

# Разработка (запуск с автоперезагрузкой)
dev:
	@echo "Starting development server..."
	@GOPATH=$$(go env GOPATH); \
	if [ -f "$$GOPATH/bin/air" ]; then \
		$$GOPATH/bin/air; \
	elif command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
		$$GOPATH/bin/air; \
	fi

# Разработка для Windows (запуск с автоперезагрузкой)
dev-win:
	@echo "Starting development server on Windows..."
	@echo "Make sure PostgreSQL is running on localhost:5432"
	@echo "Database: tax_priority, User: postgres, Password: postgres"
	@set DB_HOST=localhost& set DB_PORT=5432& set DB_USER=postgres& set DB_PASSWORD=postgres& set DB_NAME=tax_priority& set DB_SSLMODE=disable& set PORT=8081& air

# Установка инструментов для разработки
install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/google/wire/cmd/wire@latest

# Генерация Swagger документации
swagger:
	swag init -g cmd/main.go -o docs --parseDependency

# Генерация Wire кода для dependency injection
wire:
	wire ./src/wire

# Проверка всего
check: fmt vet test

# Помощь
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application (requires PostgreSQL)"
	@echo "  run-local    - Run with local PostgreSQL settings"
	@echo "  deps         - Install dependencies"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  vet          - Check code"
	@echo "  test         - Run tests"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose (includes PostgreSQL)"
	@echo "  docker-stop  - Stop Docker Compose"
	@echo "  docker-clean - Clean Docker containers and images"
	@echo "  dev          - Start development server with auto-reload (requires PostgreSQL)"
	@echo "  dev-win      - Start development server on Windows with auto-reload"
	@echo "  install-tools- Install development tools"
	@echo "  swagger      - Generate Swagger documentation"
	@echo "  wire         - Generate Wire dependency injection code"
	@echo "  check        - Run all checks (fmt, vet, test)"
	@echo "  help         - Show this help"
	@echo ""
	@echo "For quick start with database:"
	@echo "  make docker-run  - Starts API + PostgreSQL in containers"
	@echo ""
	@echo "For local development:"
	@echo "  1. Start PostgreSQL locally"
	@echo "  2. Create database 'tax_priority'"
	@echo "  3. make run-local or make dev"
	@echo ""
	@echo "For Windows development:"
	@echo "  1. Start PostgreSQL locally"
	@echo "  2. Create database 'tax_priority'"
	@echo "  3. Run start-dev.bat or make dev-win" 