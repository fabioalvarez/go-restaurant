package http

import "go-restaurant/internal/category/domain"

// CategoryResponse represents a category Response body
type CategoryResponse struct {
	ID   uint64 `json:"id" example:"1"`
	Name string `json:"name" example:"Foods"`
}

// NewCategoryResponse is a helper function to create a Response body for handling category data
func NewCategoryResponse(category *domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}
