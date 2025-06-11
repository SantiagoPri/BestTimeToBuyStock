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
	"backend/infrastructure/taskrunner"
)

type Container struct {
	StockService       *stockApp.StockService
	CategoryService    *categoryApp.CategoryService
	GameSessionService gameSessionApp.Service
	GMSessionService   gmSessionApp.Service
	TaskRunner         *taskrunner.TaskRunner
}

func NewContainer(db *gorm.DB) *Container {
	// Initialize TaskRunner with a buffer size of 100
	tr := taskrunner.New(100)
	tr.Start()

	stockRepo := stockRepo.NewStockRepository(db)
	stockService := stockApp.NewStockService(stockRepo)

	categoryRepo := categoryRepo.NewCategoryRepository(db)
	categoryService := categoryApp.NewCategoryService(categoryRepo)

	aiModel, err := ai_model.NewOpenRouterAgent()
	if err != nil {
		panic(err)
	}

	redisService := redis.NewRedisService()

	gmSessionRepository := gmSessionRepo.NewRepository(redisService)
	gmSessionService := gmSessionApp.NewService(gmSessionRepository)

	gameSessionRepository := gameSessionRepo.NewRepository(db, redisService)
	gameSessionService := gameSessionApp.NewService(
		gameSessionRepository,
		stockRepo,
		categoryRepo,
		aiModel,
		gmSessionService,
		tr,
	)

	return &Container{
		StockService:       stockService,
		CategoryService:    categoryService,
		GameSessionService: gameSessionService,
		GMSessionService:   gmSessionService,
		TaskRunner:         tr,
	}
}
