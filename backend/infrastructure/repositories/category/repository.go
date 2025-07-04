package category

import (
	"backend/domain/category"
	"backend/infrastructure/repositories"
	"backend/pkg/errors"

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
	if err := r.repo.Save(entity); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save category", err)
	}
	return nil
}

func (r *CategoryRepository) FindAll() ([]category.Category, error) {
	var entities []CategoryEntity
	err := r.repo.FindAll(&entities)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find categories", err)
	}

	categories := make([]category.Category, len(entities))
	for i, entity := range entities {
		domainCategory := ToDomain(&entity)
		categories[i] = *domainCategory
	}
	return categories, nil
}

func (r *CategoryRepository) DeleteByName(name string) error {
	if err := r.repo.DeleteByField("name", name, &CategoryEntity{}); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete category", err)
	}
	return nil
}

func (r *CategoryRepository) FindBy(filters map[string]any) (*category.Category, error) {
	var entity CategoryEntity
	err := r.repo.FindOneBy(filters, &entity)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrNotFound, "category not found")
		}
		return nil, errors.Wrap(errors.ErrInternal, "failed to find category", err)
	}
	return ToDomain(&entity), nil
}

func (r *CategoryRepository) FindPaginated(page int, limit int) ([]category.Category, int64, error) {
	var entities []CategoryEntity
	total, err := r.repo.FindPaginated(&entities, page, limit)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to find paginated categories", err)
	}

	categories := make([]category.Category, len(entities))
	for i, entity := range entities {
		domainCategory := ToDomain(&entity)
		categories[i] = *domainCategory
	}
	return categories, total, nil
}
