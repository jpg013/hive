package hive

import (
	"github.com/Code-Pundits/go-config"
	"github.com/labstack/echo/v4"
)

type Render func(echo.Context, *Response) error

var emptyResponse interface{}

var encodingJSON string = "json"

var (
	renderRegister = map[string]Render{
		encodingJSON: jsonRender,
	}
)

func getRender(cfg *config.EndpointConfig) Render {
	fallback := jsonRender

	if cfg.OutputEncoding == "" {
		return fallback
	}

	if r, ok := renderRegister[cfg.OutputEncoding]; ok {
		return r
	}

	return fallback
}

func jsonRender(c echo.Context, response *Response) error {
	if response == nil {
		return c.JSON(defaultHTTPStatus, emptyResponse)
	}
	return c.JSON(response.Status, response.Data)
}
