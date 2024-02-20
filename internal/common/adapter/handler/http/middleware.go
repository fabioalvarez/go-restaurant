package http

import (
	"github.com/gin-gonic/gin"
	"go-restaurant/internal/auth/port"
	"go-restaurant/internal/auth/util"
	cmdomain "go-restaurant/internal/common/domain"
	"go-restaurant/internal/user/domain"
	"strings"
)

const (
	// AuthorizationHeaderKey is the key for authorization header in the request
	AuthorizationHeaderKey = "authorization"
	// AuthorizationType is the accepted authorization type
	AuthorizationType = "bearer"
	// AuthorizationPayloadKey is the key for authorization payload in the context
	AuthorizationPayloadKey = "authorization_payload"
)

// authMiddleware is a middleware to check if the user is authenticated
func authMiddleware(token port.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		isEmpty := len(authorizationHeader) == 0
		if isEmpty {
			err := cmdomain.ErrEmptyAuthorizationHeader
			HandleAbort(ctx, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		isValid := len(fields) == 2
		if !isValid {
			err := cmdomain.ErrInvalidAuthorizationHeader
			HandleAbort(ctx, err)
			return
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != AuthorizationType {
			err := cmdomain.ErrInvalidAuthorizationType
			HandleAbort(ctx, err)
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			HandleAbort(ctx, err)
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

// adminMiddleware is a middleware to check if the user is an admin
func adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := util.GetAuthPayload(ctx, AuthorizationPayloadKey)

		isAdmin := payload.Role == domain.Admin
		if !isAdmin {
			err := cmdomain.ErrForbidden
			HandleAbort(ctx, err)
			return
		}

		ctx.Next()
	}
}
