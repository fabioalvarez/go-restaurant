package domain

import (
	opdomain "go-restaurant/internal/orderproduct/domain"
	pdomain "go-restaurant/internal/payment/domain"
	udomain "go-restaurant/internal/user/domain"
	"time"

	"github.com/google/uuid"
)

// Order is an entity that represents an order
type Order struct {
	ID           uint64
	UserID       uint64
	PaymentID    uint64
	CustomerName string
	TotalPrice   float64
	TotalPaid    float64
	TotalReturn  float64
	ReceiptCode  uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	User         *udomain.User
	Payment      *pdomain.Payment
	Products     []opdomain.OrderProduct
}
