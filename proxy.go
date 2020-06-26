package hive

import (
	"errors"

	"github.com/jpg013/hive/http"
)

// Proxy processes a request in a given context and returns a response and an error
type Proxy func(request *ProxyRequest) <-chan *Response

var (
	errNoBackend = errors.New("Endpoint must have at least 1 backend")
)

// func multiProxyFactory(remotes []*config.BackendConfig) Proxy {
// 	return func(p *ProxyRequest) *ProxyResponse {
// 		// wait group allows us to wait for responses from all the remotes
// 		var wg sync.WaitGroup
// 		proxyResp := NewProxyResponse()

// 		for _, r := range remotes {
// 			go func(cfg *config.BackendConfig) {
// 				defer func() {
// 					// Always call done on the wait group
// 					wg.Done()
// 				}()

// 				group := cfg.Group
// 				newBackendRequest := NewBackendRequestFactory(cfg)
// 				// add the result group to the proxy response
// 				proxyResp.AddGroup(group)

// 				// Build and fire the remote request
// 				httpResp, err := httpClient.Do(newBackendRequest(p))
// 				if err != nil {
// 					proxyResp.
// 						AddStatus(group, http.StatusInternalServerError).
// 						AddError(group, err)
// 					return
// 				}

// 				proxyResp.AddStatus(group, httpResp.StatusCode)
// 				defer httpResp.Body.Close()
// 				if err = HTTPStatusHandler(httpResp); err != nil {
// 					proxyResp.AddError(group, err)
// 					return
// 				}
// 				var data map[string]interface{}
// 				jsonBytes, err := ioutil.ReadAll(httpResp.Body)
// 				if err != nil {
// 					proxyResp.AddError(group, fmt.Errorf("error reading http response body: %s", err.Error()))
// 					return
// 				}
// 				err = json.Unmarshal(jsonBytes, &data)
// 				if err != nil {
// 					proxyResp.AddError(group, fmt.Errorf("error reading http response body: %s", err.Error()))
// 					return
// 				}

// 				proxyResp.AddData(group, data)
// 			}(r)
// 			// Add wait counter
// 			wg.Add(1)
// 		}

// 		wg.Wait()
// 		// Always set the http status flag to Status OK for multiple proxy
// 		proxyResp.Status = http.StatusOK
// 		return proxyResp
// 	}
// }

func singleProxyFactory(cfg *BackendConfig) Proxy {
	group := cfg.Group
	generateRequest := NewRequestGeneratorFactory(cfg)
	exec := http.NewExecutor()
	handler := NewResponseHandler(group)

	return func(p *ProxyRequest) <-chan *Response {
		out := make(chan *Response)

		go func() {
			out <- handler(exec.Do(generateRequest(p)))
			close(out)
		}()

		return out
	}
}

// ProxyFactory returns a new Proxy
func ProxyFactory(cfg *EndpointConfig) (Proxy, error) {
	if len(cfg.Backends) == 0 {
		return nil, errNoBackend
	}

	if len(cfg.Backends) == 1 {
		return singleProxyFactory(cfg.Backends[0]), nil
	}

	return nil, nil
}
