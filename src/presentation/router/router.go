package router

import (
	"log"
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
	if err := db.AutoMigrate(&models.FAQModel{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Инициализация FAQ обработчика через Wire
	faqHandler := wire.InitializeFAQHTTPHandler(db)

	// Регистрация маршрутов
	handlers.RegisterFAQRoutes(router, faqHandler)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Tax Priority API is running",
		})
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
