package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

var (
	// DefaultPort is the default service port param
	DefaultPort = 8080
	// DefaultMaxIdleConns is the default value for the MaxIdleConns param
	DefaultMaxIdleConns = 250
	// DefaultMaxIdleConnsPerHost is the default value for the MaxIdleConnsPerHost param
	DefaultMaxIdleConnsPerHost = 250
	// DefaultTimeout is the default value to use for the Timeout param
	DefaultTimeout = 2 * time.Second
)

// ServiceConfig defines an api service config
type ServiceConfig struct {
	// name of the service
	Name string `json:"name"`
	// set of endpoint definitions
	// Endpoints []*EndpointConfig `json:"endpoints"`
	// defafult timeout
	Timeout time.Duration `json:"timeout"`
	// default TTL for GET
	CacheTTL time.Duration `json:"cache_ttl"`
	// default host
	Host string `json:"host"`
	// port to bind the service
	Port int `json:"port"`
	// OutputEncoding defines the default encoding strategy to use for the endpoint responses
	OutputEncoding string `json:"output_encoding"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration `json:"read_timeout"`
	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration `json:"write_timeout"`
	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.
	IdleTimeout time.Duration `json:"idle_timeout"`
	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body.
	ReadHeaderTimeout time.Duration `json:"read_header_timeout"`
	// DisableKeepAlives, if true, prevents re-use of TCP connections
	// between different HTTP requests.
	DisableKeepAlives bool `json:"disable_keep_alives"`
	// DisableCompression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	// request header when the Request contains no existing
	// Accept-Encoding value. If the Transport requests gzip on
	// its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body. However, if the user
	// explicitly requested gzip it is not automatically
	// uncompressed.
	DisableCompression bool `json:"disable_compression"`
	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int `json:"max_idle_connections"`
	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost int `json:"max_idle_connections_per_host"`
	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleConnTimeout time.Duration `json:"idle_connection_timeout"`
	// ResponseHeaderTimeout, if non-zero, specifies the amount of
	// time to wait for a server's response headers after fully
	// writing the request (including its body, if any). This
	// time does not include the time to read the response body.
	ResponseHeaderTimeout time.Duration `json:"response_header_timeout"`
	// ExpectContinueTimeout, if non-zero, specifies the amount of
	// time to wait for a server's first response headers after fully
	// writing the request headers if the request has an
	// "Expect: 100-continue" header. Zero means no timeout and
	// causes the body to be sent immediately, without
	// waiting for the server to approve.
	// This time does not include the time to send the request header.
	ExpectContinueTimeout time.Duration `json:"expect_continue_timeout"`
	// DialerTimeout is the maximum amount of time a dial will wait for
	// a connect to complete. If Deadline is also set, it may fail
	// earlier.
	//
	// The default is no timeout.
	//
	// When using TCP and dialing a host name with multiple IP
	// addresses, the timeout may be divided between them.
	//
	// With or without a timeout, the operating system may impose
	// its own earlier timeout. For instance, TCP timeouts are
	// often around 3 minutes.
	DialerTimeout time.Duration `json:"dialer_timeout"`
	// DialerFallbackDelay specifies the length of time to wait before
	// spawning a fallback connection, when DualStack is enabled.
	// If zero, a default delay of 300ms is used.
	DialerFallbackDelay time.Duration `json:"dialer_fallback_delay"`
	// DialerKeepAlive specifies the keep-alive period for an active
	// network connection.
	// If zero, keep-alives are not enabled. Network protocols
	// that do not support keep-alives ignore this field.
	DialerKeepAlive time.Duration `json:"dialer_keep_alive"`

	// run service in debug mode
	Debug bool
}

// EndpointConfig defines the configuration of a single endpoint to be exposed
type EndpointConfig struct {
	// url pattern to be registered and exposed to the world
	Endpoint string `json:"endpoint"`
	// HTTP method of the endpoint (GET, POST, PUT, PATH, DELETE)
	Method string `json:"method"`
	// timeout of this endpoint
	Timeout time.Duration `json:"timeout"`
	// duration of the cache header
	CacheTTL time.Duration `json:"cache_ttl"`
	// list of query string params to be extracted from the URI
	QueryString []string `json:"querystring_params"`
	// HeadersToPass defines the list of headers to pass to the backends
	HeadersToPass []string `json:"headers_to_pass"`
}

// Backend defines how krakend should connect to the backend service (the API resource to consume)
// and how it should process the received response
// type Backend struct {
// 	// the name of the group the response should be moved to. If empty, the response is
// 	// not changed
// 	Group string `mapstructure:"group"`
// 	// HTTP method of the request to send to the backend
// 	Method string `mapstructure:"method"`
// 	// Set of hosts of the API
// 	Host []string `mapstructure:"host"`
// 	// False if the hostname should be sanitized
// 	HostSanitizationDisabled bool `mapstructure:"disable_host_sanitize"`
// 	// URL pattern to use to locate the resource to be consumed
// 	URLPattern string `mapstructure:"url_pattern"`
// 	// set of response fields to remove. If empty, the filter id not used
// 	Blacklist []string `mapstructure:"blacklist"`
// 	// set of response fields to allow. If empty, the filter id not used
// 	Whitelist []string `mapstructure:"whitelist"`
// 	// map of response fields to be renamed and their new names
// 	Mapping map[string]string `mapstructure:"mapping"`
// 	// the encoding format
// 	Encoding string `mapstructure:"encoding"`
// 	// the response to process is a collection
// 	IsCollection bool `mapstructure:"is_collection"`
// 	// name of the field to extract to the root. If empty, the formater will do nothing
// 	Target string `mapstructure:"target"`
// 	// name of the service discovery driver to use
// 	SD string `mapstructure:"sd"`
// }

// parseableServiceConfig represents the raw service config values before being parsed / sanitized
type parseableServiceConfig struct {
	// Endpoints             []*parseableEndpointConfig `json:"endpoints"`
	Timeout               string `json:"timeout"`
	CacheTTL              string `json:"cache_ttl"`
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	ReadTimeout           string `json:"read_timeout"`
	WriteTimeout          string `json:"write_timeout"`
	IdleTimeout           string `json:"idle_timeout"`
	ReadHeaderTimeout     string `json:"read_header_timeout"`
	DisableKeepAlives     bool   `json:"disable_keep_alives"`
	DisableCompression    bool   `json:"disable_compression"`
	MaxIdleConns          int    `json:"max_idle_connections"`
	IdleConnTimeout       string `json:"idle_connection_timeout"`
	ResponseHeaderTimeout string `json:"response_header_timeout"`
	ExpectContinueTimeout string `json:"expect_continue_timeout"`
	OutputEncoding        string `json:"output_encoding"`
	DialerTimeout         string `json:"dialer_timeout"`
	DialerFallbackDelay   string `json:"dialer_fallback_delay"`
	DialerKeepAlive       string `json:"dialer_keep_alive"`
	Debug                 bool
}

func (p parseableServiceConfig) normalize() ServiceConfig {
	cfg := ServiceConfig{
		Timeout:               parseDuration(p.Timeout),
		CacheTTL:              parseDuration(p.CacheTTL),
		Host:                  p.Host,
		Port:                  p.Port,
		Debug:                 p.Debug,
		ReadTimeout:           parseDuration(p.ReadTimeout),
		WriteTimeout:          parseDuration(p.WriteTimeout),
		IdleTimeout:           parseDuration(p.IdleTimeout),
		ReadHeaderTimeout:     parseDuration(p.ReadHeaderTimeout),
		DisableKeepAlives:     p.DisableKeepAlives,
		DisableCompression:    p.DisableCompression,
		MaxIdleConns:          p.MaxIdleConns,
		IdleConnTimeout:       parseDuration(p.IdleConnTimeout),
		ResponseHeaderTimeout: parseDuration(p.ResponseHeaderTimeout),
		ExpectContinueTimeout: parseDuration(p.ExpectContinueTimeout),
		DialerTimeout:         parseDuration(p.DialerTimeout),
		DialerFallbackDelay:   parseDuration(p.DialerFallbackDelay),
		DialerKeepAlive:       parseDuration(p.DialerKeepAlive),
	}

	return cfg
}

func parseDuration(v string) time.Duration {
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0
	}
	return d
}

// Init initializes service config default values
func (s *ServiceConfig) Init() error {
	if s.Port == 0 {
		s.Port = DefaultPort
	}

	if s.MaxIdleConns == 0 {
		s.MaxIdleConns = DefaultMaxIdleConns
	}

	if s.Timeout == 0 {
		s.Timeout = DefaultTimeout
	}

	return nil
}

// Parser interface declaration
type Parser interface {
	Parse(configFile string) (ServiceConfig, error)
}

// NewParser factory returns a new parser
func NewParser() Parser {
	return parser{}
}

type parser struct {
}

// Parse implements the Parser interface
func (p parser) Parse(configFile string) (cfg ServiceConfig, err error) {
	data, err := ioutil.ReadFile(configFile)
	parseableCfg := new(parseableServiceConfig)
	if err != nil {
		return cfg, fmt.Errorf("Error parsing config file: %s", err.Error())
	}
	if err = json.Unmarshal(data, &parseableCfg); err != nil {
		return cfg, fmt.Errorf("Error parsing config file: %s", err.Error())
	}

	cfg = parseableCfg.normalize()

	// Init default values
	err = cfg.Init()
	if err != nil {
		return cfg, fmt.Errorf("Error initializing config: %s", err.Error())
	}

	return cfg, err
}
