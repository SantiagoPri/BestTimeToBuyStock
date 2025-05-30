package category

import (
	"backend/domain/category"
	"time"
)

type CategoryEntity struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);unique" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (CategoryEntity) TableName() string {
	return "categories"
}

func ToDomain(e *CategoryEntity) *category.Category {
	if e == nil {
		return nil
	}
	return &category.Category{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
	}
}

func FromDomain(c *category.Category) *CategoryEntity {
	if c == nil {
		return nil
	}
	return &CategoryEntity{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
	}
}
