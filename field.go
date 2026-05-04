package httplog

import "time"

// Field represents a logger field which is contains key and value pair
type Field struct {
	Key   string
	Value any
}

// ParseFields returns slice of field as slice of any.
func ParseFields(f []Field) []any {
	results := make([]any, 0, len(f)*2)
	for i := range f {
		results = append(results, f[i].Key, f[i].Value)
	}
	return results
}

const (
	KeyFunctionName = "functionName" // KeyFunctionName is the name of the function name key
	KeyTimestamp    = "timestamp"    // KeyTimestamp is the name of the timestamp key
	KeyServiceName  = "serviceName"  // KeyServiceName is the name of the service name key
	KeyMessage      = "message"      // KeyMessage is the name of the message key
	KeyLatency      = "latency"      // KeyLatency is the name of the latency name key
	KeyRequestID    = "requestID"    // KeyRequestID is the name of the request id key
	KeyMethod       = "method"       // KeyMethod is the name of the method key
	KeyPath         = "path"         // KeyPath is the name of the path key
	KeyIP           = "ip"           // KeyClientIP is the name of the client ip key
	KeyUserAgent    = "userAgent"    // KeyUserAgent is the name of the user agent key
	KeyRequestBody  = "requestBody"  // KeyRequestBody is the name of the request body key
	KeyResponseBody = "responseBody" // KeyResponseBody is the name of the response body key
	KeyErrors       = "errors"       // KeyErrors is the name of the errors key
	KeyResponseTime = "responseTime" // KeyResponseTime is the name of the response time key
	KeyStatusCode   = "statusCode"   // KeyStatusCode is the name of the status code key
)

// NewLatencyField creates a new latency field from the given start time in milliseconds
func NewLatencyField(startTime time.Time) Field {
	return Field{KeyLatency, time.Since(startTime).Milliseconds()}
}

// NewLatencyBaseFields creates a new base field with latency fields from the given start time in milliseconds
func NewLatencyBaseFields(startTime time.Time) []Field {
	return []Field{NewLatencyField(startTime)}
}
