package hive

import (
	"encoding/json"

	"github.com/Code-Pundits/go-config"
	"github.com/jpg013/hive/transport/http"
	"github.com/labstack/echo/v4"
)

type Render func(echo.Context, *ProxyResult) error

type StreamRender func(echo.Context) func(*http.Response) error

var emptyResponse interface{}

var encodingJSON string = "json"

var (
	renderRegister = map[string]Render{
		encodingJSON: jsonRender,
	}
)

func getStreamRender(cfg *config.EndpointConfig) StreamRender {
	fallback := jsonStreamRender

	return fallback
}

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

func jsonRender(c echo.Context, response *ProxyResult) error {
	if response == nil {
		return c.JSON(http.DefaultOKStatus, emptyResponse)
	}
	return c.JSON(http.DefaultOKStatus, response.Data)
}

func jsonStreamRender(c echo.Context) func(*http.Response) error {
	c.Response().Header().Set(headerContentType, mimeApplicationJSON)
	c.Response().WriteHeader(http.DefaultOKStatus)
	enc := json.NewEncoder(c.Response())

	return func(data *http.Response) error {
		if err := enc.Encode(data); err != nil {
			return err
		}
		// Flush data to writer
		c.Response().Flush()
		return nil
	}
}
