package echo

import (
	"github.com/jpg013/hive/config"
	"github.com/labstack/echo/v4"
)

func EndpointHandler(cfg config.EndpointConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
