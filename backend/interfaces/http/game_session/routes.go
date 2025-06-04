package game_session

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	sessions := r.Group("/session")
	{
		sessions.POST("/start", h.CreateSession)
		sessions.GET("/state", h.GetSessionState)
		sessions.POST("/buy", h.BuyStock)
		sessions.POST("/sell", h.SellStock)
		sessions.POST("/advance", h.AdvanceWeek)
		sessions.POST("/end", h.EndSession)
	}

	r.GET("/leaderboard", h.GetLeaderboard)
}
