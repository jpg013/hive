package proxy

import "io"

// Response is the entity returned by the proxy
type Response struct {
	Data       map[string]interface{}
	IsComplete bool
	Metadata   Metadata
	Io         io.Reader
}
