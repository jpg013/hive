package router

import (
	"github.com/jpg013/hive/config"
	"github.com/labstack/echo/v4"
)

// Router defines the router interface
type Router interface {
	Use(...echo.MiddlewareFunc)
	RegisterEndpoint(config.EndpointConfig)
}
