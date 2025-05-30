package stock

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	stockApp "backend/application/stock"
	stockDomain "backend/domain/stock"
)

type Handler struct {
	stockService *stockApp.StockService
}

func NewHandler(stockService *stockApp.StockService) *Handler {
	return &Handler{
		stockService: stockService,
	}
}

func (h *Handler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	stocks, total, err := h.stockService.FindPaginated(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stocks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stocks": stocks,
		"pagination": gin.H{
			"total":       total,
			"currentPage": page,
			"limit":       limit,
		},
	})
}

func (h *Handler) FindOne(c *gin.Context) {
	param := c.Param("param")
	if param == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter is required"})
		return
	}

	// Check if it's an ID (19 digits)
	if matched, _ := regexp.MatchString(`^\d{19}$`, param); matched {
		stock, err := h.stockService.FindOne("id", param)
		if err != nil {
			if errors.Is(err, stockDomain.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"stock": stock})
		return
	}

	// Check if it's a ticker (3-6 uppercase characters)
	if matched, _ := regexp.MatchString(`^[A-Z]{3,6}$`, param); matched {
		stock, err := h.stockService.FindOne("ticker", param)
		if err != nil {
			if errors.Is(err, stockDomain.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"stock": stock})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter format. Must be either a 19-digit ID or a 3-6 character ticker symbol"})
}
