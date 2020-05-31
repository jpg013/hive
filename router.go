package hive

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Code-Pundits/go-config"
	logging "github.com/Code-Pundits/go-logger"
	"github.com/labstack/echo/v4"
)

// Router defines the router interface
type Router interface {
	middleware(...echo.MiddlewareFunc)
	endpoint(config.EndpointConfig)
	handler() http.Handler
}

// httpRouter implements Router interface
type httpRouter struct {
	logger      logging.Logger
	Engine      *echo.Echo
	Middlewares []echo.MiddlewareFunc
}

func (r *httpRouter) middleware(fns ...echo.MiddlewareFunc) {
	for _, fn := range fns {
		r.Middlewares = append(r.Middlewares, fn)
	}
}

func (r *httpRouter) endpoint(e config.EndpointConfig) {
	method := strings.ToTitle(e.Method)

	switch method {
	case http.MethodGet:
		r.Engine.GET(e.Endpoint, EndpointHandler(e))
	case http.MethodPost:
		r.Engine.POST(e.Endpoint, EndpointHandler(e))
	case http.MethodPut:
		r.Engine.PUT(e.Endpoint, EndpointHandler(e))
	case http.MethodPatch:
		r.Engine.PATCH(e.Endpoint, EndpointHandler(e))
	case http.MethodDelete:
		r.Engine.DELETE(e.Endpoint, EndpointHandler(e))
	default:
		r.logger.Error(fmt.Sprintf("Unsuported endpoint method: %s", method))
	}
}

func (r *httpRouter) handler() http.Handler {
	for _, middleware := range r.Middlewares {
		r.Engine.Use(middleware)
	}

	return r.Engine
}

// NewRouter factory returns a new router interface
func NewRouter(l logging.Logger) Router {
	engine := echo.New()
	engine.Debug = true

	return &httpRouter{
		logger:      l,
		Engine:      engine,
		Middlewares: make([]echo.MiddlewareFunc, 0),
	}
}
