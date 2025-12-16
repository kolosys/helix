# Context and Ctx

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain the core concepts and architecture of context handling in Helix.

## About Ctx

The `Ctx` type is Helix's unified context wrapper that provides a fluent, chainable API for accessing request data and writing responses. It wraps `http.Request` and `http.ResponseWriter` with convenient methods for common HTTP operations.

## Core Concepts

### What is Ctx?

`Ctx` is a lightweight wrapper around the standard `http.Request` and `http.ResponseWriter` that provides:

- **Fluent API**: Chainable methods for setting headers, cookies, and status codes
- **Type-Safe Accessors**: Convenient methods for path params, query params, and headers
- **Request-Scoped Storage**: Dependency injection for request-scoped values
- **Error Handling**: Automatic conversion of errors to RFC 7807 Problem responses

### Ctx vs Standard Handlers

Helix provides three handler styles:

1. **Standard `http.HandlerFunc`** - Works with any Go HTTP middleware
2. **`CtxHandler`** - Uses `Ctx` for fluent API
3. **Typed Handlers** - Generic handlers with automatic binding

## Architecture Overview

### Ctx Structure

```go
type Ctx struct {
    Request  *http.Request
    Response http.ResponseWriter

    status int              // Pending status code
    store  map[string]any  // Request-scoped storage
}
```

### Request Flow

1. Request arrives at Helix server
2. `Ctx` is created (or retrieved from pool) with `Request` and `ResponseWriter`
3. Handler receives `Ctx` and processes request
4. Handler returns error (if any) which is converted to RFC 7807 Problem
5. `Ctx` is reset and returned to pool for reuse

## Core Concepts

### Request Data Access

`Ctx` provides convenient accessors for all request data:

#### Path Parameters

```go
// String parameter
id := c.Param("id")

// Typed parameters
userID, err := c.ParamInt("id")
userID64, err := c.ParamInt64("id")
uuid, err := c.ParamUUID("id")
```

#### Query Parameters

```go
// String with default
name := c.QueryDefault("name", "World")

// Typed query parameters
page := c.QueryInt("page", 1)
limit := c.QueryInt64("limit", 20)
price := c.QueryFloat64("price", 0.0)
active := c.QueryBool("active")
tags := c.QuerySlice("tags")  // []string
```

#### Headers

```go
auth := c.Header("Authorization")
contentType := c.Header("Content-Type")
```

#### Request Body Binding

```go
var req CreateUserRequest
if err := c.Bind(&req); err != nil {
    return c.BadRequest("invalid JSON")
}
```

### Response Methods

All response methods return `*Ctx` for chaining:

#### JSON Responses

```go
// Status-specific methods
c.OK(data)              // 200 OK
c.Created(data)          // 201 Created
c.Accepted(data)         // 202 Accepted
c.NoContent()            // 204 No Content

// Generic JSON
c.JSON(http.StatusOK, data)
```

#### Other Content Types

```go
c.Text(http.StatusOK, "Hello, World!")
c.HTML(http.StatusOK, "<h1>Hello</h1>")
c.Blob(http.StatusOK, "image/png", imageData)
c.File("/path/to/file")
```

#### Error Responses

```go
c.BadRequest("invalid input")
c.Unauthorized("authentication required")
c.Forbidden("access denied")
c.NotFound("resource not found")
c.InternalServerError("server error")
```

### Chaining

Response methods return `*Ctx` for method chaining:

```go
c.SetHeader("X-Request-ID", requestID).
  SetHeader("X-Rate-Limit", "100").
  Status(http.StatusCreated).
  JSON(user)
```

### Request-Scoped Storage

`Ctx` provides a key-value store for request-scoped dependency injection:

```go
// Set a value
c.Set("user", currentUser)
c.Set("transaction", dbTx)

// Get a value
user, ok := c.Get("user")
if !ok {
    return c.Unauthorized("not authenticated")
}

// Typed getters
user := c.GetString("userID")
count := c.GetInt("itemCount")

// Panic if not found
user := c.MustGet("user")
```

## Usage Patterns

### Basic Ctx Handler

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

### Chained Response

```go
s.POST("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.BadRequest("invalid input")
    }

    user, err := userService.Create(c.Context(), req)
    if err != nil {
        return err
    }

    return c.SetHeader("Location", fmt.Sprintf("/users/%d", user.ID)).
            Created(user)
}))
```

### Request-Scoped Services

