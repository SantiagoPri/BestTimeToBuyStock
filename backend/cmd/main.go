package main

import (
	"log"
	"net/http"

	"backend/infrastructure/database"
	httpInterface "backend/interfaces/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	container := NewContainer(db)
	router := httpInterface.NewRouter(
		container.StockService,
		container.CategoryService,
		container.SnapshotService,
	)

	handler := router.SetupRoutes()

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
