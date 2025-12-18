# Middleware

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain the core concepts and architecture of the middleware package.

## About This Package

**Import Path:** `github.com/kolosys/helix/middleware`

Package middleware provides HTTP middleware for the Helix framework. Middleware functions wrap HTTP handlers to add cross-cutting concerns like logging, authentication, compression, and more.

## Architecture Overview

### Middleware Pattern

Middleware in Helix follows the standard Go HTTP middleware pattern:

```go
type Middleware func(next http.Handler) http.Handler
```

Middleware wraps an `http.Handler` and returns a new `http.Handler` that can execute code before and/or after calling the next handler.

### Execution Order

Middleware executes in the order it's added:

```go
s.Use(middleware.RequestID())    // Executes first
s.Use(middleware.Logger(...))    // Executes second
s.Use(middleware.Recover())     // Executes third (innermost)
```

The first middleware added is the outermost (executes first on request, last on response).

### Compatibility

All Helix middleware is compatible with the standard `net/http` package and can be used with any `http.Handler`.

## Core Concepts

### Using Middleware

Add middleware to your server:

```go
s := helix.New()

// Add individual middleware
s.Use(middleware.RequestID())
s.Use(middleware.Logger(middleware.LogFormatJSON))
s.Use(middleware.Recover())

// Or use middleware bundles
for _, mw := range middleware.API() {
    s.Use(mw)
}
```

### Middleware Chain

Create a reusable middleware chain:

```go
chain := middleware.Chain(
    middleware.RequestID(),
    middleware.Logger(middleware.LogFormatDev),
    middleware.Recover(),
)

s.Use(chain)
```

## Available Middleware

- **[RequestID](requestid.md)** - Generates unique request IDs for tracing
- **[Logger](logger.md)** - HTTP request logging with multiple formats
- **[Recover](recover.md)** - Panic recovery middleware
- **[CORS](cors.md)** - Cross-Origin Resource Sharing
- **[RateLimit](ratelimit.md)** - Request rate limiting
- **[BasicAuth](basicauth.md)** - HTTP Basic Authentication
- **[Compress](compress.md)** - Response compression (gzip/deflate)
- **[Timeout](timeout.md)** - Request timeout handling
- **[ETag](etag.md)** - ETag generation and conditional requests
- **[Cache](cache.md)** - HTTP cache headers
- **[Profiling](profiling.md)** - Performance profiling (build tag)

## Middleware Bundles

Pre-configured middleware sets for common scenarios:

### API Bundle

Suitable for JSON API servers:

```go
for _, mw := range middleware.API() {
    s.Use(mw)
}
```

Includes: RequestID, Logger (JSON), Recover, CORS

### API with Custom CORS

```go
corsConfig := middleware.CORSConfig{
    AllowOrigins: []string{"https://api.example.com"},
    AllowCredentials: true,
}

for _, mw := range middleware.APIWithCORS(corsConfig) {
    s.Use(mw)
}
```

### Web Bundle

Suitable for web applications:

```go
for _, mw := range middleware.Web() {
    s.Use(mw)
}
```

Includes: RequestID, Logger (dev), Recover, Compress

### Production Bundle

Suitable for production environments:

```go
for _, mw := range middleware.Production() {
    s.Use(mw)
}
```

Includes: RequestID, Logger (combined), Recover

### Development Bundle

Suitable for development (same as `helix.Default()`):

```go
for _, mw := range middleware.Development() {
    s.Use(mw)
}
```

Includes: RequestID, Logger (dev), Recover

### Secure Bundle

Security-focused middleware:

```go
for _, mw := range middleware.Secure(100, 10) {
    s.Use(mw)
}
```

Includes: RequestID, Logger (JSON), Recover, RateLimit

### Minimal Bundle

Only essential middleware:

```go
for _, mw := range middleware.Minimal() {
    s.Use(mw)
}
```

Includes: Recover

## Usage Patterns

### Conditional Middleware

Apply middleware conditionally:

```go
if os.Getenv("ENV") == "production" {
    s.Use(middleware.Logger(middleware.LogFormatJSON))
} else {
    s.Use(middleware.Logger(middleware.LogFormatDev))
}
```

### Route-Specific Middleware

Apply middleware to specific routes:

