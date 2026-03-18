package response

import (
	"encoding/json"
	"net/http"
)

type successResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type validationErrorResponse struct {
	Code    int               `json:"code"`
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Details map[string]string `json:"details"`
}

func Success(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(successResponse{
		Code:    statusCode,
		Success: true,
		Message: message,
		Data:    data,
	})
}
func Error(w http.ResponseWriter, statusCode int, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{
		Code:    statusCode,
		Success: false,
		Error:   error,
	})
}
func ValidationError(w http.ResponseWriter, statusCode int, details map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(validationErrorResponse{
		Code:    statusCode,
		Success: false,
		Error:   "Validation error",
		Details: details,
	})
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

func InternalServerError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal server error")
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}
