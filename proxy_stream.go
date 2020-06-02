package hive

// ProxyStream processes a request in a given context and returns a stream of ProxyResponses
type ProxyStream func(request *ProxyRequest) (*ProxyResponse, error)
