package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	category "backend/application/category"
	stock "backend/application/stock"
	snapshot "backend/application/stock_snapshot"
	categoryHttp "backend/interfaces/http/category"
	stockHttp "backend/interfaces/http/stock"
	snapshotHttp "backend/interfaces/http/stock_snapshot"
)

type Router struct {
	stockService    *stock.StockService
	categoryService *category.CategoryService
	snapshotService *snapshot.StockSnapshotService
}

func NewRouter(stockService *stock.StockService, categoryService *category.CategoryService, snapshotService *snapshot.StockSnapshotService) *Router {
	return &Router{
		stockService:    stockService,
		categoryService: categoryService,
		snapshotService: snapshotService,
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

	return router
}
