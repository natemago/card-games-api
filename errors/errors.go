package errors

import (
	goerrs "errors"
)

// APIError is the error structure used in the API and repositories.
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

// ErrorBuilder type defines a function that produces an error with the specified message and underlying error (if any).
type ErrorBuilder func(message string, underlyingError error) error

// ErrorTypeCheck type defines a function that checks if the supplied error is of a particular type. Used together with
// the ErrorBuilder to define a way to build and check errors of pre-defined types.
type ErrorTypeCheck func(err error) bool

// ErrorType defines an error type. It returns an ErrorBuilder function to produce errors of this type; and ErrorTypeCheck
// function that checks if a supplied error is of this type.
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

// Some base error types definitions.
var NotFoundError, IsNotFoundError = ErrorType("not-found")
var ValidationError, IsValidationError = ErrorType("validation")
var BadRequestError, IsBadRequestError = ErrorType("bad-request")
