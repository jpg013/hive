package hive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	transport "github.com/jpg013/hive/http"
)

// Response represents the http response data
type Response struct {
	Group      string                 `json:"group"`
	Data       map[string]interface{} `json:"data"`
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"status_message"`
	Errors     []string               `json:"errors"`
	IsComplete bool                   `json:"is_complete"`
}

// NewError adds a new error to the reponse
func (r *Response) NewError(err error) {
	r.Errors = append(r.Errors, err.Error())
}

// SetStatus sets the response status code and the message
func (r *Response) SetStatus(s int) {
	r.StatusCode = s
	r.Message = http.StatusText(s)
}

// NewResponse factory returns a Response Data object
func NewResponse(group string) *Response {
	return &Response{
		Group:  group,
		Errors: make([]string, 0),
		Data:   make(map[string]interface{}, 0),
	}
}

// NewResponseHandler returns a handler that accepts an http.Response / error and returns a proxy Response
func NewResponseHandler(group string) func(*http.Response, error) *Response {
	r := NewResponse(group)

	return func(resp *http.Response, err error) *Response {
		if err != nil {
			r.NewError(err)
			r.SetStatus(http.StatusInternalServerError)
			return r
		}
		defer resp.Body.Close()
		r.SetStatus(resp.StatusCode)

		// Check status handler
		if err = transport.StatusHandler(resp); err != nil {
			// Add the error to the response
			r.NewError(err)
			return r
		}

		// declare data
		var jsonData map[string]interface{}
		jsonBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			r.NewError(fmt.Errorf("error reading http response body: %s", err.Error()))
			return r
		}

		err = json.Unmarshal(jsonBytes, &jsonData)
		if err != nil {
			r.NewError(fmt.Errorf("error reading http response body: %s", err.Error()))
			return r
		}

		r.Data = jsonData

		return r
	}
}
