package main

import (
	"gorm.io/gorm"

	categoryApp "backend/application/category"
	stockApp "backend/application/stock"
	snapshotApp "backend/application/stock_snapshot"
	categoryRepo "backend/infrastructure/repositories/category"
	stockRepo "backend/infrastructure/repositories/stock"
	snapshotRepo "backend/infrastructure/repositories/stock_snapshot"
)

type Container struct {
	StockService    *stockApp.StockService
	CategoryService *categoryApp.CategoryService
	SnapshotService *snapshotApp.StockSnapshotService
}

func NewContainer(db *gorm.DB) *Container {
	stockRepo := stockRepo.NewStockRepository(db)
	stockService := stockApp.NewStockService(stockRepo)

	categoryRepo := categoryRepo.NewCategoryRepository(db)
	categoryService := categoryApp.NewCategoryService(categoryRepo)

	snapshotRepo := snapshotRepo.NewStockSnapshotRepository(db)
	snapshotService := snapshotApp.NewStockSnapshotService(snapshotRepo)

	return &Container{
		StockService:    stockService,
		CategoryService: categoryService,
		SnapshotService: snapshotService,
	}
}
