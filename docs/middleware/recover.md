# Recover Middleware

Recovers from panics and returns a 500 Internal Server Error response. Prevents the server from crashing.

## Basic Usage

```go
s.Use(middleware.Recover())
```

## Configuration

```go
s.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
    PrintStack: true,
    StackSize:  4 * 1024, // 4KB
    Output:     os.Stderr,
    Handler: func(w http.ResponseWriter, r *http.Request, err any) {
        // Custom panic handler
        logs.Error("panic recovered", logs.Any("error", err))
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    },
}))
```

## Features

- Catches panics in handlers
- Logs panic and stack trace
- Returns 500 Internal Server Error
- Prevents server crash

## Configuration Options

- **PrintStack**: Enable stack trace printing (default: `true`)
- **StackSize**: Maximum stack trace buffer size (default: 4KB)
- **Output**: Writer for panic output (default: `os.Stderr`)
- **Handler**: Custom panic handler function

## Custom Panic Handler

```go
s.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
    Handler: func(w http.ResponseWriter, r *http.Request, err any) {
        // Log panic with context
        logs.ErrorContext(r.Context(), "panic recovered",
            logs.Any("error", err),
            logs.String("path", r.URL.Path),
            logs.String("method", r.Method),
        )

        // Send custom error response
        helix.WriteProblem(w, helix.ErrInternal.WithDetailf("panic: %v", err))
    },
}))
```

## Disable Stack Trace

```go
s.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
    PrintStack: false, // Don't print stack trace
}))
```

## Important Notes

- Recover middleware should be added **last** (innermost) so it catches panics from all other middleware
- The default handler writes to `os.Stderr` - redirect if you want logs elsewhere
- Stack traces are truncated to `StackSize` bytes
