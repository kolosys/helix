# Routing

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain the core concepts and architecture of routing in Helix.

## About Routing

Helix provides a high-performance, flexible routing system built on a radix tree. It supports path parameters, route groups, modules, and resources for organizing complex APIs.

## Core Concepts

### Route Registration

Routes are registered using HTTP method helpers:

```go
s.GET("/users", listUsers)
s.POST("/users", createUser)
s.PUT("/users/{id}", updateUser)
s.PATCH("/users/{id}", patchUser)
s.DELETE("/users/{id}", deleteUser)
s.HEAD("/users", headUsers)
s.OPTIONS("/users", optionsUsers)
s.Any("/echo", echoHandler)  // All HTTP methods
```

### Path Parameters

Extract dynamic segments from URLs:

```go
// Single parameter
s.GET("/users/{id}", handler)  // /users/123 -> id="123"

// Multiple parameters
s.GET("/users/{userId}/posts/{postId}", handler)

// Catch-all parameter (must be last)
s.GET("/files/{path...}", handler)  // /files/a/b/c -> path="a/b/c"
```

### Route Matching

Helix uses a radix tree for efficient route matching:

- **Static Segments**: Exact matches (fastest)
- **Parameters**: `{name}` matches any single segment
- **Catch-All**: `{name...}` matches remaining path

Routes are matched in order of specificity (most specific first).

## Architecture Overview

### Router Structure

The router uses a radix tree (compressed trie) for O(k) route matching where k is the path length:

```
Root
├── users (static)
│   ├── / (GET handler)
│   └── {id} (param)
│       └── / (GET handler)
└── posts (static)
    └── {id} (param)
        └── /comments (static)
            └── / (GET handler)
```

### Thread Safety

The router is thread-safe:

- **Read Operations**: Concurrent route matching with `sync.RWMutex`
- **Write Operations**: Per-method locks reduce contention
- **Route Introspection**: Safe concurrent access to route list

### Performance

- **Zero Allocations**: Path parameters use `sync.Pool` for reuse
- **Fast Matching**: Radix tree provides O(k) matching
- **Pre-compilation**: Routes can be compiled before serving

## Core Concepts

### Route Groups

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

// Add middleware to existing group
api.Use(rateLimitMiddleware)
```

### Modules

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
```

### Resources

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
```

### Static Files

Serve static files and directories:

```go
// Serve directory
s.Static("/assets/", "./public")

// Serves /assets/css/style.css from ./public/css/style.css
```

## Usage Patterns

### Basic Routing

```go
s := helix.New()

// Simple routes
s.GET("/", homeHandler)
s.GET("/about", aboutHandler)
s.POST("/contact", contactHandler)
```

### RESTful API

```go
api := s.Group("/api/v1")

// Users resource
api.Resource("/users").CRUD(
    listUsers,
    createUser,
    getUser,
    updateUser,
    deleteUser,
)

// Posts resource
api.Resource("/posts").CRUD(
    listPosts,
    createPost,
    getPost,
    updatePost,
    deletePost,
)
```

### Modular Architecture

Organize large APIs into modules:

```go
// users/module.go
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

// main.go
s.Mount("/api/v1/users", &UserModule{service: userService})
s.Mount("/api/v1/posts", &PostModule{service: postService})
s.Mount("/api/v1/comments", &CommentModule{service: commentService})
```

### Route Organization

```go
// Public routes
s.GET("/", homeHandler)
s.GET("/about", aboutHandler)

// API routes
api := s.Group("/api/v1", apiMiddleware)
api.GET("/health", healthHandler)

// Authenticated routes
auth := api.Group("", authMiddleware)
auth.Resource("/users").CRUD(...)

// Admin routes
admin := auth.Group("/admin", adminMiddleware)
admin.GET("/stats", statsHandler)
admin.Resource("/users").CRUD(...)
```

## Design Decisions

### Radix Tree Router

Helix uses a radix tree instead of a simple map because:

- **Performance**: O(k) matching vs O(n) for linear search
- **Memory**: Compressed paths reduce memory usage
- **Flexibility**: Supports parameters and catch-all patterns
- **Scalability**: Handles large route tables efficiently

### Per-Method Trees

Routes are organized by HTTP method because:

- **Performance**: Smaller trees per method = faster matching
- **Concurrency**: Per-method locks reduce contention
- **Clarity**: Method-specific routing logic is simpler

### Route Groups

Groups provide shared prefixes and middleware because:

- **DRY**: Avoid repeating prefixes and middleware
- **Organization**: Logical grouping of related routes
- **Flexibility**: Can nest groups and apply middleware selectively

### Modules Interface

Modules use an interface (`RouteRegistrar`) because:

- **Flexibility**: Can mount modules on servers or groups
- **Testability**: Easy to test modules in isolation
- **Composition**: Modules can mount other modules

## Common Pitfalls

### Pitfall 1: Route Order Matters

**Problem**: More specific routes should be registered before catch-all routes.

**Solution**: Register specific routes first:

```go
// ❌ Wrong - catch-all matches first
s.GET("/files/{path...}", catchAllHandler)
s.GET("/files/users", usersHandler)  // Never matches

// ✅ Correct - specific first
s.GET("/files/users", usersHandler)
s.GET("/files/{path...}", catchAllHandler)
```

### Pitfall 2: Parameter Names Must Match

**Problem**: Parameter names in routes must match the names used in handlers.

**Solution**: Use consistent naming:

```go
// Route
s.GET("/users/{id}", handler)

// Handler
s.GET("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    id := c.Param("id")  // Must match route parameter name
    // ...
}))
```

### Pitfall 3: Base Path Applied After Routes

**Problem**: `WithBasePath` must be set before registering routes.

**Solution**: Set base path during server creation:

```go
// ❌ Wrong
s := helix.New()
s.GET("/users", handler)
s.WithBasePath("/api")  // Too late!

// ✅ Correct
s := helix.New(helix.WithBasePath("/api"))
s.GET("/users", handler)  // Route is /api/users
```

## Integration Guide

### With Middleware

Apply middleware to routes, groups, or resources:

```go
// Global middleware
s.Use(middleware.RequestID())
s.Use(middleware.Logger(middleware.LogFormatJSON))

// Group middleware
api := s.Group("/api", authMiddleware)
api.GET("/users", listUsers)  // Has auth middleware

// Resource middleware
s.Resource("/users", authMiddleware, rateLimitMiddleware).
    List(listUsers)
```

### With Dependency Injection

Modules can access services via dependency injection:

```go
type UserModule struct {
    service *UserService
    logger  *logs.Logger
}

func (m *UserModule) Register(r helix.RouteRegistrar) {
    r.GET("/", m.list)
}

func (m *UserModule) list(c *helix.Ctx) error {
    users, err := m.service.List(c.Context())
    if err != nil {
        m.logger.ErrorContext(c.Context(), "failed to list users", logs.Err(err))
        return err
    }
    return c.OK(users)
}

// Mount with dependencies
s.Mount("/users", &UserModule{
    service: userService,
    logger:  logger,
})
```

### With Validation

Combine routing with request validation:

```go
s.POST("/users", helix.Handle(func(ctx context.Context, req CreateUserRequest) (User, error) {
    // req is automatically bound and validated
    return userService.Create(ctx, req)
}))
```

## Further Reading

- [Binding Guide](./binding.md) - Path parameter binding
- [Context Guide](./context.md) - Accessing route parameters
- [Modules Guide](../examples/modular/main.md) - Modular architecture examples
- [API Reference](../api-reference/helix.md) - Complete routing API

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
