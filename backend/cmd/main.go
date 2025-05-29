package main

import (
	"log"
	"net/http"

	httpInterface "backend/interfaces/http"

	"gorm.io/gorm"
)

func main() {
	var db *gorm.DB

	router := httpInterface.NewRouter(db)
	handler := router.SetupRoutes()

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
