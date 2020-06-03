package hive

var (
	mimeApplicationJSON  string   = "application/json"
	headerContentType    string   = "Content-Type"
	defaultHeadersToSend []string = []string{headerContentType}
)

// Copy over any specified headers from Request to Backend.
// If the backend request headers is not specified, copy default headers
func copyHeadersToSend(headersToSend []string, r *ProxyRequest) map[string][]string {
	if len(headersToSend) == 0 {
		headersToSend = defaultHeadersToSend
	}

	headers := make(map[string][]string, len(headersToSend))

	for _, k := range headersToSend {
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
