package http

import (
	"bufio"
	"fmt"
	"io"
	"server"
	"strings"

	"github.com/google/logger"
)

type RequestLine struct {
	Method  string
	URI     string
	Version string
}

type Request struct {
	Line    *RequestLine
	Client  *server.Client
	Body    *bufio.Reader
	Headers map[string]string
	Params  map[string]string

	KeepAlive bool
}

func (r *Request) HasHeader(h string) bool {
	_, ok := r.Headers[h]
	return ok
}

func (rl *RequestLine) String() string {
	return fmt.Sprintf("%s %s %s", rl.Method, rl.URI, rl.Version)
}

func ReadLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Replace(line, "\r\n", "", -1), nil
}

func (r *Request) ParseRequestLine(line string) bool {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return false
	}
	r.Line = &RequestLine{Method: parts[0], URI: parts[1], Version: parts[2]}

	r.Params = make(map[string]string)
	if strings.Index(r.Line.URI, "?") != -1 {
		params := strings.Split(r.Line.URI, "?")
		if len(parts) > 2 {
			return false
		}
		for _, param := range strings.Split(params[1], "&") {
			kv := strings.Split(param, "=")
			if len(kv) != 2 {
				return false
			}
			r.Params[kv[0]] = kv[1]
		}
	}
	return true
}

func (r *Request) ParseHeaders(s []string) bool {
	r.Headers = make(map[string]string)
	for _, header := range s {
		kv := strings.SplitN(header, ":", 2)
		if len(kv) != 2 {
			return false
		}
		k := strings.Trim(kv[0], " ")
		v := strings.Trim(kv[1], " ")
		r.Headers[strings.ToLower(k)] = v
	}
	r.KeepAlive = r.Line.Version == "HTTP/1.1" || r.HasHeader("keep-alive")
	return true
}

func MakeRequest(c *server.Client) (*Request, *HttpError) {
	req := &Request{Client: c}

	header := make([]string, 0, 0)
	line, err := ReadLine(c.Reader)
	for {
		if err != nil || len(line) == 0 {
			break
		}
		header = append(header, line)
		line, err = ReadLine(c.Reader)
	}
	if err != nil {
		if err == io.EOF {
			logger.Infof("%s: Unexpected EOF <<< %s", LogInvalidRequest, (*c.Conn).RemoteAddr())
			return req, &HttpError{StatusBadRequest}
		}
		logger.Infof("%s: Read failed <<< %s", LogInvalidRequest, (*c.Conn).RemoteAddr())
		return req, &HttpError{StatusServerError}
	}

	ok := req.ParseRequestLine(header[0])
	ok = req.ParseHeaders(header[1:]) && ok
	if !ok {
		logger.Infof("%s: Bad Request <<< %s", LogInvalidRequest, (*c.Conn).RemoteAddr())
		return req, &HttpError{StatusBadRequest}
	}

	logger.Infof("%v <<< %s ", req.Line.String(), (*c.Conn).RemoteAddr())
	return req, nil
}
