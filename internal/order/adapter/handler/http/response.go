package http

import (
	"go-restaurant/internal/order/domain"
	ophttp "go-restaurant/internal/orderproduct/adapter/handler/http"
	phttp "go-restaurant/internal/payment/adapter/handler/http"
	"time"
)

// OrderResponse represents an order Response body
type OrderResponse struct {
	ID           uint64                        `json:"id" example:"1"`
	UserID       uint64                        `json:"user_id" example:"1"`
	PaymentID    uint64                        `json:"payment_type_id" example:"1"`
	CustomerName string                        `json:"customer_name" example:"John Doe"`
	TotalPrice   float64                       `json:"total_price" example:"100000"`
	TotalPaid    float64                       `json:"total_paid" example:"100000"`
	TotalReturn  float64                       `json:"total_return" example:"0"`
	ReceiptCode  string                        `json:"receipt_id" example:"4979cf6e-d215-4ff8-9d0d-b3e99bcc7750"`
	Products     []ophttp.OrderProductResponse `json:"products"`
	PaymentType  phttp.PaymentResponse         `json:"payment_type"`
	CreatedAt    time.Time                     `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt    time.Time                     `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// NewOrderResponse is a helper function to create a Response body for handling order data
func NewOrderResponse(order *domain.Order) OrderResponse {
	return OrderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		PaymentID:    order.PaymentID,
		CustomerName: order.CustomerName,
		TotalPrice:   order.TotalPrice,
		TotalPaid:    order.TotalPaid,
		TotalReturn:  order.TotalReturn,
		ReceiptCode:  order.ReceiptCode.String(),
		Products:     ophttp.NewOrderProductResponse(order.Products),
		PaymentType:  phttp.NewPaymentResponse(order.Payment),
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}
}
