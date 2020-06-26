package hive

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RenderFunc takes a context and renders the response to the client
type RenderFunc func(echo.Context) error

type StreamRender func(echo.Context) func(*Response) error

var emptyResponse interface{}

var encodingJSON string = "json"

var (
	renderRegister = map[string]RenderFunc{
		encodingJSON: jsonRender,
	}
)

func getStreamRender(cfg *EndpointConfig) StreamRender {
	fallback := jsonStreamRender

	return fallback
}

func getRender(cfg *EndpointConfig) echo.HandlerFunc {
	var handlerFunc echo.HandlerFunc

	handlerFunc = jsonRender

	for i := range cfg.AfterMiddlewares {
		fn := cfg.AfterMiddlewares[len(cfg.AfterMiddlewares)-i-1]
		handlerFunc = fn(handlerFunc)
	}

	return handlerFunc
}

func jsonRender(c echo.Context) error {
	data := c.Get("data")
	if data == nil {
		return c.JSON(http.StatusOK, emptyResponse)
	}
	return c.JSON(http.StatusOK, data)
}

func jsonStreamRender(c echo.Context) func(*Response) error {
	c.Response().Header().Set(headerContentType, mimeApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	enc := json.NewEncoder(c.Response())

	return func(data *Response) error {
		if err := enc.Encode(data); err != nil {
			return err
		}
		// Flush data to writer
		c.Response().Flush()
		return nil
	}
}
