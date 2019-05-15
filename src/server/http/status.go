package http

import "fmt"

type HttpStatus struct {
	code    int
	message string
}

var (
	StatusOK                  = &HttpStatus{200, "OK"}
	StatusBadRequest          = &HttpStatus{400, "Bad Request"}
	StatusNotFound            = &HttpStatus{404, "Not Found"}
	StatusMethodNotAllowed    = &HttpStatus{405, "Method Not Allowed"}
	StatusServerError         = &HttpStatus{500, "Internal Server Error"}
	StatusNotImplemented      = &HttpStatus{501, "Not Implemented"}
	StatusVersionNotSupported = &HttpStatus{505, "HTTP Version Not Supported"}
)

var (
	LogInvalidRequest = "Request Failed"
)

type HttpError struct {
	status *HttpStatus
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%v", e.status)
}
