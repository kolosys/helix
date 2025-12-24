# Structured Logging

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain the core concepts and architecture of the logs package.

## About This Package

**Import Path:** `github.com/kolosys/helix/logs`

Package logs provides a high-performance, context-aware structured logging library with zero-allocation hot paths, multiple output formats, and extensible hook system.

## Architecture Overview

The logs package is designed for high-performance logging in production environments:

### Core Components

1. **Logger**: Main logging interface with configurable options
2. **Entry**: Log entry with message, level, fields, and metadata
3. **Field**: Type-safe structured field builders
4. **Formatter**: Output formatting (text, JSON, pretty)
5. **Hook**: Extensibility point for custom processing
6. **Sampler**: Rate limiting for high-volume logs

### Data Flow

```text
Logger.Log()
    ↓
Level Check (skip if below threshold)
    ↓
Sampler Check (skip if sampled out)
    ↓
Create Entry (from pool)
    ↓
Add Fields (default + call-site)
    ↓
Run Hooks (pre-processing)
    ↓
Format Entry (TextFormatter/JSONFormatter)
    ↓
Write Output (sync or async)
    ↓
Return Entry to Pool
```

### Design Patterns

- **Object Pooling**: Entries and buffers use `sync.Pool` for zero allocations
- **Type-Safe Fields**: Compile-time type checking for field values
- **Formatter Interface**: Pluggable output formats
- **Hook System**: Extensible processing pipeline
- **Context Integration**: Automatic field extraction from `context.Context`

## Core Concepts

### Logger Creation

Create a logger with default settings:

```go
log := logs.New(nil)
log.Info("server started", logs.Int("port", 8080))
```

Create a logger with options:

```go
log := logs.New(&logs.Options{
    Level:           logs.DebugLevel,
    Formatter:       &logs.JSONFormatter{},
    Output:          os.Stdout,
    AddCaller:       true,
    AsyncBufferSize: 1024,
})
```

### Log Levels

Log levels in order of severity:

- `TraceLevel` - Most verbose, for detailed debugging
- `DebugLevel` - Debug information
- `InfoLevel` - General information (default)
- `WarnLevel` - Warning messages
- `ErrorLevel` - Error conditions
- `FatalLevel` - Fatal errors (exits after logging)
- `PanicLevel` - Panic errors (panics after logging)

### Structured Fields

Fields provide type-safe, structured data:

```go
logs.Info("user created",
    logs.String("user_id", "123"),
    logs.Int("age", 30),
    logs.Bool("active", true),
    logs.Duration("latency", time.Since(start)),
    logs.Err(err),
)
```

### Field Types

The package provides type-safe field builders:

- **Strings**: `String()`, `Strings()` (slice)
- **Integers**: `Int()`, `Int8()`, `Int16()`, `Int32()`, `Int64()`
- **Unsigned**: `Uint()`, `Uint8()`, `Uint16()`, `Uint32()`, `Uint64()`
- **Floats**: `Float32()`, `Float64()`
- **Booleans**: `Bool()`
- **Time**: `Time()`, `Duration()`
- **Errors**: `Err()`, `NamedErr()`
- **Any**: `Any()`, `JSON()`, `Bytes()`
- **Stringer**: `Stringer()` for types implementing `fmt.Stringer`

### Context-Aware Logging

Log with context to automatically include context fields:

```go
log.InfoContext(ctx, "request processed",
    logs.Duration("latency", time.Since(start)),
)

// Context fields are automatically extracted if set via logs.WithContext()
```

### Child Loggers

Create child loggers with additional fields:

```go
reqLog := log.With(
    logs.String("request_id", requestID),
    logs.String("user_id", userID),
)

reqLog.Info("processing request")  // Includes request_id and user_id
reqLog.Error("request failed", logs.Err(err))  // Also includes request_id and user_id
```

## Usage Patterns

### Basic Logging

```go
import "github.com/kolosys/helix/logs"

log := logs.New(nil)

log.Info("server started", logs.Int("port", 8080))
log.Warn("deprecated API used", logs.String("endpoint", "/old"))
log.Error("request failed", logs.Err(err))
```

### JSON Logging for Production

```go
log := logs.New(&logs.Options{
    Level:     logs.InfoLevel,
    Formatter: &logs.JSONFormatter{},
    Output:    os.Stdout,
})

log.Info("user created",
    logs.String("user_id", "123"),
    logs.String("email", "user@example.com"),
)
// Output: {"level":"info","time":"2024-01-15T10:30:00Z","message":"user created","user_id":"123","email":"user@example.com"}
```

### Development Logging

```go
log := logs.New(&logs.Options{
    Level: logs.DebugLevel,
    Formatter: &logs.TextFormatter{
        DisableColors: false,
        FullTimestamp: true,
    },
    AddCaller: true,
})

log.Debug("processing request",
    logs.String("method", "GET"),
    logs.String("path", "/users"),
)
```

### Context-Aware Logging

```go
// Set context fields (typically in middleware)
ctx := logs.WithContext(r.Context(),
    logs.String("request_id", requestID),
    logs.String("user_id", userID),
)

// Log with context (fields automatically included)
log.InfoContext(ctx, "request processed",
    logs.Duration("latency", time.Since(start)),
)
```

### Async Logging

Use async logging for high-throughput scenarios:

