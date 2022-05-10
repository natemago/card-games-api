package errors

import (
	goerrs "errors"
)

type APIError struct {
	Type            string
	Message         string
	UnderlyingError error
}

func (a *APIError) Error() string {
	return a.Message
}

func (a *APIError) Unwrap() error {
	return a.UnderlyingError
}

type ErrorBuilder func(message string, underlyingError error) error
type ErrorTypeCheck func(err error) bool

func ErrorType(errType string) (ErrorBuilder, ErrorTypeCheck) {
	return func(message string, underlyingError error) error {
			return &APIError{
				Type:            errType,
				Message:         message,
				UnderlyingError: underlyingError,
			}
		}, func(err error) bool {
			var errPtr *APIError = &APIError{}
			if goerrs.As(err, &errPtr) {
				if err.(*APIError).Type == errType {
					return true
				}
			}
			return false
		}
}

var NotFoundError, IsNotFoundError = ErrorType("not-found")
var ValidationError, IsValidationError = ErrorType("validation")
var BadRequestError, IsBadRequestError = ErrorType("bad-request")
