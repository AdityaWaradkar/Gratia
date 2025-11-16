package apierrors

import (
	"net/http"
)

type APIError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func New(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrBadRequest          = New(http.StatusBadRequest, "bad request")
	ErrUnauthorized        = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = New(http.StatusForbidden, "forbidden")
	ErrNotFound            = New(http.StatusNotFound, "resource not found")
	ErrInternalServerError = New(http.StatusInternalServerError, "internal server error")
)

func (e *APIError) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	body := `{"message":"` + e.Message + `"}`
	w.Write([]byte(body))
}
