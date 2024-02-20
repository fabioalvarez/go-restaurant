package http

import (
	"go-restaurant/internal/orderproduct/domain"
	phttp "go-restaurant/internal/product/adapter/handler/http"
	"time"
)

// OrderProductResponse represents an order product Response body
type OrderProductResponse struct {
	ID               uint64                `json:"id" example:"1"`
	OrderID          uint64                `json:"order_id" example:"1"`
	ProductID        uint64                `json:"product_id" example:"1"`
	Quantity         int64                 `json:"qty" example:"1"`
	Price            float64               `json:"price" example:"100000"`
	TotalNormalPrice float64               `json:"total_normal_price" example:"100000"`
	TotalFinalPrice  float64               `json:"total_final_price" example:"100000"`
	Product          phttp.ProductResponse `json:"product"`
	CreatedAt        time.Time             `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt        time.Time             `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// NewOrderProductResponse is a helper function to create a Response body for handling order product data
func NewOrderProductResponse(orderProduct []domain.OrderProduct) []OrderProductResponse {
	var orderProductResponses []OrderProductResponse

	for _, orderProduct := range orderProduct {
		orderProductResponses = append(orderProductResponses, OrderProductResponse{
			ID:               orderProduct.ID,
			OrderID:          orderProduct.OrderID,
			ProductID:        orderProduct.ProductID,
			Quantity:         orderProduct.Quantity,
			Price:            orderProduct.Product.Price,
			TotalNormalPrice: orderProduct.TotalPrice,
			TotalFinalPrice:  orderProduct.TotalPrice,
			Product:          phttp.NewProductResponse(orderProduct.Product),
			CreatedAt:        orderProduct.CreatedAt,
			UpdatedAt:        orderProduct.UpdatedAt,
		})
	}

	return orderProductResponses
}
