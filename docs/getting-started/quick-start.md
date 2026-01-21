# Quick Start

This guide will help you get started with Helix quickly. We'll use the recommended `HandleCtx` pattern which provides a clean, fluent API for building HTTP handlers.

## Hello World

Here's the simplest possible Helix server:

```go
package main

import "github.com/kolosys/helix"

func main() {
    s := helix.Default(nil)

    s.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
        return c.OK(map[string]string{"message": "Hello, World!"})
    }))

    s.Start(":8080")
}
```

Run it:

```bash
go run main.go
```

Visit `http://localhost:8080` to see your response.

## Path Parameters

Extract dynamic values from URLs using `c.Param()`:

```go
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id := c.Param("id")
    return c.OK(map[string]string{
        "id":   id,
        "name": "John Doe",
    })
}))
```

For typed parameters, use `c.ParamInt()` or `c.ParamUUID()`:

```go
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id, err := c.ParamInt("id")
    if err != nil {
        return c.BadRequest("invalid user ID")
    }
    // Use id as int...
    return c.OK(map[string]int{"id": id})
}))
```

## Query Parameters

Access query parameters with typed helpers:

```go
s.GET("/search", helix.HandleCtx(func(c *helix.Ctx) error {
    query := c.QueryDefault("q", "")           // String with default
    page := c.QueryInt("page", 1)              // Int with default
    limit := c.QueryInt("limit", 20)           // Int with default
    active := c.QueryBool("active")            // Bool
    tags := c.QuerySlice("tags")               // []string

    return c.OK(map[string]any{
        "query":  query,
        "page":   page,
        "limit":  limit,
        "active": active,
        "tags":   tags,
    })
}))
```

## Request Body Binding

Bind JSON request bodies to structs:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

s.POST("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.BadRequest("invalid request body")
    }

    // Create user...
    user := User{ID: 1, Name: req.Name, Email: req.Email}
    return c.Created(user)
}))
```

## Error Handling

Return errors from handlers and they'll be automatically converted to RFC 7807 Problem responses:

```go
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id, err := c.ParamInt("id")
    if err != nil {
        return c.BadRequest("invalid user ID")
    }

    user, err := userService.Get(c.Context(), id)
    if err != nil {
        return helix.NotFoundf("user %d not found", id)
    }

    return c.OK(user)
}))
```

Common error helpers:
- `c.BadRequest(message)` - 400 Bad Request
- `c.Unauthorized(message)` - 401 Unauthorized
- `c.Forbidden(message)` - 403 Forbidden
- `c.NotFound(message)` - 404 Not Found
- `c.InternalServerError(message)` - 500 Internal Server Error
- `helix.NotFoundf(format, args...)` - Returns a Problem error
- `helix.BadRequestf(format, args...)` - Returns a Problem error

## Route Groups

Organize routes with shared prefixes:

```go
s := helix.Default(nil)

// API v1 routes
api := s.Group("/api/v1")
api.GET("/users", helix.HandleCtx(listUsers))
api.POST("/users", helix.HandleCtx(createUser))
api.GET("/users/{id}", helix.HandleCtx(getUser))
```

## Middleware

Add middleware globally or to specific groups:

```go
import "github.com/kolosys/helix/middleware"

s := helix.Default(nil)  // Includes RequestID, Logger, Recover

// Add more middleware
s.Use(middleware.CORS())
s.Use(middleware.Compress())

// Group with middleware
admin := s.Group("/admin", middleware.BasicAuth("admin", "secret"))
admin.GET("/stats", helix.HandleCtx(getStats))
```

## Complete Example

Here's a complete example with routing, binding, and error handling:

```go
package main

import "github.com/kolosys/helix"

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    s := helix.Default(nil)

    // Health check
    s.GET("/health", helix.HandleCtx(func(c *helix.Ctx) error {
        return c.OK(map[string]string{"status": "healthy"})
    }))

    // User endpoints
    api := s.Group("/api/v1")

    api.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
        id, err := c.ParamInt("id")
        if err != nil {
            return c.BadRequest("invalid user ID")
        }
        return c.OK(User{ID: id, Name: "John Doe", Email: "john@example.com"})
    }))

    api.POST("/users", helix.HandleCtx(func(c *helix.Ctx) error {
        var req CreateUserRequest
        if err := c.Bind(&req); err != nil {
            return c.BadRequest("invalid request body")
        }
        user := User{ID: 1, Name: req.Name, Email: req.Email}
        return c.Created(user)
    }))

    s.Start(":8080")
}
```

## Next Steps

- [Handler Patterns](../core-concepts/handler-patterns.md) - Learn about different handler styles
- [Routing](../core-concepts/routing.md) - Route groups, modules, and resources
- [Middleware](../middleware/middleware.md) - Built-in middleware reference
- [Error Handling](../core-concepts/error-handling.md) - RFC 7807 Problem Details
