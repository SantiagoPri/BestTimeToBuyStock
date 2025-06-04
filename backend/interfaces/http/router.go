package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	category "backend/application/category"
	gameSession "backend/application/game_session"
	gmSession "backend/application/gm_session"
	stock "backend/application/stock"
	categoryHttp "backend/interfaces/http/category"
	gameSessionHttp "backend/interfaces/http/game_session"
	gmSessionHttp "backend/interfaces/http/gm_session"
	stockHttp "backend/interfaces/http/stock"
)

type Router struct {
	stockService       *stock.StockService
	categoryService    *category.CategoryService
	gameSessionService gameSession.Service
	gmSessionService   gmSession.Service
}

func NewRouter(
	stockService *stock.StockService,
	categoryService *category.CategoryService,
	gameSessionService gameSession.Service,
	gmSessionService gmSession.Service,
) *Router {
	return &Router{
		stockService:       stockService,
		categoryService:    categoryService,
		gameSessionService: gameSessionService,
		gmSessionService:   gmSessionService,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Apply CORS middleware
	router.Use(cors.Default())

	api := router.Group("/api")

	stockHandler := stockHttp.NewHandler(r.stockService)
	stockHttp.RegisterRoutes(api, stockHandler)

	categoryHandler := categoryHttp.NewHandler(r.categoryService)
	categoryHttp.RegisterRoutes(api, categoryHandler)

	gameSessionHandler := gameSessionHttp.NewHandler(r.gameSessionService)
	gameSessionHttp.RegisterRoutes(api, gameSessionHandler)

	gmSessionHandler := gmSessionHttp.NewHandler(r.gmSessionService)
	gmSessionHttp.RegisterRoutes(api, gmSessionHandler)

	return router
}
