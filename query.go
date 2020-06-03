package hive

import (
	"net/url"
)

// Copies the raw query string values from the proxy request
func copyQueryToSend(queryToSend []string, r *ProxyRequest) string {
	if len(queryToSend) > 0 {
		params := url.Values{}
		for _, key := range queryToSend {
			v := r.Query.Get(key)
			params.Add(key, v)
		}
		return params.Encode()
	}

	return ""
}
