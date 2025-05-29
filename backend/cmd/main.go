package main

import (
	"log"
	"net/http"

	stockApp "backend/application/stock"
	stockInfra "backend/infrastructure/repositories/stock"
	httpInterface "backend/interfaces/http"

	"gorm.io/gorm"
)

func main() {
	var db *gorm.DB

	stockRepo := stockInfra.NewStockRepository(db)
	stockService := stockApp.NewStockService(stockRepo)
	router := httpInterface.NewRouter(stockService)

	handler := router.SetupRoutes()

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
