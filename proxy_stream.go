package hive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Code-Pundits/go-config"
)

// ProxyStream processes a request in a given context and returns a stream of ProxyResponses
type ProxyStream func(request *ProxyRequest) <-chan *BackendResult

// ProxyStreamFactory returns a new Proxy Stream
func ProxyStreamFactory(cfg *config.EndpointConfig) (ProxyStream, error) {
	remotes := cfg.Backends

	return func(req *ProxyRequest) <-chan *BackendResult {
		out := make(chan *BackendResult)
		var wg sync.WaitGroup

		for _, r := range remotes {
			go func(remote *config.BackendConfig) {
				group := remote.Group
				fmt.Println(group)
				// Create the stream result
				result := &BackendResult{Group: group}

				defer func() {
					// Send the stream result to the channel
					out <- result
					// Always call done on the wait group
					wg.Done()
				}()

				// Build and fire the remote request
				httpResp, err := httpClient.Do(buildRemoteRequest(remote, req))
				if err != nil {
					result.StatusCode = http.StatusInternalServerError
					result.Errors = append(result.Errors, err.Error())
					return
				}
				result.StatusCode = httpResp.StatusCode
				defer httpResp.Body.Close()
				if err = HTTPStatusHandler(httpResp); err != nil {
					result.Errors = append(result.Errors, err.Error())
					return
				}
				var data map[string]interface{}
				jsonBytes, err := ioutil.ReadAll(httpResp.Body)
				if err != nil {
					result.Errors = append(result.Errors, fmt.Errorf("error reading http response body: %s", err.Error()).Error())
					return
				}
				err = json.Unmarshal(jsonBytes, &data)
				if err != nil {
					result.Errors = append(result.Errors, fmt.Errorf("error reading http response body: %s", err.Error()).Error())
					return
				}
				result.Data = data
			}(r)
			// Add wait counter
			wg.Add(1)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}, nil
}
