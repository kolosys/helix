# CORS Middleware

Handles Cross-Origin Resource Sharing (CORS) headers. Allows configuring which origins, methods, and headers are allowed.

## Basic Usage

```go
// Default (allows all origins)
s.Use(middleware.CORS())

// Allow all (development only)
s.Use(middleware.CORSAllowAll())
```

## Configuration

```go
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"https://example.com", "https://app.example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    ExposeHeaders:    []string{"X-Total-Count"},
    AllowCredentials: true,
    MaxAge:           86400, // 24 hours
    AllowOriginFunc: func(origin string) bool {
        // Custom origin validation
        return strings.HasSuffix(origin, ".example.com")
    },
}))
```

## Features

- Handles preflight OPTIONS requests
- Validates origin against allowed list
- Sets appropriate CORS headers
- Supports credentials

## Configuration Options

- **AllowOrigins**: List of allowed origins (use `"*"` for all)
- **AllowOriginFunc**: Custom function to validate origins
- **AllowMethods**: Allowed HTTP methods (default: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS)
- **AllowHeaders**: Allowed request headers
- **ExposeHeaders**: Headers exposed to client
- **AllowCredentials**: Allow credentials (cookies, authorization headers)
- **MaxAge**: Preflight cache duration in seconds

## Common Patterns

### Allow Specific Origins

```go
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{
        "https://app.example.com",
        "https://admin.example.com",
    },
}))
```

### Allow Subdomains

```go
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOriginFunc: func(origin string) bool {
        return strings.HasSuffix(origin, ".example.com")
    },
}))
```

### With Credentials

```go
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"https://app.example.com"}, // Must be specific, not "*"
    AllowCredentials: true,
    AllowHeaders:     []string{"Authorization", "Content-Type"},
}))
```

**Important**: When `AllowCredentials` is `true`, `AllowOrigins` cannot contain `"*"`. You must specify exact origins.

### Custom Methods and Headers

```go
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"https://example.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
    AllowHeaders: []string{
        "Authorization",
        "Content-Type",
        "X-Request-ID",
        "X-Custom-Header",
    },
    ExposeHeaders: []string{"X-Total-Count", "X-Page-Count"},
}))
```

### Preflight Caching

```go
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"https://example.com"},
    MaxAge:       86400, // Cache preflight requests for 24 hours
}))
```

## Common Pitfalls

### Credentials with Wildcard

```go
// ❌ Wrong - credentials with wildcard won't work
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"*"},
    AllowCredentials: true, // Won't work!
}))

// ✅ Correct - specific origins required
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"https://app.example.com"},
    AllowCredentials: true,
}))
```
