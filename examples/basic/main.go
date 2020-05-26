package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jpg013/hive/config"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cfgPath := filepath.Join(filepath.Dir(ex), "examples", "basic", "config.json")
	cfgFile := flag.String("c", cfgPath, "Path to the configuration filename")
	flag.Parse()

	parser := config.NewParser()
	serviceCfg, err := parser.Parse(*cfgFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(serviceCfg)
}
