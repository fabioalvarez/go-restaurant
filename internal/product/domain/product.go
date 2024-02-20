package domain

import (
	"github.com/google/uuid"
	"go-restaurant/internal/category/domain"
	"time"
)

// Product is an entity that represents a product
type Product struct {
	ID         uint64
	CategoryID uint64
	SKU        uuid.UUID
	Name       string
	Stock      int64
	Price      float64
	Image      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Category   *domain.Category
}
