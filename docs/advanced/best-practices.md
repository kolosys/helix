# Best Practices

This guide covers best practices for building production-ready applications with Helix.

## Handler Patterns

### Use HandleCtx as Your Default

`HandleCtx` is the recommended pattern for most handlers. It provides a fluent API with automatic error handling:

```go
// ✅ Recommended
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

### Always Return Errors

Handlers must return errors (or `nil`). Forgetting to return will cause the response not to be sent:

```go
// ❌ Wrong - missing return
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    c.OK(users)  // Response not sent!
}))

// ✅ Correct
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    return c.OK(users)
}))
```

### Use Typed Handlers for Complex APIs

For large APIs with many similar endpoints, typed handlers provide compile-time safety:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

s.POST("/users", helix.Handle(func(ctx context.Context, req CreateUserRequest) (User, error) {
    return userService.Create(ctx, req)
}))
```

## Error Handling

### Use Specific Problem Errors

Return specific problem errors instead of generic ones:

```go
// ❌ Generic
return helix.ErrNotFound

// ✅ Specific
return helix.NotFoundf("user %d not found", id)
```

### Create Domain-Specific Errors

Define domain-specific problem errors for your API:

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
)

// Usage
return ErrUserNotFound.WithDetailf("user %d not found", id)
```

### Don't Expose Internal Errors

Wrap internal errors with user-friendly messages:

```go
// ❌ Wrong - exposes database error
user, err := db.Query("SELECT...")
if err != nil {
    return err  // "pq: relation 'users' does not exist"
}

// ✅ Correct
user, err := db.Query("SELECT...")
if err != nil {
    log.Printf("database error: %v", err)
    return helix.ErrInternal.WithDetail("failed to retrieve user")
}
```

## Request Binding

### Use Explicit Struct Tags

Always use explicit struct tags for binding:

```go
type UpdateUserRequest struct {
    ID    int    `path:"id"`       // From URL path
    Name  string `json:"name"`     // From JSON body
    Email string `json:"email"`    // From JSON body
    Page  int    `query:"page"`    // From query string
}
```

### Implement Validation

Implement the `Validatable` interface for request validation:

```go
func (r *CreateUserRequest) Validate() error {
    v := helix.NewValidationErrors()
    
    if r.Name == "" {
        v.Add("name", "name is required")
    }
    if r.Email == "" {
        v.Add("email", "email is required")
    } else if !isValidEmail(r.Email) {
        v.Add("email", "invalid email format")
    }
    
    return v.Err()
}
```

## Project Organization

### Use Modules for Route Organization

Organize routes into modules for large applications:

```go
// users/module.go
type UserModule struct {
    service *UserService
}

func (m *UserModule) Register(r helix.RouteRegistrar) {
    r.GET("/", helix.HandleCtx(m.list))
    r.POST("/", helix.HandleCtx(m.create))
    r.GET("/{id}", helix.HandleCtx(m.get))
    r.PUT("/{id}", helix.HandleCtx(m.update))
    r.DELETE("/{id}", helix.HandleCtx(m.delete))
}

// main.go
s.Mount("/api/v1/users", &UserModule{service: userService})
```

### Group Related Routes

Use route groups for shared prefixes and middleware:

```go
api := s.Group("/api/v1")

// Public routes
api.GET("/health", helix.HandleCtx(healthCheck))

// Authenticated routes
auth := api.Group("", authMiddleware)
auth.Mount("/users", &UserModule{})
auth.Mount("/posts", &PostModule{})

// Admin routes
admin := auth.Group("/admin", adminMiddleware)
admin.GET("/stats", helix.HandleCtx(getStats))
```

## Middleware

### Use Appropriate Bundles

Use pre-configured middleware bundles:

```go
// Development
s := helix.Default(nil)  // Includes RequestID, Logger (dev), Recover

// Production
s := helix.New(nil)
for _, mw := range middleware.Production() {
    s.Use(mw)
}
```

### Order Matters

Add middleware in the correct order:

```go
// Recover should be last (innermost) to catch all panics
s.Use(middleware.RequestID())
s.Use(middleware.Logger(middleware.LogFormatJSON))
s.Use(middleware.Recover())  // Last!
```

## Performance

### Use Build() in Production

Call `Build()` to pre-compile the middleware chain:

```go
s := helix.New(nil)
// ... register routes and middleware ...
s.Build()  // Pre-compile for performance
s.Start(":8080")
```

### Don't Read Body Multiple Times

The HTTP request body can only be read once:

```go
// ❌ Wrong
body, _ := io.ReadAll(r.Body)
var req Request
c.Bind(&req)  // Fails - body already read

// ✅ Correct
var req Request
c.Bind(&req)  // Reads body automatically
```

## Security

### Use Rate Limiting

Protect your API with rate limiting:

```go
s.Use(middleware.RateLimit(100, 10))  // 100 req/sec, burst of 10
```

### Configure CORS Properly

Don't use permissive CORS in production:

```go
// ❌ Wrong for production
s.Use(middleware.CORSAllowAll())

// ✅ Correct
s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"https://app.example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
}))
```

### Set Timeouts

Always configure server timeouts:

```go
s := helix.New(&helix.Options{
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  120 * time.Second,
    GracePeriod:  30 * time.Second,
})
```

## Testing

### Use httptest for Handler Tests

Test handlers using `httptest`:

```go
func TestGetUser(t *testing.T) {
    s := helix.New(nil)
    s.GET("/users/{id}", helix.HandleCtx(getUser))
    
    req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
    rec := httptest.NewRecorder()
    
    s.ServeHTTP(rec, req)
    
    if rec.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rec.Code)
    }
}
```

### Test Error Cases

Test both success and error paths:

```go
func TestGetUser_NotFound(t *testing.T) {
    s := helix.New(nil)
    s.GET("/users/{id}", helix.HandleCtx(getUser))
    
    req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
    rec := httptest.NewRecorder()
    
    s.ServeHTTP(rec, req)
    
    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", rec.Code)
    }
}
```

## Logging

### Use Structured Logging

Use the standard library's `log/slog` for application logging:

```go
import "log/slog"

s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id := c.Param("id")
    
    user, err := userService.Get(c.Context(), id)
    if err != nil {
        slog.Error("failed to get user",
            slog.String("id", id),
            slog.Any("error", err),
        )
        return helix.NotFoundf("user not found")
    }
    
    return c.OK(user)
}))
```

### Configure Request Logging

Use appropriate log formats:

```go
// Development
s.Use(middleware.Logger(middleware.LogFormatDev))

// Production
s.Use(middleware.Logger(middleware.LogFormatJSON))
```

## Resources

- [Handler Patterns Guide](../core-concepts/handler-patterns.md)
- [Error Handling Guide](../core-concepts/error-handling.md)
- [Middleware Reference](../middleware/middleware.md)
- [Performance Tuning](./performance-tuning.md)

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
