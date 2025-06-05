package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

	// Apply CORS middleware with custom configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	api := router.Group("/api")

	stockHandler := stockHttp.NewHandler(r.stockService)
	stockHttp.RegisterRoutes(api, stockHandler)

	categoryHandler := categoryHttp.NewHandler(r.categoryService)
	categoryHttp.RegisterRoutes(api, categoryHandler)

	gameSessionHandler := gameSessionHttp.NewHandler(r.gameSessionService)
	gameSessionHttp.RegisterRoutes(api, gameSessionHandler)

	gmSessionHandler := gmSessionHttp.NewHandler(r.gmSessionService)
	gmSessionHttp.RegisterRoutes(api, gmSessionHandler)

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
