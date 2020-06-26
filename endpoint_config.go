package hive

import (
	"time"

	"github.com/labstack/echo/v4"
)

// EndpointConfig defines the configuration of a single endpoint to be exposed
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
	// StreamResponse indicates whether response should be streamed
	StreamResponse bool `json:"stream_response"`
	// List of Backends with the endpoint
	Backends []*BackendConfig
	// the encoding format
	OutputEncoding string `json:"output_encoding"`
	// Middlewares called before Backend Proxies
	BeforeMiddlewares []echo.MiddlewareFunc
	// Middlewares called after Backend Proxies
	AfterMiddlewares []echo.MiddlewareFunc
}

// Init initializes endpoint config default values
func (cfg *EndpointConfig) Init() error {
	for _, be := range cfg.Backends {
		err := be.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

type parseableEndpointConfig struct {
	Endpoint        string `json:"endpoint"`
	Method          string `json:"method"`
	ConcurrentCalls int    `json:"concurrent_calls"`
	Timeout         string `json:"timeout"`
	CacheTTL        string `json:"cache_ttl"`
	Backends        []parseableBackendConfig
	StreamResponse  bool `json:"stream_response"`
}

var defaultTimeout = 60 * time.Second

func (cfg parseableEndpointConfig) normalize() (*EndpointConfig, error) {
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		timeout = defaultTimeout
	}
	cacheTTL, err := time.ParseDuration(cfg.CacheTTL)
	if err != nil {
		cacheTTL = time.Duration(0)
	}

	backends := make([]*BackendConfig, len(cfg.Backends))
	for i, p := range cfg.Backends {
		be, err := p.normalize()

		if err == nil {
			backends[i] = be
		}
	}

	return &EndpointConfig{
		Endpoint:        cfg.Endpoint,
		Method:          cfg.Method,
		ConcurrentCalls: cfg.ConcurrentCalls,
		Timeout:         timeout,
		CacheTTL:        cacheTTL,
		Backends:        backends,
		StreamResponse:  cfg.StreamResponse,
	}, nil
}
