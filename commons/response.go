package commons

import "net/http"

// SuccessResponse represents a successful response
type SuccessResponse struct {
	Message string      `json:"message"`
	Code    uint16      `json:"statusCode"`
	Data    interface{} `json:"data"`
}

// MakeSuccessResponse is a utility method to create payload easier
func MakeSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
		Code:    http.StatusOK,
	}
}

// DefaultResponse represents a default response
type DefaultResponse struct {
	Message string `json:"message"`
	Code    uint16 `json:"statusCode"`
}

// MakeFailureResponse is a utility method to create payload easier
func MakeFailureResponse(message string, code uint16) DefaultResponse {
	return DefaultResponse{
		Message: message,
		Code:    code,
	}
}
