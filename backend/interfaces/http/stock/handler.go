package stock

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	stockApp "backend/application/stock"
	"backend/pkg/errors"
)

// Handler manages stock-related HTTP endpoints
type Handler struct {
	stockService *stockApp.StockService
}

func NewHandler(stockService *stockApp.StockService) *Handler {
	return &Handler{
		stockService: stockService,
	}
}

// @Summary List all stocks
// @Description Get a paginated list of stocks
// @Tags Stocks
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "List of stocks with pagination info"
// @Failure 500 {object} errors.Error "Internal server error"
// @Router /stocks [get]
func (h *Handler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	stocks, total, err := h.stockService.FindPaginated(page, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stocks":      stocks,
		"total":       total,
		"currentPage": page,
		"limit":       limit,
	})
}

// @Summary Get stock by ID or ticker
// @Description Get a stock by its ID (19 digits) or ticker symbol (3-6 uppercase characters)
// @Tags Stocks
// @Accept json
// @Produce json
// @Param param path string true "Stock ID or ticker" example("AAPL")
// @Success 200 {object} map[string]interface{} "Stock details"
// @Failure 400 {object} errors.Error "Invalid input - wrong format"
// @Failure 404 {object} errors.Error "Stock not found"
// @Failure 500 {object} errors.Error "Internal server error"
// @Router /stocks/{param} [get]
func (h *Handler) FindOne(c *gin.Context) {
	param := c.Param("param")
	if param == "" {
		_ = c.Error(errors.New(errors.ErrInvalidInput, "parameter is required"))
		return
	}

	// Check if it's an ID (19 digits)
	if matched, _ := regexp.MatchString(`^\d{19}$`, param); matched {
		stock, err := h.stockService.FindOne("id", param)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"stock": stock})
		return
	}

	// Check if it's a ticker (3-6 uppercase characters)
	if matched, _ := regexp.MatchString(`^[A-Z]{3,6}$`, param); matched {
		stock, err := h.stockService.FindOne("ticker", param)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"stock": stock})
		return
	}

	_ = c.Error(errors.New(errors.ErrInvalidInput, "invalid parameter format: must be either a 19-digit ID or a 3-6 character ticker symbol"))
}
