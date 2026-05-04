# go-logger

## Usage

```
go get go.kbtg.tech/733/go-logger
```

## Code Example

```go
// Initialize zap logger config
logCfg := logger.ZapConfig{
 Level:       "debug",
 OutputPaths: []string{"stdout"},
}

// Initialize zap logger instance
l, err := logger.NewZapLogger(logCfg)
if err != nil {
 panic(err)
}

// Set Global Logger with zap logger instance
logger.SetGlobalLogger(l)

// Info log without fields
l.Infof(nil, "hello world")

// Info log with fields
fields := []logger.Field{
    {Key: "key", Value: "value"},
}
l.Infof(fields, "hello world")

// Info log without field, but string format
l.Infof(nil, "log with format: %v", "custom value")

// Info log with latency field
l.Infof(logger.NewLatencyBaseFields(time.Now()), "log latency")

// Info log with context
ctx := context.Background()
l.Infocf(ctx, nil, "log with context")

// Info log with request context's go-logger
requestData := logger.RequestData{
 RequestID: "XXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
 Method:    "GET",
 Path:      "/status",
 IP:        "127.0.0.1",
 UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X",
}
requestCtx := logger.NewRequestContext(ctx, requestData)
l.Infocf(requestCtx, nil, "log with request context")

```

## Output

```
{"level":"info","timestamp":"2024-07-11T17:20:03.787+0700","caller":"server/main.go:22","functionName":"main.main","message":"log without fields"}
{"level":"info","timestamp":"2024-07-11T17:20:03.788+0700","caller":"server/main.go:27","functionName":"main.main","message":"log with fields","key":"value"}
{"level":"info","timestamp":"2024-07-11T17:20:03.788+0700","caller":"server/main.go:29","functionName":"main.main","message":"log with format: custom value"}
{"level":"info","timestamp":"2024-07-11T17:20:03.788+0700","caller":"server/main.go:31","functionName":"main.main","message":"log latency","latency":0}
{"level":"info","timestamp":"2024-07-11T17:20:03.788+0700","caller":"server/main.go:34","functionName":"main.main","message":"log with context"}
{"level":"info","timestamp":"2024-07-11T17:20:03.788+0700","caller":"server/main.go:44","functionName":"main.main","message":"log with request context","refID":"XXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX","method":"GET","path":"/status","ip":"127.0.0.1","user-agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X"}
```
