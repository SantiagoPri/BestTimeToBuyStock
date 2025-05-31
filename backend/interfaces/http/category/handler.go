package category

import (
	"net/http"
	"strconv"

	categoryApp "backend/application/category"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	categoryService *categoryApp.CategoryService
}

func NewHandler(categoryService *categoryApp.CategoryService) *Handler {
	return &Handler{
		categoryService: categoryService,
	}
}

func (h *Handler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	categories, total, err := h.categoryService.FindPaginated(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}
