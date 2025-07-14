package main

import (
	"log"
	"os"

	"tax-priority-api/src/database"
	"tax-priority-api/src/routes"
)

func main() {
	// Инициализация базы данных
	database.InitDatabase()

	// Настройка маршрутов
	router := routes.SetupRoutes()

	// Получаем порт из переменной окружения или используем 8080 по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
