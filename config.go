package hive

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type HiveConfig struct {
	Endpoints []*EndpointConfig
}

type parseableHiveConfig struct {
	Endpoints []*parseableEndpointConfig
}

func (cfg *parseableHiveConfig) normalize() (*HiveConfig, error) {
	endpoints := make([]*EndpointConfig, len(cfg.Endpoints))

	for i, k := range cfg.Endpoints {
		val, err := k.normalize()
		if err != nil {
			return nil, err
		}
		endpoints[i] = val
	}

	return &HiveConfig{endpoints}, nil
}

// // Parser interface declaration
type Parser interface {
	Parse(configFile string) (*HiveConfig, error)
}

// // NewParser factory returns a new parser
func NewParser() Parser {
	return parser{}
}

type parser struct {
}

// // Parse implements the Parser interface
func (p parser) Parse(filePath string) (cfg *HiveConfig, err error) {
	cfgPath := flag.String("p", filePath, "Path to the configuration filename")
	flag.Parse()
	data, err := ioutil.ReadFile(*cfgPath)
	parseableCfg := new(parseableHiveConfig)

	if err != nil {
		return cfg, fmt.Errorf("Error parsing config file: %s", err.Error())
	}
	if err = json.Unmarshal(data, &parseableCfg); err != nil {
		return cfg, fmt.Errorf("Error parsing config file: %s", err.Error())
	}

	cfg, err = parseableCfg.normalize()
	if err != nil {
		return cfg, fmt.Errorf("Error parsing config file: %s", err.Error())
	}

	return cfg, nil
}
