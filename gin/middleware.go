package gin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pondparinya/httplog"
)

type bufferWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bufferWriter) Write(data []byte) (int, error) {
	b.body.Write(data)
	return b.ResponseWriter.Write(data)
}

const maxBodyLogSize int = 10000

func Logger(l httplog.Logger, skipPaths []string) func(*gin.Context) {
	if l == nil {
		panic("logger is required for Logger middleware")
	}
	skipPathMap := make(map[string]struct{}, len(skipPaths))
	for _, path := range skipPaths {
		skipPathMap[path] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := skipPathMap[c.Request.URL.Path]; !ok {
			start := time.Now()

			var copyBody bytes.Buffer
			c.Request.Body = io.NopCloser(io.TeeReader(c.Request.Body, &copyBody))

			c.Request = c.Request.WithContext(httplog.NewRequestContext(c.Request.Context(), httplog.NewRequestData(c)))

			l.Infocf(c.Request.Context(), "started")

			buff := bytes.Buffer{}
			bw := bufferWriter{c.Writer, &buff}
			c.Writer = bw

			defer func() {
				contentType := c.ContentType()
				reqBody := formatRequestBody(contentType, copyBody.Bytes())

				completedReqFields := make([]httplog.Field, 0)
				completedReqFields = append(completedReqFields,
					httplog.Field{Key: httplog.KeyRequestBody, Value: reqBody},
					httplog.Field{Key: httplog.KeyStatusCode, Value: c.Writer.Status()},
					httplog.Field{Key: httplog.KeyResponseTime, Value: time.Since(start).Milliseconds()},
				)

				if bw.body.Len() > 0 {
					respBody := bw.body.String()
					if len(respBody) > maxBodyLogSize {
						respBody = respBody[:maxBodyLogSize] + "..."
					}
					completedReqFields = append(completedReqFields, httplog.Field{Key: httplog.KeyResponseBody, Value: respBody})
				}

				if len(c.Errors) > 0 {
					completedReqFields = append(completedReqFields, httplog.Field{Key: httplog.KeyErrors, Value: c.Errors.Errors()})
					l.Errorcw(c.Request.Context(), completedReqFields, "completed")
				} else {
					l.Infocw(c.Request.Context(), completedReqFields, "completed")
				}
			}()

			c.Next()
		}
	}
}

func formatRequestBody(contentType string, raw []byte) string {
	if len(raw) == 0 {
		return ""
	}

	switch {
	case bytes.Contains([]byte(contentType), []byte("application/json")):
		var compacted bytes.Buffer
		if err := json.Compact(&compacted, raw); err == nil {
			return compacted.String()
		}
	case bytes.Contains([]byte(contentType), []byte("application/xml")),
		bytes.Contains([]byte(contentType), []byte("text/xml")):
		if xml.Unmarshal(raw, &struct {
			XMLName xml.Name
			Rest    []byte `xml:",innerxml"`
		}{}) == nil {
			return string(raw)
		}
	}

	return string(raw)
}

func Recovery(l httplog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Errorcw(c.Request.Context(), []httplog.Field{
					{Key: httplog.KeyErrors, Value: err},
				}, "panic recovered")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{httplog.KeyMessage: "internal server error"})
			}
		}()
		c.Next()
	}
}
