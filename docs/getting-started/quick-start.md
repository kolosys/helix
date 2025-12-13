# Quick Start

This guide will help you get started with helix quickly with a basic example.

## Basic Usage

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/helix"
    "github.com/kolosys/helix/logs"
    "github.com/kolosys/helix/middleware"
)

func main() {
    // Basic usage example
    fmt.Println("Welcome to helix!")
    
    // TODO: Add your code here
}
```

## Common Use Cases

### Using helix

**Import Path:** `github.com/kolosys/helix`

Package helix provides a zero-dependency, context-aware, high-performance
HTTP web framework for Go with stdlib compatibility.


```go
package main

import (
    "fmt"
    "github.com/kolosys/helix"
)

func main() {
    // Example usage of helix
    fmt.Println("Using helix package")
}
```

#### Available Types
- **Ctx** - Ctx provides a unified context for HTTP handlers with fluent accessors for request data and response methods.
- **CtxHandler** - ----------------------------------------------------------------------------- CtxHandler and HandleCtx ----------------------------------------------------------------------------- CtxHandler is a handler function that uses the unified Ctx type.
- **EmptyHandler** - EmptyHandler is a handler that takes no request and returns no response.
- **ErrorHandler** - ErrorHandler is a function that handles errors from handlers. It receives the response writer, request, and error, and is responsible for writing an appropriate error response.
- **FieldError** - ----------------------------------------------------------------------------- Validation Errors ----------------------------------------------------------------------------- FieldError represents a validation error for a specific field.
- **Group** - Group represents a group of routes with a common prefix and middleware.
- **Handler** - Handler is a generic handler function that accepts a typed request and returns a typed response. The request type is automatically bound from path parameters, query parameters, headers, and JSON body. The response is automatically encoded as JSON.
- **HealthBuilder** - HealthBuilder provides a fluent interface for building health check endpoints.
- **HealthCheck** - HealthCheck is a function that checks the health of a component.
- **HealthCheckResult** - HealthCheckResult contains the result of a health check.
- **HealthResponse** - HealthResponse is the response returned by the health endpoint.
- **HealthStatus** - HealthStatus represents the health status of a component.
- **IDRequest** - IDRequest is a common request type for single-entity operations.
- **ListRequest** - ListRequest is a common request type for list operations.
- **ListResponse** - ListResponse wraps a list of entities with pagination metadata.
- **Middleware** - Middleware is a function that wraps an http.Handler to provide additional functionality. This is an alias to middleware.Middleware for convenience.
- **Module** - } func (m *UserModule) Register(r RouteRegistrar) { r.GET("/", m.list) r.POST("/", m.create) r.GET("/{id}", m.get) } // Mount the module s.Mount("/users", &UserModule{store: store})
- **ModuleFunc** - ModuleFunc is a function that implements Module.
- **NoRequestHandler** - NoRequestHandler is a handler that takes no request body, only context.
- **NoResponseHandler** - NoResponseHandler is a handler that returns no response body.
- **Option** - Option configures a Server.
- **PaginatedResponse** - PaginatedResponse wraps a list response with pagination metadata.
- **Pagination** - Pagination contains common pagination parameters. Use with struct embedding for automatic binding. Example: type ListUsersRequest struct { helix.Pagination Status string `query:"status"` }
- **Problem** - Problem represents an RFC 7807 Problem Details for HTTP APIs. See: https://tools.ietf.org/html/rfc7807
- **ResourceBuilder** - ResourceBuilder provides a fluent interface for defining REST resource routes.
- **RouteInfo** - RouteInfo contains information about a registered route.
- **RouteRegistrar** - RouteRegistrar is an interface for registering routes. Both Server and Group implement this interface.
- **Router** - Router handles HTTP request routing.
- **Server** - Server is the main HTTP server for the Helix framework.
- **TypedResourceBuilder** - TypedResourceBuilder provides a fluent interface for defining typed REST resource routes.
- **Validatable** - Validatable is an interface for types that can validate themselves.
- **ValidationErrors** - ValidationErrors collects multiple validation errors for RFC 7807 response. Implements the error interface and can be returned from Validate() methods.
- **ValidationProblem** - ValidationProblem is an RFC 7807 Problem with validation errors extension.

#### Available Functions
- **Accepted** - Accepted writes a 202 Accepted JSON response.
- **Attachment** - Attachment sets the Content-Disposition header to attachment with the given filename.
- **BadRequest** - BadRequest writes a 400 Bad Request error response.
- **Bind** - Bind binds path parameters, query parameters, headers, and JSON body to a struct. The binding sources are determined by struct tags: - `path:"name"` - binds from URL path parameters - `query:"name"` - binds from URL query parameters - `header:"name"` - binds from HTTP headers - `json:"name"` - binds from JSON body - `form:"name"` - binds from form data
- **BindAndValidate** - BindAndValidate binds and validates a request. If the bound type implements Validatable, Validate() is called after binding.
- **BindHeader** - BindHeader binds HTTP headers to a struct. Uses the `header` struct tag to determine field names.
- **BindJSON** - BindJSON binds the JSON request body to a struct.
- **BindPath** - BindPath binds URL path parameters to a struct. Uses the `path` struct tag to determine field names.
- **BindQuery** - BindQuery binds URL query parameters to a struct. Uses the `query` struct tag to determine field names.
- **Blob** - Blob writes binary data with the given content type.
- **Created** - Created writes a 201 Created JSON response.
- **Error** - Error writes an error response with the given status code and message.
- **File** - File serves a file with the given content type.
- **Forbidden** - Forbidden writes a 403 Forbidden error response.
- **FromContext** - FromContext retrieves a service from the context or falls back to global registry.
- **Get** - Get retrieves a service from the global registry. Returns the zero value and false if not found.
- **HTML** - HTML writes an HTML response with the given status code.
- **Handle** - Handle wraps a generic Handler into an http.HandlerFunc. It automatically: - Binds the request to the Req type - Calls the handler with the context and request - Encodes the response as JSON - Handles errors using RFC 7807 Problem Details
- **HandleAccepted** - HandleAccepted wraps a generic Handler into an http.HandlerFunc that returns 202 Accepted. Useful for async operations where processing happens in the background.
- **HandleCreated** - HandleCreated wraps a generic Handler into an http.HandlerFunc that returns 201 Created. This is a convenience wrapper for HandleWithStatus(http.StatusCreated, h).
- **HandleCtx** - HandleCtx wraps a CtxHandler into an http.HandlerFunc. Errors returned from the handler are automatically converted to RFC 7807 responses.
- **HandleEmpty** - HandleEmpty wraps an EmptyHandler into an http.HandlerFunc. Returns 204 No Content on success.
- **HandleErrorDefault** - HandleErrorDefault provides the default error handling logic. This can be called from custom error handlers to fall back to default behavior.
- **HandleNoRequest** - HandleNoRequest wraps a NoRequestHandler into an http.HandlerFunc. Useful for endpoints that don't need request binding (e.g., GET /users).
- **HandleNoResponse** - HandleNoResponse wraps a NoResponseHandler into an http.HandlerFunc. Returns 204 No Content on success.
- **HandleWithStatus** - HandleWithStatus wraps a generic Handler into an http.HandlerFunc with a custom success status code.
- **Inline** - Inline sets the Content-Disposition header to inline with the given filename.
- **InternalServerError** - InternalServerError writes a 500 Internal Server Error response.
- **JSON** - JSON writes a JSON response with the given status code. Uses a pooled buffer for zero-allocation in the hot path.
- **JSONPretty** - JSONPretty writes a pretty-printed JSON response with the given status code.
- **LivenessHandler** - LivenessHandler returns a simple liveness probe handler. Returns 200 OK if the server is running.
- **MustFromContext** - MustFromContext retrieves a service from context or panics.
- **MustGet** - MustGet retrieves a service from the global registry or panics.
- **NoContent** - NoContent writes a 204 No Content response.
- **NotFound** - NotFound writes a 404 Not Found error response.
- **OK** - OK writes a 200 OK JSON response.
- **Param** - Param returns the value of a path parameter. Returns an empty string if the parameter does not exist.
- **ParamInt** - ParamInt returns the value of a path parameter as an int. Returns an error if the parameter does not exist or cannot be parsed.
- **ParamInt64** - ParamInt64 returns the value of a path parameter as an int64. Returns an error if the parameter does not exist or cannot be parsed.
- **ParamUUID** - ParamUUID returns the value of a path parameter validated as a UUID. Returns an error if the parameter does not exist or is not a valid UUID format.
- **Query** - Query returns the first value of a query parameter. Returns an empty string if the parameter does not exist.
- **QueryBool** - QueryBool returns the first value of a query parameter as a bool. Returns false if the parameter does not exist or cannot be parsed. Accepts "1", "t", "T", "true", "TRUE", "True" as true. Accepts "0", "f", "F", "false", "FALSE", "False" as false.
- **QueryDefault** - QueryDefault returns the first value of a query parameter or a default value.
- **QueryFloat64** - QueryFloat64 returns the first value of a query parameter as a float64. Returns the default value if the parameter does not exist or cannot be parsed.
- **QueryInt** - QueryInt returns the first value of a query parameter as an int. Returns the default value if the parameter does not exist or cannot be parsed.
- **QueryInt64** - QueryInt64 returns the first value of a query parameter as an int64. Returns the default value if the parameter does not exist or cannot be parsed.
- **QuerySlice** - QuerySlice returns all values of a query parameter as a string slice. Returns nil if the parameter does not exist.
- **ReadinessHandler** - ReadinessHandler returns a simple readiness probe handler. Uses the provided checks to determine readiness.
- **Redirect** - Redirect redirects the request to the given URL.
- **Register** - Register registers a service in the global registry by its type.
- **Stream** - Stream streams the content from the reader to the response.
- **Text** - Text writes a plain text response with the given status code.
- **Unauthorized** - Unauthorized writes a 401 Unauthorized error response.
- **UnmarshalJSON** - UnmarshalJSON unmarshals a JSON request body into a struct. Returns an error if the request body is not a valid JSON or if the struct cannot be unmarshalled. The error will contain the stack trace of the error.
- **WithService** - WithService returns a new context with the service added. This is useful for request-scoped services like database transactions.
- **WriteProblem** - WriteProblem writes a Problem response to the http.ResponseWriter.
- **WriteValidationProblem** - WriteValidationProblem writes a ValidationProblem response to the http.ResponseWriter.

For detailed API documentation, see the [helix API Reference](../api-reference/helix.md).

### Using logs

**Import Path:** `github.com/kolosys/helix/logs`

Package logs provides a high-performance, context-aware structured logging library.

Features:
  - Zero-allocation hot paths using sync.Pool
  - Context-aware logging with context.Context
  - Type-safe field builders
  - Multiple output formats (text, JSON, pretty)
  - Sampling for high-volume logs
  - Async logging option
  - Hook system for extensibility
  - Built-in caller information
  - Chained/fluent API

Basic usage:

	log := logs.New()
	log.Info("server started", logs.Int("port", 8080))

With context:

	log.InfoContext(ctx, "request processed", logs.Duration("latency", time.Since(start)))


```go
package main

