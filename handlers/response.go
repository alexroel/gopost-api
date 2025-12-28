package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type AppError struct {
	Message string
	Code    int
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(message string, code int) *AppError {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

func RespondError(w http.ResponseWriter, err *AppError) {
	RespondJSON(w, err.Code, ErrorResponse{
		Error:   http.StatusText(err.Code),
		Message: err.Message,
		Code:    err.Code,
	})
}

func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
