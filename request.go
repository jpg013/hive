package hive

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Code-Pundits/go-config"
	"github.com/labstack/echo/v4"
)

const requestParamsAsterisk string = "*"

// ProxyRequest contains the data to send to the remote backend
type ProxyRequest struct {
	Method  string
	Body    io.ReadCloser
	Params  map[string]string
	Headers map[string][]string
	Query   url.Values
	URL     *url.URL
}

// NewRequest takes a context and creates a new ProxyRequest
func NewRequest(c echo.Context) *ProxyRequest {
	req := c.Request()

	return &ProxyRequest{
		Method:  req.Method,
		Body:    cloneRequestBody(req),
		Headers: cloneRequestHeaders(req),
	}
}

func cloneRequestHeaders(req *http.Request) map[string][]string {
	headers := make(map[string][]string, len(req.Header))
	for k, vs := range req.Header {
		tmp := make([]string, len(vs))
		copy(tmp, vs)
		headers[k] = tmp
	}
	return headers
}

// Check out https://github.com/golang/go/issues/36095 for how I came by this
func cloneRequestBody(req *http.Request) io.ReadCloser {
	var b bytes.Buffer
	b.ReadFrom(req.Body)
	req.Body = ioutil.NopCloser(&b)
	return ioutil.NopCloser(bytes.NewReader(b.Bytes()))
}

// RequestGeneratorFunc takes a proxy request and converts it into an http request
type RequestGeneratorFunc func(*ProxyRequest) *http.Request

// NewRequestGeneratorFactory returns a RequestGeneratorFunc
func NewRequestGeneratorFactory(cfg *config.BackendConfig) RequestGeneratorFunc {
	return func(p *ProxyRequest) *http.Request {
		reqURL := url.URL{
			Scheme:   cfg.Scheme,
			Host:     cfg.Host,
			Path:     cfg.Path,
			RawQuery: copyQueryToSend(cfg.Query, p),
		}

		// Make the http request
		request, err := http.NewRequest(strings.ToTitle(cfg.Method), reqURL.String(), p.Body)

		if err != nil {
			panic(err)
		}

		request.Header = copyHeadersToSend(cfg.HeadersToSend, p)

		return request
	}
}
