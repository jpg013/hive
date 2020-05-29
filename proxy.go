package hive

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Proxy processes a request in a given context and returns a response and an error
type Proxy func(request *Request) (*Response, error)

func buildBackendRequest(be *BackendConfig, r *Request) *http.Request {
	reqURL := url.URL{
		Scheme: be.Scheme,
		Host:   be.Host,
		Path:   be.Path,
	}

	// Prepare query params
	if len(be.Query) > 0 {
		params := url.Values{}
		for _, key := range be.Query {
			v := r.Query.Get(key)
			params.Add(key, v)
		}
		reqURL.RawQuery = params.Encode()
	}

	// Make the http request
	backendRequest, err := http.NewRequest(strings.ToTitle(be.Method), reqURL.String(), r.Body)

	if err != nil {
		panic(err)
	}

	// Copy over any specified headers from Request to Backend.
	// If the backend request headers is not specified, copy default headers
	if len(be.HeadersToSend) == 0 {
		be.HeadersToSend = defaultHeadersToSend
	}

	for _, k := range be.HeadersToSend {
		if k == requestParamsAsterisk {
			for name, vs := range r.Headers {
				tmp := make([]string, len(vs))
				copy(tmp, vs)
				backendRequest.Header[name] = tmp
			}
			break
		}

		vs, ok := r.Headers[k]
		if ok {
			tmp := make([]string, len(vs))
			copy(tmp, vs)
			backendRequest.Header[k] = tmp
		}
	}

	return backendRequest
}

func ProxyFactory(e *EndpointConfig) Proxy {
	// Just assume 1 be for now
	return func(r *Request) (*Response, error) {
		be := e.Backends[0]
		request := buildBackendRequest(be, r)
		resp, err := httpClient.Do(request)

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		err = HTTPStatusHandler(resp)

		// We received an error status code from the backend
		if err != nil {
			return &Response{
				Status: resp.StatusCode,
				Errors: []string{err.Error()},
			}, nil
		}

		var data map[string]interface{}
		jsonBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonBytes, &data)
		if err != nil {
			return nil, err
		}

		return &Response{
			Data:   data,
			Group:  be.Group,
			Status: resp.StatusCode,
			Errors: make([]string, 0),
		}, nil
	}
}
