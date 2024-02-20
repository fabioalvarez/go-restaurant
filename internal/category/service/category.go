package service

import (
	"context"
	"errors"
	"go-restaurant/internal/category/domain"
	"go-restaurant/internal/category/port"
	cmdomain "go-restaurant/internal/common/domain"
	cmport "go-restaurant/internal/common/port"
	"go-restaurant/internal/common/util"
)

/*CategoryService implements port.CategoryService interface
 * and provides access to the category repository
 * and cache service
 */
type CategoryService struct {
	repo  port.CategoryRepository
	cache cmport.CacheRepository
}

// NewCategoryService creates a new category service instance
func NewCategoryService(repo port.CategoryRepository, cache cmport.CacheRepository) *CategoryService {
	return &CategoryService{
		repo,
		cache,
	}
}

// CreateCategory creates a new category
func (cs *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	category, err := cs.repo.CreateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, cmdomain.ErrConflictingData) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("category", category.ID)
	categorySerialized, err := cmutil.Serialize(category)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = cs.cache.DeleteByPrefix(ctx, "categories:*")
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return category, nil
}

// GetCategory retrieves a category by id
func (cs *CategoryService) GetCategory(ctx context.Context, id uint64) (*domain.Category, error) {
	var category *domain.Category

	cacheKey := cmutil.GenerateCacheKey("category", id)
	cachedCategory, err := cs.cache.Get(ctx, cacheKey)
	if err == nil {
		err := cmutil.Deserialize(cachedCategory, &category)
		if err != nil {
			return nil, cmdomain.ErrInternal
		}
		return category, nil
	}

	category, err = cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	categorySerialized, err := cmutil.Serialize(category)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return category, nil
}

// ListCategories retrieves a list of categories
func (cs *CategoryService) ListCategories(ctx context.Context, skip, limit uint64) ([]domain.Category, error) {
	var categories []domain.Category

	params := cmutil.GenerateCacheKeyParams(skip, limit)
	cacheKey := cmutil.GenerateCacheKey("categories", params)

	cachedCategories, err := cs.cache.Get(ctx, cacheKey)
	if err == nil {
		err := cmutil.Deserialize(cachedCategories, &categories)
		if err != nil {
			return nil, cmdomain.ErrInternal
		}

		return categories, nil
	}

	categories, err = cs.repo.ListCategories(ctx, skip, limit)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	categoriesSerialized, err := cmutil.Serialize(categories)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = cs.cache.Set(ctx, cacheKey, categoriesSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return categories, nil
}

// UpdateCategory updates a category
func (cs *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	existingCategory, err := cs.repo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	emptyData := category.Name == ""
	sameData := existingCategory.Name == category.Name
	if emptyData || sameData {
		return nil, cmdomain.ErrNoUpdatedData
	}

	_, err = cs.repo.UpdateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, cmdomain.ErrConflictingData) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("category", category.ID)

	err = cs.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	categorySerialized, err := cmutil.Serialize(category)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = cs.cache.DeleteByPrefix(ctx, "categories:*")
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return category, nil
}

// DeleteCategory deletes a category
func (cs *CategoryService) DeleteCategory(ctx context.Context, id uint64) error {
	_, err := cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return err
		}
		return cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("category", id)

	err = cs.cache.Delete(ctx, cacheKey)
	if err != nil {
		return cmdomain.ErrInternal
	}

	err = cs.cache.DeleteByPrefix(ctx, "categories:*")
	if err != nil {
		return cmdomain.ErrInternal
	}

	return cs.repo.DeleteCategory(ctx, id)
}
