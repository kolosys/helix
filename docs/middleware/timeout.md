# Timeout Middleware

Adds a timeout to request processing. Cancels the request context if the timeout is exceeded.

## Basic Usage

```go
// 30 second timeout
s.Use(middleware.Timeout(30 * time.Second))
```

## Configuration

```go
s.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
    Timeout: 30 * time.Second,
    Handler: func(w http.ResponseWriter, r *http.Request) {
        // Custom timeout response
        helix.ServiceUnavailable(w, "Request timeout")
    },
    SkipFunc: func(r *http.Request) bool {
        // Don't timeout long-running operations
        return r.URL.Path == "/long-task"
    },
}))
```

## Features

- Context-based timeout
- Automatic cancellation
- Custom timeout handler
- Prevents hanging requests

## How It Works

The middleware:

1. Creates a context with timeout
2. Replaces the request context
3. Processes the request in a goroutine
4. If timeout occurs before completion, sends timeout response
5. If headers already written, cannot send timeout response

## Custom Timeout Response

```go
s.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
    Timeout: 30 * time.Second,
    Handler: func(w http.ResponseWriter, r *http.Request) {
        helix.WriteProblem(w, helix.ErrGatewayTimeout.WithDetail(
            "Request processing exceeded timeout",
        ))
    },
}))
```

## Skipping Timeout

Skip timeout for long-running operations:

```go
s.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
    Timeout: 30 * time.Second,
    SkipFunc: func(r *http.Request) bool {
        // Don't timeout file uploads or long tasks
        return r.URL.Path == "/upload" || r.URL.Path == "/long-task"
    },
}))
```

## Context Cancellation

Handlers should respect context cancellation:

```go
s.GET("/process", helix.HandleCtx(func(c *helix.Ctx) error {
    ctx := c.Context()

    // Check if context is cancelled
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Do work that respects context
    result, err := longRunningTask(ctx)
    if err != nil {
        return err
    }

    return c.OK(result)
}))
```

## Important Notes

- Timeout applies to the entire request processing time
- If response headers are already written, timeout response cannot be sent
- Handlers should check `context.Context` for cancellation
- Use `SkipFunc` for endpoints that legitimately take longer

## Example

```go
s := helix.New()

// Global timeout: 30 seconds
s.Use(middleware.Timeout(30 * time.Second))

// Shorter timeout for API routes
api := s.Group("/api", middleware.Timeout(10 * time.Second))

// No timeout for file uploads
s.POST("/upload", uploadHandler) // Not in timeout group
```
