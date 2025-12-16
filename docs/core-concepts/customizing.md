# Customizing Helix

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain how to customize and extend Helix.

## About Customization

Helix is designed to be highly customizable while maintaining sensible defaults. You can customize server configuration, error handling, middleware, and more through the functional options pattern and extension points.

## Core Concepts

### Functional Options Pattern

Helix uses the functional options pattern for configuration. This provides:

- **Type Safety**: Compile-time checking of configuration
- **Flexibility**: Options can be applied in any order
- **Extensibility**: Easy to add new options without breaking changes
- **Readability**: Clear, self-documenting configuration

### Server Options

All server customization happens through `Option` functions passed to `helix.New()` or `helix.Default()`:

```go
s := helix.New(
    helix.WithAddr(":3000"),
    helix.WithReadTimeout(30 * time.Second),
    helix.WithWriteTimeout(30 * time.Second),
    helix.WithIdleTimeout(120 * time.Second),
    helix.WithGracePeriod(30 * time.Second),
    helix.WithBasePath("/api/v1"),
    helix.WithTLS("cert.pem", "key.pem"),
    helix.WithErrorHandler(customErrorHandler),
    helix.HideBanner(),
)
```

## Architecture Overview

### Configuration Flow

1. **Server Creation**: Options are applied during `helix.New()` or `helix.Default()`
2. **Route Registration**: Routes and middleware are registered
3. **Build Phase**: `Build()` pre-compiles middleware chain (called automatically)
4. **Server Start**: Server begins accepting connections

### Extension Points

Helix provides several extension points:

- **Error Handlers**: Custom error handling logic
- **Middleware**: Request/response processing
- **Lifecycle Hooks**: Startup and shutdown callbacks
- **Custom Banners**: Server startup messages

## Core Concepts

### Server Configuration

#### Address and Timeouts

```go
s := helix.New(
    helix.WithAddr(":3000"),                    // Listen address
    helix.WithReadTimeout(30 * time.Second),    // Request read timeout
    helix.WithWriteTimeout(30 * time.Second),   // Response write timeout
    helix.WithIdleTimeout(120 * time.Second),  // Keep-alive timeout
    helix.WithGracePeriod(30 * time.Second),    // Shutdown grace period
    helix.WithMaxHeaderBytes(1 << 20),          // Max header size (1MB)
)
```

#### TLS Configuration

```go
// Simple TLS with certificate files
s := helix.New(
    helix.WithTLS("cert.pem", "key.pem"),
)

// Advanced TLS configuration
s := helix.New(
    helix.WithTLSConfig(&tls.Config{
        MinVersion: tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        },
    }),
)
```

#### Base Path

Set a prefix for all routes:

```go
s := helix.New(
    helix.WithBasePath("/api/v1"),
)

// Route "/users" becomes "/api/v1/users"
s.GET("/users", handler)
```

### Error Handling

#### Custom Error Handler

Replace the default RFC 7807 error handler:

```go
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func customErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
    // Custom error handling logic
    switch e := err.(type) {
    case helix.Problem:
        // Handle Problem errors
        helix.WriteProblem(w, e)
    case *helix.ValidationErrors:
        // Handle validation errors
        helix.WriteValidationProblem(w, e)
    default:
        // Handle other errors
        helix.InternalServerError(w, "internal server error")
    }
}

s := helix.New(
    helix.WithErrorHandler(customErrorHandler),
)
```

#### Error Handler Middleware

Error handlers are injected as middleware, so they have access to the full request context:

```go
func loggingErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
    // Log error with request context
    logs.ErrorContext(r.Context(), "handler error",
        logs.Err(err),
        logs.String("path", r.URL.Path),
        logs.String("method", r.Method),
    )

    // Use default handler for actual response
    helix.HandleError(w, r, err)
}
```

### Middleware Customization

#### Global Middleware

Add middleware that applies to all routes:

```go
s := helix.New()

// Add middleware
s.Use(middleware.RequestID())
s.Use(middleware.Logger(middleware.LogFormatJSON))
s.Use(middleware.Recover())
s.Use(customMiddleware)

// Middleware executes in order added
```

#### Route-Specific Middleware

Apply middleware to specific routes or groups:

```go
// Group with middleware
admin := s.Group("/admin", authMiddleware, adminOnlyMiddleware)
admin.GET("/stats", getStats)

// Resource with middleware
s.Resource("/users", authMiddleware).
    List(listUsers).
    Create(createUser)
```

#### Custom Middleware

Create custom middleware that works with any `http.Handler`:

```go
func requestLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap response writer to capture status
        rw := &responseWriter{ResponseWriter: w}

        next.ServeHTTP(rw, r)

        duration := time.Since(start)
        logs.Info("request completed",
            logs.String("method", r.Method),
            logs.String("path", r.URL.Path),
            logs.Int("status", rw.status),
            logs.Duration("duration", duration),
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### Lifecycle Hooks

#### Startup Hooks

Execute code when the server starts:

```go
s.OnStart(func(s *helix.Server) {
    logs.Info("server starting",
        logs.String("addr", s.Addr()),
        logs.String("version", helix.Version),
    )

    // Initialize services, connect to databases, etc.
    db.Connect()
    cache.Connect()
})
```

#### Shutdown Hooks

Execute code during graceful shutdown:

```go
s.OnStop(func(ctx context.Context, s *helix.Server) {
    logs.Info("server shutting down")

    // Cleanup resources
    db.Close()
    cache.Close()

    // Wait for in-flight requests
    // Context is already cancelled, so operations should respect it
    <-ctx.Done()
})
```

### Banner Customization

#### Hide Banner

```go
s := helix.New(helix.HideBanner())
```

#### Custom Banner

```go
s := helix.New(
    helix.WithCustomBanner(`
    ╔═══════════════════════════╗
    ║   My Awesome API v1.0.0   ║
    ╚═══════════════════════════╝
    `),
)
```

## Usage Patterns

### Development Server

```go
s := helix.Default(
    helix.WithAddr(":8080"),
    helix.HideBanner(),  // Cleaner output in dev
)

