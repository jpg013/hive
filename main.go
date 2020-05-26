package main

import "github.com/jpg013/hive/config"

// Hive is the shitz
type Hive interface {
	RunServer(config.ServiceConfig)
	// RegisterEndpoint(config.EndpoizntConfig) error
}

func New() Hive {
	return &Router{}
}

func main() {

}
