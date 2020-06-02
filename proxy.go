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
type Proxy func(request *ProxyRequest) *ProxyResponse

var (
	errNoBackend = errors.New("Endpoint must have at least 1 backend")
)

func multiProxyFactory(remotes []*config.BackendConfig) Proxy {
	return func(req *ProxyRequest) *ProxyResponse {
		// wait group allows us to wait for responses from all the remotes
		var wg sync.WaitGroup
		proxyResp := NewProxyResponse()

		for _, r := range remotes {
			go func(remote *config.BackendConfig) {
				defer func() {
					// Always call done on the wait group
					wg.Done()
				}()

				group := remote.Group
				// add the result group to the proxy response
				proxyResp.AddGroup(group)

				// Build and fire the remote request
				httpResp, err := httpClient.Do(buildRemoteRequest(remote, req))
				if err != nil {
					proxyResp.
						AddStatus(group, http.StatusInternalServerError).
						AddError(group, err)
					return
				}

				proxyResp.AddStatus(group, httpResp.StatusCode)
				defer httpResp.Body.Close()
				if err = HTTPStatusHandler(httpResp); err != nil {
					proxyResp.AddError(group, err)
					return
				}
				var data map[string]interface{}
				jsonBytes, err := ioutil.ReadAll(httpResp.Body)
				if err != nil {
					proxyResp.AddError(group, fmt.Errorf("error reading http response body: %s", err.Error()))
					return
				}
				err = json.Unmarshal(jsonBytes, &data)
				if err != nil {
					if err != nil {
						proxyResp.AddError(group, fmt.Errorf("error reading http response body: %s", err.Error()))
						return
					}
				}

				proxyResp.AddData(group, data)
			}(r)
			// Add wait counter
			wg.Add(1)
		}

		wg.Wait()
		// Always set the http status flag to Status OK for multiple proxy
		proxyResp.Status = http.StatusOK
		return proxyResp
	}
}

func singleProxyFactory(remote *config.BackendConfig) Proxy {
	group := remote.Group

	return func(req *ProxyRequest) *ProxyResponse {
		proxyResp := NewProxyResponse()
		httpResp, err := httpClient.Do(buildRemoteRequest(remote, req))

		if err != nil {
			return proxyResp.
				AddStatus(group, http.StatusInternalServerError).
				AddError(group, err)
		}

		proxyResp.AddStatus(group, httpResp.StatusCode)
		defer httpResp.Body.Close()
		if err = HTTPStatusHandler(httpResp); err != nil {
			return proxyResp.AddError(group, err)
		}
		var data map[string]interface{}
		jsonBytes, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return proxyResp.AddError(group, fmt.Errorf("error reading http response body: %s", err.Error()))
		}
		err = json.Unmarshal(jsonBytes, &data)
		if err != nil {
			if err != nil {
				return proxyResp.AddError(group, fmt.Errorf("error reading http response body: %s", err.Error()))
			}
		}

		return proxyResp.AddData(group, data)
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
