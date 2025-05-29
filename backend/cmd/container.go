package main

import (
	"gorm.io/gorm"

	stockApp "backend/application/stock"
	stockRepo "backend/infrastructure/repositories/stock"
)

type Container struct {
	StockService *stockApp.StockService
}

func NewContainer(db *gorm.DB) *Container {
	stockRepo := stockRepo.NewStockRepository(db)
	stockService := stockApp.NewStockService(stockRepo)

	return &Container{
		StockService: stockService,
	}
}
