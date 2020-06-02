package hive

import (
	"encoding/json"
	"fmt"

	"github.com/Code-Pundits/go-config"
	"github.com/labstack/echo/v4"
)

func EndpointHandler(cfg *config.EndpointConfig) echo.HandlerFunc {
	proxy, err := ProxyFactory(cfg)
	render := getRender(cfg)

	if err != nil {
		panic(err)
	}

	return func(c echo.Context) error {
		request := NewRequest(c)
		resp := proxy(request)

		if err != nil {
			return err
		}

		return render(c, resp)
	}
}

func StreamingEndpointHandler(cfg *config.EndpointConfig) echo.HandlerFunc {
	proxyStream, err := ProxyStreamFactory(cfg)

	if err != nil {
		panic(err)
	}

	return func(c echo.Context) error {
		request := NewRequest(c)
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(defaultHTTPStatus)
		enc := json.NewEncoder(c.Response())
		for chunk := range proxyStream(request) {
			fmt.Println("What the shit?")
			if err := enc.Encode(chunk); err != nil {
				return nil
			}
			c.Response().Flush()
		}
		return nil
	}
}

func GetHandler(cfg *config.EndpointConfig) echo.HandlerFunc {
	if cfg.StreamResponse == true {
		return StreamingEndpointHandler(cfg)
	}

	return EndpointHandler(cfg)
}
