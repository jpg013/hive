package main

import (
	"fmt"

	logging "github.com/Code-Pundits/go-logger"
	"github.com/jpg013/hive"
)

func main() {
	// Create new hive instance
	h := hive.New(logging.NewLogger())
	h.LoadConfig("./examples/basic/config.json")
	err := h.RunServer()
	fmt.Println(err)
}
