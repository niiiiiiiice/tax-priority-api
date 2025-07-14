package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"tax-priority-api/src/handlers"
)

// SetupRoutes настраивает все маршруты приложения
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Настройка CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // В продакшене указать конкретные домены
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler()
	productHandler := handlers.NewProductHandler()

	// API группа
	api := r.Group("/api")
	{
		// Пользователи
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/count", userHandler.GetUserCount)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.POST("/batch", userHandler.GetUsersBatch)
			users.POST("/bulk-delete", userHandler.BulkDeleteUsers)
			users.PATCH("/:id/status", userHandler.ChangeUserStatus)
			users.GET("/by-role/:role", userHandler.GetUsersByRole)
			users.GET("/search", userHandler.SearchUsers)
		}

		// Продукты
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/count", productHandler.GetProductCount)
			products.GET("/:id", productHandler.GetProduct)
			products.POST("", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
			products.POST("/batch", productHandler.GetProductsBatch)
			products.POST("/bulk-delete", productHandler.BulkDeleteProducts)
			products.GET("/category/:category", productHandler.GetProductsByCategory)
			products.PATCH("/:id/stock", productHandler.UpdateProductStock)
			products.GET("/search", productHandler.SearchProducts)
		}
	}

	// Здоровье приложения
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})

	return r
}
