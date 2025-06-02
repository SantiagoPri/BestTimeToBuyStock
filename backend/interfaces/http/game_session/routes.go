package game_session

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	sessions := r.Group("/session")
	{
		sessions.POST("/start", h.CreateSession)
		sessions.GET("/state", h.GetSessionState)
		sessions.POST("/update", h.UpdateSessionState)
	}

	r.GET("/leaderboard", h.GetLeaderboard)
}
