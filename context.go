package httplog

import (
	"context"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	// ContextKey is a type of context key
	ContextKey string

	// RequestData is a type of request data that will be set inside the context, should be set from request handler
	RequestData struct {
		RequestID string
		Method    string
		Path      string
		IP        string
		UserAgent string
	}
)

const (
	// KeyRequestData is a request data context key which is used to set or get the request data from context
	KeyRequestData ContextKey = "_go-logger/request_data"
)

// NewRequestContext returns a new Context that carries request data
func NewRequestContext(ctx context.Context, request RequestData) context.Context {
	return context.WithValue(ctx, KeyRequestData, request)
}

// NewRequestData returns request data from request value.
func NewRequestData(v any) RequestData {
	if r, ok := v.(*http.Request); ok {
		return RequestData{
			RequestID: requestIDOrNew(r.Header.Get("X-Request-ID")),
			Method:    r.Method,
			Path:      r.URL.Path,
			IP:        r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}
	}

	if gc, ok := v.(*gin.Context); ok {
		requestid.New()(gc)
		return RequestData{
			RequestID: requestIDOrNew(requestid.Get(gc)),
			Method:    gc.Request.Method,
			Path:      gc.Request.URL.Path,
			IP:        gc.ClientIP(),
			UserAgent: gc.Request.UserAgent(),
		}
	}

	if ctx, ok := v.(context.Context); ok {
		return getRequestDataFromContext(ctx)
	}

	return RequestData{RequestID: uuid.New().String()}
}

func requestIDOrNew(id string) string {
	if id == "" {
		return uuid.New().String()
	}
	return id
}

// getRequestDataFromContext returns request data from context.
func getRequestDataFromContext(ctx context.Context) RequestData {
	if ctx == nil {
		return RequestData{}
	}
	v, ok := ctx.Value(KeyRequestData).(RequestData)
	if !ok {
		return RequestData{}
	}
	return v
}

// GetRequestIDFromContext returns request id from context.
func GetRequestIDFromContext(ctx context.Context) (string, error) {
	v := getRequestDataFromContext(ctx)
	return v.RequestID, nil
}

// ParseRequestContext returns request data from context as slice of any.
func ParseRequestContext(ctx context.Context) []any {
	v := getRequestDataFromContext(ctx)
	result := make([]any, 0)
	if v.RequestID != "" {
		result = append(result, KeyRequestID, v.RequestID)
	}
	if v.Method != "" {
		result = append(result, KeyMethod, v.Method)
	}
	if v.Path != "" {
		result = append(result, KeyPath, v.Path)
	}
	if v.IP != "" {
		result = append(result, KeyIP, v.IP)
	}
	if v.UserAgent != "" {
		result = append(result, KeyUserAgent, v.UserAgent)
	}
	return result
}
