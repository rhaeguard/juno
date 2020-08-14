package commons

import "net/http"

type SuccessResponse struct {
	Message string      `json:"message"`
	Code    uint16      `json:"statusCode"`
	Data    interface{} `json:"data"`
}

func MakeSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
		Code:    http.StatusOK,
	}
}

type DefaultResponse struct {
	Message string `json:"message"`
	Code    uint16 `json:"statusCode"`
}

func MakeFailureResponse(message string, code uint16) DefaultResponse {
	return DefaultResponse{
		Message: message,
		Code:    code,
	}
}
