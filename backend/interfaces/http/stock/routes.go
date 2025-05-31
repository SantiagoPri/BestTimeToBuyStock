package stock

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	stocks := r.Group("/stocks")
	{
		stocks.GET("", h.FindAll)
		stocks.GET("/:param", h.FindOne)
		// stocks.GET("/:ticker", h.FindBy)
	}
}
