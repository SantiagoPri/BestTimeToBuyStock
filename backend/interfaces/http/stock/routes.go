package stock

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	stocks := r.Group("/stocks")
	{
		stocks.GET("", h.FindAll)
		// stocks.GET("/:ticker", h.FindBy)
	}
}
