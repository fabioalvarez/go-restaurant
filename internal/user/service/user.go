package service

import (
	"context"
	"errors"
	cmdomain "go-restaurant/internal/common/domain"
	cmport "go-restaurant/internal/common/port"
	cmutil "go-restaurant/internal/common/util"
	"go-restaurant/internal/user/domain"
	"go-restaurant/internal/user/port"
)

/*UserService implements port.UserService interface
 * and provides access to the user repository
 * and cache service
 */
type UserService struct {
	repo  port.UserRepository
	cache cmport.CacheRepository
}

// NewUserService creates a new user service instance
func NewUserService(repo port.UserRepository, cache cmport.CacheRepository) *UserService {
	return &UserService{
		repo,
		cache,
	}
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := cmutil.HashPassword(user.Password)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	user.Password = hashedPassword

	user, err = us.repo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, cmdomain.ErrConflictingData) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("user", user.ID)
	userSerialized, err := cmutil.Serialize(user)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return user, nil
}

// GetUser gets a user by ID
func (us *UserService) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	var user *domain.User

	cacheKey := cmutil.GenerateCacheKey("user", id)
	cachedUser, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := cmutil.Deserialize(cachedUser, &user)
		if err != nil {
			return nil, cmdomain.ErrInternal
		}
		return user, nil
	}

	user, err = us.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	userSerialized, err := cmutil.Serialize(user)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return user, nil
}

// ListUsers lists all users
func (us *UserService) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	var users []domain.User

	params := cmutil.GenerateCacheKeyParams(skip, limit)
	cacheKey := cmutil.GenerateCacheKey("users", params)

	cachedUsers, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := cmutil.Deserialize(cachedUsers, &users)
		if err != nil {
			return nil, cmdomain.ErrInternal
		}
		return users, nil
	}

	users, err = us.repo.ListUsers(ctx, skip, limit)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	usersSerialized, err := cmutil.Serialize(users)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, usersSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return users, nil
}

// UpdateUser updates a user's name, email, and password
func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	emptyData := user.Name == "" &&
		user.Email == "" &&
		user.Password == "" &&
		user.Role == ""
	sameData := existingUser.Name == user.Name &&
		existingUser.Email == user.Email &&
		existingUser.Role == user.Role
	if emptyData || sameData {
		return nil, cmdomain.ErrNoUpdatedData
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = cmutil.HashPassword(user.Password)
		if err != nil {
			return nil, cmdomain.ErrInternal
		}
	}

	user.Password = hashedPassword

	_, err = us.repo.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, cmdomain.ErrConflictingData) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("user", user.ID)

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	userSerialized, err := cmutil.Serialize(user)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return err
		}
		return cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("user", id)

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return cmdomain.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return cmdomain.ErrInternal
	}

	return us.repo.DeleteUser(ctx, id)
}
