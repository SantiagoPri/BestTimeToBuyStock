package game_session

import (
	"backend/application/game_session"
	"backend/pkg/errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler manages game session HTTP endpoints
type Handler struct {
	service game_session.Service
}

func NewHandler(service game_session.Service) *Handler {
	return &Handler{service: service}
}

// @Description Request body for creating a new session
type createSessionRequest struct {
	// @Description User's display name
	// @Required
	Username string `json:"username" binding:"required" example:"john_doe"`
	// @Description List of exactly 3 stock categories
	// @Required
	// @MinItems 3
	// @MaxItems 3
	Categories []string `json:"categories" binding:"required,len=3" example:"['tech','healthcare','energy']"`
}

// @Description Response for session creation
type createSessionResponse struct {
	// @Description Unique session identifier
	SessionID string `json:"sessionId" example:"abc123def456"`
}

type updateStateRequest struct {
	Status string  `json:"status" binding:"required"`
	Cash   float64 `json:"cash" binding:"required"`
}

// @Description Request body for trading stocks
type tradeRequest struct {
	// @Description Stock ticker symbol
	// @Required
	Ticker string `json:"ticker" binding:"required" example:"AAPL"`
	// @Description Number of shares to trade
	// @Required
	// @Minimum 1
	Quantity int `json:"quantity" binding:"required" example:"100"`
}

// @Summary Create a new game session
// @Description Creates a new game session for a user with selected stock categories
// @Tags Game Session
// @Accept json
// @Produce json
// @Param request body createSessionRequest true "Session creation parameters"
// @Success 201 {object} createSessionResponse "Session created successfully"
// @Failure 400 {object} errors.Error "Invalid input - Username missing or categories != 3"
// @Failure 500 {object} errors.Error "Internal server error"
// @Router /sessions [post]
func (h *Handler) CreateSession(c *gin.Context) {
	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(errors.ErrInvalidInput, "invalid request body", err))
		return
	}

	sessionID, err := h.service.Create(req.Username, req.Categories)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createSessionResponse{SessionID: sessionID})
}

// @Summary Get current session state
// @Description Retrieves the current state of a game session including cash, holdings, and status
// @Tags Game Session
// @Produce json
// @Security BearerAuth
// @Success 200 {object} game_session.GameSession "Current session state"
// @Failure 401 {object} errors.Error "Unauthorized - Invalid session"
// @Failure 404 {object} errors.Error "Session not found"
// @Failure 503 {object} errors.Error "Session is no longer active"
// @Router /sessions/state [get]
func (h *Handler) GetSessionState(c *gin.Context) {
	sessionID := extractBearerToken(c)
	if sessionID == "" {
		_ = c.Error(errors.New(errors.ErrUnauthorized, "missing or invalid session token"))
		return
	}

	state, err := h.service.GetState(sessionID)
	if err != nil {
		log.Printf("Handler: Error getting state for session %s: %v", sessionID, err)
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, state)
}

// @Summary Get leaderboard
// @Description Retrieves top 10 finished sessions ordered by total balance
// @Tags Game Session
// @Produce json
// @Success 200 {array} game_session.GameSession "Leaderboard entries"
// @Failure 500 {object} errors.Error "Internal server error"
// @Router /leaderboard [get]
func (h *Handler) GetLeaderboard(c *gin.Context) {
	leaderboard, err := h.service.GetLeaderboard()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// @Summary Buy stocks
// @Description Purchase a specified quantity of a stock in the current session
// @Tags Trading
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tradeRequest true "Buy order details"
// @Success 200 "Purchase successful"
// @Failure 400 {object} errors.Error "Invalid input - Missing ticker or quantity"
// @Failure 401 {object} errors.Error "Unauthorized - Invalid session"
// @Failure 404 {object} errors.Error "Stock not found"
// @Failure 422 {object} errors.Error "Insufficient funds"
// @Router /sessions/buy [post]
func (h *Handler) BuyStock(c *gin.Context) {
	sessionID := extractBearerToken(c)
	if sessionID == "" {
		_ = c.Error(errors.New(errors.ErrUnauthorized, "missing or invalid session token"))
		return
	}

	var req tradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(errors.ErrInvalidInput, "invalid request body", err))
		return
	}

	if err := h.service.Buy(sessionID, req.Ticker, req.Quantity); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Sell stocks
// @Description Sell a specified quantity of a stock in the current session
// @Tags Trading
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tradeRequest true "Sell order details"
// @Success 200 "Sale successful"
// @Failure 400 {object} errors.Error "Invalid input - Missing ticker or quantity"
// @Failure 401 {object} errors.Error "Unauthorized - Invalid session"
// @Failure 404 {object} errors.Error "Stock not found or insufficient holdings"
// @Router /sessions/sell [post]
func (h *Handler) SellStock(c *gin.Context) {
	sessionID := extractBearerToken(c)
	if sessionID == "" {
		_ = c.Error(errors.New(errors.ErrUnauthorized, "missing or invalid session token"))
		return
	}

	var req tradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(errors.ErrInvalidInput, "invalid request body", err))
		return
	}

	if err := h.service.Sell(sessionID, req.Ticker, req.Quantity); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Advance to next week
// @Description Advances the session to the next week, updating stock prices
// @Tags Game Session
// @Security BearerAuth
// @Success 200 "Advanced to next week"
// @Failure 401 {object} errors.Error "Unauthorized - Invalid session"
// @Failure 400 {object} errors.Error "Cannot advance beyond week 5"
// @Router /sessions/advance [post]
func (h *Handler) AdvanceWeek(c *gin.Context) {
	sessionID := extractBearerToken(c)
	if sessionID == "" {
		_ = c.Error(errors.New(errors.ErrUnauthorized, "missing or invalid session token"))
		return
	}

	if err := h.service.AdvanceWeek(sessionID); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary End session
// @Description Ends the current session, selling all holdings at current prices
// @Tags Game Session
// @Security BearerAuth
// @Success 202 {object} game_session.GameSession "Session ended successfully"
// @Failure 401 {object} errors.Error "Unauthorized - Invalid session"
// @Failure 400 {object} errors.Error "Can only end session in week 5"
// @Router /sessions/end [post]
func (h *Handler) EndSession(c *gin.Context) {
	sessionID := extractBearerToken(c)
	if sessionID == "" {
		_ = c.Error(errors.New(errors.ErrUnauthorized, "missing or invalid session token"))
		return
	}

	gameSession, err := h.service.EndSession(sessionID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusAccepted, gameSession)
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
