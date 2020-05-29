package hive

import "net/url"

type BackendConfig struct {
	// Scheme is the URL scheme (http, https)
	Scheme string `json:"scheme"`
	// HTTP method of the request to send to the backend
	Method string `json:"method"`
	// host for the API
	Host string `json:"host"`
	// URL
	URL *url.URL `json:"url"`
	// Path is the remote path
	Path string `json:"path"`
	// the name of the group for the response
	Group string `json:"group"`
	// Backend URL Path Params
	Params map[string]string
	// List of query string params to be sent to backend
	Query []string
	// HeadersToSend defines the list of headers to send to the backend
	HeadersToSend []string `json:"headers_to_send"`
}

type parseableBackendConfig struct {
	Scheme        string            `json:"scheme"`
	Method        string            `json:"method"`
	Host          string            `json:"host"`
	URL           string            `json:"url"`
	Path          string            `json:"path"`
	Group         string            `json:"group"`
	Params        map[string]string `json:"params"`
	Query         []string          `json:"query"`
	HeadersToSend []string          `json:"headers_to_send"`
}

func (cfg *parseableBackendConfig) normalize() (*BackendConfig, error) {
	u, err := url.Parse(cfg.URL)

	return &BackendConfig{
		Scheme:        cfg.Scheme,
		Method:        cfg.Method,
		Host:          cfg.Host,
		URL:           u,
		Path:          cfg.Path,
		Group:         cfg.Group,
		Params:        cfg.Params,
		Query:         cfg.Query,
		HeadersToSend: cfg.HeadersToSend,
	}, err
}
