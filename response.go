package hive

type Response struct {
	Group  string                 `json:"group"`
	Data   map[string]interface{} `json:"data"`
	Status int                    `json:"status"`
	Errors []string               `json:"errors"`
}
