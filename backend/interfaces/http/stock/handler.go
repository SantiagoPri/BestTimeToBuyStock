package stock

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	stock "backend/application/stock"
)

type Handler struct {
	stockService *stock.StockService
}

func NewHandler(stockService *stock.StockService) *Handler {
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

	c.JSON(http.StatusOK, gin.H{"stocks": stocks, "total": total})
}
