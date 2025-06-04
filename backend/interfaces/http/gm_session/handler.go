package gm_session

import (
	"backend/application/gm_session"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service gm_session.Service
}

func NewHandler(service gm_session.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetWeekData(c *gin.Context) {
	sessionID := extractBearerToken(c)
	if sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid session token"})
		return
	}

	weekStr := c.Param("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid week number"})
		return
	}

	weekData, err := h.service.GetWeekData(sessionID, week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, weekData)
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
