package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BadRequestError ...
const (
	BadRequestError   = "BadRequest"
	UnauthorizedError = "Unauthorized"
	ForbiddenError    = "Forbidden"
	NotFoundError     = "NotFound"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	System  string `json:"system"`
}

// ErrorResponse  ...
type ErrorResponse struct {
	Error Error `json:"error"`
}

// BadRequestErrorJSONResponse ...
var (
	UnauthorizedJSONResponse = NewErrorJSONResponse(UnauthorizedError, http.StatusUnauthorized)
)

func NewErrorJSONResponse(
	errorCode string,
	statusCode int,
) func(c *gin.Context, message string) {
	return func(c *gin.Context, message string) {
		c.JSON(statusCode, ErrorResponse{Error: Error{
			Code:    errorCode,
			Message: message,
			System:  "bk-apigateway-sdks",
		}})
	}
}
