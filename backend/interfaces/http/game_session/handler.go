package game_session

import (
	"backend/application/game_session"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service game_session.Service
}

func NewHandler(service game_session.Service) *Handler {
	return &Handler{service: service}
}

type createSessionRequest struct {
	Username   string   `json:"username" binding:"required"`
	Categories []string `json:"categories" binding:"required,len=3"`
}

type createSessionResponse struct {
	SessionID string `json:"sessionId"`
}

type updateStateRequest struct {
	Status string  `json:"status" binding:"required"`
	Cash   float64 `json:"cash" binding:"required"`
}

func (h *Handler) CreateSession(c *gin.Context) {
	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := h.service.Create(req.Username, req.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createSessionResponse{SessionID: sessionID})
}

func (h *Handler) GetSessionState(c *gin.Context) {
	sessionID := extractBearerToken(c)
	if sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid session token"})
		return
	}

	state, err := h.service.GetState(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	leaderboard, err := h.service.GetLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
