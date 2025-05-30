package stock_snapshot

type Repository interface {
	Save(snapshot *StockSnapshot) error
	FindAll() ([]StockSnapshot, error)
	FindBy(filters map[string]any) (*StockSnapshot, error)
	FindByCategory(category string) ([]StockSnapshot, error)
}
