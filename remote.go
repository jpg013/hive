package hive

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/Code-Pundits/go-config"
)

func buildRawQuery(remote config.BackendConfig, r *Request) string {
	if len(remote.Query) > 0 {
		params := url.Values{}
		for _, key := range remote.Query {
			v := r.Query.Get(key)
			params.Add(key, v)
		}
		return params.Encode()
	}

	return ""
}

// Copy over any specified headers from Request to Backend.
// If the backend request headers is not specified, copy default headers
func copyRemoteHeaders(remote config.BackendConfig, r *Request) map[string][]string {
	if len(remote.HeadersToSend) == 0 {
		remote.HeadersToSend = defaultHeadersToSend
	}

	headers := make(map[string][]string, len(remote.HeadersToSend))

	for _, k := range remote.HeadersToSend {
		if k == requestParamsAsterisk {
			for name, vs := range r.Headers {
				tmp := make([]string, len(vs))
				copy(tmp, vs)
				headers[name] = tmp
			}
			break
		}

		vs, ok := r.Headers[k]
		if ok {
			tmp := make([]string, len(vs))
			copy(tmp, vs)
			headers[k] = tmp
		}
	}

	return headers
}

func buildRemoteRequest(remote config.BackendConfig, r *Request) *http.Request {
	reqURL := url.URL{
		Scheme:   remote.Scheme,
		Host:     remote.Host,
		Path:     remote.Path,
		RawQuery: buildRawQuery(remote, r),
	}

	// Make the http request
	backendRequest, err := http.NewRequest(strings.ToTitle(remote.Method), reqURL.String(), r.Body)

	if err != nil {
		panic(err)
	}

	backendRequest.Header = copyRemoteHeaders(remote, r)

	return backendRequest
}
