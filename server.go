package hive

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Code-Pundits/go-config"
)

var onceTransportConfig sync.Once

// NewServer returns a http.Server ready to serve the injected handler
func NewServer(cfg config.ServiceConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
}

// RunServer is the main function that is called to configure and run an http server
func RunServer(ctx context.Context, cfg config.ServiceConfig, handler http.Handler) error {
	InitHTTPDefaultTransport(cfg)

	done := make(chan error)
	s := NewServer(cfg, handler)

	go func() {
		done <- s.ListenAndServe()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	}
}

// InitHTTPDefaultTransport ensures the default HTTP transport is configured just once per execution
func InitHTTPDefaultTransport(cfg config.ServiceConfig) {
	onceTransportConfig.Do(func() {
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:       cfg.DialerTimeout,
				KeepAlive:     cfg.DialerKeepAlive,
				FallbackDelay: cfg.DialerFallbackDelay,
				DualStack:     true,
			}).DialContext,
			DisableCompression:    cfg.DisableCompression,
			DisableKeepAlives:     cfg.DisableKeepAlives,
			MaxIdleConns:          cfg.MaxIdleConns,
			MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
			IdleConnTimeout:       cfg.IdleConnTimeout,
			ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
			ExpectContinueTimeout: cfg.ExpectContinueTimeout,
			TLSHandshakeTimeout:   10 * time.Second,
		}
	})
}
