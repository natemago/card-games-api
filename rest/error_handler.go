package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/natemago/card-games-api/errors"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()
		if err == nil {
			return
		}

		statusCode := http.StatusInternalServerError

		if errors.IsBadRequestError(err.Err) || errors.IsValidationError(err.Err) {
			statusCode = http.StatusBadRequest
		} else if errors.IsNotFoundError(err.Err) {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, &ErrorResponse{
			Message: err.Error(),
		})
	}
}
