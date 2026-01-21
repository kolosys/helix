# Performance Tuning

This guide covers performance optimization techniques for Helix applications.

## Built-in Performance Features

Helix is designed for high performance out of the box:

### Zero Allocations in Hot Paths

Helix uses `sync.Pool` for frequently created objects:

- **Path Parameters**: Reused for each request
- **Ctx Instances**: Pooled and reset between requests
- **JSON Encoding Buffers**: Pooled for response encoding

### Radix Tree Router

The router uses a radix tree (compressed trie) for O(k) route matching where k is the path length:

```
Root
├── users (static)
│   ├── / (GET handler)
│   └── {id} (param)
│       └── / (GET handler)
└── posts (static)
    └── {id} (param)
```

### Pre-compiled Middleware Chain

Call `Build()` after registering all routes to pre-compile the middleware chain:

```go
s := helix.New(nil)

// Register routes and middleware
s.Use(middleware.RequestID())
s.Use(middleware.Logger(middleware.LogFormatJSON))
s.Use(middleware.Recover())

s.GET("/users", helix.HandleCtx(listUsers))
s.POST("/users", helix.HandleCtx(createUser))

// Pre-compile for production
s.Build()

s.Start(":8080")
```

## Optimization Techniques

### Use HandleCtx for Most Handlers

`HandleCtx` is optimized for common use cases with minimal overhead:

```go
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id := c.Param("id")
    return c.OK(user)
}))
```

### Avoid Unnecessary Allocations

#### Use Typed Parameters

```go
// ✅ Better - uses typed accessor
id, err := c.ParamInt("id")

// ❌ Worse - parses twice
idStr := c.Param("id")
id, _ := strconv.Atoi(idStr)
```

#### Reuse Slices for Query Parameters

```go
// QuerySlice returns existing slice when possible
tags := c.QuerySlice("tags")
```

### Binding Cache

Struct reflection information is cached automatically:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// First call: reflects and caches struct info
// Subsequent calls: uses cached info
req, err := helix.Bind[CreateUserRequest](r)
```

### Skip Unnecessary Middleware

Use Skip functions to bypass middleware for specific routes:

```go
s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: middleware.LogFormatJSON,
    Skip: func(r *http.Request) bool {
        // Don't log health checks
        return r.URL.Path == "/health"
    },
}))
```

## Server Configuration

### Tune Timeouts

Configure appropriate timeouts for your use case:

```go
s := helix.New(&helix.Options{
    ReadTimeout:  30 * time.Second,  // Max time to read request
    WriteTimeout: 30 * time.Second,  // Max time to write response
    IdleTimeout:  120 * time.Second, // Keep-alive timeout
    GracePeriod:  30 * time.Second,  // Shutdown grace period
})
```

### Connection Limits

For high-traffic scenarios, consider connection limits at the load balancer level rather than in the application.

## Middleware Performance

### Compression Trade-offs

Compression reduces bandwidth but increases CPU usage:

```go
// Enable compression for large responses
s.Use(middleware.CompressWithConfig(middleware.CompressConfig{
    Level: gzip.DefaultCompression,
    MinSize: 1024, // Only compress responses > 1KB
}))
```

### Rate Limiting

Token bucket rate limiting is efficient but adds overhead:

```go
// Only enable where needed
api := s.Group("/api")
api.Use(middleware.RateLimit(100, 10))  // 100 req/sec, burst 10
```

## Benchmarking

### Use Go's Built-in Benchmarks

```go
func BenchmarkListUsers(b *testing.B) {
    s := helix.New(nil)
    s.GET("/users", helix.HandleCtx(listUsers))
    
    req := httptest.NewRequest(http.MethodGet, "/users", nil)
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        rec := httptest.NewRecorder()
        s.ServeHTTP(rec, req)
    }
}
```

Run benchmarks with:

```bash
go test -bench=. -benchmem -benchtime=5s
```

### Sample Output

```
BenchmarkListUsers-8    500000    2450 ns/op    256 B/op    4 allocs/op
```

### Compare Before/After

```bash
# Before changes
go test -bench=. -benchmem > before.txt

# After changes
go test -bench=. -benchmem > after.txt

# Compare
benchstat before.txt after.txt
```

## Profiling

### CPU Profiling

Use the profiling middleware (requires build tag):

```go
//go:build profiling

s.Use(middleware.Profiling("/debug/pprof"))
```

Or use `net/http/pprof` directly:

```go
import _ "net/http/pprof"

// Profiles available at /debug/pprof/
```

### Memory Profiling

```bash
# Collect heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof

# Analyze
go tool pprof heap.prof
```

### Tracing

```bash
# Collect trace
curl http://localhost:8080/debug/pprof/trace?seconds=5 > trace.out

# View trace
go tool trace trace.out
```

## Common Performance Issues

### Issue: High Memory Usage

**Cause**: Unbounded request body reading

**Solution**: Limit body size

```go
s.Use(middleware.BodyLimit(1024 * 1024))  // 1MB limit
```

### Issue: Slow Route Matching

**Cause**: Too many routes with overlapping patterns

**Solution**: Organize routes with groups

```go
// Better organization
api := s.Group("/api/v1")
api.Mount("/users", &UserModule{})
api.Mount("/posts", &PostModule{})
```

### Issue: Memory Leaks

**Cause**: Not closing response bodies in HTTP clients

**Solution**: Always close response bodies

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()
```

### Issue: Goroutine Leaks

**Cause**: Context not respected in background operations

**Solution**: Always check context cancellation

```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-resultChan:
    return result, nil
}
```

## Production Checklist

- [ ] Call `s.Build()` before starting the server
- [ ] Configure appropriate timeouts
- [ ] Enable compression for text responses
- [ ] Add rate limiting to public endpoints
- [ ] Use JSON logging format
- [ ] Skip logging for health check endpoints
- [ ] Profile under realistic load before deployment
- [ ] Monitor memory usage and goroutine counts

## Resources

- [Go Performance Best Practices](https://github.com/dgryski/go-perfbook)
- [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)
- [pprof Documentation](https://pkg.go.dev/runtime/pprof)
- [Best Practices Guide](./best-practices.md)

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
