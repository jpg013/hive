package hive

import (
	"context"

	"github.com/Code-Pundits/go-config"
	logging "github.com/Code-Pundits/go-logger"
)

// Hive is the shitz
type Hive interface {
	RunServer() error
	RegisterEndpoint(*config.EndpointConfig)
	LoadConfig(string) error
}

// hive implements the Hive interface
type hive struct {
	config config.Configuration
	router Router
	logger logging.Logger
}

func (h *hive) LoadConfig(path string) error {
	cfg, err := config.NewParser().Parse(path)

	if err != nil {
		return err
	}

	h.config = cfg

	// Register each endpoint with config
	for _, endpoint := range cfg.Endpoints {
		h.RegisterEndpoint(endpoint)
	}

	return nil
}

func (h *hive) RunServer() error {
	return RunServer(context.Background(), h.config.ServiceConfig, h.router.handler())
}

func (h *hive) RegisterEndpoint(e *config.EndpointConfig) {
	h.router.endpoint(e)
}

// New returns a Hive implementation
func New(logger logging.Logger) Hive {
	router := NewRouter(logger)

	return &hive{logger: logger, router: router}
}
