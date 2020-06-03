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

func StreamingEndpointHandler(cfg *config.EndpointConfig) echo.HandlerFunc {
	return nil
	// proxyStream, err := ProxyStreamFactory(cfg)
	// if err != nil {
	// 	panic(err)
	// }
	// render := getStreamRender(cfg)
	// return func(c echo.Context) error {
	// 	request := NewRequest(c)
	// 	re := render(c)

	// 	for data := range proxyStream(request) {
	// 		if err := re(data); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// }
}

func GetHandler(cfg *config.EndpointConfig) echo.HandlerFunc {
	if cfg.StreamResponse == true {
		return StreamingEndpointHandler(cfg)
	}

	return EndpointHandler(cfg)
}
