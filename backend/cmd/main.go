package main

import (
	"log"
	"net/http"

	stockApp "backend/application/stock"
	"backend/infrastructure/database"
	stockInfra "backend/infrastructure/repositories/stock"
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

	stockRepo := stockInfra.NewStockRepository(db)
	stockService := stockApp.NewStockService(stockRepo)
	router := httpInterface.NewRouter(stockService)

	handler := router.SetupRoutes()

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