import (
    "fmt"
    "github.com/kolosys/helix/logs"
)

func main() {
    // Example usage of logs
    fmt.Println("Using logs package")
}
```

#### Available Types
- **AlwaysSampler** - AlwaysSampler always allows logging.
- **Builder** - Builder provides a fluent/chainable API for building log entries. It accumulates fields and then emits a log entry when a level method is called.
- **CompositeSampler** - CompositeSampler combines multiple samplers with AND logic.
- **CountSampler** - CountSampler logs every Nth occurrence.
- **Entry** - Entry represents a log entry.
- **ErrorBuilder** - ErrorBuilder provides a fluent API for logging errors.
- **ErrorHook** - ErrorHook collects errors for inspection.
- **Field** - Field represents a structured log field.
- **FieldType** - FieldType represents the type of a field value.
- **FileHook** - FileHook writes entries to a file.
- **FilterHook** - FilterHook conditionally fires another hook.
- **FirstNSampler** - FirstNSampler logs only the first N occurrences.
- **Formatter** - Formatter formats log entries.
- **FuncHook** - FuncHook wraps a function as a hook.
- **Hook** - Hook is called when a log entry is written.
- **JSONFormatter** - JSONFormatter formats logs as JSON.
- **Level** - Level represents a log level.
- **LevelHook** - LevelHook fires only for specific levels.
- **LevelSampler** - LevelSampler applies different samplers per level.
- **Logger** - Logger is the main logging interface.
- **MetricsHook** - MetricsHook tracks log counts by level.
- **NamedFormatter** - NamedFormatter is a formatter that includes the logger name in output.
- **NeverSampler** - NeverSampler never allows logging.
- **NoopFormatter** - NoopFormatter discards all output.
- **OncePerSampler** - OncePerSampler logs a message only once per duration.
- **Option** - Option configures a Logger.
- **PrettyFormatter** - PrettyFormatter formats logs with colors and alignment for development.
- **RandomSampler** - RandomSampler samples a percentage of logs.
- **RateSampler** - RateSampler limits logs to a certain rate per message.
- **Sampler** - Sampler determines if a log entry should be emitted.
- **TextFormatter** - TextFormatter formats logs as text.
- **WriterHook** - WriterHook writes entries to an io.Writer.

#### Available Functions
- **CtxDebug** - CtxDebug logs at debug level using the logger from context.
- **CtxError** - CtxError logs at error level using the logger from context.
- **CtxInfo** - CtxInfo logs at info level using the logger from context.
- **CtxTrace** - CtxTrace logs at trace level using the logger from context.
- **CtxWarn** - CtxWarn logs at warn level using the logger from context.
- **Debug** - Debug logs at debug level using the default logger.
- **Debugf** - Debugf logs a formatted message at debug level.
- **Error** - Error logs at error level using the default logger.
- **Errorf** - Errorf logs a formatted message at error level.
- **Fatal** - Fatal logs at fatal level using the default logger and exits.
- **Fatalf** - Fatalf logs a formatted message at fatal level and exits.
- **Info** - Info logs at info level using the default logger.
- **Infof** - Infof logs a formatted message at info level.
- **Must** - Must logs and panics if error is not nil. Useful for initialization code. db := log.Must(sql.Open("postgres", dsn))
- **Panic** - Panic logs at panic level using the default logger and panics.
- **Panicf** - Panicf logs a formatted message at panic level and panics.
- **Print** - Print logs a message at info level.
- **Printf** - Printf logs a formatted message at info level.
- **Println** - Println logs a message at info level.
- **SetDefault** - SetDefault sets the default logger.
- **SetDefaultFormatter** - SetDefaultFormatter sets the default formatter.
- **SetDefaultLevel** - SetDefaultLevel sets the default level.
- **Trace** - Trace logs at trace level using the default logger.
- **Tracef** - Tracef logs a formatted message at trace level.
- **Warn** - Warn logs at warn level using the default logger.
- **Warnf** - Warnf logs a formatted message at warn level.
- **WithContextFields** - WithFields adds fields to the context that will be included in all logs.
- **WithLogger** - WithLogger attaches a logger to the context.
- **WithRequestID** - WithRequestID adds a request ID to the context.
- **WithTraceID** - WithTraceID adds a trace ID to the context.
- **WithUserID** - WithUserID adds a user ID to the context.

For detailed API documentation, see the [logs API Reference](../api-reference/logs.md).

### Using middleware

**Import Path:** `github.com/kolosys/helix/middleware`

Package middleware provides HTTP middleware for the Helix framework.


```go
package main

