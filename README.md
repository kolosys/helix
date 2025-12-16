# Helix 🧬

![GoVersion](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Zero Dependencies](https://img.shields.io/badge/Zero-Dependencies-green.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/kolosys/helix.svg)](https://pkg.go.dev/github.com/kolosys/helix)
[![Go Report Card](https://goreportcard.com/badge/github.com/kolosys/helix)](https://goreportcard.com/report/github.com/kolosys/helix)

```
    __ __    ___
   / // /__ / (_)_ __
  / _  / -_) / /\ \ /
 /_//_/\__/_/_//_\_\
 Developer friendly HTTP framework
```

**Helix** is a zero-dependency, high-performance HTTP web framework for Go with a focus on developer experience, type safety, and stdlib compatibility. Built by [Kolosys](https://github.com/kolosys) for enterprise-grade applications.

## Features

- **Zero Dependencies** - Built entirely on Go's standard library
- **High Performance** - Zero-allocation hot paths using `sync.Pool`
- **Type-Safe Handlers** - Generic handlers with automatic request binding and response encoding
- **RFC 7807 Problem Details** - Standardized error responses out of the box
- **Modular Architecture** - First-class support for organizing routes into modules
- **Fluent API** - Chainable context methods for clean handler code
- **Middleware Ecosystem** - Comprehensive built-in middleware suite
- **Dependency Injection** - Type-safe service registry with request-scoped support
- **Health Checks** - Built-in Kubernetes-ready liveness and readiness probes
- **Structured Logging** - High-performance logging with JSON and text formatters
- **Graceful Shutdown** - Context-aware shutdown with configurable grace period
- **stdlib Compatible** - Works with any `http.Handler` middleware

## Installation

```bash
go get github.com/kolosys/helix
```

Requires Go 1.24 or later.

## Quick Start

```go
package main

import (
    "context"
    "net/http"

    "github.com/kolosys/helix"
)

func main() {
    // Create server with default middleware (RequestID, Logger, Recover)
    s := helix.Default()

    // Simple handler
    s.GET("/", func(w http.ResponseWriter, r *http.Request) {
        helix.OK(w, map[string]string{"message": "Hello, World!"})
    })

    // Handler with Ctx for cleaner API
    s.GET("/hello", helix.HandleCtx(func(c *helix.Ctx) error {
        name := c.QueryDefault("name", "World")
        return c.OK(map[string]string{"message": "Hello, " + name + "!"})
    }))

    // Typed handler with automatic binding
    s.GET("/users/{id}", helix.Handle(func(ctx context.Context, req struct {
        ID int `path:"id"`
    }) (User, error) {
        return getUser(ctx, req.ID)
    }))

    s.Start(":8080")
}

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func getUser(ctx context.Context, id int) (User, error) {
    return User{ID: id, Name: "John Doe"}, nil
}
```

## Server Creation

### Basic Server

```go
s := helix.New()
```

### Default Server (with middleware)

Includes `RequestID`, `Logger` (dev format), and `Recover` middleware:

```go
s := helix.Default()
```

### Server with Options

```go
s := helix.New(
    helix.WithAddr(":3000"),
    helix.WithReadTimeout(30 * time.Second),
    helix.WithWriteTimeout(30 * time.Second),
    helix.WithIdleTimeout(120 * time.Second),
    helix.WithGracePeriod(30 * time.Second),
    helix.WithBasePath("/api/v1"),
    helix.WithTLS("cert.pem", "key.pem"),
    helix.WithErrorHandler(customErrorHandler),
    helix.HideBanner(),
)
```

## Routing

### HTTP Methods

```go
s.GET("/users", listUsers)
s.POST("/users", createUser)
s.PUT("/users/{id}", updateUser)
s.PATCH("/users/{id}", patchUser)
s.DELETE("/users/{id}", deleteUser)
s.HEAD("/users", headUsers)
s.OPTIONS("/users", optionsUsers)
s.Any("/echo", echoHandler)  // All methods
```

### Path Parameters

```go
// Single parameter
s.GET("/users/{id}", handler)      // Param(r, "id")

// Multiple parameters
s.GET("/users/{userId}/posts/{postId}", handler)

// Catch-all parameter
s.GET("/files/{path...}", handler) // Matches /files/a/b/c
```

### Static Files

```go
s.Static("/assets/", "./public")
```

## Handlers

Helix provides multiple handler types for different use cases:

### Standard Handler

Works with any `http.HandlerFunc`:

```go
s.GET("/", func(w http.ResponseWriter, r *http.Request) {
    helix.JSON(w, http.StatusOK, data)
})
```

### Context Handler (`HandleCtx`)

Uses the unified `Ctx` type for fluent API:

```go
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    // Access path params
    id := c.Param("id")

    // Access query params
    page := c.QueryInt("page", 1)
    search := c.QueryDefault("q", "")

    // Access headers
    auth := c.Header("Authorization")

    // Bind JSON body
    var input CreateUserInput
    if err := c.Bind(&input); err != nil {
        return c.BadRequest("invalid input")
    }

    // Return response
    return c.OK(users)
}))
```

### Typed Handler (`Handle`)

Generic handlers with automatic request binding and JSON response:

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
    // req is automatically bound from JSON body
    user, err := userService.Create(ctx, req.Name, req.Email)
    if err != nil {
        return User{}, err
    }
    return user, nil // Automatically encoded as JSON
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

// No request, no response
helix.HandleEmpty(func(ctx context.Context) error {
    return pingService(ctx)
})
```

## Request Binding

Bind request data to structs using struct tags:

```go
type UpdateUserRequest struct {
    // From path parameters
    ID int `path:"id"`

    // From query parameters
    Include string `query:"include"`

    // From headers
    APIKey string `header:"X-API-Key"`

    // From JSON body
    Name  string `json:"name"`
    Email string `json:"email"`

    // From form data
    Avatar string `form:"avatar"`
}
```

### Binding Functions

```go
// Bind all sources
user, err := helix.Bind[UpdateUserRequest](r)

// Bind specific sources
query, err := helix.BindQuery[QueryParams](r)
path, err := helix.BindPath[PathParams](r)
headers, err := helix.BindHeader[HeaderParams](r)
body, err := helix.BindJSON[CreateRequest](r)

// Bind and validate
user, err := helix.BindAndValidate[CreateUserRequest](r)
```

### Parameter Helpers

```go
// Path parameters
id := helix.Param(r, "id")
userID, err := helix.ParamInt(r, "id")
uuid, err := helix.ParamUUID(r, "id")

// Query parameters
name := helix.Query(r, "name")
name := helix.QueryDefault(r, "name", "default")
page := helix.QueryInt(r, "page", 1)
active := helix.QueryBool(r, "active")
tags := helix.QuerySlice(r, "tags")
price := helix.QueryFloat64(r, "price", 0.0)
```

## Response Helpers

### JSON Responses

```go
helix.JSON(w, http.StatusOK, data)
helix.JSONPretty(w, http.StatusOK, data, "  ")
helix.OK(w, data)           // 200
helix.Created(w, data)      // 201
helix.Accepted(w, data)     // 202
helix.NoContent(w)          // 204
```

### Other Content Types

```go
helix.Text(w, http.StatusOK, "Hello, World!")
helix.HTML(w, http.StatusOK, "<h1>Hello</h1>")
helix.Blob(w, http.StatusOK, "image/png", imageData)
helix.Stream(w, "application/octet-stream", reader)
helix.File(w, r, "/path/to/file")
```

### Error Responses

```go
helix.Error(w, http.StatusBadRequest, "invalid input")
helix.BadRequest(w, "invalid input")
helix.Unauthorized(w, "authentication required")
helix.Forbidden(w, "access denied")
helix.NotFound(w, "user not found")
helix.InternalServerError(w, "something went wrong")
```

### Content Disposition

```go
helix.Attachment(w, "report.pdf")  // Force download
helix.Inline(w, "image.png")       // Display inline
helix.Redirect(w, r, "/new-url", http.StatusFound)
```

## Problem Details (RFC 7807)

Helix uses [RFC 7807](https://tools.ietf.org/html/rfc7807) Problem Details for standardized error responses:

```go
// Return a problem from a handler
return helix.ErrNotFound.WithDetail("user 123 not found")

// Or create custom problems
return helix.NewProblem(
    http.StatusConflict,
    "duplicate_email",
    "Email Already Exists",
).WithDetail("The email address is already registered")
```

### Sentinel Errors

```go
helix.ErrBadRequest          // 400
helix.ErrUnauthorized        // 401
helix.ErrForbidden           // 403
helix.ErrNotFound            // 404
helix.ErrMethodNotAllowed    // 405
helix.ErrConflict            // 409
helix.ErrGone                // 410
helix.ErrUnprocessableEntity // 422
helix.ErrTooManyRequests     // 429
helix.ErrInternal            // 500
helix.ErrNotImplemented      // 501
helix.ErrBadGateway          // 502
helix.ErrServiceUnavailable  // 503
helix.ErrGatewayTimeout      // 504
```

### Convenience Functions

```go
return helix.NotFoundf("user %d not found", id)
return helix.BadRequestf("invalid email: %s", email)
return helix.Conflictf("username %q already taken", username)
```

### Response Format

```json
{
  "type": "about:blank#not_found",
  "title": "Not Found",
  "status": 404,
  "detail": "user 123 not found",
  "instance": "/users/123"
}
```

## Validation

Implement the `Validatable` interface for automatic validation:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
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
    if r.Age < 0 || r.Age > 150 {
        v.Addf("age", "age must be between 0 and 150, got %d", r.Age)
    }

    return v.Err() // Returns nil if no errors
}
```

Validation errors are returned as RFC 7807 with field-level details:

```json
{
  "type": "about:blank#unprocessable_entity",
  "title": "Unprocessable Entity",
  "status": 422,
  "detail": "One or more validation errors occurred",
  "instance": "/users",
  "errors": [
    { "field": "name", "message": "name is required" },
    { "field": "email", "message": "invalid email format" }
  ]
}
```

## Middleware

### Using Middleware

```go
// Global middleware
s.Use(middleware.RequestID())
s.Use(middleware.Logger(middleware.LogFormatDev))
s.Use(middleware.Recover())

// Works with any func(http.Handler) http.Handler
s.Use(thirdPartyMiddleware)
```

### Built-in Middleware

#### Request ID

```go
middleware.RequestID()  // Generates X-Request-ID header
```

#### Logger

```go
middleware.Logger(middleware.LogFormatDev)      // Colorized development
middleware.Logger(middleware.LogFormatJSON)     // JSON format
middleware.Logger(middleware.LogFormatCombined) // Apache combined
middleware.Logger(middleware.LogFormatCommon)   // Apache common
middleware.Logger(middleware.LogFormatShort)    // Short format
middleware.Logger(middleware.LogFormatTiny)     // Minimal format

// Custom format with tokens
middleware.LoggerWithFormat(":method :url :status :response-time")

// Advanced configuration
middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format:     middleware.LogFormatJSON,
    Output:     os.Stdout,
    Skip:       func(r *http.Request) bool { return r.URL.Path == "/health" },
    TimeFormat: time.RFC3339,
    Fields:     map[string]string{"api_version": "header:X-API-Version"},
})
```

#### Recover

```go
middleware.Recover()  // Recovers from panics, returns 500
```

#### CORS

```go
middleware.CORS()  // Default permissive config

middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"https://example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    ExposeHeaders:    []string{"X-Total-Count"},
    AllowCredentials: true,
    MaxAge:           86400,
})

middleware.CORSAllowAll()  // Allow everything (dev only)
```

#### Rate Limiting

```go
middleware.RateLimit(100, 10)  // 100 req/sec, burst of 10

middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Rate:     100,
    Burst:    10,
    KeyFunc:  func(r *http.Request) string { return r.Header.Get("X-API-Key") },
    Handler:  customRateLimitHandler,
    SkipFunc: func(r *http.Request) bool { return r.URL.Path == "/health" },
})
```

#### Basic Auth

```go
middleware.BasicAuth(map[string]string{
    "admin": "secret",
    "user":  "password",
})
```

#### Compression

```go
middleware.Compress()  // Gzip compression
```

#### Timeout

```go
middleware.Timeout(30 * time.Second)
```

#### ETag

```go
middleware.ETag()  // Automatic ETag generation
```

#### Cache

```go
middleware.Cache(time.Hour)  // HTTP cache headers
```

### Middleware Bundles

Pre-configured middleware sets for common scenarios:

```go
// API server (RequestID, Logger JSON, Recover, CORS)
for _, mw := range middleware.API() {
    s.Use(mw)
}

// Web application (RequestID, Logger Dev, Recover, Compress)
for _, mw := range middleware.Web() {
    s.Use(mw)
}

// Production (RequestID, Logger Combined, Recover)
for _, mw := range middleware.Production() {
    s.Use(mw)
}

// Development (same as helix.Default())
for _, mw := range middleware.Development() {
    s.Use(mw)
}

// Secure (RequestID, Logger JSON, Recover, RateLimit)
for _, mw := range middleware.Secure(100, 10) {
    s.Use(mw)
}

// Minimal (Recover only)
for _, mw := range middleware.Minimal() {
    s.Use(mw)
}
```

### Middleware Chain

```go
chain := middleware.Chain(
    middleware.RequestID(),
    middleware.Logger(middleware.LogFormatDev),
    middleware.Recover(),
)
s.Use(chain)
```

## Route Groups

Organize routes with shared prefixes and middleware:

```go
// Create a group
api := s.Group("/api/v1")
api.GET("/users", listUsers)
api.POST("/users", createUser)

// Group with middleware
admin := s.Group("/admin", authMiddleware, adminOnlyMiddleware)
admin.GET("/stats", getStats)
admin.DELETE("/users/{id}", deleteUser)

// Nested groups
v2 := api.Group("/v2")
v2.GET("/users", listUsersV2)

// Add middleware to group after creation
api.Use(rateLimitMiddleware)
```

## Modules

Modules provide a clean way to organize routes into separate files or packages:

```go
// Define a module
type UserModule struct {
    service *UserService
}

func (m *UserModule) Register(r helix.RouteRegistrar) {
    r.GET("/", m.list)
    r.POST("/", m.create)
    r.GET("/{id}", m.get)
    r.PUT("/{id}", m.update)
    r.DELETE("/{id}", m.delete)
}

// Mount the module
s.Mount("/users", &UserModule{service: userService})

// Mount with middleware
s.Mount("/users", &UserModule{}, authMiddleware)

// Mount using a function
s.MountFunc("/posts", func(r helix.RouteRegistrar) {
    r.GET("/", listPosts)
    r.POST("/", createPost)
})

// Mount within a group
api := s.Group("/api/v1")
api.Mount("/users", &UserModule{})
```

## Resources

REST resource builder for CRUD operations:

```go
// Fluent resource definition
s.Resource("/users").
    List(listUsers).      // GET /users
    Create(createUser).   // POST /users
    Get(getUser).         // GET /users/{id}
    Update(updateUser).   // PUT /users/{id}
    Patch(patchUser).     // PATCH /users/{id}
    Delete(deleteUser)    // DELETE /users/{id}

// All CRUD in one call
s.Resource("/posts").CRUD(listPosts, createPost, getPost, updatePost, deletePost)

// Read-only resource
s.Resource("/articles").ReadOnly(listArticles, getArticle)

// Custom actions
s.Resource("/users").
    Get(getUser).
    Custom("POST", "/{id}/activate", activateUser).
    Custom("POST", "/{id}/deactivate", deactivateUser)

// Resource with middleware
s.Resource("/admin/users", authMiddleware, adminMiddleware).
    CRUD(list, create, get, update, delete)

// Typed resources
helix.TypedResource[User](s, "/users").
    List(listHandler).
    Create(createHandler).
    Get(getHandler).
    Update(updateHandler).
    Delete(deleteHandler)
```

## Dependency Injection

Type-safe service registry with global and request-scoped support:

### Global Services

```go
// Register services at startup
userService := NewUserService(db)
emailService := NewEmailService(smtp)

helix.Register(userService)
helix.Register(emailService)

// Access in handlers
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    svc := helix.MustGet[*UserService]()
    users, err := svc.List(c.Context())
    if err != nil {
        return err
    }
    return c.OK(users)
}))

// Safe access (returns ok bool)
svc, ok := helix.Get[*UserService]()
```

### Request-Scoped Services

```go
// Add service to request context
ctx := helix.WithService(r.Context(), txn)

// Retrieve from context (falls back to global)
svc, ok := helix.FromContext[*Transaction](ctx)
svc := helix.MustFromContext[*Transaction](ctx)

// Middleware that provides request-scoped services
s.Use(helix.ProvideMiddleware(func(r *http.Request) *Transaction {
    return db.BeginTx(r.Context())
}))
```

## Pagination

Built-in pagination helpers:

### Pagination Type

```go
type ListUsersRequest struct {
    helix.Pagination  // Embeds Page, Limit, Sort, Order, Cursor
    Status string `query:"status"`
}

s.GET("/users", helix.Handle(func(ctx context.Context, req ListUsersRequest) (helix.PaginatedResponse[User], error) {
    page := req.GetPage()                        // Default: 1
    limit := req.GetLimit(20, 100)               // Default 20, max 100
    offset := req.GetOffset(limit)               // Calculate SQL offset
    sort := req.GetSort("created_at", []string{"created_at", "name"})
    order := req.GetOrder()                      // "asc" or "desc"

    users, total, err := userService.List(ctx, limit, offset, sort, order)
    if err != nil {
        return helix.PaginatedResponse[User]{}, err
    }

    return helix.NewPaginatedResponse(users, total, page, limit), nil
}))
```

### Ctx Pagination

```go
s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    p := c.BindPagination(20, 100)  // defaultLimit, maxLimit

    users, total, err := userService.List(c.Context(), p.GetPage(), p.GetLimit(20, 100))
    if err != nil {
        return err
    }

    return c.Paginated(users, total, p.GetPage(), p.GetLimit(20, 100))
}))
```

### Response Format

```json
{
  "items": [...],
  "total": 150,
  "page": 2,
  "limit": 20,
  "total_pages": 8,
  "has_more": true
}
```

## Health Checks

Kubernetes-ready health check endpoints:

### Comprehensive Health Check

```go
health := helix.Health().
    Version("1.0.0").
    Timeout(5 * time.Second).
    CheckFunc("database", func(ctx context.Context) error {
        return db.PingContext(ctx)
    }).
    CheckFunc("redis", func(ctx context.Context) error {
        return redis.Ping(ctx).Err()
    }).
    Check("external_api", func(ctx context.Context) helix.HealthCheckResult {
        start := time.Now()
        err := callExternalAPI(ctx)
        return helix.HealthCheckResult{
            Status:  helix.HealthStatusUp,
            Latency: time.Since(start),
            Details: map[string]any{"endpoint": "api.example.com"},
        }
    })

s.GET("/health", health.Handler())
```

### Simple Probes

```go
// Liveness probe (is the process running?)
s.GET("/health/live", helix.LivenessHandler())

// Readiness probe (is the service ready to accept traffic?)
s.GET("/health/ready", helix.ReadinessHandler(
    func(ctx context.Context) error { return db.PingContext(ctx) },
    func(ctx context.Context) error { return cache.Ping(ctx) },
))
```

### Health Response

```json
{
  "status": "up",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "components": {
    "database": {
      "status": "up",
      "latency_ms": 2
    },
    "redis": {
      "status": "up",
      "latency_ms": 1
    }
  }
}
```

## Logging

Helix includes a high-performance structured logging library:

```go
import "github.com/kolosys/helix/logs"

// Basic usage
logs.Info("server started", logs.Int("port", 8080))
logs.Error("request failed", logs.Err(err), logs.String("path", "/api"))

// Create custom logger
log := logs.New(
    logs.WithLevel(logs.DebugLevel),
    logs.WithFormatter(&logs.JSONFormatter{}),
    logs.WithCaller(),
    logs.WithAsync(1024),
)

// Child logger with fields
reqLog := log.With(logs.String("request_id", "abc123"))
reqLog.Info("processing request")

// Context-aware logging
log.InfoContext(ctx, "user authenticated", logs.String("user_id", "123"))
```

### Field Types

```go
logs.String("key", "value")
logs.Int("count", 42)
logs.Int64("id", 123456789)
logs.Float64("price", 99.99)
logs.Bool("active", true)
logs.Duration("latency", time.Since(start))
logs.Time("timestamp", time.Now())
logs.Err(err)
logs.Any("data", complexStruct)
logs.Stringer("ip", net.IP{...})
```

### Log Levels

```go
logs.TraceLevel  // Most verbose
logs.DebugLevel
logs.InfoLevel   // Default
logs.WarnLevel
logs.ErrorLevel
logs.FatalLevel  // Exits after logging
logs.PanicLevel  // Panics after logging
```

## Lifecycle Hooks

```go
s.OnStart(func(s *helix.Server) {
    log.Printf("Server starting on %s", s.Addr())
})

s.OnStop(func(ctx context.Context, s *helix.Server) {
    log.Println("Server shutting down...")
    db.Close()
})
```

## Route Introspection

```go
// Get all registered routes
routes := s.Routes()
for _, r := range routes {
    fmt.Printf("%s %s\n", r.Method, r.Pattern)
}

// Print routes to writer
s.PrintRoutes(os.Stdout)
```

## Configuration Options

| Option                  | Description                           | Default    |
| ----------------------- | ------------------------------------- | ---------- |
| `WithAddr(addr)`        | Server listen address                 | `:8080`    |
| `WithReadTimeout(d)`    | Maximum duration for reading request  | `30s`      |
| `WithWriteTimeout(d)`   | Maximum duration for writing response | `30s`      |
| `WithIdleTimeout(d)`    | Maximum time to wait for next request | `120s`     |
| `WithGracePeriod(d)`    | Shutdown grace period                 | `30s`      |
| `WithMaxHeaderBytes(n)` | Maximum size of request headers       | Go default |
| `WithBasePath(path)`    | Base path prefix for all routes       | `""`       |
| `WithTLS(cert, key)`    | Enable TLS with certificate files     | Disabled   |
| `WithTLSConfig(cfg)`    | Custom TLS configuration              | `nil`      |
| `WithErrorHandler(h)`   | Custom error handler                  | RFC 7807   |
| `HideBanner()`          | Hide startup banner                   | Shown      |
| `WithCustomBanner(b)`   | Custom startup banner                 | Default    |

## Examples

See the [examples](./examples) directory for complete working examples:

- **[basic](./examples/basic)** - Simple handlers and routing
- **[crud](./examples/crud)** - CRUD operations with typed handlers
- **[groups](./examples/groups)** - Route grouping and nested groups
- **[middleware](./examples/middleware)** - Middleware usage patterns
- **[modular](./examples/modular)** - Modular architecture with DI
- **[resource](./examples/resource)** - REST resource builder
- **[validation](./examples/validation)** - Request validation

## Performance

Helix is designed for high performance:

- **Zero allocations** in hot paths using `sync.Pool` for contexts and parameters
- **Pre-compiled middleware chains** via `s.Build()`
- **Radix tree router** for efficient route matching
- **Buffer pooling** for JSON encoding
- **Minimal reflection** - binding info is cached

### Pre-compile for Production

```go
// Call Build() after registering all routes and middleware
s.Build()

// Then start the server
s.Start(":8080")
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Helix is released under the [MIT License](LICENSE).

---

Built with ❤️ by [Kolosys](https://github.com/kolosys)
