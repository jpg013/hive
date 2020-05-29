package hive

import "net/http"

var (
	defaultHeadersToSend []string = []string{"Content-Type"}
)

var httpClient = &http.Client{}
