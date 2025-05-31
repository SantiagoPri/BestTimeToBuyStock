package category

import (
	"backend/domain/category"
	"backend/infrastructure/repositories"

	"gorm.io/gorm"
)

// Compile-time assertion to ensure CategoryRepository implements category.Repository
var _ category.Repository = (*CategoryRepository)(nil)

type CategoryRepository struct {
	repo repositories.Repository
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		repo: repositories.NewBaseRepository(db, &CategoryEntity{}),
	}
}

func (r *CategoryRepository) Save(c *category.Category) error {
	entity := FromDomain(c)
	return r.repo.Save(entity)
}

func (r *CategoryRepository) FindAll() ([]category.Category, error) {
	var entities []CategoryEntity
	err := r.repo.FindAll(&entities)
	if err != nil {
		return nil, err
	}

	categories := make([]category.Category, len(entities))
	for i, entity := range entities {
		domainCategory := ToDomain(&entity)
		categories[i] = *domainCategory
	}
	return categories, nil
}

func (r *CategoryRepository) DeleteByName(name string) error {
	return r.repo.DeleteByField("name", name, &CategoryEntity{})
}

func (r *CategoryRepository) FindBy(filters map[string]any) (*category.Category, error) {
	var entity CategoryEntity
	err := r.repo.FindOneBy(filters, &entity)
	if err != nil {
		return nil, err
	}
	return ToDomain(&entity), nil
}

func (r *CategoryRepository) FindPaginated(page int, limit int) ([]category.Category, int64, error) {
	var entities []CategoryEntity
	total, err := r.repo.FindPaginated(&entities, page, limit)
	if err != nil {
		return nil, 0, err
	}

	categories := make([]category.Category, len(entities))
	for i, entity := range entities {
		domainCategory := ToDomain(&entity)
		categories[i] = *domainCategory
	}
	return categories, total, nil
}
