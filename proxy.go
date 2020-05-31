package hive

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Code-Pundits/go-config"
)

// Proxy processes a request in a given context and returns a response and an error
type Proxy func(request *Request) (*Response, error)

func ProxyFactory(cfg config.EndpointConfig) Proxy {
	return func(r *Request) (*Response, error) {
		be := cfg.Backends[0]
		resp, err := httpClient.Do(buildRemoteRequest(be, r))

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		err = HTTPStatusHandler(resp)

		// We received an error status code from the backend
		if err != nil {
			return &Response{
				Status: resp.StatusCode,
				Errors: []string{err.Error()},
			}, nil
		}

		var data map[string]interface{}
		jsonBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonBytes, &data)
		if err != nil {
			return nil, err
		}

		return &Response{
			Data:   data,
			Group:  be.Group,
			Status: resp.StatusCode,
			Errors: make([]string, 0),
		}, nil
	}
}
