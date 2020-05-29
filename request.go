package hive

import (
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

const requestParamsAsterisk string = "*"

type RequestGeneratorFunc func(c echo.Context) *Request

// Request contains the data to send to the backend
type Request struct {
	Method  string
	Body    io.ReadCloser
	Params  map[string]string
	Headers map[string][]string
	Query   url.Values
	URL     *url.URL
}

func NewRequest(c echo.Context) *Request {
	req := c.Request()

	// httpReq.Header
	return &Request{
		Method:  req.Method,
		Body:    req.Body,
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
