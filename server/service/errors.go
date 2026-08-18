package service

import "net/http"

type ErrorKind int

const (
	ErrInternal ErrorKind = iota
	ErrBadRequest
	ErrUnauthorized
	ErrForbidden
	ErrNotFound
	ErrConflict
)

// AppError is the unified application error used across the service layer.
// Status carries an optional extra field for responses such as
// {"error": "...", "status": "pending"}.
type AppError struct {
	Kind   ErrorKind
	Msg    string
	Status string
}

func (e *AppError) Error() string { return e.Msg }

func (e *AppError) HTTPStatus() int {
	switch e.Kind {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func NewError(kind ErrorKind, msg string) *AppError {
	return &AppError{Kind: kind, Msg: msg}
}
