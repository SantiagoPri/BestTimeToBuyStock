package gm_session

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	gm := r.Group("/game")
	{
		gm.GET("/week/:week", h.GetWeekData)
	}
}
