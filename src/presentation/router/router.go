package router

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tax-priority-api/src/infrastructure/persistence"
	"tax-priority-api/src/infrastructure/persistence/models"
	"tax-priority-api/src/presentation/handlers"
	"tax-priority-api/src/wire"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	// Подключение к базе данных
	db, err := persistence.Connect(persistence.NewDatabaseConfig())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Миграции
	if err := db.AutoMigrate(&models.FAQModel{}, &models.TestimonialModel{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Инициализация фабрики обработчиков
	handlerFactory := wire.InitializeHandlerFactory(db)

	// Создание обработчиков через фабрику
	faqHandler := handlerFactory.CreateFAQHandler()
	wsHandler := handlerFactory.CreateWebSocketHandler()
	testimonialHandler := handlerFactory.CreateTestimonialHandler()

	// Запуск WebSocket хаба в горутине
	go wsHandler.GetHub().Run(context.Background())

	// Регистрация маршрутов
	handlers.RegisterFAQRoutes(router, faqHandler)
	handlers.RegisterTestimonialRoutes(router, testimonialHandler)
	RegisterWebSocketRoutes(router, wsHandler)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Tax Priority API is running",
			"endpoints": gin.H{
				"swagger":   "/swagger",
				"websocket": "/ws",
				"ws_test":   "/ws/test-page",
				"api_docs":  "/swagger/index.html",
			},
		})
	})

	// WebSocket test page
	router.GET("/ws/test-page", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Ошибка получения рабочей директории:", err)
			return
		}

		fmt.Println("Текущая директория:", wd)

		// Путь к файлу относительно рабочей директории
		filePath := filepath.Join(wd, "websocket-test.html")

		c.File(filePath)
	})

	// Swagger documentation
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	return router
}

// RegisterWebSocketRoutes регистрирует WebSocket маршруты
func RegisterWebSocketRoutes(r *gin.Engine, handler *handlers.WebSocketHandler) {
	ws := r.Group("/ws")
	{
		ws.GET("", handler.HandleWebSocket)
		ws.GET("/stats", handler.GetWebSocketStats)
		ws.GET("/info", handler.GetConnectionInfo)
		ws.POST("/test", handler.SendTestNotification)
		ws.POST("/broadcast", handler.BroadcastMessage)
	}
}
