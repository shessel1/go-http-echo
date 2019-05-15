package main

import (
	"io/ioutil"

	"server/http"

	"github.com/google/logger"
)

func main() {
	logger.Init("", true, true, ioutil.Discard)
	http.Server(":8080")
}
