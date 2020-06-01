package hive

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var httpAllowedStatusCodes map[uint]bool = map[uint]bool{
	http.StatusOK:      true,
	http.StatusCreated: true,
}

// HTTPResponseError is the error to be returned by the HTTPStatusHandler
type HTTPResponseError struct {
	Code int    `json:"http_status_code"`
	Msg  string `json:"http_body,omitempty"`
	name string
}

// Error returns the error message
func (r HTTPResponseError) Error() string {
	return r.Msg
}

// Name returns the name of the error
func (r HTTPResponseError) Name() string {
	return r.name
}

// StatusCode returns the status code returned by the backend
func (r HTTPResponseError) StatusCode() int {
	return r.Code
}

// HTTPStatusHandler checks the http.Response for invalid responses
func HTTPStatusHandler(resp *http.Response) error {
	_, ok := httpAllowedStatusCodes[uint(resp.StatusCode)]

	if !ok {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			body = []byte{}
		}
		resp.Body.Close()

		// Copy the ReadCloser
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		// Try to parse the message from the body
		var data map[string]string
		var msg string
		err = json.Unmarshal(body, &data)
		if err != nil {
			msg = data["message"]
		}
		if msg == "" {
			msg = string(body)
		}
		return HTTPResponseError{
			Code: resp.StatusCode,
			Msg:  msg,
		}
	}

	return nil
}
