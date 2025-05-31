package category

import (
	"backend/domain/category"
	"strconv"
	"time"
)

type CategoryEntity struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100)" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (CategoryEntity) TableName() string {
	return "categories"
}

func ToDomain(e *CategoryEntity) *category.Category {
	if e == nil {
		return nil
	}
	return &category.Category{
		ID:        strconv.FormatUint(uint64(e.ID), 10),
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
	}
}

func FromDomain(c *category.Category) *CategoryEntity {
	if c == nil {
		return nil
	}
	id, err := strconv.ParseUint(c.ID, 10, 64)
	if err != nil {
		id = 0 // For create operations
	}
	return &CategoryEntity{
		ID:        uint(id),
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
	}
}
