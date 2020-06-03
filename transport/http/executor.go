package http

import (
	"net/http"
)

var DefaultOKStatus int = http.StatusOK
var DefaultErrorStatus int = http.StatusInternalServerError

var defaultClient = &http.Client{}

// Executor exposes a Do method for http request
type Executor interface {
	Do(string, *http.Request) *Response
}

type httpExecutor struct {
	client *http.Client
}

func (exec *httpExecutor) Do(group string, req *http.Request) *Response {
	handler := NewResponseHandler(group)

	resp, err := exec.client.Do(req)

	return handler(resp, err)
}

func NewExecutor() Executor {
	return &httpExecutor{
		client: defaultClient,
	}
}
