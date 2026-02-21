package api

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func badRequestError(format string, args ...interface{}) *HTTPError {
	return &HTTPError{Code: http.StatusBadRequest, Message: fmt.Sprintf(format, args...)}
}

func unauthorizedError(format string, args ...interface{}) *HTTPError {
	return &HTTPError{Code: http.StatusUnauthorized, Message: fmt.Sprintf(format, args...)}
}

func forbiddenError(format string, args ...interface{}) *HTTPError {
	return &HTTPError{Code: http.StatusForbidden, Message: fmt.Sprintf(format, args...)}
}

func notFoundError(format string, args ...interface{}) *HTTPError {
	return &HTTPError{Code: http.StatusNotFound, Message: fmt.Sprintf(format, args...)}
}

func conflictError(format string, args ...interface{}) *HTTPError {
	return &HTTPError{Code: http.StatusConflict, Message: fmt.Sprintf(format, args...)}
}

func internalServerError(format string, args ...interface{}) *HTTPError {
	return &HTTPError{Code: http.StatusInternalServerError, Message: fmt.Sprintf(format, args...)}
}

func validationError(err error) *HTTPError {
	return &HTTPError{Code: http.StatusUnprocessableEntity, Message: err.Error()}
}
