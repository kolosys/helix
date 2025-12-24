# RequestID Middleware

Generates or propagates a unique request ID for each request. Useful for request tracing and logging.

## Basic Usage

```go
s.Use(middleware.RequestID())
```

## Configuration

```go
s.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
    Header:       "X-Request-ID",
    Generator:    func() string { return customID() },
    TargetHeader: "X-Request-ID",
}))
```

## Features

- Reads existing request ID from header (if present)
- Generates random 16-byte hex string if not present
- Stores request ID in request context
- Sets request ID in response header

## Accessing Request ID

```go
requestID := middleware.GetRequestIDFromRequest(r)
// or
requestID := middleware.GetRequestID(r.Context())
```

## Configuration Options

- **Header**: Header name to read/write request ID (default: `X-Request-ID`)
- **Generator**: Function that generates a new request ID (default: random 16-byte hex)
- **TargetHeader**: Header name to set on response (default: same as Header)

## Example

```go
s := helix.New(nil)

// Use default RequestID middleware
s.Use(middleware.RequestID())

// Access in handler
s.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
    requestID := middleware.GetRequestIDFromRequest(c.Request)
    logs.Info("request", logs.String("request_id", requestID))
    return c.OK(map[string]string{"request_id": requestID})
}))
```

## Custom Generator

```go
s.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
    Generator: func() string {
        // Generate UUID or custom ID format
        return generateUUID()
    },
}))
```
