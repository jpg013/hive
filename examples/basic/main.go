package main

import (
	"fmt"

	"github.com/Code-Pundits/go-config"
	logging "github.com/Code-Pundits/go-logger"
	"github.com/jpg013/hive"
)

func main() {
	logger := logging.NewLogger()
	cfg, _ := config.NewParser().Parse("./examples/basic/config.json")

	// Create new hive instance
	h := hive.New(logger, cfg)

	h.LoadConfig("./examples/basic/config.json")

	// h.RegisterEndpoint(&hive.Endpoint{
	// 	Endpoint: "/authenticate",
	// 	Method:   "POST",
	// 	Backends: []*hive.Backend{{
	// 		Scheme: "http",
	// 		Method: "POST",
	// 		Host:   "127.0.0.1:9001",
	// 		Path:   "/auth",
	// 	}},
	// })

	// h.RegisterEndpoint(&hive.Endpoint{
	// 	Endpoint: "/hash_password",
	// 	Method:   "POST",
	// 	Backends: []*hive.Backend{{
	// 		Scheme: "http",
	// 		Method: "POST",
	// 		Host:   "127.0.0.1:9001",
	// 		Path:   "/password/hash",
	// 	}},
	// })

	err := h.RunServer()
	fmt.Println(err)
}