import (
    "fmt"
    "github.com/kolosys/helix/middleware"
)

func main() {
    // Example usage of middleware
    fmt.Println("Using middleware package")
}
```

#### Available Types
- **BasicAuthConfig** - BasicAuthConfig configures the BasicAuth middleware.
- **CORSConfig** - CORSConfig configures the CORS middleware.
- **CacheConfig** - CacheConfig configures the Cache middleware.
- **CompressConfig** - CompressConfig configures the Compress middleware.
- **ETagConfig** - ETagConfig configures the ETag middleware.
- **LogEntry** - LogEntry represents a JSON log entry.
- **LogFormat** - LogFormat represents a predefined log format.
- **LoggerConfig** - LoggerConfig configures the Logger middleware.
- **Middleware** - Middleware is a function that wraps an http.Handler to provide additional functionality.
- **RateLimitConfig** - RateLimitConfig configures the RateLimit middleware.
- **RecoverConfig** - RecoverConfig configures the Recover middleware.
- **RequestIDConfig** - RequestIDConfig configures the RequestID middleware.
- **TimeoutConfig** - TimeoutConfig configures the Timeout middleware.
- **TokenExtractor** - TokenExtractor is a function that extracts a value from the request. It receives the request and the captured request body (if body capture is enabled).

#### Available Functions
- **ETagFromContent** - ETagFromContent generates an ETag from content.
- **ETagFromString** - ETagFromString generates an ETag from a string.
- **ETagFromVersion** - ETagFromVersion generates an ETag from a version number.
- **GetRequestID** - GetRequestID retrieves the request ID from the context. Returns an empty string if no request ID is set.
- **GetRequestIDFromRequest** - GetRequestIDFromRequest retrieves the request ID from the request context.
- **SetCacheControl** - SetCacheControl sets the Cache-Control header on the response.
- **SetExpires** - SetExpires sets the Expires header on the response.
- **SetLastModified** - SetLastModified sets the Last-Modified header on the response.

For detailed API documentation, see the [middleware API Reference](../api-reference/middleware.md).

## Step-by-Step Tutorial

### Step 1: Import the Package

First, import the necessary packages in your Go file:

```go
import (
    "fmt"
    "github.com/kolosys/helix"
    "github.com/kolosys/helix/logs"
    "github.com/kolosys/helix/middleware"
)
```

### Step 2: Initialize

Set up the basic configuration:

```go
func main() {
    // Initialize your application
    fmt.Println("Initializing helix...")
}
```

### Step 3: Use the Library

Implement your specific use case:

```go
func main() {
    // Your implementation here
}
```

## Running Your Code

To run your Go program:

```bash
go run main.go
```

To build an executable:

```bash
go build -o myapp
./myapp
```

## Configuration Options

helix can be configured to suit your needs. Check the [Core Concepts](../core-concepts/) section for detailed information about configuration options.

## Error Handling

Always handle errors appropriately:

```go
result, err := someFunction()
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

## Best Practices

- Always handle errors returned by library functions
- Check the API documentation for detailed parameter information
- Use meaningful variable and function names
- Add comments to document your code

## Complete Example

Here's a complete working example:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/helix"
    "github.com/kolosys/helix/logs"
    "github.com/kolosys/helix/middleware"
)

func main() {
    fmt.Println("Starting helix application...")
    
    // Add your implementation here
    
    fmt.Println("Application completed successfully!")
}
```

## Next Steps

Now that you've seen the basics, explore:

- **[Core Concepts](../core-concepts/)** - Understanding the library architecture
- **[API Reference](../api-reference/)** - Complete API documentation
- **[Examples](../examples/README.md)** - More detailed examples
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced patterns

## Getting Help

If you run into issues:

1. Check the [API Reference](../api-reference/)
2. Browse the [Examples](../examples/README.md)
3. Visit the [GitHub Issues](https://github.com/kolosys/helix/issues) page

