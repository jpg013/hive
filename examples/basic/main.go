package main

import (
	"fmt"
	"net/http"

	"github.com/Code-Pundits/go-config"
	logging "github.com/Code-Pundits/go-logger"
	"github.com/Code-Pundits/go-middleware"
	"github.com/jpg013/hive"
	"github.com/labstack/echo/v4"
)

/*
	Handler func!
	func (c echo.Context) error {

	}
*/

func BullshitMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("IN BULLSHIT MIDDLEWARE")
		return next(c)
	}
}

func HandleAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Response().Status != http.StatusOK {
			return next(c)
		}

		data, ok := c.Get("data").(map[string]*hive.Response)

		if !ok || data == nil {
			return next(c)
		}

		// Reset the data, we don't want to send the auth token to the client
		c.Set("data", nil)

		authResponse, ok := data["auth"]

		if !ok {
			return next(c)
		}

		fmt.Println("BITCH FUCK YEA!!!")
		fmt.Println(authResponse.Data["refresh_token"].(string))
		fmt.Println(authResponse.Data["token"].(string))

		return next(c)

		// Stat
		// Data       map[string]interface{} `json:"data"`
		// Reset the data, we don't want to send the auth token to the client
		// c.Set("data", nil)
		// cookie := new(http.Cookie)
		// cookie.Name = "auth"
		// cookie.Value = "jon"
		// cookie.Expires = time.Now().Add(24 * time.Hour)
		// c.SetCookie(cookie)
		// c.
		// return next(c)
	}
}

func main() {
	logger := logging.
		NewLogger().
		WithLevel(logging.InfoLevel).
		WithTransports(
			logging.NewStdOutTransport(logging.StdOutTransportConfig{Level: logging.InfoLevel}),
		).
		WithDefaults(
			&logging.FieldPair{Name: "Component", Value: "hive"},
		)

	// Create new hive instance
	h := hive.New(logger)
	cfg, _ := config.NewParser().Parse("./examples/basic/config.json")

	h.RegisterEndpoint(&hive.EndpointConfig{
		Endpoint:       "/authenticate",
		Method:         "POST",
		OutputEncoding: "json",
		AfterMiddlewares: []echo.MiddlewareFunc{
			BullshitMiddleware,
			HandleAuthentication,
		},
		Backends: []*hive.BackendConfig{
			{
				Group:  "auth",
				Scheme: "http",
				Method: "POST",
				Host:   "127.0.0.1:9001",
				Path:   "/auth",
			},
		},
	})

	h.UseMiddleware(middleware.LogRequest(cfg, logger))
	h.RunServer(cfg.ServiceConfig)
}

// {
// 	"group": "password_hash",
// 	"scheme": "http",
// 	"method": "POST",
// 	"host":   "127.0.0.1:9001",
// 	"path":   "/password/hash"
// },

// Backends []*BackendConfig
// // the encoding format
// OutputEncoding string `json:"output_encoding"`
