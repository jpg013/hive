package hive

// Response represents the backend remote response object
type Response struct {
	Group      string                 `json:"group"`
	Data       map[string]interface{} `json:"data"`
	Status     int                    `json:"status"`
	Errors     map[string]string      `json:"errors"`
	IsComplete bool                   `json:"is_complete"`
}

// NewResponse factory returns a new response
func NewResponse() *Response {
	return &Response{
		Errors: make(map[string]string, 0),
		Data:   make(map[string]interface{}, 0),
	}
}
