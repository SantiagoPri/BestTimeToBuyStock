package gm_session

import (
	"backend/application/gm_session"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler manages GM session HTTP endpoints
type Handler struct {
	service gm_session.Service
}

func NewHandler(service gm_session.Service) *Handler {
	return &Handler{service: service}
}

// @Summary Get week data
// @Description Get the game master's data for a specific week
// @Tags GM Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param week path int true "Week number (1-5)" minimum(1) maximum(5)
// @Success 200 {object} gm_session.GMWeekData "Week data including stock prices and news"
// @Failure 400 {object} errors.Error "Invalid week number"
// @Failure 401 {object} errors.Error "Unauthorized - Invalid session"
// @Failure 404 {object} errors.Error "Week data not found"
// @Failure 500 {object} errors.Error "Internal server error"
// @Router /gm/week/{week} [get]
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
