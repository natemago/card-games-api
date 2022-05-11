package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the general structure of the API errors.
type ErrorResponse struct {
	Message string `json:"message"`
}

// ErrorHandler builds new error handling middleware to be attached to Gin router.
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()
		if err == nil {
			return
		}

		statusCode := http.StatusInternalServerError

		if IsBadRequestError(err.Err) || IsValidationError(err.Err) {
			statusCode = http.StatusBadRequest
		} else if IsNotFoundError(err.Err) {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, &ErrorResponse{
			Message: err.Error(),
		})
	}
}
