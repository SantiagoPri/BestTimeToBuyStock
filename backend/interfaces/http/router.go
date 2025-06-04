package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	category "backend/application/category"
	gameSession "backend/application/game_session"
	gmSession "backend/application/gm_session"
	stock "backend/application/stock"
	snapshot "backend/application/stock_snapshot"
	categoryHttp "backend/interfaces/http/category"
	gameSessionHttp "backend/interfaces/http/game_session"
	gmSessionHttp "backend/interfaces/http/gm_session"
	stockHttp "backend/interfaces/http/stock"
	snapshotHttp "backend/interfaces/http/stock_snapshot"
)

type Router struct {
	stockService       *stock.StockService
	categoryService    *category.CategoryService
	snapshotService    *snapshot.StockSnapshotService
	gameSessionService gameSession.Service
	gmSessionService   gmSession.Service
}

func NewRouter(
	stockService *stock.StockService,
	categoryService *category.CategoryService,
	snapshotService *snapshot.StockSnapshotService,
	gameSessionService gameSession.Service,
	gmSessionService gmSession.Service,
) *Router {
	return &Router{
		stockService:       stockService,
		categoryService:    categoryService,
		snapshotService:    snapshotService,
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

	snapshotHandler := snapshotHttp.NewHandler(r.snapshotService)
	snapshotHttp.RegisterRoutes(api, snapshotHandler)

	gameSessionHandler := gameSessionHttp.NewHandler(r.gameSessionService)
	gameSessionHttp.RegisterRoutes(api, gameSessionHandler)

	gmSessionHandler := gmSessionHttp.NewHandler(r.gmSessionService)
	gmSessionHttp.RegisterRoutes(api, gmSessionHandler)

	return router
}
