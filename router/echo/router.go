package echo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jpg013/hive/config"
	"github.com/jpg013/hive/logging"
	"github.com/jpg013/hive/router"
	"github.com/labstack/echo/v4"
)

// echoRouter implements Echo Router interface
type echoRouter struct {
	logger      logging.Logger
	Engine      *echo.Echo
	Middlewares []echo.MiddlewareFunc
}

func (r *echoRouter) Use(fns ...echo.MiddlewareFunc) {
	for _, fn := range fns {
		r.Middlewares = append(r.Middlewares, fn)
	}
}

func (r *echoRouter) RegisterEndpoint(cfg config.EndpointConfig) {
	method := strings.ToTitle(cfg.Method)

	switch method {
	case http.MethodGet:
		r.Engine.GET(cfg.Endpoint, EndpointHandler(cfg))
	case http.MethodPost:
		r.Engine.POST(cfg.Endpoint, EndpointHandler(cfg))
	case http.MethodPut:
		r.Engine.PUT(cfg.Endpoint, EndpointHandler(cfg))
	case http.MethodPatch:
		r.Engine.PATCH(cfg.Endpoint, EndpointHandler(cfg))
	case http.MethodDelete:
		r.Engine.DELETE(cfg.Endpoint, EndpointHandler(cfg))
	default:
		r.logger.Error(fmt.Sprintf("Unsuported endpoint method: %s", method))
	}
}

// New factory returns a new echo router
func New(logger logging.Logger) router.Router {
	return &echoRouter{
		logger:      logger,
		Engine:      echo.New(),
		Middlewares: make([]echo.MiddlewareFunc, 0),
	}
}
