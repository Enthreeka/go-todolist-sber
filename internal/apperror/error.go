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
	ErrNotFound            = NewError("Tasks not found", errors.New("not_found"))

	ErrHashPasswordsNotEqual = NewError("Invalid password", errors.New("hashes_not_equal"))
	ErrDataNotValid          = NewError("Provided data is not valid", errors.New("not_valid"))
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
	case errors.Is(err, ErrHashPasswordsNotEqual):
		return http.StatusForbidden
	case errors.Is(err, ErrDataNotValid):
		return http.StatusUnprocessableEntity
	}

	return http.StatusInternalServerError
}
