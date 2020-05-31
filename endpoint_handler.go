package hive

import (
	"net/http"

	"github.com/Code-Pundits/go-config"
	"github.com/labstack/echo/v4"
)

func EndpointHandler(cfg config.EndpointConfig) echo.HandlerFunc {
	proxy := ProxyFactory(cfg)

	return func(c echo.Context) error {
		request := NewRequest(c)
		resp, err := proxy(request)

		if err != nil {
			return err
		}

		c.JSON(http.StatusPartialContent, "bitch is this working!")
		return c.JSON(resp.Status, resp)
	}
}
