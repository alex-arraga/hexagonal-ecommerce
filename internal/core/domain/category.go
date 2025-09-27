package domain

import (
	"time"
)

// Category is an entity that represents a category of product
type Category struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCategory(name string) (*Category, error) {
	if len(name) == 0 {
		return nil, ErrCategoryNameIsRequire
	}

	now := time.Now()
	return &Category{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Category) UpdateCategory(name string) {
	c.Name = name
}