```go
log := logs.New(&logs.Options{
    AsyncBufferSize: 1024, // Buffer size
    Formatter:       &logs.JSONFormatter{},
})

// Logs are written asynchronously
log.Info("high volume log", logs.Int("count", 1000))
```

### Sampling

Use sampling to reduce log volume:

```go
sampler := logs.NewRateSampler(100, time.Second) // 100 logs per second

log := logs.New(&logs.Options{
    Sampler:   sampler,
    Formatter: &logs.JSONFormatter{},
})

// Only 100 logs per second will be written
for i := 0; i < 10000; i++ {
    log.Info("high volume", logs.Int("i", i))
}
```

### Hooks

Add hooks for custom processing (metrics, alerting, etc.):

```go
type metricsHook struct{}

func (h *metricsHook) Levels() []logs.Level {
    return []logs.Level{logs.ErrorLevel, logs.FatalLevel}
}

func (h *metricsHook) Fire(entry *logs.Entry) error {
    metrics.IncrementErrorCounter(entry.Level.String())
    return nil
}

log := logs.New(&logs.Options{
    Hooks: []logs.Hook{&metricsHook{}},
})
```

## Design Decisions

### Zero-Allocation Hot Paths

The logs package uses several techniques to avoid allocations:

- **Object Pooling**: Entries and buffers use `sync.Pool`
- **Type-Safe Fields**: Fields store values directly (no boxing)
- **Small Integer Optimization**: Common integers (< 100) use pre-allocated strings
- **Buffer Reuse**: Formatters reuse buffers from pool

### Type-Safe Fields

Fields are type-safe rather than using `map[string]interface{}` because:

- **Performance**: No reflection or type assertions in hot path
- **Type Safety**: Compile-time checking prevents errors
- **Memory**: Direct value storage (no boxing)
- **Clarity**: Explicit types improve readability

### Formatter Interface

Formatters are pluggable because:

- **Flexibility**: Easy to add custom formats
- **Performance**: Can optimize for specific formats
- **Testing**: Easy to test formatting logic
- **Extensibility**: Third-party formatters possible

### Context Integration

Context-aware logging provides:

- **Automatic Fields**: Request-scoped data automatically included
- **Propagation**: Fields propagate through call chain
- **Middleware Integration**: Easy to add request context in middleware
- **Zero Boilerplate**: No need to manually pass fields

## Common Pitfalls

### Pitfall 1: Not Closing Async Logger

**Problem**: Async loggers must be closed to flush pending logs.

**Solution**: Always close async loggers:

```go
log := logs.New(&logs.Options{AsyncBufferSize: 1024})
defer log.Close() // Flushes pending logs

log.Info("message")
// Logger must be closed to ensure message is written
```

### Pitfall 2: Logging After Close

**Problem**: Logging after `Close()` may panic or lose logs.

**Solution**: Ensure logger lifecycle matches application:

```go
log := logs.New(&logs.Options{AsyncBufferSize: 1024})

// In shutdown handler
s.OnStop(func(ctx context.Context, s *helix.Server) {
    log.Close() // Flush logs before shutdown
    // Don't log after this point
})
```

### Pitfall 3: High Allocation Fields

**Problem**: Creating fields with `Any()` for simple types causes allocations.

**Solution**: Use type-specific field builders:

```go
// ❌ Wrong - causes allocation
log.Info("message", logs.Any("count", 42))

// ✅ Correct - zero allocation
log.Info("message", logs.Int("count", 42))
```

### Pitfall 4: Not Setting Log Level

**Problem**: Default level is `InfoLevel`, so debug logs are ignored.

**Solution**: Set appropriate level for environment:

```go
// Development
log := logs.New(&logs.Options{Level: logs.DebugLevel})

// Production
log := logs.New(&logs.Options{Level: logs.InfoLevel})
```

## Integration Guide

### With Helix Framework

Helix middleware automatically uses the logs package:

```go
import "github.com/kolosys/helix/middleware"

s := helix.New(nil)
s.Use(middleware.Logger(middleware.LogFormatJSON))
```

### With Context

Add request-scoped fields in middleware:

```go
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := generateRequestID()
        ctx := logs.WithContext(r.Context(),
            logs.String("request_id", requestID),
        )
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### With Metrics

Use hooks to send logs to metrics systems:

```go
type prometheusHook struct {
    counter *prometheus.CounterVec
}

func (h *prometheusHook) Levels() []logs.Level {
    return []logs.Level{logs.ErrorLevel, logs.FatalLevel}
}

func (h *prometheusHook) Fire(entry *logs.Entry) error {
    h.counter.WithLabelValues(entry.Level.String()).Inc()
    return nil
}

log := logs.New(&logs.Options{
    Hooks: []logs.Hook{&prometheusHook{counter: errorCounter}},
})
```

### With External Services

Send logs to external services via hooks:

```go
type cloudWatchHook struct {
    client *cloudwatchlogs.Client
}

func (h *cloudWatchHook) Levels() []logs.Level {
    return []logs.Level{}  // All levels
}

func (h *cloudWatchHook) Fire(entry *logs.Entry) error {
    // Send to CloudWatch
    return h.client.PutLogEvents(...)
}
```

## Further Reading

- [API Reference](../api-reference/logs.md) - Complete logs API documentation
- [Middleware Guide](../middleware/middleware.md) - Logger middleware usage
- [Best Practices](../advanced/best-practices.md) - Logging best practices
- [Examples](../examples/README.md) - Practical logging examples

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
