package apierrors

import "net/http"

// APIError represents a structured error response
type APIError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

// New creates a new APIError
func New(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// Predefined errors
var (
	ErrBadRequest          = New(http.StatusBadRequest, "bad request")
	ErrUnauthorized        = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = New(http.StatusForbidden, "forbidden")
	ErrNotFound            = New(http.StatusNotFound, "resource not found")
	ErrInternalServerError = New(http.StatusInternalServerError, "internal server error")
)

// WriteJSON writes the APIError as a JSON response
func (e *APIError) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	jsonStr := `{"message":"` + e.Message + `"}` // simple JSON without extra dependency
	w.Write([]byte(jsonStr))
}
