package http

import (
	"github.com/go-playground/validator/v10"
	"go-restaurant/internal/payment/domain"
)

// PaymentTypeValidator is a custom validator for validating payment types
var PaymentTypeValidator validator.Func = func(fl validator.FieldLevel) bool {
	paymentType := fl.Field().Interface().(domain.PaymentType)

	switch paymentType {
	case "CASH", "E-WALLET", "EDC":
		return true
	default:
		return false
	}
}
