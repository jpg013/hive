package hive

import "github.com/labstack/echo/v4"

var emptyResponse interface{}

func renderJson(c echo.Context, resp *Response) {
	if resp == nil {
		return
	}
}
