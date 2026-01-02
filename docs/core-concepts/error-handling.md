# Error Handling

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain the core concepts and architecture of error handling in Helix.

## About Error Handling

Helix uses [RFC 7807 Problem Details for HTTP APIs](https://tools.ietf.org/html/rfc7807) as the standard format for error responses. This provides a consistent, machine-readable way to communicate errors that integrates seamlessly with modern API clients.

## Core Concepts

### RFC 7807 Problem Details

RFC 7807 defines a standard JSON format for HTTP API errors. Helix implements this standard with extensions for validation errors.

### Problem Structure

A Problem response contains:

- **type** (string): URI reference identifying the problem type
- **title** (string): Short, human-readable summary
- **status** (int): HTTP status code
- **detail** (string, optional): Human-readable explanation
- **instance** (string, optional): URI identifying the specific occurrence

### Error Flow

1. Handler returns an `error`
2. Helix checks if error is a `Problem`
3. If not, converts to appropriate `Problem`
4. Writes RFC 7807 JSON response
5. Sets appropriate HTTP status code

## Architecture Overview

### Problem Type

```go
type Problem struct {
    Type     string `json:"type"`
    Title    string `json:"title"`
    Status   int    `json:"status"`
    Detail   string `json:"detail,omitempty"`
    Instance string `json:"instance,omitempty"`
    Err      error  `json:"-"`  // Original error (not serialized)
}
```

### Sentinel Errors

Helix provides pre-defined `Problem` values for common HTTP errors:

- `ErrBadRequest` (400)
- `ErrUnauthorized` (401)
- `ErrForbidden` (403)
- `ErrNotFound` (404)
- `ErrMethodNotAllowed` (405)
- `ErrConflict` (409)
- `ErrGone` (410)
- `ErrUnprocessableEntity` (422)
- `ErrTooManyRequests` (429)
- `ErrInternal` (500)
- `ErrNotImplemented` (501)
- `ErrBadGateway` (502)
- `ErrServiceUnavailable` (503)
- `ErrGatewayTimeout` (504)

## Core Concepts

### Returning Errors from Handlers

Handlers return `error`, which is automatically converted to Problem responses:

```go
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id := c.Param("id")
    user, err := userService.Get(c.Context(), id)
    if err != nil {
        return helix.NotFoundf("user %s not found", id)
    }
    return c.OK(user)
}))
```

### Using Sentinel Errors

Use pre-defined errors for common cases:

```go
// Simple error
return helix.ErrNotFound

// With detail
return helix.ErrNotFound.WithDetail("user 123 not found")

// With formatted detail
return helix.NotFoundf("user %d not found", id)
```

### Creating Custom Problems

Create custom problems for domain-specific errors:

```go
var ErrDuplicateEmail = helix.NewProblem(
    http.StatusConflict,
    "duplicate_email",
    "Email Already Exists",
)

// In handler
if exists {
    return ErrDuplicateEmail.WithDetailf("email %s is already registered", email)
}
```

### Validation Errors

Validation errors are automatically converted to RFC 7807 with field-level details:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (r *CreateUserRequest) Validate() error {
    v := helix.NewValidationErrors()

    if r.Name == "" {
        v.Add("name", "name is required")
    }
    if r.Email == "" {
        v.Add("email", "email is required")
    } else if !strings.Contains(r.Email, "@") {
        v.Add("email", "invalid email format")
    }

    return v.Err()  // Returns nil if no errors
}

// Handler automatically gets validation errors
s.POST("/users", helix.Handle(func(ctx context.Context, req CreateUserRequest) (User, error) {
    // Validation happens automatically if req implements Validatable
    return userService.Create(ctx, req)
}))
```

### Validation Error Response Format

Validation errors include an `errors` array with field-level details:

```json
{
  "type": "about:blank#unprocessable_entity",
  "title": "Unprocessable Entity",
  "status": 422,
  "detail": "One or more validation errors occurred",
  "instance": "/users",
  "errors": [
    {
      "field": "name",
      "message": "name is required"
    },
    {
      "field": "email",
      "message": "invalid email format"
    }
  ]
}
```

## Usage Patterns

### Basic Error Handling

```go
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id, err := c.ParamInt("id")
    if err != nil {
        return helix.BadRequestf("invalid user id: %v", err)
    }

    user, err := userService.Get(c.Context(), id)
    if err == ErrUserNotFound {
        return helix.NotFoundf("user %d not found", id)
    }
    if err != nil {
        return err  // Wrapped as ErrInternal
    }

    return c.OK(user)
}))
```

### Error Wrapping

Wrap errors with context for better debugging:

```go
user, err := userService.Get(ctx, id)
if err != nil {
    return fmt.Errorf("failed to get user %d: %w", id, err)
    // Automatically converted to ErrInternal with detail
}
```

### Domain-Specific Errors

Define domain errors that map to HTTP problems:

```go
var (
    ErrUserNotFound = helix.NewProblem(
        http.StatusNotFound,
        "user_not_found",
        "User Not Found",
    )

    ErrEmailExists = helix.NewProblem(
        http.StatusConflict,
        "email_exists",
        "Email Already Exists",
    )

    ErrInvalidCredentials = helix.NewProblem(
        http.StatusUnauthorized,
        "invalid_credentials",
        "Invalid Credentials",
    )
)

