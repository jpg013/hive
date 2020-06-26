package hive

import (
	"context"
	"os"

	stdLog "log"

	"github.com/Code-Pundits/go-config"
	logging "github.com/Code-Pundits/go-logger"

	"github.com/jpg013/hive/http"
	"github.com/labstack/echo/v4"
	gocolor "github.com/labstack/gommon/color"
)

const (
	// Version of Echo
	Version = "VERSION_HOLDER"
	banner  = `
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ %s
High performance, minimalist Go web framework
____________________________________O/_______
                                    O\
`
)

// Hive is the shitz
type Hive interface {
	RunServer(config.ServiceConfig) error
	RegisterEndpoint(*EndpointConfig)
	UseMiddleware(...echo.MiddlewareFunc)
}

// hive implements the Hive interface
type hive struct {
	router    Router
	logger    logging.Logger
	StdLogger *stdLog.Logger
	colorer   *gocolor.Color
}

func (h *hive) UseMiddleware(fns ...echo.MiddlewareFunc) {
	h.router.middleware(fns...)
}

func (h *hive) RunServer(cfg config.ServiceConfig) error {
	h.colorer.SetOutput(os.Stdout)
	h.colorer.Printf(h.colorer.Cyan(banner), h.colorer.Red("v"+Version))

	return http.RunServer(context.Background(), cfg, h.router.handler())
}

func (h *hive) RegisterEndpoint(e *EndpointConfig) {
	h.router.endpoint(e)
}

// New returns a Hive implementation
func New(logger logging.Logger) Hive {
	router := NewRouter(logger)

	return &hive{
		logger:  logger,
		router:  router,
		colorer: gocolor.New(),
	}
}
