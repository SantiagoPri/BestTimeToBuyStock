package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"backend/infrastructure/database"
	"backend/infrastructure/redis"
	httpInterface "backend/interfaces/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
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
		container.SnapshotService,
		container.GameSessionService,
		container.GMSessionService,
	)

	handler := router.SetupRoutes()

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
