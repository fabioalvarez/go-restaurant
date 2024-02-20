package http

import (
	"go-restaurant/internal/category/adapter/handler/http"
	"go-restaurant/internal/product/domain"
	"time"
)

// ProductResponse represents a product Response body
type ProductResponse struct {
	ID        uint64                `json:"id" example:"1"`
	SKU       string                `json:"sku" example:"9a4c25d3-9786-492c-b084-85cb75c1ee3e"`
	Name      string                `json:"name" example:"Chiki Ball"`
	Stock     int64                 `json:"stock" example:"100"`
	Price     float64               `json:"price" example:"5000"`
	Image     string                `json:"image" example:"https://example.com/chiki-ball.png"`
	Category  http.CategoryResponse `json:"category"`
	CreatedAt time.Time             `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time             `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// NewProductResponse is a helper function to create a Response body for handling product data
func NewProductResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		SKU:       product.SKU.String(),
		Name:      product.Name,
		Stock:     product.Stock,
		Price:     product.Price,
		Image:     product.Image,
		Category:  http.NewCategoryResponse(product.Category),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
