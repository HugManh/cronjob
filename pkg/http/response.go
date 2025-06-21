package http

import "github.com/HugManh/cronjob/pkg/constant"

func Null() interface{} {
	return nil
}

type ApiResponse[T any] struct {
	ResponseKey     string `json:"response_key"`
	ResponseMessage string `json:"response_message"`
	Data            T      `json:"data"`
}

func BuildResponse[T any](responseStatus constant.ResponseStatus, data T) ApiResponse[T] {
	return BuildResponse_(responseStatus.GetResponseStatus(), responseStatus.GetResponseMessage(), data)
}

func BuildResponse_[T any](status string, message string, data T) ApiResponse[T] {
	return ApiResponse[T]{
		ResponseKey:     status,
		ResponseMessage: message,
		Data:            data,
	}
}
