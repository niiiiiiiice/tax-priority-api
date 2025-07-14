$API_URL = "http://localhost:8081"

Write-Host "Testing Swagger endpoints..." -ForegroundColor Green
Write-Host "================================" -ForegroundColor Green

# Test health endpoint
Write-Host "1. Testing health endpoint..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "$API_URL/health" -Method Get
    Write-Host "Health Status: $($healthResponse.status)" -ForegroundColor Green
    Write-Host "Message: $($healthResponse.message)" -ForegroundColor Green
}
catch {
    Write-Host "Health endpoint failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Test Swagger JSON
Write-Host "2. Testing Swagger JSON..." -ForegroundColor Yellow
try {
    $swaggerResponse = Invoke-RestMethod -Uri "$API_URL/swagger/doc.json" -Method Get
    Write-Host "API Title: $($swaggerResponse.info.title)" -ForegroundColor Green
    Write-Host "API Version: $($swaggerResponse.info.version)" -ForegroundColor Green
    Write-Host "API Description: $($swaggerResponse.info.description)" -ForegroundColor Green
}
catch {
    Write-Host "Swagger JSON failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Test Swagger UI redirect
Write-Host "3. Testing Swagger UI redirect..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$API_URL/swagger" -Method Get -MaximumRedirection 0
    Write-Host "Swagger redirect Status Code: $($response.StatusCode)" -ForegroundColor Green
}
catch {
    if ($_.Exception.Response.StatusCode -eq 302) {
        Write-Host "Swagger redirect works (302 Found)" -ForegroundColor Green
    }
    else {
        Write-Host "Swagger redirect failed: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""

# Test Swagger UI page
Write-Host "4. Testing Swagger UI page..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$API_URL/swagger/index.html" -Method Get
    Write-Host "Swagger UI Status Code: $($response.StatusCode)" -ForegroundColor Green
}
catch {
    Write-Host "Swagger UI failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "Swagger UI is available at: $API_URL/swagger (redirects to /swagger/index.html)" -ForegroundColor Cyan
Write-Host "Swagger JSON is available at: $API_URL/swagger/doc.json" -ForegroundColor Cyan

Write-Host ""
Write-Host "Test completed!" -ForegroundColor Green 