package http

import (
	"go-restaurant/internal/payment/domain"
)

// PaymentResponse represents a payment Response body
type PaymentResponse struct {
	ID   uint64             `json:"id" example:"1"`
	Name string             `json:"name" example:"Tunai"`
	Type domain.PaymentType `json:"type" example:"CASH"`
	Logo string             `json:"logo" example:"https://example.com/cash.png"`
}

// NewPaymentResponse is a helper function to create a Response body for handling payment data
func NewPaymentResponse(payment *domain.Payment) PaymentResponse {
	return PaymentResponse{
		ID:   payment.ID,
		Name: payment.Name,
		Type: payment.Type,
		Logo: payment.Logo,
	}
}
