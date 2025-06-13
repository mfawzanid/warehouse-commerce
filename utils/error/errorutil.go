package errorutil

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	Error   = "error"
	Message = "message"
)

type CustomError struct {
	ErrorType     error
	OriginalError error
}

func NewErrorCode(errorType, originalError error) *CustomError {
	return &CustomError{
		ErrorType:     errorType,
		OriginalError: originalError,
	}
}

func GetErrorType(err error) error {
	if e, ok := err.(*CustomError); ok {
		return e.ErrorType
	}
	return err
}

func GetOriginalError(err error) error {
	if e, ok := err.(*CustomError); ok {
		return e.OriginalError
	}
	return err
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorType.Error(), e.OriginalError.Error())
}

var (
	ErrBadRequest      = errors.New("bad request")
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrUniqueViolation = errors.New("unique violation")
)

func CombineHTTPErrorMessage(httpStatusCode int, err error) string {
	return fmt.Sprintf("%s: %s", http.StatusText(httpStatusCode), err.Error())
}
