# Helix Framework Overview

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain the core concepts and architecture of the Helix package.

## About This Package

**Import Path:** `github.com/kolosys/helix`

Package helix provides a zero-dependency, context-aware, high-performance HTTP web framework for Go with stdlib compatibility.

## Architecture Overview

Helix is built on Go's standard library (`net/http`) and provides a thin, high-performance layer on top. The architecture follows these principles:

- **Zero Dependencies**: Built entirely on Go's standard library
- **Performance First**: Zero-allocation hot paths using `sync.Pool`
- **Type Safety**: Generic handlers with compile-time type checking
- **Developer Experience**: Fluent APIs and sensible defaults
- **Standards Compliant**: RFC 7807 Problem Details, stdlib compatibility

### Core Components

1. **Server**: Main HTTP server with middleware support
2. **Router**: Radix tree-based route matching
3. **Ctx**: Unified context wrapper for handlers
4. **Binding**: Type-safe request data binding
5. **Problem**: RFC 7807 error responses
6. **Middleware**: Request/response processing pipeline

### Data Flow

```text
HTTP Request
    ↓
Middleware Chain (RequestID, Logger, etc.)
    ↓
Router (matches route, extracts parameters)
    ↓
Handler (Ctx handler, typed handler, or http.HandlerFunc)
    ↓
Response (JSON, Problem, etc.)
    ↓
Middleware Chain (response processing)
    ↓
HTTP Response
```

### Design Patterns

- **Functional Options**: Configuration via option functions
- **Middleware Chain**: Composable request/response processing
- **Dependency Injection**: Type-safe service registry
- **Module Pattern**: Organize routes into modules
- **Resource Pattern**: REST resource builder

## Core Concepts

### Server Creation

Helix provides two ways to create a server:

```go
// Basic server (no middleware)
s := helix.New()

// Default server (includes RequestID, Logger, Recover)
s := helix.Default()
```

### Handler Types

Helix supports three handler styles:

#1. **Standard Handler**: Works with any `http.HandlerFunc`

```go
s.GET("/", func(w http.ResponseWriter, r *http.Request) {
    helix.OK(w, data)
})
```

#2. **Ctx Handler**: Uses `Ctx` for fluent API

```go
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    return c.OK(users)
}))
```

#3. **Typed Handler**: Generic handlers with automatic binding

```go
s.POST("/users", helix.Handle(func(ctx context.Context, req CreateRequest) (User, error) {
    return userService.Create(ctx, req)
}))
```

### Request Binding

Bind request data using struct tags:

```go
type CreateUserRequest struct {
    ID    int    `path:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Page  int    `query:"page"`
}

req, err := helix.Bind[CreateUserRequest](r)
```

### Error Handling

Errors are automatically converted to RFC 7807 Problem responses:

```go
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    user, err := userService.Get(c.Context(), c.Param("id"))
    if err != nil {
        return helix.NotFoundf("user not found")
    }
    return c.OK(user)
}))
```

### Routing

Routes support path parameters, groups, modules, and resources:

```go
// Path parameters
s.GET("/users/{id}", handler)

// Groups
api := s.Group("/api/v1")
api.GET("/users", handler)

// Modules
s.Mount("/users", &UserModule{})

// Resources
s.Resource("/users").CRUD(list, create, get, update, delete)
```

## Design Decisions

### Zero Dependencies

Helix has zero external dependencies because:

- **Reliability**: Fewer dependencies = fewer breaking changes
- **Performance**: No dependency overhead
- **Compatibility**: Works with any Go version that supports stdlib
- **Simplicity**: Easier to understand and maintain

### Performance Considerations

- **Object Pooling**: `Ctx` and path parameters use `sync.Pool`
- **Radix Tree**: O(k) route matching where k is path length
- **Pre-compilation**: Middleware chain compiled before serving
- **Zero Allocations**: Hot paths avoid allocations

### API Design Choices

- **Functional Options**: Type-safe, extensible configuration
- **Generic Handlers**: Compile-time type checking
- **Fluent API**: Chainable methods for readability
- **Error Returns**: Explicit error handling

### Backward Compatibility

Helix prioritizes:

- **stdlib Compatibility**: Works with any `http.Handler` middleware
- **API Stability**: Breaking changes are avoided when possible
- **Migration Path**: Easy to migrate from other frameworks

## Usage Patterns

### Basic Server

```go
package main

