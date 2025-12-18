# Compress Middleware

Compresses HTTP responses using gzip or deflate compression.

## Basic Usage

```go
s.Use(middleware.Compress())
```

## Configuration

```go
s.Use(middleware.CompressWithConfig(middleware.CompressConfig{
    Level:   gzip.BestCompression, // 1-9, or gzip.DefaultCompression
    MinSize: 1024,                  // Only compress if >= 1KB
    Types: []string{
        "text/",
        "application/json",
        "application/javascript",
    },
    SkipFunc: func(r *http.Request) bool {
        // Don't compress small responses
        return false
    },
}))
```

## Features

- Gzip and deflate support
- Automatic content type detection
- Minimum size threshold
- Pooled writers for performance
- Respects `Accept-Encoding` header

## Compression Levels

- `gzip.DefaultCompression` (-1) - Default compression level
- `gzip.BestSpeed` (1) - Fastest compression
- `gzip.BestCompression` (9) - Best compression ratio
- Levels 2-8 - Balance between speed and compression

## Minimum Size

Only compress responses above a certain size:

```go
s.Use(middleware.CompressWithConfig(middleware.CompressConfig{
    MinSize: 1024, // Only compress if >= 1KB
}))
```

Small responses may not benefit from compression due to overhead.

## Content Types

Specify which content types to compress:

```go
s.Use(middleware.CompressWithConfig(middleware.CompressConfig{
    Types: []string{
        "text/html",
        "text/css",
        "text/javascript",
        "application/json",
        "application/javascript",
        "application/xml",
        "image/svg+xml",
    },
}))
```

Use prefix matching with trailing slash:

```go
s.Use(middleware.CompressWithConfig(middleware.CompressConfig{
    Types: []string{
        "text/",        // All text/* types
        "application/", // All application/* types
    },
}))
```

## Skipping Compression

Skip compression for specific requests:

```go
s.Use(middleware.CompressWithConfig(middleware.CompressConfig{
    SkipFunc: func(r *http.Request) bool {
        // Don't compress already compressed content
        return r.URL.Path == "/compressed-file.gz"
    },
}))
```

## Accept-Encoding

The middleware respects the `Accept-Encoding` header:

- If client accepts `gzip`, uses gzip compression
- If client accepts `deflate`, uses deflate compression
- If client doesn't accept compression, response is not compressed

## Performance

The middleware uses object pooling for compression writers to minimize allocations:

- Gzip writers are pooled and reused
- Deflate writers are pooled and reused
- Buffers are reused to reduce memory allocations

## Example

```go
s := helix.New()

// Basic compression
s.Use(middleware.Compress())

// High compression for text content
s.Use(middleware.CompressWithConfig(middleware.CompressConfig{
    Level:   gzip.BestCompression,
    MinSize: 512,
    Types: []string{
        "text/",
        "application/json",
    },
}))
```
