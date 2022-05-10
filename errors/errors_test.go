package errors

import (
	"fmt"
	"testing"
)

func TestNotFoundError(t *testing.T) {
	err := NotFoundError("record not found", nil)
	if !IsNotFoundError(err) {
		t.Error("Expected to be a 'not-found-error'.")
	}
	if err.Error() != "record not found" {
		t.Error("Expected to get correct error message from APIError.")
	}

	err = fmt.Errorf("generic-error")
	if IsNotFoundError(err) {
		t.Error("Generic error should not be a 'not-found-error'.")
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError("invalid value", nil)
	if !IsValidationError(err) {
		t.Error("Expected to be a 'validation-error'.")
	}
	if err.Error() != "invalid value" {
		t.Error("Expected to get correct error message from APIError.")
	}

	err = fmt.Errorf("generic-error")
	if IsValidationError(err) {
		t.Error("Generic error should not be a 'validation-error'.")
	}
}

func TestBadRequestError(t *testing.T) {
	err := BadRequestError("bad parameter", nil)
	if !IsBadRequestError(err) {
		t.Error("Expected to be a 'bad-request-error'.")
	}
	if err.Error() != "bad parameter" {
		t.Error("Expected to get correct error message from APIError.")
	}

	err = fmt.Errorf("generic-error")
	if IsBadRequestError(err) {
		t.Error("Generic error should not be a 'bad-request-error'.")
	}
}
