package logging

import (
	"os"
)

// NewStdOutTransport returns a transport that writes to Stdout
func NewStdOutTransport(cfg StdOutTransportConfig) Transport {
	return Transport{
		Level: cfg.Level,
		Out:   os.Stdout,
	}
}
