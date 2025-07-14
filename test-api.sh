#!/bin/bash

# Скрипт для тестирования Tax Priority API
# Убедитесь, что API запущен на localhost:8080

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api"

echo "🚀 Testing Tax Priority API"
echo "================================"

# Проверка здоровья API
echo "1. Health Check:"
curl -s "$BASE_URL/health"
echo ""

# Тест пользователей (покажет ошибку без БД, но проверит маршруты)
echo "2. Users endpoint (without database):"
curl -s "$API_URL/users"
echo ""

echo "3. Users count endpoint:"
curl -s "$API_URL/users/count"
echo ""

echo "4. Users with sorting (createdAt -> created_at):"
curl -s "$API_URL/users?_sort=createdAt&_order=desc"
echo ""

# Тест продуктов
echo "5. Products endpoint (without database):"
curl -s "$API_URL/products"
echo ""

echo "6. Products count endpoint:"
curl -s "$API_URL/products/count"
echo ""

echo "7. Products with sorting (createdAt -> created_at):"
curl -s "$API_URL/products?_sort=createdAt&_order=desc"
echo ""

echo "================================"
echo "✅ API is responding to requests"
echo "❌ Database connection errors are expected without PostgreSQL"
echo ""
echo "To run with database:"
echo "  make docker-run"
echo ""
echo "To test with real data, first start PostgreSQL and create the database:"
echo "  createdb tax_priority"
echo "  make run-local" 