// Development middleware
s.Use(middleware.Logger(middleware.LogFormatDev))
s.Use(middleware.CORSAllowAll())  // Allow all origins in dev
```

### Production Server

```go
s := helix.New(
    helix.WithAddr(":8080"),
    helix.WithReadTimeout(30 * time.Second),
    helix.WithWriteTimeout(30 * time.Second),
    helix.WithIdleTimeout(120 * time.Second),
    helix.WithGracePeriod(30 * time.Second),
    helix.WithTLS("cert.pem", "key.pem"),
)

// Production middleware bundle
for _, mw := range middleware.Production() {
    s.Use(mw)
}

// Lifecycle hooks
s.OnStart(func(s *helix.Server) {
    // Health checks, metrics, etc.
})

s.OnStop(func(ctx context.Context, s *helix.Server) {
    // Graceful shutdown
})
```

### API Server with Base Path

```go
s := helix.New(
    helix.WithBasePath("/api/v1"),
    helix.WithErrorHandler(apiErrorHandler),
)

// API middleware bundle
for _, mw := range middleware.API() {
    s.Use(mw)
}

// Routes are automatically prefixed
s.GET("/users", listUsers)  // Becomes /api/v1/users
```

## Design Decisions

### Functional Options Pattern

Helix uses functional options instead of a config struct because:

- **Backward Compatibility**: New options don't break existing code
- **Type Safety**: Compile-time checking prevents invalid configurations
- **Flexibility**: Options can be conditionally applied
- **Readability**: Self-documenting configuration code

### Error Handler Injection

Error handlers are injected as middleware rather than stored in Server because:

- **Context Access**: Full access to request context
- **Middleware Chain**: Can be combined with other middleware
- **Flexibility**: Can be different per route group if needed

### Lifecycle Hooks

Lifecycle hooks are separate from server options because:

- **Timing**: Hooks execute at specific lifecycle events
- **Multiple Hooks**: Can register multiple hooks for the same event
- **Context**: Shutdown hooks receive cancellation context

## Common Pitfalls

### Pitfall 1: Options Applied After Routes

**Problem**: Some options (like `WithBasePath`) must be set before registering routes.

**Solution**: Set all options during server creation:

```go
// ❌ Wrong - base path set after routes
s := helix.New()
s.GET("/users", handler)
s.WithBasePath("/api")  // Too late!

// ✅ Correct - set options first
s := helix.New(helix.WithBasePath("/api"))
s.GET("/users", handler)  // Route is /api/users
```

### Pitfall 2: Error Handler Not Handling All Types

**Problem**: Custom error handlers must handle all error types or use a default case.

**Solution**: Always include a default case:

```go
func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
    switch e := err.(type) {
    case helix.Problem:
        helix.WriteProblem(w, e)
    case *helix.ValidationErrors:
        helix.WriteValidationProblem(w, e)
    default:
        // Always handle unknown errors
        helix.WriteProblem(w, helix.ErrInternal.WithDetailf("%v", err))
    }
}
```

### Pitfall 3: Not Respecting Context in Shutdown Hooks

**Problem**: Shutdown hooks receive a cancelled context. Operations must respect cancellation.

**Solution**: Always check context in shutdown hooks:

```go
s.OnStop(func(ctx context.Context, s *helix.Server) {
    // ❌ Wrong - doesn't respect cancellation
    db.Close()  // Blocks indefinitely

    // ✅ Correct - respects context
    done := make(chan struct{})
    go func() {
        db.Close()
        close(done)
    }()

    select {
    case <-done:
        // Closed successfully
    case <-ctx.Done():
        // Timeout - force close
        db.ForceClose()
    }
})
```

## Integration Guide

### With Dependency Injection

Customize service registration:

```go
s := helix.New()

// Register global services
helix.Register(userService)
helix.Register(emailService)

// Provide request-scoped services via middleware
s.Use(helix.ProvideMiddleware(func(r *http.Request) *Transaction {
    return db.BeginTx(r.Context())
}))
```

### With Logging

Integrate with Helix's logging package:

```go
import "github.com/kolosys/helix/logs"

s := helix.New()

// Configure logger
log := logs.New(
    logs.WithLevel(logs.InfoLevel),
    logs.WithFormatter(&logs.JSONFormatter{}),
    logs.WithCaller(),
)

// Use in lifecycle hooks
s.OnStart(func(s *helix.Server) {
    log.Info("server starting", logs.String("addr", s.Addr()))
})
```

### With Metrics

Add metrics middleware:

```go
func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        rw := &responseWriter{ResponseWriter: w}
        next.ServeHTTP(rw, r)

        duration := time.Since(start)
        metrics.RecordRequest(r.Method, r.URL.Path, rw.status, duration)
    })
}

s.Use(metricsMiddleware)
```

## Further Reading

- [Error Handling](./error-handling.md) - Error handling patterns
- [Middleware Guide](../middleware/middleware.md) - Built-in middleware
- [Best Practices](../advanced/best-practices.md) - Recommended patterns
- [API Reference](../api-reference/helix.md) - Complete options API

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
