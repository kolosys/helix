# Cache Middleware

Sets HTTP cache headers (`Cache-Control`, `Expires`, `Vary`) for responses.

## Basic Usage

```go
// Simple cache with max-age
s.Use(middleware.Cache(3600)) // 1 hour

// Public cache
s.Use(middleware.CachePublic(3600))

// Private cache
s.Use(middleware.CachePrivate(3600))

// Immutable content
s.Use(middleware.CacheImmutable(86400)) // 1 day

// Disable caching
s.Use(middleware.NoCache())
```

## Configuration

```go
s.Use(middleware.CacheWithConfig(middleware.CacheConfig{
    MaxAge:            3600, // 1 hour
    SMaxAge:           1800, // 30 minutes for shared caches
    Public:            true,
    MustRevalidate:    true,
    Immutable:         false,
    StaleWhileRevalidate: 60,
    StaleIfError:      300,
    VaryHeaders:        []string{"Accept-Encoding", "Accept-Language"},
    SkipFunc: func(r *http.Request) bool {
        // Don't cache authenticated requests
        return r.Header.Get("Authorization") != ""
    },
}))
```

## Cache Directives

### Public and Private

```go
// Public: can be cached by any cache
s.Use(middleware.CachePublic(3600))

// Private: only for single user
s.Use(middleware.CachePrivate(3600))
```

### Max-Age

```go
// Cache for 1 hour
s.Use(middleware.Cache(3600))

// Cache for 1 day
s.Use(middleware.Cache(86400))
```

### No Cache

```go
// Disable caching
s.Use(middleware.NoCache())
```

This sets: `Cache-Control: no-cache, no-store, must-revalidate`

### Immutable Content

For content that never changes:

```go
s.Use(middleware.CacheImmutable(86400 * 365)) // 1 year
```

Sets: `Cache-Control: public, max-age=31536000, immutable`

## Advanced Configuration

### Stale-While-Revalidate

Serve stale content while revalidating:

```go
s.Use(middleware.CacheWithConfig(middleware.CacheConfig{
    MaxAge:            3600,
    StaleWhileRevalidate: 60, // Serve stale for 60 seconds while revalidating
}))
```

### Stale-If-Error

Serve stale content if there's an error:

```go
s.Use(middleware.CacheWithConfig(middleware.CacheConfig{
    MaxAge:        3600,
    StaleIfError:  300, // Serve stale for 5 minutes if error
}))
```

### Vary Header

Specify headers that affect caching:

```go
s.Use(middleware.CacheWithConfig(middleware.CacheConfig{
    MaxAge: 3600,
    VaryHeaders: []string{
        "Accept-Encoding",
        "Accept-Language",
        "User-Agent",
    },
}))
```

## Helper Functions

Set cache headers directly:

```go
// Set Cache-Control header directly
middleware.SetCacheControl(w, "public, max-age=3600")

// Set Expires header
middleware.SetExpires(w, time.Now().Add(time.Hour))

// Set Last-Modified header
middleware.SetLastModified(w, time.Now())
```

## Skipping Cache

Skip cache headers for specific requests:

```go
s.Use(middleware.CacheWithConfig(middleware.CacheConfig{
    MaxAge: 3600,
    SkipFunc: func(r *http.Request) bool {
        // Don't cache authenticated or dynamic requests
        return r.Header.Get("Authorization") != "" ||
               r.URL.Path == "/api/dynamic"
    },
}))
```

## Common Patterns

### Static Assets

```go
// Cache static assets for 1 year
s.Static("/assets/", "./public")
s.Group("/assets", middleware.CacheImmutable(86400 * 365))
```

### API Responses

```go
// Short cache for API responses
api := s.Group("/api", middleware.Cache(300)) // 5 minutes
```

### User-Specific Content

```go
// Private cache for user content
s.Group("/user", middleware.CachePrivate(3600))
```

## Example

```go
s := helix.New(nil)

// Global: no cache by default
s.Use(middleware.NoCache())

// Public API: 5 minute cache
api := s.Group("/api/public", middleware.CachePublic(300))

// Static assets: 1 year cache
s.Group("/static", middleware.CacheImmutable(86400 * 365))

// User content: private cache
s.Group("/user", middleware.CachePrivate(3600))
```

## Best Practices

- Use `public` for cacheable content
- Use `private` for user-specific content
- Use `immutable` for versioned static assets
- Use `no-cache` or `no-store` for dynamic content
- Set appropriate `max-age` based on content volatility
- Use `Vary` header when content varies by headers
