package category

type Repository interface {
	Save(category *Category) error
	FindAll() ([]Category, error)
	FindBy(filters map[string]any) (*Category, error)
	DeleteByName(name string) error
	FindPaginated(page int, limit int) ([]Category, int64, error)
}
