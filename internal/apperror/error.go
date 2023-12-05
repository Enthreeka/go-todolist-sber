package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Msg string `json:"message"`
	Err error  `json:"-"`
}

var (
	ErrUniqueViolation     = NewError("Violation must be unique", errors.New("non_unique_value"))
	ErrForeignKeyViolation = NewError("Foreign Key Violation", errors.New("foreign_key_violation "))
	ErrNoRows              = NewError("No rows in result set", errors.New("no_rows"))
)

func (a *AppError) Error() string {
	return fmt.Sprintf("%s", a.Msg)
}

func NewError(msg string, err error) *AppError {
	return &AppError{
		Msg: msg,
		Err: err,
	}
}

func ParseHTTPErrStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrUniqueViolation):
		return http.StatusBadRequest
	case errors.Is(err, ErrForeignKeyViolation):
		return http.StatusBadRequest
	case errors.Is(err, ErrNoRows):
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
