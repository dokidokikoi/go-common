package errors

import "net/http"

type APIError struct {
	Code       int    `json:"code"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func Success(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusOK,
		Message:    message,
	}
}

func ServerFailed(message string, code int) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Code:       code,
	}
}

func ClientFailed(message string, code int) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Code:       code,
	}
}

func ClientAuthnFailed(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}

}
