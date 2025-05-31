package repositories

import (
	"strconv"

	"gorm.io/gorm"
)

type Repository interface {
	Save(entity any) error
	FindAll(out any) error
	FindByField(field string, value any, out any) error
	DeleteByField(field string, value any, model any) error
	FindOneBy(filters map[string]any, out any) error
	FindPaginated(out any, page int, limit int) (int64, error)
}

type BaseRepository struct {
	db    *gorm.DB
	model any
}

func NewBaseRepository(db *gorm.DB, model any) *BaseRepository {
	return &BaseRepository{
		db:    db,
		model: model,
	}
}

// ParseIDFilter converts string IDs in filters to uint for GORM queries
func ParseIDFilter(filters map[string]any, idFields ...string) map[string]any {
	result := make(map[string]any)
	for k, v := range filters {
		if strVal, ok := v.(string); ok {
			// Check if this is an ID field that needs conversion
			for _, idField := range idFields {
				if k == idField {
					if id, err := strconv.ParseUint(strVal, 10, 64); err == nil {
						result[k] = uint(id)
					} else {
						result[k] = 0 // Default for invalid ID
					}
					break
				}
			}
			if _, exists := result[k]; !exists {
				result[k] = v // Keep original value if not an ID field
			}
		} else {
			result[k] = v // Keep non-string values as is
		}
	}
	return result
}

func (r *BaseRepository) Save(entity any) error {
	return r.db.Save(entity).Error
}

func (r *BaseRepository) FindAll(out any) error {
	return r.db.Find(out).Error
}

func (r *BaseRepository) FindByField(field string, value any, out any) error {
	return r.db.Where(field+" = ?", value).First(out).Error
}

func (r *BaseRepository) DeleteByField(field string, value any, model any) error {
	return r.db.Where(field+" = ?", value).Delete(model).Error
}

func (r *BaseRepository) FindOneBy(filters map[string]any, out any) error {
	return r.db.Where(filters).First(out).Error
}

func (r *BaseRepository) FindPaginated(out any, page int, limit int) (int64, error) {
	var total int64

	// Count total records for pagination
	if err := r.db.Model(r.model).Count(&total).Error; err != nil {
		return 0, err
	}

	offset := (page - 1) * limit

	err := r.db.
		Limit(limit).
		Offset(offset).
		Find(out).
		Error

	if err != nil {
		return 0, err
	}

	return total, nil
}
