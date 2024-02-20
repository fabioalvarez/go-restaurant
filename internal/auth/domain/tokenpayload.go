package domain

import (
	"github.com/google/uuid"
	udomain "go-restaurant/internal/user/domain"
	"time"
)

// TokenPayload is an entity that represents the payload of the token
type TokenPayload struct {
	ID        uuid.UUID
	UserID    uint64
	Role      udomain.UserRole
	IssuedAt  time.Time
	ExpiredAt time.Time
}
