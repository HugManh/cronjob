package httpx

import "time"

type ResponseCode string

const (
	CodeSuccess        ResponseCode = "SUCCESS"
	CodeDataNotFound   ResponseCode = "DATA_NOT_FOUND"
	CodeUnknownError   ResponseCode = "UNKNOWN_ERROR"
	CodeInvalidRequest ResponseCode = "INVALID_REQUEST"
	CodeUnauthorized   ResponseCode = "UNAUTHORIZED"
)

type APIResponse[T any] struct {
	Success   bool              `json:"success"`
	Data      T                 `json:"data,omitempty"`
	Message   string            `json:"message,omitempty"`
	Code      ResponseCode      `json:"code,omitempty"`
	Errors    map[string]string `json:"errors,omitempty"`
	Meta      map[string]any    `json:"meta,omitempty"`
	Timestamp string            `json:"timestamp"`
}

func NewSuccessResponse[T any](code ResponseCode, message string, data T, meta map[string]any) APIResponse[T] {
	return APIResponse[T]{
		Success:   true,
		Data:      data,
		Message:   message,
		Code:      code,
		Meta:      meta,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func NewErrorResponse(code ResponseCode, message string, errors map[string]string) APIResponse[any] {
	return APIResponse[any]{
		Success:   false,
		Message:   message,
		Code:      code,
		Errors:    errors,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
