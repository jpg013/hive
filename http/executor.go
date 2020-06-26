package http

import (
	"net/http"
)

var defaultClient = &http.Client{}

// Executor exposes a Do method for http request
type Executor interface {
	Do(*http.Request) (*http.Response, error)
}

type httpExecutor struct {
	client *http.Client
}

func (exec *httpExecutor) Do(req *http.Request) (*http.Response, error) {
	resp, err := exec.client.Do(req)

	return resp, err
}

func NewExecutor() Executor {
	return &httpExecutor{
		client: defaultClient,
	}
}