```go
// Global middleware
s.Use(middleware.RequestID())

// Route-specific middleware
api := s.Group("/api", middleware.Logger(middleware.LogFormatJSON))
admin := s.Group("/admin", middleware.BasicAuth("admin", "secret"))
```

### Custom Middleware

Create custom middleware:

```go
func customMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Before handler
        start := time.Now()

        // Call next handler
        next.ServeHTTP(w, r)

        // After handler
        duration := time.Since(start)
        logs.Info("request completed", logs.Duration("duration", duration))
    })
}

s.Use(customMiddleware)
```

## Design Decisions

### Standard Library Compatibility

All middleware uses `net/http` types because:

- **Compatibility**: Works with any HTTP framework
- **Simplicity**: No framework-specific types
- **Flexibility**: Easy to integrate with existing code

### Functional Options Pattern

Middleware uses configuration structs because:

- **Type Safety**: Compile-time checking
- **Extensibility**: Easy to add new options
- **Clarity**: Self-documenting configuration

### Performance Considerations

- **Object Pooling**: Compress middleware pools gzip/deflate writers
- **Pre-computation**: CORS headers pre-computed for performance
- **Lazy Evaluation**: Skip functions prevent unnecessary work

## Common Pitfalls

### Pitfall 1: Middleware Order Matters

**Problem**: Middleware executes in the order added. Recover should be innermost.

**Solution**: Add Recover last (or use bundles):

```go
// ❌ Wrong - Recover won't catch panics in Logger
s.Use(middleware.Recover())
s.Use(middleware.Logger(...))

// ✅ Correct - Recover catches all panics
s.Use(middleware.Logger(...))
s.Use(middleware.Recover())
```

### Pitfall 2: Reading Request Body Multiple Times

**Problem**: HTTP request bodies can only be read once. Logger middleware with body capture reads the body.

**Solution**: Only enable body capture when needed:

```go
// ❌ Wrong - body read twice
s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    CaptureBody: true, // Reads body
}))
// Handler tries to read body again - fails!

// ✅ Correct - don't capture body, or read manually
s.Use(middleware.Logger(middleware.LogFormatJSON)) // No body capture
```

### Pitfall 3: CORS with Credentials

**Problem**: CORS with credentials requires specific origin (can't use `*`).

**Solution**: Specify exact origins:

```go
// ❌ Wrong - credentials with wildcard
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"*"},
    AllowCredentials: true, // Won't work!
}))

// ✅ Correct - specific origins
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"https://app.example.com"},
    AllowCredentials: true,
}))
```

## Integration Guide

### With Helix Framework

Middleware integrates seamlessly with Helix:

```go
s := helix.Default() // Already includes RequestID, Logger, Recover

// Add more middleware
s.Use(middleware.CORS())
s.Use(middleware.Compress())
```

### With Standard Library

Use Helix middleware with standard library:

```go
mux := http.NewServeMux()
mux.HandleFunc("/", handler)

handler := middleware.RequestID()(
    middleware.Logger(middleware.LogFormatJSON)(
        middleware.Recover()(mux),
    ),
)

http.ListenAndServe(":8080", handler)
```

### With Logging

Logger middleware works with Helix logs package:

```go
s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: middleware.LogFormatJSON,
    Output: logFile,
    Fields: map[string]string{
        "service": "api",
        "version": "1.0.0",
    },
}))
```

## Further Reading

- [RequestID Middleware](requestid.md) - Request ID generation
- [Logger Middleware](logger.md) - HTTP request logging
- [Recover Middleware](recover.md) - Panic recovery
- [CORS Middleware](cors.md) - Cross-Origin Resource Sharing
- [RateLimit Middleware](ratelimit.md) - Request rate limiting
- [BasicAuth Middleware](basicauth.md) - HTTP Basic Authentication
- [Compress Middleware](compress.md) - Response compression
- [Timeout Middleware](timeout.md) - Request timeout
- [ETag Middleware](etag.md) - ETag generation
- [Cache Middleware](cache.md) - Cache headers
- [Profiling Middleware](profiling.md) - Performance profiling
- [API Reference](../api-reference/middleware.md) - Complete middleware API documentation
- [Core Concepts](../core-concepts/) - Framework concepts
- [Examples](../examples/middlewares/main.md) - Middleware usage examples
- [Best Practices](../advanced/best-practices.md) - Recommended patterns

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
