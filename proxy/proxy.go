package proxy

// Proxy processes a request in a given context and returns a response and an error
type Proxy func(request *Request) (*Response, error)
