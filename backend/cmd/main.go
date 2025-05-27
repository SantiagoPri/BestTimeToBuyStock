package main

import (
	"log"
	"net/http"

	"backend/internal/stocks/transport"
)

func main() {
	router := transport.SetupRouter()

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
