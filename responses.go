package repository

import "net/http"

type Err struct {
	Message string `json:"message,omitempty"`
	Code    int64  `json:"code,omitempty"`
}

type Response struct {
	Result string      `json:"result"`
	Error  *Err        `json:"error,omitempty"`
	Data   interface{} `json:"data"`
}

var (
	Ok                  = Response{Result: "ok"}
	ErrorUnauthorized   = Error(&Err{Code: http.StatusUnauthorized, Message: "Unauthorized"})
	ErrorInternalServer = Error(&Err{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
	ErrorForbidden      = Error(&Err{Code: http.StatusForbidden, Message: "Forbidden"})
	ErrorBadRequest     = Error(&Err{Code: http.StatusBadRequest, Message: "Bad Request"})
	ErrorNotFound       = Error(&Err{Code: http.StatusNotFound, Message: "Not Found"})
)

func Error(err *Err) Response {
	return Response{
		Result: "error",
		Error:  err,
	}
}

func (r Response) WithData(data interface{}) Response {
	r.Data = data
	return r
}
