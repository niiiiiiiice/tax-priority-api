Write-Host "Starting Tax Priority API development server..." -ForegroundColor Green
Write-Host "Make sure PostgreSQL is running on localhost:5432" -ForegroundColor Yellow
Write-Host "Database: tax_priority, User: postgres, Password: postgres" -ForegroundColor Yellow
Write-Host ""

$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "postgres"
$env:DB_NAME = "tax_priority"
$env:DB_SSLMODE = "disable"
$env:PORT = "8080"

Write-Host "Environment variables set:" -ForegroundColor Cyan
Write-Host "DB_HOST=$env:DB_HOST" -ForegroundColor Gray
Write-Host "DB_PORT=$env:DB_PORT" -ForegroundColor Gray
Write-Host "DB_USER=$env:DB_USER" -ForegroundColor Gray
Write-Host "DB_NAME=$env:DB_NAME" -ForegroundColor Gray
Write-Host "PORT=$env:PORT" -ForegroundColor Gray
Write-Host ""

Write-Host "Starting Air..." -ForegroundColor Green
air 