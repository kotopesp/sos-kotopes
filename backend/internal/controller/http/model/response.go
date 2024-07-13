package model

type status string

const (
	success = status("success")
	err     = status("error")
)

type (
	Response struct {
		Status status `json:"status"`
		Data   any    `json:"data,omitempty"`
	}
)

func OKResponse(data any) Response {
	return Response{
		Status: success,
		Data:   data,
	}
}

func ErrorResponse(data any) Response {
	return Response{
		Status: err,
		Data:   data,
	}
}
