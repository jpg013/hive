package hive

import "net/http"

// BackendResult represents the backend remote response data
type BackendResult struct {
	Group      string                 `json:"group"`
	Data       map[string]interface{} `json:"data"`
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"status_message"`
	Errors     []string               `json:"errors"`
	IsComplete bool                   `json:"is_complete"`
}

// NewBackendResult factory returns a new proxy response
func NewBackendResult(group string) *BackendResult {
	return &BackendResult{
		Group:  group,
		Errors: make([]string, 0),
		Data:   make(map[string]interface{}, 0),
	}
}

// ProxyResponse is a collection of BackendResults that is returned by the proxy
type ProxyResponse struct {
	Data       map[string]*BackendResult `json:"data"`
	IsComplete bool                      `json:"is_complete"`
	Status     int                       `json:"status"`
}

// NewProxyResponse factory returns a new proxy response
func NewProxyResponse() *ProxyResponse {
	return &ProxyResponse{
		IsComplete: false,
		Status:     0,
		Data:       make(map[string]*BackendResult, 0),
	}
}

// AddData adds a backend reponse to the specified remote result group
func (p *ProxyResponse) AddData(group string, data map[string]interface{}) *ProxyResponse {
	if _, ok := p.Data[group]; !ok {
		p.Data[group] = NewBackendResult(group)
	}

	p.Data[group].Data = data
	p.Data[group].IsComplete = true
	return p
}

// AddError adds an error specified remote result group
func (p *ProxyResponse) AddError(group string, err error) *ProxyResponse {
	if _, ok := p.Data[group]; !ok {
		p.Data[group] = NewBackendResult(group)
	}

	// Add the error to the remote result
	p.Data[group].Errors = append(p.Data[group].Errors, err.Error())
	return p
}

// AddStatus adds a status code specified to remote status
func (p *ProxyResponse) AddStatus(group string, status int) *ProxyResponse {
	if _, ok := p.Data[group]; !ok {
		p.Data[group] = NewBackendResult(group)
	}

	p.Data[group].StatusCode = status
	p.Data[group].Message = http.StatusText(status)
	return p
}

// AddGroup adds a remote result group
func (p *ProxyResponse) AddGroup(group string) *ProxyResponse {
	if _, ok := p.Data[group]; !ok {
		p.Data[group] = NewBackendResult(group)
	}

	return p
}
