package stock_snapshot

import (
	"errors"
	"net/http"
	"strconv"

	snapshotApp "backend/application/stock_snapshot"
	"backend/domain/stock_snapshot"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	snapshotService *snapshotApp.StockSnapshotService
}

func NewHandler(snapshotService *snapshotApp.StockSnapshotService) *Handler {
	return &Handler{
		snapshotService: snapshotService,
	}
}

func (h *Handler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	snapshots, total, err := h.snapshotService.FindPaginated(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snapshots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"snapshots": snapshots,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

func (h *Handler) FindByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	snapshot, err := h.snapshotService.FindByID(id)
	if err != nil {
		if errors.Is(err, stock_snapshot.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Snapshot not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch snapshot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"snapshot": snapshot})
}