// In service layer
func (s *UserService) Get(ctx context.Context, id int) (User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err == sql.ErrNoRows {
        return User{}, ErrUserNotFound.WithDetailf("user %d not found", id)
    }
    return user, err
}
```

### Error Handling in Middleware

Middleware can handle errors before they reach handlers:

```go
func errorLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Wrap response writer to capture errors
        rw := &errorResponseWriter{ResponseWriter: w}
        next.ServeHTTP(rw, r)

        if rw.err != nil {
            logs.ErrorContext(r.Context(), "request error",
                logs.Err(rw.err),
                logs.String("path", r.URL.Path),
                logs.Int("status", rw.status),
            )
        }
    })
}
```

## Design Decisions

### RFC 7807 Standard

Helix uses RFC 7807 because:

- **Standardization**: Widely adopted standard for HTTP APIs
- **Machine Readable**: Easy for clients to parse and handle
- **Extensibility**: Can add custom fields for domain-specific errors
- **Tooling**: Works with API documentation tools and client generators

### Automatic Error Conversion

Errors are automatically converted to Problems because:

- **Consistency**: All errors follow the same format
- **Developer Experience**: No need to manually convert errors
- **Type Safety**: Sentinel errors are type-safe
- **Flexibility**: Can still return raw errors for custom handling

### Validation Error Extension

Validation errors extend RFC 7807 with an `errors` array because:

- **Field-Level Details**: Clients can highlight specific form fields
- **Multiple Errors**: Return all validation errors at once
- **User Experience**: Better UX than single error message
- **Standards Compliant**: Still valid RFC 7807 (allows extensions)

## Common Pitfalls

### Pitfall 1: Not Returning Errors

**Problem**: If a handler doesn't return an error, the response may not be sent correctly.

**Solution**: Always return errors (or `nil`):

```go
// ❌ Wrong
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    users, _ := userService.List(c.Context())
    c.OK(users)  // Missing return
}))

// ✅ Correct
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    users, err := userService.List(c.Context())
    if err != nil {
        return err
    }
    return c.OK(users)
}))
```

### Pitfall 2: Exposing Internal Errors

**Problem**: Returning raw database or internal errors exposes implementation details.

**Solution**: Wrap errors with user-friendly messages:

```go
// ❌ Wrong - exposes database error
user, err := db.Query("SELECT...")
if err != nil {
    return err  // "pq: relation 'users' does not exist"
}

// ✅ Correct - user-friendly error
user, err := db.Query("SELECT...")
if err != nil {
    logs.ErrorContext(ctx, "database error", logs.Err(err))
    return helix.ErrInternal.WithDetail("failed to retrieve user")
}
```

### Pitfall 3: Not Handling All Error Types

**Problem**: Custom error handlers must handle all error types.

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
        logs.ErrorContext(r.Context(), "unhandled error", logs.Err(err))
        helix.WriteProblem(w, helix.ErrInternal.WithDetail("an error occurred"))
    }
}
```

## Integration Guide

### With Logging

Log errors with context:

```go
import "log"

s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    user, err := userService.Get(c.Context(), c.Param("id"))
    if err != nil {
        log.Printf("failed to get user %s: %v", c.Param("id"), err)
        return helix.NotFoundf("user not found")
    }
    return c.OK(user)
}))
```

### With Monitoring

Track error rates and types:

```go
func errorTrackingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rw := &trackingResponseWriter{ResponseWriter: w}
        next.ServeHTTP(rw, r)

        if rw.status >= 400 {
            metrics.IncrementErrorCounter(r.Method, r.URL.Path, rw.status)
        }
    })
}
```

### With Custom Error Handlers

Replace default error handling:

```go
func apiErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
    var problem helix.Problem

    switch e := err.(type) {
    case helix.Problem:
        problem = e
    case *helix.ValidationErrors:
        helix.WriteValidationProblem(w, e)
        return
    default:
        // Log unexpected errors
        logs.ErrorContext(r.Context(), "unexpected error", logs.Err(err))
        problem = helix.ErrInternal.WithDetail("an error occurred")
    }

    // Add request ID to problem
    if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
        problem = problem.WithDetailf("%s (request-id: %s)", problem.Detail, requestID)
    }

    helix.WriteProblem(w, problem)
}

s := helix.New(&helix.Options{
    ErrorHandler: apiErrorHandler,
})
```

## Further Reading

- [Binding Guide](./binding.md) - Validation error handling
- [Context Guide](./context.md) - Error handling in Ctx handlers
- [RFC 7807 Specification](https://tools.ietf.org/html/rfc7807) - Problem Details standard
- [API Reference](../api-reference/helix.md) - Complete error handling API

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
