package main

import (
	"gorm.io/gorm"

	categoryApp "backend/application/category"
	gameSessionApp "backend/application/game_session"
	gmSessionApp "backend/application/gm_session"
	stockApp "backend/application/stock"
	"backend/infrastructure/ai_model"
	"backend/infrastructure/redis"
	categoryRepo "backend/infrastructure/repositories/category"
	gameSessionRepo "backend/infrastructure/repositories/game_session"
	gmSessionRepo "backend/infrastructure/repositories/gm_session"
	stockRepo "backend/infrastructure/repositories/stock"
)

type Container struct {
	StockService       *stockApp.StockService
	CategoryService    *categoryApp.CategoryService
	GameSessionService gameSessionApp.Service
	GMSessionService   gmSessionApp.Service
}

func NewContainer(db *gorm.DB) *Container {
	stockRepo := stockRepo.NewStockRepository(db)
	stockService := stockApp.NewStockService(stockRepo)

	categoryRepo := categoryRepo.NewCategoryRepository(db)
	categoryService := categoryApp.NewCategoryService(categoryRepo)

	aiModel, err := ai_model.NewOpenRouterAgent()
	if err != nil {
		panic(err)
	}

	redisService := redis.NewRedisService()

	// Create GM session repository and service
	gmSessionRepository := gmSessionRepo.NewRepository(redisService)
	gmSessionService := gmSessionApp.NewService(gmSessionRepository)

	// Create game session repository and service with all dependencies
	gameSessionRepository := gameSessionRepo.NewRepository(db, redisService)
	gameSessionService := gameSessionApp.NewService(
		gameSessionRepository,
		stockRepo,
		aiModel,
		gmSessionService,
	)

	return &Container{
		StockService:       stockService,
		CategoryService:    categoryService,
		GameSessionService: gameSessionService,
		GMSessionService:   gmSessionService,
	}
}
