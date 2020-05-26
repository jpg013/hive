package proxy

import (
	"io"
	"net/url"
)

// Request contains the data to send to the backend
type Request struct {
	Method  string
	URL     *url.URL
	Query   url.Values
	Path    string
	Body    io.ReadCloser
	Params  map[string]string
	Headers map[string][]string
}
