package http

import (
	"bufio"
	"fmt"
	"net"
	"server"

	"github.com/google/logger"
)

func HandleRequest(c *server.Client) {
	req, err := MakeRequest(c)
	resp := MakeResponse(req)
	if err != nil {
		resp.Status = &ResponseStatus{Status: err.status}
	}
	resp.Status = &ResponseStatus{Status: StatusOK}
	resp.WriteHeader()
}

func Connect(c *net.Conn) {
	client := &server.Client{c, bufio.NewReader(*c), bufio.NewWriter(*c)}
	defer func() {
		client.Writer.Flush()
		(*client.Conn).Close()
	}()
	HandleRequest(client)
}

func Server(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	logger.Infof("Server listening on %s\n", l.Addr())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go Connect(&c)
	}
}
