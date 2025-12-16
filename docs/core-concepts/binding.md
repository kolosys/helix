# Request Binding

> **Note**: This is a developer-maintained documentation page. The content here is not auto-generated and should be updated manually to explain the core concepts and architecture of request binding in Helix.

## About Request Binding

Request binding in Helix provides a type-safe, declarative way to extract and validate data from HTTP requests. Using struct tags, you can bind data from multiple sources (path parameters, query parameters, headers, JSON body, form data) into a single struct with automatic type conversion.

## Core Concepts

### Binding Sources

Helix supports binding from five different sources:

- **Path Parameters** (`path:"name"`) - URL path segments like `/users/{id}`
- **Query Parameters** (`query:"name"`) - URL query string like `?page=1&limit=10`
- **Headers** (`header:"name"`) - HTTP request headers
- **JSON Body** (`json:"name"`) - Request body as JSON
- **Form Data** (`form:"name"`) - URL-encoded or multipart form data

### Struct Tags

Binding is controlled through struct tags that specify the source and parameter name:

```go
type UpdateUserRequest struct {
    // Path parameter: /users/{id}
    ID int `path:"id"`

    // Query parameter: ?include=profile,posts
    Include string `query:"include"`

    // Header: X-API-Key: abc123
    APIKey string `header:"X-API-Key"`

    // JSON body field
    Name  string `json:"name"`
    Email string `json:"email"`

    // Form data field
    Avatar string `form:"avatar"`
}
```

### Required Fields

Mark fields as required by adding `required` to the tag options:

```go
type CreateUserRequest struct {
    Name  string `json:"name" required:"true"`
    Email string `json:"email" required:"true"`
    Age   int    `json:"age"`  // Optional
}
```

When a required field is missing, binding returns `ErrRequiredField`.

## Usage Patterns

### Basic Binding

Use `Bind` to bind all sources at once:

```go
s.PUT("/users/{id}", helix.HandleCtx(func(c *helix.Ctx) error {
    var req UpdateUserRequest
    if err := helix.Bind[UpdateUserRequest](c.Request); err != nil {
        return c.BadRequest("invalid request")
    }

    // Use req.ID, req.Name, req.Email, etc.
    return c.OK(req)
}))
```

### Source-Specific Binding

Bind from a single source when you only need specific data:

```go
// Bind only query parameters
params, err := helix.BindQuery[QueryParams](r)

// Bind only path parameters
pathParams, err := helix.BindPath[PathParams](r)

// Bind only headers
headers, err := helix.BindHeader[HeaderParams](r)

// Bind only JSON body
body, err := helix.BindJSON[CreateRequest](r)
```

### Binding with Validation

Use `BindAndValidate` to automatically validate after binding:

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

    return v.Err()
}