import (
    "github.com/kolosys/helix"
)

func main() {
    s := helix.Default()

    s.GET("/", func(w http.ResponseWriter, r *http.Request) {
        helix.OK(w, map[string]string{"message": "Hello, World!"})
    })

    s.Start(":8080")
}
```

### RESTful API

```go
s := helix.Default()

api := s.Group("/api/v1")

// Users resource
api.Resource("/users").CRUD(
    listUsers,
    createUser,
    getUser,
    updateUser,
    deleteUser,
)
```

### Modular Architecture

```go
// users/module.go
type UserModule struct {
    service *UserService
}

func (m *UserModule) Register(r helix.RouteRegistrar) {
    r.GET("/", m.list)
    r.POST("/", m.create)
    r.GET("/{id}", m.get)
}

// main.go
s.Mount("/api/v1/users", &UserModule{service: userService})
```

### Production Server

```go
s := helix.New(
    helix.WithAddr(":8080"),
    helix.WithReadTimeout(30 * time.Second),
    helix.WithWriteTimeout(30 * time.Second),
    helix.WithGracePeriod(30 * time.Second),
)

// Production middleware
for _, mw := range middleware.Production() {
    s.Use(mw)
}

s.OnStart(func(s *helix.Server) {
    // Initialize services
})

s.OnStop(func(ctx context.Context, s *helix.Server) {
    // Cleanup
})

s.Start()
```

## Common Pitfalls

### Pitfall 1: Not Returning Errors

**Problem**: Handlers must return errors for proper response handling.

**Solution**: Always return errors (or `nil`):

```go
// ❌ Wrong
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    c.OK(users)  // Missing return
}))

// ✅ Correct
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    return c.OK(users)
}))
```

### Pitfall 2: Reading Body Multiple Times

**Problem**: HTTP request bodies can only be read once.

**Solution**: Let binding handle the body:

```go
// ❌ Wrong
body, _ := io.ReadAll(r.Body)
var req CreateRequest
helix.BindJSON[CreateRequest](r)  // Fails

// ✅ Correct
var req CreateRequest
helix.BindJSON[CreateRequest](r)  // Reads body automatically
```

### Pitfall 3: Route Order

**Problem**: More specific routes must be registered before catch-all routes.

**Solution**: Register specific routes first:

```go
// ❌ Wrong
s.GET("/files/{path...}", catchAll)
s.GET("/files/users", users)  // Never matches

// ✅ Correct
s.GET("/files/users", users)
s.GET("/files/{path...}", catchAll)
```

## Integration Guide

### With Standard Library

Helix is fully compatible with `net/http`:

```go
// Use stdlib middleware
s.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Your middleware
        next.ServeHTTP(w, r)
    })
})

// Use stdlib handlers
s.Handle("/legacy", legacyHandler)
```

### With Other Frameworks

Helix can be integrated with other frameworks:

```go
// Mount Helix server as sub-router
mux := http.NewServeMux()
helixServer := helix.Default()
mux.Handle("/api/", helixServer)
```

### With Database Libraries

Helix works with any database library:

```go
// Using sql.DB
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    rows, err := db.QueryContext(c.Context(), "SELECT...")
    // ...
}))

// Using gorm
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    var users []User
    db.WithContext(c.Context()).Find(&users)
    return c.OK(users)
}))
```

## Further Reading

- [Binding Guide](./binding.md) - Request data binding
- [Context Guide](./context.md) - Using Ctx for handlers
- [Routing Guide](./routing.md) - Route organization
- [Error Handling](./error-handling.md) - Error handling patterns
- [Customizing](./customizing.md) - Framework customization
- [API Reference](../api-reference/helix.md) - Complete API documentation
- [Examples](../examples/README.md) - Practical examples

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
