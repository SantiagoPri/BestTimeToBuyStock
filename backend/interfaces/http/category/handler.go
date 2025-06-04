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

// @Summary List all categories
// @Description Get a paginated list of stock categories
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "List of categories with pagination info"
// @Failure 500 {object} errors.Error "Internal server error"
// @Router /categories [get]
func (h *Handler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	categories, total, err := h.categoryService.FindPaginated(page, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}
