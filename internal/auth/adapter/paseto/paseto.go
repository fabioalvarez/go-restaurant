package paseto

import (
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"go-restaurant/internal/auth/domain"
	"go-restaurant/internal/auth/port"
	"go-restaurant/internal/common/adapter/config"
	cmdomain "go-restaurant/internal/common/domain"
	udomain "go-restaurant/internal/user/domain"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

/*Token implements port.TokenService interface
 * and provides access to the paseto library
 */
type Token struct {
	paseto       *paseto.V2
	symmetricKey []byte
	duration     time.Duration
}

// New creates a new paseto instance
func New(config *config.Token) (port.TokenService, error) {
	symmetricKey := config.SymmetricKey
	durationStr := config.Duration

	validSymmetricKey := len(symmetricKey) == chacha20poly1305.KeySize
	if !validSymmetricKey {
		return nil, cmdomain.ErrInvalidTokenSymmetricKey
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, err
	}

	return &Token{
		paseto.NewV2(),
		[]byte(symmetricKey),
		duration,
	}, nil
}

// CreateToken creates a new paseto token
func (pt *Token) CreateToken(user *udomain.User) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", cmdomain.ErrTokenCreation
	}

	payload := domain.TokenPayload{
		ID:        id,
		UserID:    user.ID,
		Role:      user.Role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(pt.duration),
	}

	token, err := pt.paseto.Encrypt(pt.symmetricKey, payload, nil)
	if err != nil {
		return "", cmdomain.ErrTokenCreation
	}

	return token, nil
}

// VerifyToken verifies the paseto token
func (pt *Token) VerifyToken(token string) (*domain.TokenPayload, error) {
	var payload domain.TokenPayload

	err := pt.paseto.Decrypt(token, pt.symmetricKey, &payload, nil)
	if err != nil {
		return nil, cmdomain.ErrInvalidToken
	}

	isExpired := time.Now().After(payload.ExpiredAt)
	if isExpired {
		return nil, cmdomain.ErrExpiredToken
	}

	return &payload, nil
}
