package util

import (
	"github.com/gin-gonic/gin"
	"go-restaurant/internal/auth/domain"
)

// GetAuthPayload is a helper function to get the auth payload from the context
func GetAuthPayload(ctx *gin.Context, key string) *domain.TokenPayload {
	return ctx.MustGet(key).(*domain.TokenPayload)
}
