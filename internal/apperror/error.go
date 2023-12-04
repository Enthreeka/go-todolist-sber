package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Msg string `json:"message"`
	Err error  `json:"-"`
}

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

	}

	return http.StatusInternalServerError
}
