package category

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	categories := r.Group("/categories")
	{
		categories.GET("", h.FindAll)
	}
}
