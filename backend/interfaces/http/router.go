package http

import (
	"github.com/gin-gonic/gin"

	stock "backend/application/stock"
	stockHttp "backend/interfaces/http/stock"
)

type Router struct {
	stockService *stock.StockService
}

func NewRouter(stockService *stock.StockService) *Router {
	return &Router{
		stockService: stockService,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	stockHandler := stockHttp.NewHandler(r.stockService)
	stockHttp.RegisterRoutes(router, stockHandler)

	return router
}
