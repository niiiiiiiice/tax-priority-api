package main

import (
	"log"
	"os"

	_ "tax-priority-api/docs" // Swagger docs
	"tax-priority-api/src/infrastructure/persistence"
	"tax-priority-api/src/presentation/router"
)

// @title Tax Priority API
// @version 1.0
// @description REST API для управления FAQ в системе Tax Priority
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:5000
// @BasePath /
// @schemes http https

// @securityDefinitions.oauth2.authorizationCode OAuth2AuthorizationCode
// @authorizationUrl http://localhost:8080/realms/master/protocol/openid-connect/auth
// @tokenUrl http://localhost:8080/realms/master/protocol/openid-connect/token
// @scope.api:read Read access to API
// @scope.api:write Write access to API

func main() {
	persistence.Connect(persistence.NewDatabaseConfig())

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
