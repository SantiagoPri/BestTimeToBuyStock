package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"backend/infrastructure/database"
	"backend/infrastructure/redis"
	httpInterface "backend/interfaces/http"

	_ "backend/docs" // This line is important!

	"github.com/joho/godotenv"
)

// @title           Best Time To Buy Stock API
// @version         1.0
// @description     A stock trading game API where users can practice trading with historical data.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the session ID.

func generateSwaggerDocs() error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}
	swagPath := filepath.Join(gopath, "bin", "swag")

	cmd := exec.Command(swagPath, "init", "-g", "cmd/main.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Failed to generate Swagger docs: %v\nOutput: %s", err, output)
		return err
	}
	log.Printf("Swagger documentation generated successfully")
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Generate Swagger documentation
	if err := generateSwaggerDocs(); err != nil {
		log.Printf("Warning: Could not generate Swagger docs: %v", err)
	}

	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// Check Redis availability
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redis.Ping(ctx); err != nil {
		log.Fatalf("Redis is not available: %v", err)
	}
	log.Printf("Redis connection verified")

	container := NewContainer(db)
	router := httpInterface.NewRouter(
		container.StockService,
		container.CategoryService,
		container.GameSessionService,
		container.GMSessionService,
	)

	handler := router.SetupRoutes()

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
