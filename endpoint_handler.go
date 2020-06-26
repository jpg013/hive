package hive

import (
	"github.com/labstack/echo/v4"
)

func EndpointHandler(cfg *EndpointConfig) echo.HandlerFunc {
	proxy, err := ProxyFactory(cfg)
	render := getRender(cfg)

	if err != nil {
		panic(err)
	}

	return func(c echo.Context) error {
		request := NewRequest(c)
		result := make(map[string]*Response)
		for resp := range proxy(request) {
			result[resp.Group] = resp
		}

		// Set the context data
		c.Set("data", result)

		return render(c)
	}
}

func StreamingEndpointHandler(cfg *EndpointConfig) echo.HandlerFunc {
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
