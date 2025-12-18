# RateLimit Middleware

Limits the number of requests per client using a token bucket algorithm.

## Basic Usage

```go
// 100 requests per second, burst of 10
s.Use(middleware.RateLimit(100, 10))
```

## Configuration

```go
s.Use(middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Rate:     100,  // requests per second
    Burst:    10,   // maximum burst
    KeyFunc:  func(r *http.Request) string {
        // Rate limit by API key instead of IP
        return r.Header.Get("X-API-Key")
    },
    Handler: func(w http.ResponseWriter, r *http.Request) {
        // Custom rate limit response
        helix.TooManyRequests(w, "Rate limit exceeded")
    },
    SkipFunc: func(r *http.Request) bool {
        // Skip rate limiting for health checks
        return r.URL.Path == "/health"
    },
    CleanupInterval: time.Minute,
    ExpirationTime:  5 * time.Minute,
}))
```

## Features

- Token bucket algorithm
- Per-client rate limiting (by IP or custom key)
- Automatic cleanup of expired entries
- Rate limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `Retry-After`)

## Token Bucket Algorithm

The token bucket algorithm allows:

- **Rate**: Steady rate of requests per second
- **Burst**: Maximum number of requests that can be made in quick succession

Example: Rate=100, Burst=10 means:

- Up to 100 requests per second sustained
- Up to 10 requests can be made immediately (burst)
- After burst is exhausted, requests are limited to the rate

## Rate Limiting by IP

Default behavior limits by client IP:

```go
s.Use(middleware.RateLimit(100, 10))
```

The middleware extracts the client IP from:

1. `X-Forwarded-For` header (first IP)
2. `X-Real-IP` header
3. `RemoteAddr` from request

## Rate Limiting by API Key

Limit by API key instead of IP:

```go
s.Use(middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Rate: 100,
    Burst: 10,
    KeyFunc: func(r *http.Request) string {
        return r.Header.Get("X-API-Key")
    },
}))
```

## Rate Limiting by User

Limit by authenticated user:

```go
s.Use(middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Rate: 100,
    Burst: 10,
    KeyFunc: func(r *http.Request) string {
        // Extract user ID from context or header
        userID := getUserID(r)
        return userID
    },
}))
```

## Custom Rate Limit Response

```go
s.Use(middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Rate: 100,
    Burst: 10,
    Handler: func(w http.ResponseWriter, r *http.Request) {
        helix.WriteProblem(w, helix.ErrTooManyRequests.WithDetail(
            "Rate limit exceeded. Please try again later.",
        ))
    },
}))
```

## Skipping Rate Limits

Skip rate limiting for specific requests:

```go
s.Use(middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Rate: 100,
    Burst: 10,
    SkipFunc: func(r *http.Request) bool {
        // Don't rate limit health checks or metrics
        return r.URL.Path == "/health" || r.URL.Path == "/metrics"
    },
}))
```

## Rate Limit Headers

The middleware sets these headers:

- `X-RateLimit-Limit`: Maximum requests per second
- `X-RateLimit-Remaining`: Remaining requests in current window
- `Retry-After`: Seconds to wait before retrying (when limit exceeded)

## Cleanup

The middleware automatically cleans up expired rate limit entries:

```go
s.Use(middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Rate:            100,
    Burst:           10,
    CleanupInterval: time.Minute,  // How often to run cleanup
    ExpirationTime:  5 * time.Minute, // When to expire entries
}))
```

## Example

```go
s := helix.New()

// Global rate limit: 100 req/s, burst of 10
s.Use(middleware.RateLimit(100, 10))

// Stricter rate limit for API routes
api := s.Group("/api", middleware.RateLimit(50, 5))

// Even stricter for admin routes
admin := api.Group("/admin", middleware.RateLimit(10, 2))
```
