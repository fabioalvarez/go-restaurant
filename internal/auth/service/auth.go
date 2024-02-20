package service

import (
	"context"
	"errors"
	"go-restaurant/internal/auth/port"
	"go-restaurant/internal/common/domain"
	cmutil "go-restaurant/internal/common/util"
	uport "go-restaurant/internal/user/port"
)

/*AuthService implements port.AuthService interface
 * and provides access to the user repository
 * and token service
 */
type AuthService struct {
	repo uport.UserRepository
	ts   port.TokenService
}

// NewAuthService creates a new auth service instance
func NewAuthService(repo uport.UserRepository, ts port.TokenService) *AuthService {
	return &AuthService{
		repo,
		ts,
	}
}

// Login gives a registered user an access token if the credentials are valid
func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", domain.ErrInternal
	}

	err = cmutil.ComparePassword(password, user.Password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	return accessToken, nil
}
