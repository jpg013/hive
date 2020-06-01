package hive

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

const requestParamsAsterisk string = "*"

// RequestGeneratorFunc takes a context and creates a new ProxyRequest
type RequestGeneratorFunc func(c echo.Context) *ProxyRequest

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

	// httpReq.Header
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

// func (rf *requestFactory) parseHeaders(c echo.Context) map[string][]string {
// 	req := c.Request()
// 	headers := make(map[string][]string, len(rf.headersToSend)+1)
// 	headers["X-Forwarded-For"] = []string{c.RealIP()}

// 	for _, k := range rf.headersToSend {
// 		if k == requestParamsAsterisk {
// 			headers = req.Header
// 			break
// 		}

// 		v := req.Header.Get(k)

// 		if k != "" {
// 			headers[k] = []string{v}
// 		}
// 	}

// 	return headers
// }

// type requestFactory struct {
// 	headersToSend []string
// 	queryParams   []string
// 	params        []string
// }

// func (rf *requestFactory) parseQueryParams(c echo.Context) map[string][]string {
// 	req := c.Request()
// 	// Copy the query string params
// 	query := make(map[string][]string, len(rf.queryParams))

// 	for _, k := range rf.queryParams {
// 		if k == requestParamsAsterisk {
// 			query = req.URL.Query()
// 			break
// 		}
// 		query[k] = []string{c.QueryParam(k)}
// 	}

// 	return query
// }

// func (rf *requestFactory) parseParams(c echo.Context) map[string]string {
// 	params := make(map[string]string, len(rf.params))

// 	for _, k := range rf.params {
// 		if k == requestParamsAsterisk {
// 			for _, n := range c.ParamNames() {
// 				params[n] = c.Param(n)
// 			}

// 			break
// 		}

// 		params[k] = c.Param(k)
// 	}

// 	return params
// }

// func (rf *requestFactory) parseHeaders(c echo.Context) map[string][]string {
// 	req := c.Request()
// 	headers := make(map[string][]string, len(rf.headersToSend)+1)
// 	headers["X-Forwarded-For"] = []string{c.RealIP()}

// 	for _, k := range rf.headersToSend {
// 		if k == requestParamsAsterisk {
// 			headers = req.Header
// 			break
// 		}

// 		v := req.Header.Get(k)

// 		if k != "" {
// 			headers[k] = []string{v}
// 		}
// 	}

// 	return headers
// }

// func (r *requestFactory) New(c echo.Context) *Request {
// 	req := c.Request()

// 	return &Request{
// 		Method:  req.Method,
// 		Query:   r.parseQueryParams(c),
// 		Headers: r.parseHeaders(c),
// 		Params:  r.parseParams(c),
// 		Body:    req.Body,
// 	}
// }

// func NewRequestFactory(e *Endpoint) RequestGeneratorFunc {
// 	headersToSend := e.HeadersToSend
// 	if len(headersToSend) == 0 {
// 		headersToSend = defaultHeadersToSend
// 	}
// 	queryParams := e.QueryString

// 	rf := &requestFactory{headersToSend, queryParams, e.Params}

// 	return func(c echo.Context) *Request {
// 		req := c.Request()

// 		return &Request{
// 			Method:  req.Method,
// 			Query:   rf.parseQueryParams(c),
// 			Headers: rf.parseHeaders(c),
// 			Params:  rf.parseParams(c),
// 			Body:    req.Body,
// 		}
// 	}
// }
