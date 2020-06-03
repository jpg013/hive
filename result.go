package hive

import "github.com/jpg013/hive/transport/http"

// ProxyResult is a collection of http responses that is returned by the proxy
type ProxyResult struct {
	Data       map[string]*http.Response `json:"data"`
	IsComplete bool                      `json:"is_complete"`
	Status     int                       `json:"status"`
}

// NewProxyResult factory returns a new proxy response
func NewProxyResult() *ProxyResult {
	return &ProxyResult{
		IsComplete: false,
		Status:     0,
		Data:       make(map[string]*http.Response, 0),
	}
}