```go
// Middleware that provides request-scoped service
func dbMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tx := db.BeginTx(r.Context())
        defer tx.Rollback()

        c := helix.NewCtx(w, r)
        c.Set("transaction", tx)

        // Convert back to handler
        next.ServeHTTP(w, r)

        if err := tx.Commit(); err != nil {
            // Handle error
        }
    })
}

// Handler using request-scoped service
s.POST("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    tx := c.MustGet("transaction").(*sql.Tx)

    // Use transaction
    _, err := tx.Exec("INSERT INTO users...")
    return err
}))
```

### Pagination

```go
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    p := c.BindPagination(20, 100)  // defaultLimit, maxLimit

    users, total, err := userService.List(
        c.Context(),
        p.GetPage(),
        p.GetLimit(20, 100),
    )
    if err != nil {
        return err
    }

    return c.Paginated(users, total, p.GetPage(), p.GetLimit(20, 100))
}))
```

## Design Decisions

### Performance Considerations

- **Object Pooling**: `Ctx` instances are pooled using `sync.Pool` to reduce allocations
- **Zero Allocations**: Hot path methods avoid allocations where possible
- **Lazy Initialization**: Request-scoped storage (`store`) is only allocated when needed

### API Design Choices

- **Fluent API**: Chainable methods improve readability and reduce boilerplate
- **Error Returns**: Handlers return `error` which is automatically converted to RFC 7807
- **Context Integration**: `Ctx.Context()` provides access to `context.Context` for cancellation

### Why Not Extend http.Request?

Helix chose to create a separate `Ctx` type rather than extending `http.Request` because:

- **Compatibility**: Works with any `http.Handler` middleware
- **Clarity**: Explicit API that's easy to understand
- **Flexibility**: Can add Helix-specific features without breaking stdlib compatibility

## Common Pitfalls

### Pitfall 1: Not Returning Errors

**Problem**: If your handler doesn't return an error, the response won't be sent.

**Solution**: Always return an error (or `nil`):

```go
// ❌ Wrong - response never sent
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    c.OK(users)
    // Missing return!
}))

// ✅ Correct
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    return c.OK(users)
}))
```

### Pitfall 2: Reading Body Multiple Times

**Problem**: HTTP request bodies can only be read once. If you read it manually and then call `c.Bind()`, binding will fail.

**Solution**: Let `Ctx.Bind()` handle the body:

```go
// ❌ Wrong
body, _ := io.ReadAll(c.Request.Body)
var req CreateRequest
c.Bind(&req)  // Fails - body already read

// ✅ Correct
var req CreateRequest
c.Bind(&req)  // Reads body automatically
```

### Pitfall 3: Status Code Not Applied

**Problem**: Setting status with `c.Status()` doesn't write the response immediately.

**Solution**: Status is applied when you write the response body:

```go
// ✅ Correct - status applied when JSON is written
c.Status(http.StatusCreated).JSON(user)

// ✅ Also correct - explicit status in method
c.Created(user)  // Always 201
```

## Integration Guide

### With Middleware

`Ctx` works seamlessly with standard middleware:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract auth from request
        token := r.Header.Get("Authorization")

        // Validate and add to context
        user := validateToken(token)
        ctx := context.WithValue(r.Context(), "user", user)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Handler accesses user from context
s.GET("/profile", helix.HandleCtx(func(c *helix.Ctx) error {
    user := c.Context().Value("user").(*User)
    return c.OK(user)
}))
```

### With Dependency Injection

Use `Ctx.Set()` and `Ctx.Get()` for request-scoped services:

```go
func provideDBMiddleware(db *sql.DB) middleware.Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tx := db.BeginTx(r.Context())
            defer tx.Rollback()

            c := helix.NewCtx(w, r)
            c.Set("db", tx)

            // Convert back - this is a limitation of current design
            // Better to use helix.ProvideMiddleware
            next.ServeHTTP(w, r)

            tx.Commit()
        })
    }
}
```

### With Typed Handlers

You can mix `Ctx` handlers with typed handlers:

```go
// Ctx handler for complex logic
s.GET("/users/{id}/posts", helix.HandleCtx(func(c *helix.Ctx) error {
    userID := c.Param("id")
    p := c.BindPagination(10, 50)
    // Complex pagination logic...
    return c.Paginated(posts, total, p.GetPage(), p.GetLimit(10, 50))
}))

// Typed handler for simple CRUD
s.POST("/users", helix.Handle(func(ctx context.Context, req CreateUserRequest) (User, error) {
    return userService.Create(ctx, req)
}))
```

## Further Reading

- [Binding Guide](./binding.md) - Request binding with struct tags
- [Error Handling](./error-handling.md) - Error handling patterns
- [Routing Guide](./routing.md) - Route registration and organization
- [API Reference](../api-reference/helix.md) - Complete Ctx API documentation

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
