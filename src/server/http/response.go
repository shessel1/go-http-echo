package http

import (
	"fmt"

	"github.com/google/logger"
)

type ResponseStatus struct {
	Version string
	Status  *HttpStatus
}

func (r *ResponseStatus) String() string {
	ver := r.Version
	if len(r.Version) == 0 {
		ver = "HTTP/1.0"
	}
	return fmt.Sprintf("%s %d %s", ver, r.Status.code, r.Status.message)
}

type Response struct {
	Status  *ResponseStatus
	Req     *Request
	Headers map[string]string
}

func MakeResponse(r *Request) *Response {
	return &Response{
		Req:     r,
		Headers: make(map[string]string)}
}

func (r *Response) WriteHeader() {
	logger.Infof("%v >>> %s", r.Status.String(), (*r.Req.Client.Conn).RemoteAddr())
	w := r.Req.Client.Writer
	w.WriteString(fmt.Sprintf("%s\r\n", r.Status.String()))
	for k, v := range r.Headers {
		w.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	w.WriteString("\r\n")
	w.Flush()
}
