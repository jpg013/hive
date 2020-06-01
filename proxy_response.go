package hive

// ProxyResponse represents the backend remote response object
type ProxyResponse struct {
	Data       map[string]interface{} `json:"data"`
	Status     int                    `json:"status"`
	Errors     map[string]string      `json:"errors"`
	IsComplete bool                   `json:"is_complete"`
}

// NewResponse factory returns a new proxy response
func NewResponse() *ProxyResponse {
	return &ProxyResponse{
		Errors: make(map[string]string, 0),
		Data:   make(map[string]interface{}, 0),
	}
}
