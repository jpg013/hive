package logging

import (
	"io"
	"sync/atomic"
)

// Transport represents type with Level and ioWriter
type Transport struct {
	Level Level
	Out   io.Writer
}

// TransportConfig represents config type to be passed to Transport factories
type TransportConfig struct {
	FilePath   string
	FileName   string
	Level      Level
	KafkaHosts []string
}

// StdOutTransportConfig represents config for standard out transport
type StdOutTransportConfig struct {
	Level Level
}

// GetLogLevel coerces transport log level to int32 and returns value
func (t *Transport) GetLogLevel() Level {
	return Level(atomic.LoadUint32((*uint32)(&t.Level)))
}
