package stock_snapshot

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	snapshots := r.Group("/snapshots")
	{
		snapshots.GET("", h.FindAll)
		snapshots.GET("/:id", h.FindByID)
	}
}
