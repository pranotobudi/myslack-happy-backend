package common

import "fmt"

// response
// meta payload

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func (r Response) String() string {
	return fmt.Sprintf("meta: %v\n data:%v\n", r.Meta, r.Data)
}

type Meta struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}

func (m Meta) String() string {
	return fmt.Sprintf("code: %v\n Status:%v\n Message:%v", m.Code, m.Status, m.Message)
}

type M map[string]interface{}

func ResponseFormatter(code int, status string, message interface{}, data interface{}) Response {
	meta := Meta{
		Code:    code,
		Status:  status,
		Message: message,
	}
	response := Response{
		Meta: meta,
		Data: data,
	}
	return response
}
func ResponseErrorFormatter(code int, err error) Response {
	response := ResponseFormatter(code, "error", "invalid request", err.Error())
	return response
}
