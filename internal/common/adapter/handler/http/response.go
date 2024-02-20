package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-restaurant/internal/common/domain"
	"net/http"
)

// Response represents a Response body format
type Response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// NewResponse is a helper function to create a response body
func NewResponse(success bool, message string, data any) Response {
	return Response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// Meta represents metadata for a paginated response
type Meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

// NewMeta is a helper function to create metadata for a paginated response
func NewMeta(total, limit, skip uint64) Meta {
	return Meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	domain.ErrInternal:                   http.StatusInternalServerError,
	domain.ErrDataNotFound:               http.StatusNotFound,
	domain.ErrConflictingData:            http.StatusConflict,
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrUnauthorized:               http.StatusUnauthorized,
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrInvalidToken:               http.StatusUnauthorized,
	domain.ErrExpiredToken:               http.StatusUnauthorized,
	domain.ErrForbidden:                  http.StatusForbidden,
	domain.ErrNoUpdatedData:              http.StatusBadRequest,
	domain.ErrInsufficientStock:          http.StatusBadRequest,
	domain.ErrInsufficientPayment:        http.StatusBadRequest,
}

// ValidationError sends an error response for some specific request validation error
func ValidationError(ctx *gin.Context, err error) {
	errMsg := ParseError(err)
	errRsp := NewErrorResponse(errMsg)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// HandleError determines the status code of an error and returns a JSON response with the error message and status code
func HandleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := ParseError(err)
	errRsp := NewErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// HandleAbort sends an error response and aborts the request with the specified status code and error message
func HandleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := ParseError(err)
	errRsp := NewErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// ParseError parses error messages from the error object and returns a slice of error messages
func ParseError(err error) []string {
	var errMsg []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, err.Error())
		}
	} else {
		errMsg = append(errMsg, err.Error())
	}

	return errMsg
}

// ErrorResponse represents an error response body format
type ErrorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// NewErrorResponse is a helper function to create an error response body
func NewErrorResponse(errMsg []string) ErrorResponse {
	return ErrorResponse{
		Success:  false,
		Messages: errMsg,
	}
}

// HandleSuccess sends a success response with the specified status code and optional data
func HandleSuccess(ctx *gin.Context, data any) {
	rsp := NewResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}
