package logging

// Fields represent that values that can be in the log message.
type Fields struct {
	TraceID          string  `json:"traceId,omitempty"`
	ID               string  `json:"id,omitempty"`
	ParentID         string  `json:"parentId,omitempty"`
	Timestamp        string  `json:"timestamp,omitempty"`
	Hostname         string  `json:"hostname,omitempty"`
	LogType          string  `json:"type,omitempty"`
	Facility         string  `json:"facility,omitempty"`
	Severity         string  `json:"severity,omitempty"`
	Namespace        string  `json:"namespace,omitempty"`
	Component        string  `json:"component,omitempty"`
	Version          string  `json:"version,omitempty"`
	Message          string  `json:"message,omitempty"`
	Stack            string  `json:"stackTrace,omitempty"`
	LineNumber       int     `json:"lineNumber,omitempty"`
	HTTPMethod       string  `json:"httpMethod,omitempty"`
	HTTPRoute        string  `json:"httpRoute,omitempty"`
	Parameters       string  `json:"parameters,omitempty"`
	HTTPResponseCode int     `json:"httpResponseCode,omitempty"`
	Duration         float64 `json:"duration,omitempty"`
	Environment      string  `json:"environment,omitempty"`
	Subenvironment   string  `json:"subenvironment,omitempty"`
	BuildNumber      int     `json:"buildNumber,omitempty"`
}

// FieldPair represents a key / value for logger field.
type FieldPair struct {
	Name  string
	Value interface{}
}

// CopyDataFields loops over field pairs and tries sets value for Field if zero value
func CopyDataFields(data *Fields, args ...*FieldPair) *Fields {
	for _, arg := range args {
		if arg.Name == "TraceID" && isZeroString(data.TraceID) {
			data.TraceID = arg.Value.(string)
		}
		if arg.Name == "ID" && isZeroString(data.ID) {
			data.ID = arg.Value.(string)
		}
		if arg.Name == "ParentID" && isZeroString(data.ParentID) {
			data.ParentID = arg.Value.(string)
		}
		if arg.Name == "Timestamp" && isZeroString(data.Timestamp) {
			data.Timestamp = arg.Value.(string)
		}
		if arg.Name == "Hostname" && isZeroString(data.Hostname) {
			data.Hostname = arg.Value.(string)
		}
		if arg.Name == "LogType" && isZeroString(data.LogType) {
			data.LogType = arg.Value.(string)
		}
		if arg.Name == "Facility" && isZeroString(data.Facility) {
			data.Facility = arg.Value.(string)
		}
		if arg.Name == "Severity" && isZeroString(data.Severity) {
			data.Severity = arg.Value.(string)
		}
		if arg.Name == "Namespace" && isZeroString(data.Namespace) {
			data.Namespace = arg.Value.(string)
		}
		if arg.Name == "Component" && isZeroString(data.Component) {
			data.Component = arg.Value.(string)
		}
		if arg.Name == "Version" && isZeroString(data.Version) {
			data.Version = arg.Value.(string)
		}
		if arg.Name == "Message" && isZeroString(data.Message) {
			data.Message = arg.Value.(string)
		}
		if arg.Name == "Stack" && isZeroString(data.Stack) {
			data.Stack = arg.Value.(string)
		}
		if arg.Name == "LineNumber" && isZeroInt(data.LineNumber) {
			data.LineNumber = arg.Value.(int)
		}
		if arg.Name == "HTTPMethod" && isZeroString(data.HTTPMethod) {
			data.HTTPMethod = arg.Value.(string)
		}
		if arg.Name == "HTTPRoute" && isZeroString(data.HTTPRoute) {
			data.HTTPRoute = arg.Value.(string)
		}
		if arg.Name == "Parameters" && isZeroString(data.Parameters) {
			data.Parameters = arg.Value.(string)
		}
		if arg.Name == "HTTPResponseCode" && isZeroInt(data.HTTPResponseCode) {
			data.HTTPResponseCode = arg.Value.(int)
		}
		if arg.Name == "Duration" && isZeroFloat(data.Duration) {
			data.Duration = arg.Value.(float64)
		}
		if arg.Name == "Environment" && isZeroString(data.Parameters) {
			data.Environment = arg.Value.(string)
		}
		if arg.Name == "Subenvironment" && isZeroString(data.Parameters) {
			data.Subenvironment = arg.Value.(string)
		}
		if arg.Name == "BuildNumber" && isZeroInt(data.BuildNumber) {
			data.BuildNumber = arg.Value.(int)
		}
	}

	return data
}
