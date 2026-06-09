package response

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	DB_ERR            ErrorCode = "DB_ERROR"
	INVALID_INPUT_ERR ErrorCode = "INVALID_INPUT_ERROR"
	SERVICE_ERR       ErrorCode = "SERVICE_ERROR"
	NOT_FOUND_ERR     ErrorCode = "NOT_FOUND_ERROR"
	FORBIDDEN_ERR     ErrorCode = "FORBIDDEN_ERROR"
)

type AppError struct {
	Status ErrorCode `json:"status"`
	Reason string    `json:"error"`
}

func (ap *AppError) Error() string {
	return fmt.Sprintf("Error: %s", ap.Reason)
}

type Response[T any] struct {
	Data T
	Err  *AppError
}

func (r Response[T]) IsError() bool {
	return r.Err != nil
}

func (r Response[T]) GetHttpCode() int {
	if !r.IsError() {
		return http.StatusOK
	}

	switch r.Err.Status {
	case DB_ERR:
		return http.StatusServiceUnavailable
	case INVALID_INPUT_ERR:
		return http.StatusBadRequest
	case NOT_FOUND_ERR:
		return http.StatusNotFound
	case FORBIDDEN_ERR:
		return http.StatusForbidden
	case SERVICE_ERR:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func Success[T any](data T) Response[T] {
	return Response[T]{
		Data: data,
	}
}

func Error[T any](status ErrorCode, format string, params ...any) Response[T] {
	return Response[T]{
		Err: &AppError{
			Status: status,
			Reason: fmt.Sprintf(format, params...),
		},
	}
}

func PropagateError[V, U any](u Response[U]) Response[V] {
	if !u.IsError() {
		// This function should only be called on error responses
		panic("PropagateError called on non-error response")
	}

	return Response[V]{
		Err: &AppError{
			Status: u.Err.Status,
			Reason: u.Err.Reason,
		},
	}
}
