package service

import (
	"context"
	"errors"
	cmdomain "go-restaurant/internal/common/domain"
	cmport "go-restaurant/internal/common/port"
	cmutil "go-restaurant/internal/common/util"
	"go-restaurant/internal/payment/domain"
	"go-restaurant/internal/payment/port"
)

/*PaymentService implements port.PaymentService interface
 * and provides access to the payment repository
 * and cache service
 */
type PaymentService struct {
	repo  port.PaymentRepository
	cache cmport.CacheRepository
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(repo port.PaymentRepository, cache cmport.CacheRepository) *PaymentService {
	return &PaymentService{
		repo,
		cache,
	}
}

// CreatePayment creates a new payment
func (ps *PaymentService) CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	payment, err := ps.repo.CreatePayment(ctx, payment)
	if err != nil {
		if errors.Is(err, cmdomain.ErrConflictingData) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("payment", payment.ID)
	paymentSerialized, err := cmutil.Serialize(payment)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, paymentSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = ps.cache.DeleteByPrefix(ctx, "payments:*")
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return payment, nil
}

// GetPayment retrieves a payment by id
func (ps *PaymentService) GetPayment(ctx context.Context, id uint64) (*domain.Payment, error) {
	var payment *domain.Payment

	cacheKey := cmutil.GenerateCacheKey("payment", id)
	cachedPayment, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		err := cmutil.Deserialize(cachedPayment, &payment)
		if err != nil {
			return nil, cmdomain.ErrInternal
		}

		return payment, nil
	}

	payment, err = ps.repo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	paymentSerialized, err := cmutil.Serialize(payment)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, paymentSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return payment, nil
}

// ListPayments retrieves a list of payments
func (ps *PaymentService) ListPayments(ctx context.Context, skip, limit uint64) ([]domain.Payment, error) {
	var payments []domain.Payment

	params := cmutil.GenerateCacheKeyParams(skip, limit)
	cacheKey := cmutil.GenerateCacheKey("payments", params)

	cachedPayments, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		err := cmutil.Deserialize(cachedPayments, &payments)
		if err != nil {
			return nil, cmdomain.ErrInternal
		}

		return payments, nil
	}

	payments, err = ps.repo.ListPayments(ctx, skip, limit)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	paymentsSerialized, err := cmutil.Serialize(payments)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, paymentsSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return payments, nil

}

// UpdatePayment updates a payment
func (ps *PaymentService) UpdatePayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	existingPayment, err := ps.repo.GetPaymentByID(ctx, payment.ID)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	emptyData := payment.Name == "" && payment.Type == "" && payment.Logo == ""
	sameData := existingPayment.Name == payment.Name && existingPayment.Type == payment.Type && existingPayment.Logo == payment.Logo
	if emptyData || sameData {
		return nil, cmdomain.ErrNoUpdatedData
	}

	_, err = ps.repo.UpdatePayment(ctx, payment)
	if err != nil {
		if errors.Is(err, cmdomain.ErrConflictingData) {
			return nil, err
		}
		return nil, cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("payment", payment.ID)

	err = ps.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	paymentSerialized, err := cmutil.Serialize(payment)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, paymentSerialized, 0)
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	err = ps.cache.DeleteByPrefix(ctx, "payments:*")
	if err != nil {
		return nil, cmdomain.ErrInternal
	}

	return payment, nil
}

// DeletePayment deletes a payment
func (ps *PaymentService) DeletePayment(ctx context.Context, id uint64) error {
	_, err := ps.repo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, cmdomain.ErrDataNotFound) {
			return err
		}
		return cmdomain.ErrInternal
	}

	cacheKey := cmutil.GenerateCacheKey("payment", id)

	err = ps.cache.Delete(ctx, cacheKey)
	if err != nil {
		return cmdomain.ErrInternal
	}

	err = ps.cache.DeleteByPrefix(ctx, "payments:*")
	if err != nil {
		return cmdomain.ErrInternal
	}

	return ps.repo.DeletePayment(ctx, id)
}
