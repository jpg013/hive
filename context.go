package hive

import "io"

// Context represents a request context
type Context interface {
	Body() io.ReadCloser
	Header() map[string][]string
	Method() string
	Set(string, interface{})
	Get(string) interface{}
}
