package errs

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func NewNotFoundError(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func NewUnexpectedError(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func NewValidationError(msg string) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: msg}
}
