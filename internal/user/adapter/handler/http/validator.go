package http

import (
	"github.com/go-playground/validator/v10"
	"go-restaurant/internal/user/domain"
)

// UserRoleValidator is a custom validator for validating user roles
var UserRoleValidator validator.Func = func(fl validator.FieldLevel) bool {
	userRole := fl.Field().Interface().(domain.UserRole)

	switch userRole {
	case "admin", "cashier":
		return true
	default:
		return false
	}
}
