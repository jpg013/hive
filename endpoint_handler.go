package hive

import (
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
