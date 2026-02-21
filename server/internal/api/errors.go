package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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
	if ve, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, 0, len(ve))
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				msgs = append(msgs, field+" is required")
			case "email":
				msgs = append(msgs, field+" must be a valid email address")
			case "min":
				msgs = append(msgs, field+" must be at least "+fe.Param()+" characters")
			case "max":
				msgs = append(msgs, field+" must be at most "+fe.Param()+" characters")
			default:
				msgs = append(msgs, field+" is invalid")
			}
		}
		return &HTTPError{Code: http.StatusUnprocessableEntity, Message: strings.Join(msgs, "; ")}
	}
	return &HTTPError{Code: http.StatusUnprocessableEntity, Message: err.Error()}
}
