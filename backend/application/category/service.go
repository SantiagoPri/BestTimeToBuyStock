package category

import "backend/domain/category"

type CategoryService struct {
	repo category.Repository
}

func NewCategoryService(repo category.Repository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) FindPaginated(page, limit int) ([]category.Category, int64, error) {
	return s.repo.FindPaginated(page, limit)
}
