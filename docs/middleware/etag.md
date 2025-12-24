# ETag Middleware

Generates ETag headers for responses and handles `If-None-Match` requests. Returns 304 Not Modified when appropriate.

## Basic Usage

```go
// Strong ETags
s.Use(middleware.ETag())

// Weak ETags
s.Use(middleware.ETagWeak())
```

## Configuration

```go
s.Use(middleware.ETagWithConfig(middleware.ETagConfig{
    Weak: false, // Strong ETags
    SkipFunc: func(r *http.Request) bool {
        // Skip ETag for dynamic content
        return r.URL.Path == "/api/dynamic"
    },
}))
```

## Features

- Automatic ETag generation from response body
- Handles `If-None-Match` header
- Returns 304 Not Modified when appropriate
- Supports weak ETags

## How It Works

1. Middleware buffers the response body
2. Generates ETag hash from body content
3. Checks `If-None-Match` header from request
4. If ETag matches, returns 304 Not Modified
5. Otherwise, sets ETag header and sends full response

## Strong vs Weak ETags

**Strong ETags** (default):

- Indicate exact byte-for-byte match
- Format: `"abc123..."`
- Use when content is identical

**Weak ETags**:

- Indicate semantic equivalence
- Format: `W/"abc123..."`
- Use when content is semantically same but may differ

```go
// Strong ETags (default)
s.Use(middleware.ETag())

// Weak ETags
s.Use(middleware.ETagWeak())
```

## Skipping ETag

Skip ETag generation for dynamic content:

```go
s.Use(middleware.ETagWithConfig(middleware.ETagConfig{
    SkipFunc: func(r *http.Request) bool {
        // Skip for dynamic or personalized content
        return r.URL.Path == "/api/user/profile" ||
               r.URL.Path == "/api/dynamic"
    },
}))
```

## Helper Functions

Generate ETags manually:

```go
// Generate ETag from content
etag := middleware.ETagFromContent([]byte("content"), false)

// Generate ETag from string
etag := middleware.ETagFromString("content", false)

// Generate ETag from version number
etag := middleware.ETagFromVersion(123, false)
```

## Client Usage

Clients can use ETags for caching:

```http
GET /api/users HTTP/1.1
If-None-Match: "abc123def456"
```

If content hasn't changed, server responds:

```http
HTTP/1.1 304 Not Modified
ETag: "abc123def456"
```

## Example

```go
s := helix.New(nil)

// Enable ETag for all GET requests
s.Use(middleware.ETag())

// Handler that benefits from ETag
s.GET("/api/users", helix.HandleCtx(func(c *helix.Ctx) error {
    users, err := userService.List(c.Context())
    if err != nil {
        return err
    }
    // ETag is automatically generated and checked
    return c.OK(users)
}))
```

## Best Practices

- Use ETags for cacheable content (GET requests)
- Skip ETags for dynamic or personalized content
- Use weak ETags when content is semantically equivalent but may differ
- Combine with Cache middleware for optimal caching
