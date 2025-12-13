# Helix Examples

This directory contains examples demonstrating various features of the Helix web framework.

## Examples

### [basic](./basic)

The simplest example showing how to create a server, register routes, and handle requests.

**Features demonstrated:**

- Creating a server with `helix.Default()`
- Basic routing with `GET`, `POST`, etc.
- Path parameters
- Query parameters
- Lifecycle hooks (`OnStart`, `OnStop`)
- Error responses with RFC 7807 Problem Details

```bash
cd basic && go run main.go
```

### [crud](./crud)

A complete CRUD API example using typed handlers.

**Features demonstrated:**

- Typed request/response handlers with `helix.Handle[Req, Res]`
- `HandleCreated` for 201 responses
- Automatic request binding from path, query, and JSON body
- Field-level validation with `ValidationErrors`
- Different handler types (`Handle`, `HandleCreated`, `HandleNoResponse`)
- Thread-safe in-memory storage

```bash
cd crud && go run main.go
```

### [middleware](./middleware)

Demonstrates the use of various middleware.

**Features demonstrated:**

- Global middleware (`Use`) - no type casting needed!
- Route group middleware - accepts middleware directly
- RequestID middleware
- Logger middleware
- Recover middleware
- CORS middleware
- Compression middleware
- Basic authentication
- Rate limiting
- Request timeouts

```bash
cd middleware && go run main.go
```

### [groups](./groups)

Route grouping and API versioning example.

**Features demonstrated:**

- Route groups with common prefixes
- Nested groups
- API versioning (v1, v2)
- Group-level middleware (no casting required)
- `HandleCreated` for POST endpoints
- Different response types per version

```bash
cd groups && go run main.go
```

### [resource](./resource)

RESTful resource builder pattern for CRUD operations.

**Features demonstrated:**

- `Resource()` fluent builder
- Standard CRUD methods (`List`, `Create`, `Get`, `Update`, `Delete`)
- Custom resource actions (`publish`, `unpublish`)
- Nested resources (comments under articles)
- Resources within groups

```bash
cd resource && go run main.go
```

### [validation](./validation)

Request binding and validation examples.

**Features demonstrated:**

- Multi-source binding (path, query, header, JSON body)
- Struct tag binding (`path`, `query`, `header`, `json`)
- Field-level validation with `ValidationErrors`
- RFC 7807 responses with `errors` array for validation failures
- Type conversion (string, int, int64, float64, bool, []string)
- Default values and range validation
- Manual binding with `Ctx` methods

```bash
cd validation && go run main.go
```

### [modular](./modular)

Advanced example demonstrating modular architecture and dependency injection.

**Features demonstrated:**

- Module pattern for organizing routes (`Module` interface)
- `Mount()` and `MountFunc()` for mounting modules
- Service registration with `helix.Register[T]()`
- Service retrieval with `helix.MustGet[T]()`
- Built-in pagination with `Pagination` struct and `BindPagination()`
- `c.Paginated()` helper for paginated responses
- Health check builder with liveness/readiness probes
- Pre-compiled middleware chain with `s.Build()`

```bash
cd modular && go run main.go
```

## Running Examples

Each example can be run directly:

```bash
# From the examples directory
cd <example-name>
go run main.go
```

The server will start on `:8080` by default.

## Testing with curl

```bash
# Basic hello
curl http://localhost:8080/

# Get with query parameter
curl http://localhost:8080/hello?name=Helix

# Path parameters
curl http://localhost:8080/users/123

# Create resource (POST with JSON)
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}'

# Protected endpoint with Basic Auth
curl -u admin:secret http://localhost:8080/admin/dashboard

# Search with query parameters
curl "http://localhost:8080/search?q=widget&page=1&limit=10&tags=electronics"
```
