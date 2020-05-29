package hive

import (
	"time"

	"github.com/labstack/echo/v4"
)

// Endpoint defines the configuration of a single endpoint to be exposed
type EndpointConfig struct {
	// url pattern to be registered and exposed to the world
	Endpoint string `json:"endpoint"`
	// HTTP method of the endpoint (GET, POST, PUT, etc)
	Method string `json:"method"`
	// number of concurrent calls this endpoint must send to the backends
	ConcurrentCalls int `json:"concurrent_calls"`
	// timeout of this endpoint
	Timeout time.Duration `json:"timeout"`
	// duration of the cache header
	CacheTTL time.Duration `json:"cache_ttl"`
	// List of path params to extract from the URI
	// Params []string `json:"params"`
	// // list of query string params to be extracted from the URI
	// QueryString []string `json:"querystring_params"`
	// HeadersToSend defines the list of headers to send to the backend
	// HeadersToSend []string `json:"headers_to_send"`
	// List of middlewares applied to this endpoint
	Middlewares []echo.MiddlewareFunc
	// List of Backends with the endpoint
	Backends []*BackendConfig
}

type parseableEndpointConfig struct {
	Endpoint        string `json:"endpoint"`
	Method          string `json:"method"`
	ConcurrentCalls int    `json:"concurrent_calls"`
	Timeout         string `json:"timeout"`
	CacheTTL        string `json:"cache_ttl"`
	Backends        []*parseableBackendConfig
}

func (cfg *parseableEndpointConfig) normalize() (*EndpointConfig, error) {
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return nil, err
	}
	cacheTTL, err := time.ParseDuration(cfg.CacheTTL)
	if err != nil {
		return nil, err
	}

	backends := make([]*BackendConfig, len(cfg.Backends))

	for i, p := range cfg.Backends {
		be, err := p.normalize()

		if err != nil {
			return nil, err
		}
		backends[i] = be
	}

	return &EndpointConfig{
		Endpoint:        cfg.Endpoint,
		Method:          cfg.Method,
		ConcurrentCalls: cfg.ConcurrentCalls,
		Timeout:         timeout,
		CacheTTL:        cacheTTL,
		Backends:        backends,
	}, nil
}

func EndpointHandler(cfg *EndpointConfig) echo.HandlerFunc {
	// cacheControlHeaderValue := fmt.Sprintf("public, max-age=%d", int(cfg.CacheTTL.Seconds()))
	// isCacheEnabled := cfg.CacheTTL.Seconds() != 0
	// requestGenerator := NewRequestFactory(cfg)
	proxy := ProxyFactory(cfg)

	return func(c echo.Context) error {
		request := NewRequest(c)
		resp, err := proxy(request)

		if err != nil {
			return err
		}

		return c.JSON(resp.Status, resp)
	}
}
