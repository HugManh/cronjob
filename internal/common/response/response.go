package response

import (
	"time"

	"github.com/HugManh/cronjob/pkg/constant"
)

type SendResponse interface {
	SuccessMsgResponse(message string)
	SuccessDataResponse(message string, data any)
	BadRequestError(message string, err error)
	ForbiddenError(message string, err error)
	UnauthorizedError(message string, err error)
	NotFoundError(message string, err error)
	InternalServerError(message string, err error)
	MixedError(err error)
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Total      int `json:"total"`
	Limit      int `json:"limit"`
	TotalPages int `json:"totalPages"`
}

type Meta struct {
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

type ApiResponse[T any] struct {
	Success   bool              `json:"success"`
	Data      T                 `json:"data,omitempty"`
	Message   string            `json:"message,omitempty"`
	Code      string            `json:"code,omitempty"`
	Errors    map[string]string `json:"errors,omitempty"`
	Meta      map[string]any    `json:"meta,omitempty"`
	Timestamp string            `json:"timestamp"`
}

func Null() interface{} {
	return nil
}

// Success trả về response thành công, có data, meta (có thể nil)
func Success[T any](responseStatus constant.ResponseStatus, data T, meta map[string]any) ApiResponse[T] {
	return BuildResponse_(
		true,
		responseStatus.GetResponseMessage(),
		responseStatus.GetResponseStatus(),
		data,
		nil,
		meta,
	)
}

// Fail trả về response lỗi, có errors (có thể nil)
func Fail[T any](responseStatus constant.ResponseStatus, errors map[string]string) ApiResponse[T] {
	return BuildResponse_(
		false,
		responseStatus.GetResponseMessage(),
		responseStatus.GetResponseStatus(),
		*new(T), // giá trị zero của T (ví dụ nil cho pointer, 0 cho int,...)
		errors,
		nil,
	)
}

func BuildResponse_[T any](
	success bool,
	message string,
	code string,
	data T,
	errors map[string]string,
	meta map[string]any,
) ApiResponse[T] {
	return ApiResponse[T]{
		Success:   success,
		Data:      data,
		Message:   message,
		Code:      code,
		Errors:    errors,
		Meta:      meta,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
