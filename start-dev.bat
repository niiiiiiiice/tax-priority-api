@echo off
echo Starting Tax Priority API development server...
echo Make sure PostgreSQL is running on localhost:5432
echo Database: tax_priority, User: postgres, Password: postgres
echo.

set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=postgres
set DB_NAME=tax_priority
set DB_SSLMODE=disable
set PORT=8081

echo Environment variables set:
echo DB_HOST=%DB_HOST%
echo DB_PORT=%DB_PORT%
echo DB_USER=%DB_USER%
echo DB_NAME=%DB_NAME%
echo PORT=%PORT%
echo.

echo Starting Air...
air 