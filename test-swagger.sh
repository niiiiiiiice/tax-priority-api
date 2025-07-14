#!/bin/bash

# Test script for Swagger endpoints
API_URL="http://localhost:8081"

echo "Testing Swagger endpoints..."
echo "================================"

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "${API_URL}/health" | jq .

echo -e "\n2. Testing Swagger JSON..."
curl -s "${API_URL}/swagger/doc.json" | jq '.info'

echo -e "\n3. Testing Swagger UI availability..."
curl -s -o /dev/null -w "%{http_code}" "${API_URL}/swagger/index.html"

echo -e "\n\nSwagger UI is available at: ${API_URL}/swagger/index.html"
echo "Swagger JSON is available at: ${API_URL}/swagger/doc.json"

echo -e "\nTest completed!" 