package hive

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Code-Pundits/go-config"
)

// Proxy processes a request in a given context and returns a response and an error
type Proxy func(request *ProxyRequest) (*ProxyResponse, error)

var (
	errNoBackend = errors.New("Endpoint must have at least 1 backend")
)

func multiProxyFactory(remotes []*config.BackendConfig) Proxy {
	return func(req *ProxyRequest) (*ProxyResponse, error) {
		var wait sync.WaitGroup
		proxyResponse := NewResponse()

		for _, r := range remotes {
			go func(remote *config.BackendConfig) {
				defer func() {
					wait.Done()
				}()
				group := remote.Group
				httpResp, err := httpClient.Do(buildRemoteRequest(remote, req))
				if err != nil {
					panic(err)
				}

				if err != nil {
					proxyResponse.Errors[group] = err.Error()
					return
				}

				defer httpResp.Body.Close()

				if err = HTTPStatusHandler(httpResp); err != nil {
					proxyResponse.Errors[group] = err.Error()
					return
				}

				var data map[string]interface{}
				jsonBytes, err := ioutil.ReadAll(httpResp.Body)
				if err != nil {
					proxyResponse.Errors[group] = fmt.Errorf("error reading http response: %s", err.Error()).Error()
					return
				}
				err = json.Unmarshal(jsonBytes, &data)
				if err != nil {
					proxyResponse.Errors[group] = fmt.Errorf("error parsing http response JSON: %s", err.Error()).Error()
					return
				}
				proxyResponse.Data[group] = data
			}(r)
			// Add wait counter
			wait.Add(1)
		}

		wait.Wait()
		// Set the IsComplete flag to true
		proxyResponse.IsComplete = true
		// Always set the http status flag to Status OK for multiple proxy
		proxyResponse.Status = http.StatusOK
		return proxyResponse, nil
	}
}

func singleProxyFactory(remote *config.BackendConfig) Proxy {
	group := remote.Group

	return func(req *ProxyRequest) (*ProxyResponse, error) {
		httpResp, err := httpClient.Do(buildRemoteRequest(remote, req))

		if err != nil {
			return nil, err
		}

		defer httpResp.Body.Close()
		if err = HTTPStatusHandler(httpResp); err != nil {
			return &ProxyResponse{
				Status: httpResp.StatusCode,
				Errors: map[string]string{group: err.Error()},
			}, nil
		}

		proxyResponse := NewResponse()
		proxyResponse.Status = httpResp.StatusCode

		var data map[string]interface{}
		jsonBytes, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonBytes, &data)
		if err != nil {
			return nil, err
		}

		proxyResponse.Data[group] = data
		proxyResponse.IsComplete = true

		return proxyResponse, nil
	}
}

// ProxyFactory returns a new Proxy
func ProxyFactory(cfg *config.EndpointConfig) (Proxy, error) {
	if len(cfg.Backends) == 0 {
		return nil, errNoBackend
	}

	if len(cfg.Backends) == 1 {
		return singleProxyFactory(cfg.Backends[0]), nil
	}

	return multiProxyFactory(cfg.Backends), nil
}
