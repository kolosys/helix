# Handler Patterns

Helix supports three handler patterns to accommodate different use cases. Choose the pattern that best fits your needs.

## Recommended: HandleCtx

The `HandleCtx` pattern is recommended for most applications. It provides a fluent API with automatic error handling.

### Why Use HandleCtx

- Clean, readable code with fluent method chaining
- Automatic error conversion to RFC 7807 Problem responses
- Access to all request data through a unified context
- Works seamlessly with all Helix features

### Basic Usage

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

### Request Data Access

```go
s.POST("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    // Path parameters
    id := c.Param("id")
    idInt, err := c.ParamInt("id")
    uuid, err := c.ParamUUID("id")

    // Query parameters
    name := c.Query("name")
    page := c.QueryInt("page", 1)
    active := c.QueryBool("active")
    tags := c.QuerySlice("tags")

    // Headers
    auth := c.Header("Authorization")

    // JSON body
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.BadRequest("invalid request")
    }

    // Response
    return c.Created(user)
}))
```

### Response Methods

```go
c.OK(data)                    // 200 OK with JSON
c.Created(data)               // 201 Created with JSON
c.Accepted(data)              // 202 Accepted with JSON
c.NoContent()                 // 204 No Content
c.JSON(status, data)          // Custom status with JSON
c.Text(status, text)          // Plain text response
c.HTML(status, html)          // HTML response
c.Redirect(url, code)         // HTTP redirect
c.File(path)                  // Serve file
```

### Error Responses

```go
c.BadRequest(message)         // 400
c.Unauthorized(message)       // 401
c.Forbidden(message)          // 403
c.NotFound(message)           // 404
c.InternalServerError(message) // 500

// Or return Problem errors directly
return helix.NotFoundf("user %d not found", id)
return helix.BadRequestf("invalid email: %s", email)
```

---

## Advanced: Typed Handlers

Typed handlers use Go generics for automatic request binding and type-safe responses. Use them when you want maximum type safety and reduced boilerplate.

### Why Use Typed Handlers

- Automatic request binding via struct tags
- Compile-time type checking for request/response types
- Less boilerplate for CRUD operations
- Automatic validation when implementing `Validatable`

### When to Use

- Large APIs with many similar endpoints
- When you want compile-time type safety for request/response
- CRUD operations with consistent patterns
- APIs where request structure is well-defined

### Basic Usage

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

s.POST("/users", helix.Handle(func(ctx context.Context, req CreateUserRequest) (User, error) {
    user, err := userService.Create(ctx, req.Name, req.Email)
    if err != nil {
        return User{}, err
    }
    return user, nil
}))
```

### Multi-Source Binding

Bind from multiple sources using struct tags:

```go
type UpdateUserRequest struct {
    ID    int    `path:"id"`              // From URL path
    Page  int    `query:"page"`           // From query string
    Token string `header:"Authorization"` // From header
    Name  string `json:"name"`            // From JSON body
    Email string `json:"email"`           // From JSON body
}

s.PUT("/users/{id}", helix.Handle(func(ctx context.Context, req UpdateUserRequest) (User, error) {
    // req.ID comes from path, req.Name/Email from JSON body
    return userService.Update(ctx, req.ID, req.Name, req.Email)
}))
```

### Handler Variants

```go
// Returns 201 Created
helix.HandleCreated(handler)

// Returns 202 Accepted
helix.HandleAccepted(handler)

// Custom status code
helix.HandleWithStatus(http.StatusCreated, handler)

// No request body (GET endpoints)
helix.HandleNoRequest(func(ctx context.Context) ([]User, error) {
    return userService.List(ctx)
})

// No response body (DELETE endpoints)
helix.HandleNoResponse(func(ctx context.Context, req DeleteRequest) error {
    return userService.Delete(ctx, req.ID)
})

// No request, no response (health checks)
helix.HandleEmpty(func(ctx context.Context) error {
    return pingService(ctx)
})
```

### Validation

Implement the `Validatable` interface for automatic validation:

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
    }
    return v.Err()
}

// Validation is called automatically
s.POST("/users", helix.Handle(func(ctx context.Context, req CreateUserRequest) (User, error) {
    // req is already validated
    return userService.Create(ctx, req)
}))
```

---

## Compatibility: http.HandlerFunc

Standard `http.HandlerFunc` handlers provide maximum compatibility with the Go standard library and third-party middleware.

### Why Use http.HandlerFunc

- Integrating with existing stdlib-based code
- Maximum control over response handling
- Using third-party handlers directly
- When you need raw access to `http.ResponseWriter`

### When to Use

- Wrapping existing handlers
- Serving static files
- WebSocket upgrades
- Proxying requests
- Any case where you need direct response control

### Basic Usage

```go
s.GET("/", func(w http.ResponseWriter, r *http.Request) {
    helix.OK(w, map[string]string{"status": "ok"})
})
```

### Using Helix Helpers

```go
s.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := helix.Param(r, "id")
    user, err := userService.Get(r.Context(), id)
    if err != nil {
        helix.NotFound(w, "user not found")
        return
    }
    helix.OK(w, user)
})
```

### Response Helpers

```go
helix.JSON(w, status, data)   // JSON response
helix.OK(w, data)             // 200 OK
helix.Created(w, data)        // 201 Created
helix.NoContent(w)            // 204 No Content
helix.Text(w, status, text)   // Plain text
helix.HTML(w, status, html)   // HTML
helix.BadRequest(w, message)  // 400 error
helix.NotFound(w, message)    // 404 error
```

---

## Comparison Table

| Feature | HandleCtx | Typed Handlers | http.HandlerFunc |
|---------|-----------|----------------|------------------|
| Fluent API | Yes | No | No |
| Auto error handling | Yes | Yes | No |
| Auto request binding | Manual | Automatic | Manual |
| Type safety | Runtime | Compile-time | Runtime |
| Learning curve | Low | Medium | Lowest |
| Boilerplate | Low | Lowest | Highest |
| stdlib compatible | Via adapter | Via adapter | Native |

## Choosing a Pattern

```
Start here
    │
    ▼
Need stdlib compatibility? ──Yes──► http.HandlerFunc
    │
    No
    │
    ▼
Many similar CRUD endpoints? ──Yes──► Typed Handlers
    │
    No
    │
    ▼
Use HandleCtx (recommended)
```

## Mixing Patterns

You can mix patterns in the same application:

```go
s := helix.Default(nil)

// HandleCtx for most routes
s.GET("/health", helix.HandleCtx(healthCheck))

// Typed handlers for CRUD
s.POST("/users", helix.Handle(createUser))
s.GET("/users", helix.HandleNoRequest(listUsers))

// http.HandlerFunc for special cases
s.GET("/ws", websocketHandler)
```

## Next Steps

- [Quick Start](../getting-started/quick-start.md) - Get started with HandleCtx
- [Binding](./binding.md) - Request binding in detail
- [Error Handling](./error-handling.md) - RFC 7807 Problem Details
