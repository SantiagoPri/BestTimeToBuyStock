package stock

type Repository interface {
	FindPaginated(page int, limit int) ([]Stock, int64, error)
	FindBy(filters map[string]any) (*Stock, error)
}