// In handler
req, err := helix.BindAndValidate[CreateUserRequest](r)
if err != nil {
    // Returns ValidationErrors as RFC 7807 Problem
    return err
}
```

### Using Ctx for Binding

The `Ctx` type provides a convenient `Bind` method:

```go
s.POST("/users", helix.HandleCtx(func(c *helix.Ctx) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.BadRequest("invalid input")
    }

    // Process request...
    return c.Created(user)
}))
```

## Supported Types

Binding automatically converts string values to these Go types:

- `string` - Direct assignment
- `int`, `int8`, `int16`, `int32`, `int64` - Parsed as base-10 integers
- `uint`, `uint8`, `uint16`, `uint32`, `uint64` - Parsed as unsigned integers
- `float32`, `float64` - Parsed as floating-point numbers
- `bool` - Parsed as boolean (also accepts "yes", "no", "on", "off")
- `[]string` - Comma-separated values for query parameters
- Pointer types (`*string`, `*int`, etc.) - Optional fields

### Type Conversion Examples

```go
type FilterRequest struct {
    // Query: ?page=2
    Page int `query:"page"`

    // Query: ?active=true
    Active bool `query:"active"`

    // Query: ?tags=go,rust,typescript
    Tags []string `query:"tags"`

    // Query: ?price=99.99
    Price float64 `query:"price"`

    // Optional query parameter
    Sort *string `query:"sort"`
}
```

## Design Decisions

### Performance Considerations

- **Caching**: Struct reflection information is cached using `sync.Map` to avoid repeated reflection overhead
- **Zero Allocations**: Binding uses value semantics and avoids unnecessary allocations in hot paths
- **Lazy JSON Parsing**: JSON body is only parsed if there are `json:` tagged fields

### API Design Choices

- **Generic Functions**: `Bind[T]()` provides type safety at compile time
- **Multiple Sources**: Single struct can bind from multiple sources simultaneously
- **Validation Integration**: `Validatable` interface allows seamless validation after binding

### Error Handling

Binding errors are specific and actionable:

- `ErrBindingFailed` - General binding failure
- `ErrRequiredField` - Required field missing
- `ErrInvalidJSON` - Invalid JSON body
- `ErrInvalidFieldValue` - Type conversion failed
- `ErrUnsupportedType` - Field type not supported

## Common Pitfalls

### Pitfall 1: Reading Body Multiple Times

**Problem**: HTTP request bodies can only be read once. If you read the body manually before binding, binding will fail.

**Solution**: Let binding handle the body, or use `io.NopCloser` to wrap a pre-read body:

```go
// ❌ Wrong - body already read
body, _ := io.ReadAll(r.Body)
var req CreateRequest
err := helix.BindJSON[CreateRequest](r) // Fails!

// ✅ Correct - let binding read the body
var req CreateRequest
err := helix.BindJSON[CreateRequest](r)

// ✅ Alternative - if you need to read body manually
body, _ := io.ReadAll(r.Body)
r.Body = io.NopCloser(bytes.NewReader(body))
var req CreateRequest
err := json.Unmarshal(body, &req)
```

### Pitfall 2: Case Sensitivity

**Problem**: Header names and query parameter names are case-sensitive in HTTP, but Go struct field names are case-insensitive for JSON.

**Solution**: Always specify explicit names in tags:

```go
// ❌ Wrong - relies on default JSON field name
type Request struct {
    APIKey string `json:"apiKey"` // Won't match "X-API-Key" header
}

// ✅ Correct - explicit header name
type Request struct {
    APIKey string `header:"X-API-Key"`
}
```

### Pitfall 3: Pointer vs Value Types

**Problem**: Using value types for optional fields makes it impossible to distinguish between "not provided" and "zero value".

**Solution**: Use pointers for optional fields:

```go
// ❌ Wrong - can't tell if age was 0 or not provided
type Request struct {
    Age int `query:"age"`
}

// ✅ Correct - nil means not provided
type Request struct {
    Age *int `query:"age"`
}
```

## Integration Guide

### With Typed Handlers

Binding works seamlessly with typed handlers:

```go
s.POST("/users", helix.Handle(func(ctx context.Context, req CreateUserRequest) (User, error) {
    // req is automatically bound from JSON body
    return userService.Create(ctx, req)
}))
```

### With Validation

Combine binding with the `Validatable` interface:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (r *CreateUserRequest) Validate() error {
    // Validation logic
}

// Automatic validation
req, err := helix.BindAndValidate[CreateUserRequest](r)
```

### With Middleware

Binding can be used in middleware to extract common request data:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var auth AuthHeader
        if err := helix.BindHeader[AuthHeader](r); err != nil {
            helix.Unauthorized(w, "missing authorization")
            return
        }

        // Validate auth and add to context
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Further Reading

- [Context Guide](./context.md) - Using Ctx for request handling
- [Error Handling](./error-handling.md) - Handling binding errors
- [Validation Examples](../examples/validation/main.md) - Complete validation examples
- [API Reference](../api-reference/helix.md) - Complete binding API documentation

---

_This documentation should be updated by package maintainers to reflect the actual architecture and design patterns used._
