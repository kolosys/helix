# helix API

Complete API documentation for the helix package.

**Import Path:** `github.com/kolosys/helix`

## Package Documentation

Package helix provides a zero-dependency, context-aware, high-performance
HTTP web framework for Go with stdlib compatibility.


## Constants

**MIMETextPlain, MIMETextHTML, MIMETextCSS, MIMETextCSV, MIMETextJavaScript, MIMETextXML, MIMETextPlainCharsetUTF8, MIMETextHTMLCharsetUTF8, MIMETextCSSCharsetUTF8, MIMETextCSVCharsetUTF8, MIMETextJavaScriptCharsetUTF8, MIMETextXMLCharsetUTF8, MIMEApplicationJSON, MIMEApplicationXML, MIMEApplicationJavaScript, MIMEApplicationXHTMLXML, MIMEApplicationJSONCharsetUTF8, MIMEApplicationXMLCharsetUTF8, MIMEApplicationJavaScriptCharsetUTF8, MIMEApplicationProblemJSON, MIMEApplicationForm, MIMEApplicationProtobuf, MIMEApplicationMsgPack, MIMEApplicationOctetStream, MIMEApplicationPDF, MIMEApplicationZip, MIMEApplicationGzip, MIMEMultipartForm, MIMEImagePNG, MIMEImageSVG, MIMEImageJPEG, MIMEImageGIF, MIMEImageWebP, MIMEImageICO, MIMEImageAVIF, MIMEAudioMPEG, MIMEAudioWAV, MIMEAudioOGG, MIMEVideoMP4, MIMEVideoWebM, MIMEVideoOGG**

MIME type constants for HTTP Content-Type headers.
Base types (without charset) are used for content-type detection/matching.
CharsetUTF8 variants are used for setting response headers.


```go
const MIMETextPlain = "text/plain"	// Text types - base (for matching)

const MIMETextHTML = "text/html"
const MIMETextCSS = "text/css"
const MIMETextCSV = "text/csv"
const MIMETextJavaScript = "text/javascript"
const MIMETextXML = "text/xml"
const MIMETextPlainCharsetUTF8 = "text/plain; charset=utf-8"	// Text types - with charset (for responses)

const MIMETextHTMLCharsetUTF8 = "text/html; charset=utf-8"
const MIMETextCSSCharsetUTF8 = "text/css; charset=utf-8"
const MIMETextCSVCharsetUTF8 = "text/csv; charset=utf-8"
const MIMETextJavaScriptCharsetUTF8 = "text/javascript; charset=utf-8"
const MIMETextXMLCharsetUTF8 = "text/xml; charset=utf-8"
const MIMEApplicationJSON = "application/json"	// Application types - base (for matching)

const MIMEApplicationXML = "application/xml"
const MIMEApplicationJavaScript = "application/javascript"
const MIMEApplicationXHTMLXML = "application/xhtml+xml"
const MIMEApplicationJSONCharsetUTF8 = "application/json; charset=utf-8"	// Application types - with charset (for responses)

const MIMEApplicationXMLCharsetUTF8 = "application/xml; charset=utf-8"
const MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"
const MIMEApplicationProblemJSON = "application/problem+json"	// Application types - no charset needed

const MIMEApplicationForm = "application/x-www-form-urlencoded"
const MIMEApplicationProtobuf = "application/x-protobuf"
const MIMEApplicationMsgPack = "application/msgpack"
const MIMEApplicationOctetStream = "application/octet-stream"
const MIMEApplicationPDF = "application/pdf"
const MIMEApplicationZip = "application/zip"
const MIMEApplicationGzip = "application/gzip"
const MIMEMultipartForm = "multipart/form-data"
const MIMEImagePNG = "image/png"	// Image types

const MIMEImageSVG = "image/svg+xml"
const MIMEImageJPEG = "image/jpeg"
const MIMEImageGIF = "image/gif"
const MIMEImageWebP = "image/webp"
const MIMEImageICO = "image/x-icon"
const MIMEImageAVIF = "image/avif"
const MIMEAudioMPEG = "audio/mpeg"	// Audio types

const MIMEAudioWAV = "audio/wav"
const MIMEAudioOGG = "audio/ogg"
const MIMEVideoMP4 = "video/mp4"	// Video types

const MIMEVideoWebM = "video/webm"
const MIMEVideoOGG = "video/ogg"
```

**Version**



```go
const Version = "0.1.0"	// Version of Hexix

```

## Variables

**ErrBindingFailed, ErrUnsupportedType, ErrInvalidJSON, ErrRequiredField, ErrBodyAlreadyRead, ErrInvalidFieldValue**

Binding errors


```go
var ErrBindingFailed = errors.New("helix: binding failed")
var ErrUnsupportedType = errors.New("helix: unsupported type for binding")
var ErrInvalidJSON = errors.New("helix: invalid JSON body")
var ErrRequiredField = errors.New("helix: required field missing")
var ErrBodyAlreadyRead = errors.New("helix: request body already read")
var ErrInvalidFieldValue = errors.New("helix: invalid field value")
```

**ErrBadRequest, ErrUnauthorized, ErrForbidden, ErrNotFound, ErrMethodNotAllowed, ErrConflict, ErrGone, ErrUnprocessableEntity, ErrTooManyRequests, ErrInternal, ErrNotImplemented, ErrBadGateway, ErrServiceUnavailable, ErrGatewayTimeout**

Sentinel errors for common HTTP error responses.


```go
var ErrBadRequest = NewProblem(http.StatusBadRequest, "bad_request", "Bad Request")	// ErrBadRequest represents a 400 Bad Request error.

var ErrUnauthorized = NewProblem(http.StatusUnauthorized, "unauthorized", "Unauthorized")	// ErrUnauthorized represents a 401 Unauthorized error.

var ErrForbidden = NewProblem(http.StatusForbidden, "forbidden", "Forbidden")	// ErrForbidden represents a 403 Forbidden error.

var ErrNotFound = NewProblem(http.StatusNotFound, "not_found", "Not Found")	// ErrNotFound represents a 404 Not Found error.

var ErrMethodNotAllowed = NewProblem(http.StatusMethodNotAllowed, "method_not_allowed", "Method Not Allowed")	// ErrMethodNotAllowed represents a 405 Method Not Allowed error.

var ErrConflict = NewProblem(http.StatusConflict, "conflict", "Conflict")	// ErrConflict represents a 409 Conflict error.

var ErrGone = NewProblem(http.StatusGone, "gone", "Gone")	// ErrGone represents a 410 Gone error.

var ErrUnprocessableEntity = NewProblem(http.StatusUnprocessableEntity, "unprocessable_entity", "Unprocessable Entity")	// ErrUnprocessableEntity represents a 422 Unprocessable Entity error.

var ErrTooManyRequests = NewProblem(http.StatusTooManyRequests, "too_many_requests", "Too Many Requests")	// ErrTooManyRequests represents a 429 Too Many Requests error.

var ErrInternal = NewProblem(http.StatusInternalServerError, "internal_error", "Internal Server Error")	// ErrInternal represents a 500 Internal Server Error.

var ErrNotImplemented = NewProblem(http.StatusNotImplemented, "not_implemented", "Not Implemented")	// ErrNotImplemented represents a 501 Not Implemented error.

var ErrBadGateway = NewProblem(http.StatusBadGateway, "bad_gateway", "Bad Gateway")	// ErrBadGateway represents a 502 Bad Gateway error.

var ErrServiceUnavailable = NewProblem(http.StatusServiceUnavailable, "service_unavailable", "Service Unavailable")	// ErrServiceUnavailable represents a 503 Service Unavailable error.

var ErrGatewayTimeout = NewProblem(http.StatusGatewayTimeout, "gateway_timeout", "Gateway Timeout")	// ErrGatewayTimeout represents a 504 Gateway Timeout error.

```

## Types

### Ctx
Ctx provides a unified context for HTTP handlers with fluent accessors for request data and response methods.

#### Example Usage

```go
// Create a new Ctx
ctx := Ctx{
    Request: &/* value */{},
    Response: /* value */,
}
```

#### Type Definition

```go
type Ctx struct {
    Request *http.Request
    Response http.ResponseWriter
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Request | `*http.Request` |  |
| Response | `http.ResponseWriter` |  |

### Constructor Functions

### NewCtx

NewCtx creates a new Ctx from an http.Request and http.ResponseWriter.

```go
func NewCtx(w http.ResponseWriter, r *http.Request) *Ctx
```

**Parameters:**
- `w` (http.ResponseWriter)
- `r` (*http.Request)

**Returns:**
- *Ctx

## Methods

### Accepted

Accepted writes a 202 Accepted JSON response.

```go
func Accepted(w http.ResponseWriter, v any) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `v` (any)

**Returns:**
- error

### AddHeader

AddHeader adds a response header value and returns the Ctx for chaining.

```go
func (*Ctx) AddHeader(key, value string) *Ctx
```

**Parameters:**
- `key` (string)
- `value` (string)

**Returns:**
- *Ctx

### Attachment

Attachment sets the Content-Disposition header to attachment.

```go
func Attachment(w http.ResponseWriter, filename string)
```

**Parameters:**
- `w` (http.ResponseWriter)
- `filename` (string)

**Returns:**
  None

### BadRequest

BadRequest writes a 400 Bad Request error response.

```go
func (*Ctx) BadRequest(message string) error
```

**Parameters:**
- `message` (string)

**Returns:**
- error

### Bind

Bind binds the request body to the given struct using JSON decoding.

```go
func (*Ctx) Bind(v any) error
```

**Parameters:**
- `v` (any)

**Returns:**
- error

### BindJSON

BindJSON is an alias for Bind.

```go
func BindJSON(r *http.Request) (T, error)
```

**Parameters:**
- `r` (*http.Request)

**Returns:**
- T
- error

### BindPagination

BindPaginationCtx extracts pagination from the Ctx with defaults.

```go
func (*Ctx) BindPagination(defaultLimit, maxLimit int) Pagination
```

**Parameters:**
- `defaultLimit` (int)
- `maxLimit` (int)

**Returns:**
- Pagination

### Blob

Blob writes binary data with the given content type.

```go
func (*Ctx) Blob(status int, contentType string, data []byte) error
```

**Parameters:**
- `status` (int)
- `contentType` (string)
- `data` ([]byte)

**Returns:**
- error

### Context

Context returns the request's context.Context.

```go
func (*Ctx) Context() context.Context
```

**Parameters:**
  None

**Returns:**
- context.Context

### Created

Created writes a 201 Created JSON response.

```go
func Created(w http.ResponseWriter, v any) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `v` (any)

**Returns:**
- error

### CreatedMessage

CreatedMessage writes a 201 Created response with a message and ID.

```go
func (*Ctx) CreatedMessage(message string, id any) error
```

**Parameters:**
- `message` (string)
- `id` (any)

**Returns:**
- error

### DeletedMessage

DeletedMessage writes a 200 OK response indicating deletion.

```go
func (*Ctx) DeletedMessage(message string) error
```

**Parameters:**
- `message` (string)

**Returns:**
- error

### File

File serves a file.

```go
func File(w http.ResponseWriter, r *http.Request, path string)
```

**Parameters:**
- `w` (http.ResponseWriter)
- `r` (*http.Request)
- `path` (string)

**Returns:**
  None

### Forbidden

Forbidden writes a 403 Forbidden error response.

```go
func Forbidden(w http.ResponseWriter, message string) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `message` (string)

**Returns:**
- error

### Get

Get retrieves a value from the request-scoped store.

```go
func (**ast.IndexExpr) Get(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### GetInt

GetInt retrieves an int value from the request-scoped store.

```go
func (*Ctx) GetInt(key string) int
```

**Parameters:**
- `key` (string)

**Returns:**
- int

### GetString

GetString retrieves a string value from the request-scoped store.

```go
func (*Ctx) GetString(key string) string
```

**Parameters:**
- `key` (string)

**Returns:**
- string

### HTML

HTML writes an HTML response with the given status code.

```go
func HTML(w http.ResponseWriter, status int, html string) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `status` (int)
- `html` (string)

**Returns:**
- error

### Header

Header returns the value of a request header.

```go
func (*Ctx) Header(name string) string
```

**Parameters:**
- `name` (string)

**Returns:**
- string

### Inline

Inline sets the Content-Disposition header to inline.

```go
func (*Ctx) Inline(filename string) *Ctx
```

**Parameters:**
- `filename` (string)

**Returns:**
- *Ctx

### InternalServerError

InternalServerError writes a 500 Internal Server Error response.

```go
func InternalServerError(w http.ResponseWriter, message string) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `message` (string)

**Returns:**
- error

### JSON

JSON writes a JSON response with the given status code.

```go
func JSON(w http.ResponseWriter, status int, v any) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `status` (int)
- `v` (any)

**Returns:**
- error

### MustGet

MustGet retrieves a value from the request-scoped store or panics if not found.

```go
func (*Ctx) MustGet(key string) any
```

**Parameters:**
- `key` (string)

**Returns:**
- any

### NoContent

NoContent writes a 204 No Content response.

```go
func (*Ctx) NoContent() error
```

**Parameters:**
  None

**Returns:**
- error

### NotFound

NotFound writes a 404 Not Found error response.

```go
func NotFound(w http.ResponseWriter, message string) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `message` (string)

**Returns:**
- error

### OK

OK writes a 200 OK JSON response.

```go
func OK(w http.ResponseWriter, v any) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `v` (any)

**Returns:**
- error

### OKMessage

OKMessage writes a 200 OK response with a message.

```go
func (*Ctx) OKMessage(message string) error
```

**Parameters:**
- `message` (string)

**Returns:**
- error

### Paginated

Paginated writes a paginated JSON response.

```go
func (*Ctx) Paginated(items any, total, page, limit int) error
```

**Parameters:**
- `items` (any)
- `total` (int)
- `page` (int)
- `limit` (int)

**Returns:**
- error

### Param

Param returns the value of a path parameter.

```go
func Param(r *http.Request, name string) string
```

**Parameters:**
- `r` (*http.Request)
- `name` (string)

**Returns:**
- string

### ParamInt

ParamInt returns the value of a path parameter as an int.

```go
func ParamInt(r *http.Request, name string) (int, error)
```

**Parameters:**
- `r` (*http.Request)
- `name` (string)

**Returns:**
- int
- error

### ParamInt64

ParamInt64 returns the value of a path parameter as an int64.

```go
func (*Ctx) ParamInt64(name string) (int64, error)
```

**Parameters:**
- `name` (string)

**Returns:**
- int64
- error

### ParamUUID

ParamUUID returns the value of a path parameter validated as a UUID.

```go
func (*Ctx) ParamUUID(name string) (string, error)
```

**Parameters:**
- `name` (string)

**Returns:**
- string
- error

### Problem

Problem writes an RFC 7807 Problem response.

```go
func (*Ctx) Problem(p Problem) error
```

**Parameters:**
- `p` (Problem)

**Returns:**
- error

### Query

Query returns the first value of a query parameter.

```go
func Query(r *http.Request, name string) string
```

**Parameters:**
- `r` (*http.Request)
- `name` (string)

**Returns:**
- string

### QueryBool

QueryBool returns the first value of a query parameter as a bool.

```go
func QueryBool(r *http.Request, name string) bool
```

**Parameters:**
- `r` (*http.Request)
- `name` (string)

**Returns:**
- bool

### QueryDefault

QueryDefault returns the first value of a query parameter or a default value.

```go
func (*Ctx) QueryDefault(name, defaultVal string) string
```

**Parameters:**
- `name` (string)
- `defaultVal` (string)

**Returns:**
- string

### QueryFloat64

QueryFloat64 returns the first value of a query parameter as a float64.

```go
func QueryFloat64(r *http.Request, name string, defaultVal float64) float64
```

**Parameters:**
- `r` (*http.Request)
- `name` (string)
- `defaultVal` (float64)

**Returns:**
- float64

### QueryInt

QueryInt returns the first value of a query parameter as an int.

```go
func (*Ctx) QueryInt(name string, defaultVal int) int
```

**Parameters:**
- `name` (string)
- `defaultVal` (int)

**Returns:**
- int

### QueryInt64

QueryInt64 returns the first value of a query parameter as an int64.

```go
func QueryInt64(r *http.Request, name string, defaultVal int64) int64
```

**Parameters:**
- `r` (*http.Request)
- `name` (string)
- `defaultVal` (int64)

**Returns:**
- int64

### QuerySlice

QuerySlice returns all values of a query parameter as a string slice.

```go
func (*Ctx) QuerySlice(name string) []string
```

**Parameters:**
- `name` (string)

**Returns:**
- []string

### Redirect

Redirect redirects the request to the given URL.

```go
func (*Ctx) Redirect(url string, code int)
```

**Parameters:**
- `url` (string)
- `code` (int)

**Returns:**
  None

### Reset

Reset resets the Ctx for reuse from a pool.

```go
func (*Ctx) Reset(w http.ResponseWriter, r *http.Request)
```

**Parameters:**
- `w` (http.ResponseWriter)
- `r` (*http.Request)

**Returns:**
  None

### SendJSON

SendJSON writes a JSON response with the pending status code (or 200 if not set).

```go
func (*Ctx) SendJSON(v any) error
```

**Parameters:**
- `v` (any)

**Returns:**
- error

### Set

Set stores a value in the request-scoped store.

```go
func (*Ctx) Set(key string, value any)
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
  None

### SetCookie

SetCookie sets a cookie on the response and returns the Ctx for chaining.

```go
func (*Ctx) SetCookie(cookie *http.Cookie) *Ctx
```

**Parameters:**
- `cookie` (*http.Cookie)

**Returns:**
- *Ctx

### SetHeader

SetHeader sets a response header and returns the Ctx for chaining.

```go
func (*Ctx) SetHeader(key, value string) *Ctx
```

**Parameters:**
- `key` (string)
- `value` (string)

**Returns:**
- *Ctx

### Status

Status sets the pending status code for the response and returns the Ctx for chaining. The status is applied when a response body is written.

```go
func (*Ctx) Status(code int) *Ctx
```

**Parameters:**
- `code` (int)

**Returns:**
- *Ctx

### Text

Text writes a plain text response with the given status code.

```go
func Text(w http.ResponseWriter, status int, text string) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `status` (int)
- `text` (string)

**Returns:**
- error

### Unauthorized

Unauthorized writes a 401 Unauthorized error response.

```go
func Unauthorized(w http.ResponseWriter, message string) error
```

**Parameters:**
- `w` (http.ResponseWriter)
- `message` (string)

**Returns:**
- error

### CtxHandler
----------------------------------------------------------------------------- CtxHandler and HandleCtx ----------------------------------------------------------------------------- CtxHandler is a handler function that uses the unified Ctx type.

#### Example Usage

```go
// Example usage of CtxHandler
var value CtxHandler
// Initialize with appropriate value
```

#### Type Definition

```go
type CtxHandler func(c *Ctx) error
```

### EmptyHandler
EmptyHandler is a handler that takes no request and returns no response.

#### Example Usage

```go
// Example usage of EmptyHandler
var value EmptyHandler
// Initialize with appropriate value
```

#### Type Definition

```go
type EmptyHandler func(ctx context.Context) error
```

### ErrorHandler
ErrorHandler is a function that handles errors from handlers. It receives the response writer, request, and error, and is responsible for writing an appropriate error response.

#### Example Usage

```go
// Example usage of ErrorHandler
var value ErrorHandler
// Initialize with appropriate value
```

#### Type Definition

```go
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)
```

### FieldError
----------------------------------------------------------------------------- Validation Errors ----------------------------------------------------------------------------- FieldError represents a validation error for a specific field.

#### Example Usage

```go
// Create a new FieldError
fielderror := FieldError{
    Field: "example",
    Message: "example",
}
```

#### Type Definition

```go
type FieldError struct {
    Field string `json:"field"`
    Message string `json:"message"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Field | `string` |  |
| Message | `string` |  |

### Group
Group represents a group of routes with a common prefix and middleware.

#### Example Usage

```go
// Create a new Group
group := Group{

}
```

#### Type Definition

```go
type Group struct {
}
```

## Methods

### Any

Any registers a handler for all HTTP methods.

```go
func (*Group) Any(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### DELETE

DELETE registers a handler for DELETE requests.

```go
func (*Group) DELETE(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### GET

GET registers a handler for GET requests.

```go
func (*Server) GET(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Group

Group creates a nested group with the given prefix. The prefix is appended to the parent group's prefix. Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.

```go
func (*Group) Group(prefix string, mw ...any) *Group
```

**Parameters:**
- `prefix` (string)
- `mw` (...any)

**Returns:**
- *Group

### HEAD

HEAD registers a handler for HEAD requests.

```go
func (*Server) HEAD(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Handle

Handle registers a handler for the given method and pattern.

```go
func (*Group) Handle(method, pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `method` (string)
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Mount

Mount mounts a module at the given prefix within a group.

```go
func (*Group) Mount(prefix string, m Module, mw ...any)
```

**Parameters:**
- `prefix` (string)
- `m` (Module)
- `mw` (...any)

**Returns:**
  None

### MountFunc

MountFunc mounts a function as a module at the given prefix within a group.

```go
func (*Group) MountFunc(prefix string, fn func(r RouteRegistrar), mw ...any)
```

**Parameters:**
- `prefix` (string)
- `fn` (func(r RouteRegistrar))
- `mw` (...any)

**Returns:**
  None

### OPTIONS

OPTIONS registers a handler for OPTIONS requests.

```go
func (*Group) OPTIONS(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### PATCH

PATCH registers a handler for PATCH requests.

```go
func (*Group) PATCH(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### POST

POST registers a handler for POST requests.

```go
func (*Group) POST(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### PUT

PUT registers a handler for PUT requests.

```go
func (*Server) PUT(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Resource

Resource creates a new ResourceBuilder for the given pattern within this group. The pattern is relative to the group's prefix. Optional middleware can be applied to all routes in the resource. Accepts helix.Middleware, middleware.Middleware, or func(http.Handler) http.Handler.

```go
func (*Server) Resource(pattern string, mw ...any) *ResourceBuilder
```

**Parameters:**
- `pattern` (string)
- `mw` (...any)

**Returns:**
- *ResourceBuilder

### Static

Static serves static files from the given file system root.

```go
func (*Group) Static(pattern, root string)
```

**Parameters:**
- `pattern` (string)
- `root` (string)

**Returns:**
  None

### Use

Use adds middleware to the group. Middleware is applied to all routes registered on this group. Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.

```go
func (*Server) Use(mw ...any)
```

**Parameters:**
- `mw` (...any)

**Returns:**
  None

### Handler
Handler is a generic handler function that accepts a typed request and returns a typed response. The request type is automatically bound from path parameters, query parameters, headers, and JSON body. The response is automatically encoded as JSON.

#### Example Usage

```go
// Example usage of Handler
var value Handler
// Initialize with appropriate value
```

#### Type Definition

```go
type Handler func(ctx context.Context, req Req) (Res, error)
```

### HealthBuilder
HealthBuilder provides a fluent interface for building health check endpoints.

#### Example Usage

```go
// Create a new HealthBuilder
healthbuilder := HealthBuilder{

}
```

#### Type Definition

```go
type HealthBuilder struct {
}
```

### Constructor Functions

### Health

Health creates a new HealthBuilder.

```go
func Health() *HealthBuilder
```

**Parameters:**
  None

**Returns:**
- *HealthBuilder

## Methods

### Check

Check adds a health check for a named component.

```go
func (*HealthBuilder) Check(name string, check HealthCheck) *HealthBuilder
```

**Parameters:**
- `name` (string)
- `check` (HealthCheck)

**Returns:**
- *HealthBuilder

### CheckFunc

CheckFunc adds a simple health check that returns an error.

```go
func (*HealthBuilder) CheckFunc(name string, check func(ctx context.Context) error) *HealthBuilder
```

**Parameters:**
- `name` (string)
- `check` (func(ctx context.Context) error)

**Returns:**
- *HealthBuilder

### Handler

Handler returns an http.HandlerFunc for the health check endpoint.

```go
func (*HealthBuilder) Handler() http.HandlerFunc
```

**Parameters:**
  None

**Returns:**
- http.HandlerFunc

### Timeout

Timeout sets the timeout for health checks.

```go
func (*HealthBuilder) Timeout(d time.Duration) *HealthBuilder
```

**Parameters:**
- `d` (time.Duration)

**Returns:**
- *HealthBuilder

### Version

Version sets the application version shown in health responses.

```go
func (*HealthBuilder) Version(v string) *HealthBuilder
```

**Parameters:**
- `v` (string)

**Returns:**
- *HealthBuilder

### HealthCheck
HealthCheck is a function that checks the health of a component.

#### Example Usage

```go
// Example usage of HealthCheck
var value HealthCheck
// Initialize with appropriate value
```

#### Type Definition

```go
type HealthCheck func(ctx context.Context) HealthCheckResult
```

### HealthCheckResult
HealthCheckResult contains the result of a health check.

#### Example Usage

```go
// Create a new HealthCheckResult
healthcheckresult := HealthCheckResult{
    Status: HealthStatus{},
    Message: "example",
    Latency: /* value */,
    Details: map[],
}
```

#### Type Definition

```go
type HealthCheckResult struct {
    Status HealthStatus `json:"status"`
    Message string `json:"message,omitempty"`
    Latency time.Duration `json:"latency_ms,omitempty"`
    Details map[string]any `json:"details,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Status | `HealthStatus` |  |
| Message | `string` |  |
| Latency | `time.Duration` |  |
| Details | `map[string]any` |  |

### HealthResponse
HealthResponse is the response returned by the health endpoint.

#### Example Usage

```go
// Create a new HealthResponse
healthresponse := HealthResponse{
    Status: HealthStatus{},
    Timestamp: /* value */,
    Version: "example",
    Components: map[],
}
```

#### Type Definition

```go
type HealthResponse struct {
    Status HealthStatus `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Version string `json:"version,omitempty"`
    Components map[string]HealthCheckResult `json:"components,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Status | `HealthStatus` |  |
| Timestamp | `time.Time` |  |
| Version | `string` |  |
| Components | `map[string]HealthCheckResult` |  |

### HealthStatus
HealthStatus represents the health status of a component.

#### Example Usage

```go
// Example usage of HealthStatus
var value HealthStatus
// Initialize with appropriate value
```

#### Type Definition

```go
type HealthStatus string
```

### IDRequest
IDRequest is a common request type for single-entity operations.

#### Example Usage

```go
// Create a new IDRequest
idrequest := IDRequest{
    ID: 42,
}
```

#### Type Definition

```go
type IDRequest struct {
    ID int `path:"id"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `int` |  |

### ListRequest
ListRequest is a common request type for list operations.

#### Example Usage

```go
// Create a new ListRequest
listrequest := ListRequest{
    Page: 42,
    Limit: 42,
    Sort: "example",
    Order: "example",
    Search: "example",
}
```

#### Type Definition

```go
type ListRequest struct {
    Page int `query:"page"`
    Limit int `query:"limit"`
    Sort string `query:"sort"`
    Order string `query:"order"`
    Search string `query:"search"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Page | `int` |  |
| Limit | `int` |  |
| Sort | `string` |  |
| Order | `string` |  |
| Search | `string` |  |

### ListResponse
ListResponse wraps a list of entities with pagination metadata.

#### Example Usage

```go
// Create a new ListResponse
listresponse := ListResponse{
    Items: [],
    Total: 42,
    Page: 42,
    Limit: 42,
}
```

#### Type Definition

```go
type ListResponse struct {
    Items []Entity `json:"items"`
    Total int `json:"total"`
    Page int `json:"page,omitempty"`
    Limit int `json:"limit,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Items | `[]Entity` |  |
| Total | `int` |  |
| Page | `int` |  |
| Limit | `int` |  |

### Middleware
Middleware is a function that wraps an http.Handler to provide additional functionality. This is an alias to middleware.Middleware for convenience.

#### Example Usage

```go
// Example usage of Middleware
var value Middleware
// Initialize with appropriate value
```

#### Type Definition

```go
type Middleware middleware.Middleware
```

### Constructor Functions

### ProvideMiddleware

ProvideMiddleware returns middleware that injects services into the request context. Services added this way are request-scoped.

```go
func ProvideMiddleware(factory func(r *http.Request) T) Middleware
```

**Parameters:**
- `factory` (func(r *http.Request) T)

**Returns:**
- Middleware

### Module
} func (m *UserModule) Register(r RouteRegistrar) { r.GET("/", m.list) r.POST("/", m.create) r.GET("/{id}", m.get) } // Mount the module s.Mount("/users", &UserModule{store: store})

#### Example Usage

```go
// Example implementation of Module
type MyModule struct {
    // Add your fields here
}

func (m MyModule) Register(param1 RouteRegistrar)  {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Module interface {
    Register(r RouteRegistrar)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### ModuleFunc
ModuleFunc is a function that implements Module.

#### Example Usage

```go
// Example usage of ModuleFunc
var value ModuleFunc
// Initialize with appropriate value
```

#### Type Definition

```go
type ModuleFunc func(r RouteRegistrar)
```

## Methods

### Register

Register implements Module.

```go
func Register(service T)
```

**Parameters:**
- `service` (T)

**Returns:**
  None

### NoRequestHandler
NoRequestHandler is a handler that takes no request body, only context.

#### Example Usage

```go
// Example usage of NoRequestHandler
var value NoRequestHandler
// Initialize with appropriate value
```

#### Type Definition

```go
type NoRequestHandler func(ctx context.Context) (Res, error)
```

### NoResponseHandler
NoResponseHandler is a handler that returns no response body.

#### Example Usage

```go
// Example usage of NoResponseHandler
var value NoResponseHandler
// Initialize with appropriate value
```

#### Type Definition

```go
type NoResponseHandler func(ctx context.Context, req Req) error
```

### Option
Option configures a Server.

#### Example Usage

```go
// Example usage of Option
var value Option
// Initialize with appropriate value
```

#### Type Definition

```go
type Option func(*Server)
```

### Constructor Functions

### HideBanner

HideBanner hides the banner and sets the banner to an empty string.

```go
func HideBanner() Option
```

**Parameters:**
  None

**Returns:**
- Option

### WithAddr

WithAddr sets the address the server listens on. Default is ":8080".

```go
func WithAddr(addr string) Option
```

**Parameters:**
- `addr` (string)

**Returns:**
- Option

### WithBasePath

WithBasePath sets a base path prefix for all routes. All registered routes will be prefixed with this path. For example, with base path "/api/v1", a route "/users" becomes "/api/v1/users". The base path should start with "/" but should not end with "/" (it will be normalized).

```go
func WithBasePath(path string) Option
```

**Parameters:**
- `path` (string)

**Returns:**
- Option

### WithCustomBanner

WithCustomBanner sets a custom banner for the server.

```go
func WithCustomBanner(banner string) Option
```

**Parameters:**
- `banner` (string)

**Returns:**
- Option

### WithErrorHandler

WithErrorHandler sets a custom error handler for the server. The error handler will be called whenever an error occurs in a handler. If not set, the default error handling (RFC 7807 Problem Details) is used.

```go
func WithErrorHandler(handler ErrorHandler) Option
```

**Parameters:**
- `handler` (ErrorHandler)

**Returns:**
- Option

### WithGracePeriod

WithGracePeriod sets the maximum duration to wait for active connections to finish during graceful shutdown.

```go
func WithGracePeriod(d time.Duration) Option
```

**Parameters:**
- `d` (time.Duration)

**Returns:**
- Option

### WithIdleTimeout

WithIdleTimeout sets the maximum amount of time to wait for the next request when keep-alives are enabled.

```go
func WithIdleTimeout(d time.Duration) Option
```

**Parameters:**
- `d` (time.Duration)

**Returns:**
- Option

### WithMaxHeaderBytes

WithMaxHeaderBytes sets the maximum size of request headers.

```go
func WithMaxHeaderBytes(n int) Option
```

**Parameters:**
- `n` (int)

**Returns:**
- Option

### WithReadTimeout

WithReadTimeout sets the maximum duration for reading the entire request.

```go
func WithReadTimeout(d time.Duration) Option
```

**Parameters:**
- `d` (time.Duration)

**Returns:**
- Option

### WithTLS

WithTLS configures the server to use TLS with the provided certificate and key files.

```go
func WithTLS(certFile, keyFile string) Option
```

**Parameters:**
- `certFile` (string)
- `keyFile` (string)

**Returns:**
- Option

### WithTLSConfig

WithTLSConfig sets a custom TLS configuration for the server.

```go
func WithTLSConfig(config *tls.Config) Option
```

**Parameters:**
- `config` (*tls.Config)

**Returns:**
- Option

### WithWriteTimeout

WithWriteTimeout sets the maximum duration before timing out writes of the response.

```go
func WithWriteTimeout(d time.Duration) Option
```

**Parameters:**
- `d` (time.Duration)

**Returns:**
- Option

### PaginatedResponse
PaginatedResponse wraps a list response with pagination metadata.

#### Example Usage

```go
// Create a new PaginatedResponse
paginatedresponse := PaginatedResponse{
    Items: [],
    Total: 42,
    Page: 42,
    Limit: 42,
    TotalPages: 42,
    HasMore: true,
    NextCursor: "example",
}
```

#### Type Definition

```go
type PaginatedResponse struct {
    Items []T `json:"items"`
    Total int `json:"total"`
    Page int `json:"page"`
    Limit int `json:"limit"`
    TotalPages int `json:"total_pages"`
    HasMore bool `json:"has_more"`
    NextCursor string `json:"next_cursor,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Items | `[]T` |  |
| Total | `int` |  |
| Page | `int` |  |
| Limit | `int` |  |
| TotalPages | `int` |  |
| HasMore | `bool` |  |
| NextCursor | `string` |  |

### Constructor Functions

### NewCursorResponse

NewCursorResponse creates a new cursor-based paginated response.

```go
func NewCursorResponse(items []T, total int, nextCursor string) *ast.IndexExpr
```

**Parameters:**
- `items` ([]T)
- `total` (int)
- `nextCursor` (string)

**Returns:**
- *ast.IndexExpr

### NewPaginatedResponse

NewPaginatedResponse creates a new paginated response.

```go
func NewPaginatedResponse(items []T, total, page, limit int) *ast.IndexExpr
```

**Parameters:**
- `items` ([]T)
- `total` (int)
- `page` (int)
- `limit` (int)

**Returns:**
- *ast.IndexExpr

### Pagination
Pagination contains common pagination parameters. Use with struct embedding for automatic binding. Example: type ListUsersRequest struct { helix.Pagination Status string `query:"status"` }

#### Example Usage

```go
// Create a new Pagination
pagination := Pagination{
    Page: 42,
    Limit: 42,
    Sort: "example",
    Order: "example",
    Cursor: "example",
}
```

#### Type Definition

```go
type Pagination struct {
    Page int `query:"page"`
    Limit int `query:"limit"`
    Sort string `query:"sort"`
    Order string `query:"order"`
    Cursor string `query:"cursor"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Page | `int` |  |
| Limit | `int` |  |
| Sort | `string` |  |
| Order | `string` |  |
| Cursor | `string` |  |

### Constructor Functions

### BindPagination

BindPagination extracts pagination from the request with defaults.

```go
func (*Ctx) BindPagination(defaultLimit, maxLimit int) Pagination
```

**Parameters:**
- `defaultLimit` (int)
- `maxLimit` (int)

**Returns:**
- Pagination

## Methods

### GetLimit

GetLimit returns the limit with a default and maximum.

```go
func (Pagination) GetLimit(defaultLimit, maxLimit int) int
```

**Parameters:**
- `defaultLimit` (int)
- `maxLimit` (int)

**Returns:**
- int

### GetOffset

GetOffset calculates the offset for SQL queries.

```go
func (Pagination) GetOffset(limit int) int
```

**Parameters:**
- `limit` (int)

**Returns:**
- int

### GetOrder

GetOrder returns the order (asc/desc) with a default of desc.

```go
func (Pagination) GetOrder() string
```

**Parameters:**
  None

**Returns:**
- string

### GetPage

GetPage returns the page number with a default of 1.

```go
func (Pagination) GetPage() int
```

**Parameters:**
  None

**Returns:**
- int

### GetSort

GetSort returns the sort field with a default.

```go
func (Pagination) GetSort(defaultSort string, allowed []string) string
```

**Parameters:**
- `defaultSort` (string)
- `allowed` ([]string)

**Returns:**
- string

### IsAscending

IsAscending returns true if the order is ascending.

```go
func (Pagination) IsAscending() bool
```

**Parameters:**
  None

**Returns:**
- bool

### Problem
Problem represents an RFC 7807 Problem Details for HTTP APIs. See: https://tools.ietf.org/html/rfc7807

#### Example Usage

```go
// Create a new Problem
problem := Problem{
    Type: "example",
    Title: "example",
    Status: 42,
    Detail: "example",
    Instance: "example",
    Err: error{},
}
```

#### Type Definition

```go
type Problem struct {
    Type string `json:"type"`
    Title string `json:"title"`
    Status int `json:"status"`
    Detail string `json:"detail,omitempty"`
    Instance string `json:"instance,omitempty"`
    Err error `json:"-"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Type | `string` | Type is a URI reference that identifies the problem type. |
| Title | `string` | Title is a short, human-readable summary of the problem type. |
| Status | `int` | Status is the HTTP status code for this problem. |
| Detail | `string` | Detail is a human-readable explanation specific to this occurrence of the problem. |
| Instance | `string` | Instance is a URI reference that identifies the specific occurrence of the problem. |
| Err | `error` | Err is the error that caused the problem. |

### Constructor Functions

### BadGatewayf

BadGatewayf creates a 502 Bad Gateway Problem with a formatted detail message.

```go
func BadGatewayf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### BadRequestf

BadRequestf creates a 400 Bad Request Problem with a formatted detail message.

```go
func BadRequestf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### Conflictf

Conflictf creates a 409 Conflict Problem with a formatted detail message.

```go
func Conflictf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### Forbiddenf

Forbiddenf creates a 403 Forbidden Problem with a formatted detail message.

```go
func Forbiddenf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### GatewayTimeoutf

GatewayTimeoutf creates a 504 Gateway Timeout Problem with a formatted detail message.

```go
func GatewayTimeoutf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### Gonef

Gonef creates a 410 Gone Problem with a formatted detail message.

```go
func Gonef(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### Internalf

Internalf creates a 500 Internal Server Error Problem with a formatted detail message.

```go
func Internalf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### MethodNotAllowedf

MethodNotAllowedf creates a 405 Method Not Allowed Problem with a formatted detail message.

```go
func MethodNotAllowedf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### NewProblem

NewProblem creates a new Problem with the given status, type, and title.

```go
func NewProblem(status int, problemType, title string) Problem
```

**Parameters:**
- `status` (int)
- `problemType` (string)
- `title` (string)

**Returns:**
- Problem

### NotFoundf

NotFoundf creates a 404 Not Found Problem with a formatted detail message.

```go
func NotFoundf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### NotImplementedf

NotImplementedf creates a 501 Not Implemented Problem with a formatted detail message.

```go
func NotImplementedf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### ProblemFromStatus

ProblemFromStatus creates a Problem from an HTTP status code.

```go
func ProblemFromStatus(status int) Problem
```

**Parameters:**
- `status` (int)

**Returns:**
- Problem

### ServiceUnavailablef

ServiceUnavailablef creates a 503 Service Unavailable Problem with a formatted detail message.

```go
func ServiceUnavailablef(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### TooManyRequestsf

TooManyRequestsf creates a 429 Too Many Requests Problem with a formatted detail message.

```go
func TooManyRequestsf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### Unauthorizedf

Unauthorizedf creates a 401 Unauthorized Problem with a formatted detail message.

```go
func Unauthorizedf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### UnprocessableEntityf

UnprocessableEntityf creates a 422 Unprocessable Entity Problem with a formatted detail message.

```go
func UnprocessableEntityf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

## Methods

### Error

Error implements the error interface.

```go
func (*ValidationErrors) Error() string
```

**Parameters:**
  None

**Returns:**
- string

### WithDetail



```go
func (Problem) WithDetail(detail string) Problem
```

**Parameters:**
- `detail` (string)

**Returns:**
- Problem

### WithDetailf

WithDetailf returns a copy of the Problem with the given detail message.

```go
func (Problem) WithDetailf(format string, args ...any) Problem
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
- Problem

### WithErr

WithStack returns a copy of the Problem with the given stack trace.

```go
func (Problem) WithErr(err error) Problem
```

**Parameters:**
- `err` (error)

**Returns:**
- Problem

### WithInstance

WithInstance returns a copy of the Problem with the given instance URI.

```go
func (Problem) WithInstance(instance string) Problem
```

**Parameters:**
- `instance` (string)

**Returns:**
- Problem

### WithType

WithType returns a copy of the Problem with the given type URI.

```go
func (Problem) WithType(problemType string) Problem
```

**Parameters:**
- `problemType` (string)

**Returns:**
- Problem

### ResourceBuilder
ResourceBuilder provides a fluent interface for defining REST resource routes.

#### Example Usage

```go
// Create a new ResourceBuilder
resourcebuilder := ResourceBuilder{

}
```

#### Type Definition

```go
type ResourceBuilder struct {
}
```

## Methods

### CRUD

CRUD registers all standard CRUD handlers in one call. Handlers: list, create, get, update, delete

```go
func (*ResourceBuilder) CRUD(list, create, get, update, delete http.HandlerFunc) *ResourceBuilder
```

**Parameters:**
- `list` (http.HandlerFunc)
- `create` (http.HandlerFunc)
- `get` (http.HandlerFunc)
- `update` (http.HandlerFunc)
- `delete` (http.HandlerFunc)

**Returns:**
- *ResourceBuilder

### Create

Create registers a POST handler for creating resources (e.g., POST /users).

```go
func (**ast.IndexExpr) Create(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### Custom

Custom registers a handler with a custom method and path suffix. The suffix is appended to the base pattern. Example: Custom("POST", "/{id}/archive", archiveHandler) for POST /users/{id}/archive

```go
func (**ast.IndexExpr) Custom(method, suffix string, handler http.HandlerFunc) **ast.IndexExpr
```

**Parameters:**
- `method` (string)
- `suffix` (string)
- `handler` (http.HandlerFunc)

**Returns:**
- **ast.IndexExpr

### Delete

Delete registers a DELETE handler for deleting a resource (e.g., DELETE /users/{id}).

```go
func (**ast.IndexExpr) Delete(h *ast.IndexExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexExpr)

**Returns:**
- **ast.IndexExpr

### Destroy

Destroy is an alias for Delete.

```go
func (*ResourceBuilder) Destroy(handler http.HandlerFunc) *ResourceBuilder
```

**Parameters:**
- `handler` (http.HandlerFunc)

**Returns:**
- *ResourceBuilder

### Get

Get registers a GET handler for a single resource (e.g., GET /users/{id}).

```go
func (*Ctx) Get(key string) (any, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- any
- bool

### Index

Index is an alias for List.

```go
func (*ResourceBuilder) Index(handler http.HandlerFunc) *ResourceBuilder
```

**Parameters:**
- `handler` (http.HandlerFunc)

**Returns:**
- *ResourceBuilder

### List

List registers a GET handler for the collection (e.g., GET /users).

```go
func (**ast.IndexExpr) List(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### Patch

Patch registers a PATCH handler for partial updates (e.g., PATCH /users/{id}).

```go
func (**ast.IndexExpr) Patch(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### ReadOnly

ReadOnly registers only read handlers (list and get).

```go
func (*ResourceBuilder) ReadOnly(list, get http.HandlerFunc) *ResourceBuilder
```

**Parameters:**
- `list` (http.HandlerFunc)
- `get` (http.HandlerFunc)

**Returns:**
- *ResourceBuilder

### Show

Show is an alias for Get.

```go
func (*ResourceBuilder) Show(handler http.HandlerFunc) *ResourceBuilder
```

**Parameters:**
- `handler` (http.HandlerFunc)

**Returns:**
- *ResourceBuilder

### Store

Store is an alias for Create.

```go
func (*ResourceBuilder) Store(handler http.HandlerFunc) *ResourceBuilder
```

**Parameters:**
- `handler` (http.HandlerFunc)

**Returns:**
- *ResourceBuilder

### Update

Update registers a PUT handler for updating a resource (e.g., PUT /users/{id}).

```go
func (**ast.IndexExpr) Update(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### RouteInfo
RouteInfo contains information about a registered route.

#### Example Usage

```go
// Create a new RouteInfo
routeinfo := RouteInfo{
    Method: "example",
    Pattern: "example",
}
```

#### Type Definition

```go
type RouteInfo struct {
    Method string
    Pattern string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Method | `string` |  |
| Pattern | `string` |  |

### RouteRegistrar
RouteRegistrar is an interface for registering routes. Both Server and Group implement this interface.

#### Example Usage

```go
// Example implementation of RouteRegistrar
type MyRouteRegistrar struct {
    // Add your fields here
}

func (m MyRouteRegistrar) GET(param1 string, param2 http.HandlerFunc)  {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) POST(param1 string, param2 http.HandlerFunc)  {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) PUT(param1 string, param2 http.HandlerFunc)  {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) PATCH(param1 string, param2 http.HandlerFunc)  {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) DELETE(param1 string, param2 http.HandlerFunc)  {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) OPTIONS(param1 string, param2 http.HandlerFunc)  {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) HEAD(param1 string, param2 http.HandlerFunc)  {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) Handle(param1 string, param2 http.HandlerFunc)  {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) Group(param1 string, param2 ...any) *Group {
    // Implement your logic here
    return
}

func (m MyRouteRegistrar) Resource(param1 string, param2 ...any) *ResourceBuilder {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type RouteRegistrar interface {
    GET(pattern string, handler http.HandlerFunc)
    POST(pattern string, handler http.HandlerFunc)
    PUT(pattern string, handler http.HandlerFunc)
    PATCH(pattern string, handler http.HandlerFunc)
    DELETE(pattern string, handler http.HandlerFunc)
    OPTIONS(pattern string, handler http.HandlerFunc)
    HEAD(pattern string, handler http.HandlerFunc)
    Handle(method, pattern string, handler http.HandlerFunc)
    Group(prefix string, mw ...any) *Group
    Resource(pattern string, mw ...any) *ResourceBuilder
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### Router
Router handles HTTP request routing.

#### Example Usage

```go
// Create a new Router
router := Router{

}
```

#### Type Definition

```go
type Router struct {
}
```

## Methods

### Handle

Handle registers a new route with the given method and pattern.

```go
func (*Server) Handle(method, pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `method` (string)
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Routes

Routes returns all registered routes.

```go
func (*Router) Routes() []RouteInfo
```

**Parameters:**
  None

**Returns:**
- []RouteInfo

### ServeHTTP

ServeHTTP implements http.Handler.

```go
func (*Router) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

**Parameters:**
- `w` (http.ResponseWriter)
- `req` (*http.Request)

**Returns:**
  None

### Server
Server is the main HTTP server for the Helix framework.

#### Example Usage

```go
// Create a new Server
server := Server{

}
```

#### Type Definition

```go
type Server struct {
}
```

### Constructor Functions

### Default

Default creates a new Server with sensible defaults for development. It includes RequestID, Logger (dev format), and Recover middleware.

```go
func Default(opts ...Option) *Server
```

**Parameters:**
- `opts` (...Option)

**Returns:**
- *Server

### New

New creates a new Server with the provided options.

```go
func New(opts ...Option) *Server
```

**Parameters:**
- `opts` (...Option)

**Returns:**
- *Server

## Methods

### Addr

Addr returns the address the server is configured to listen on.

```go
func (*Server) Addr() string
```

**Parameters:**
  None

**Returns:**
- string

### Any

Any registers a handler for all HTTP methods.

```go
func (*Group) Any(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Build

Build pre-compiles the middleware chain for optimal performance. This is called automatically before the server starts, but can be called manually after all routes and middleware are registered.

```go
func (*Server) Build()
```

**Parameters:**
  None

**Returns:**
  None

### CONNECT

CONNECT registers a handler for CONNECT requests.

```go
func (*Server) CONNECT(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### DELETE

DELETE registers a handler for DELETE requests.

```go
func (*Group) DELETE(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### GET

GET registers a handler for GET requests.

```go
func (*Server) GET(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Group

Group creates a new route group with the given prefix. The prefix is prepended to all routes registered on the group. Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.

```go
func (*Group) Group(prefix string, mw ...any) *Group
```

**Parameters:**
- `prefix` (string)
- `mw` (...any)

**Returns:**
- *Group

### HEAD

HEAD registers a handler for HEAD requests.

```go
func (*Group) HEAD(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Handle

Handle registers a handler for the given method and pattern.

```go
func (*Group) Handle(method, pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `method` (string)
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Mount

Mount mounts a module at the given prefix. The module's routes are prefixed with the given path.

```go
func (*Group) Mount(prefix string, m Module, mw ...any)
```

**Parameters:**
- `prefix` (string)
- `m` (Module)
- `mw` (...any)

**Returns:**
  None

### MountFunc

MountFunc mounts a function as a module at the given prefix.

```go
func (*Group) MountFunc(prefix string, fn func(r RouteRegistrar), mw ...any)
```

**Parameters:**
- `prefix` (string)
- `fn` (func(r RouteRegistrar))
- `mw` (...any)

**Returns:**
  None

### OPTIONS

OPTIONS registers a handler for OPTIONS requests.

```go
func (*Group) OPTIONS(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### OnStart

OnStart registers a function to be called when the server starts. Multiple functions can be registered and will be called in order.

```go
func (*Server) OnStart(fn func(s *Server))
```

**Parameters:**
- `fn` (func(s *Server))

**Returns:**
  None

### OnStop

OnStop registers a function to be called when the server stops. Multiple functions can be registered and will be called in order. The context passed to the function has the grace period as its deadline.

```go
func (*Server) OnStop(fn func(ctx context.Context, s *Server))
```

**Parameters:**
- `fn` (func(ctx context.Context, s *Server))

**Returns:**
  None

### PATCH

PATCH registers a handler for PATCH requests.

```go
func (*Group) PATCH(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### POST

POST registers a handler for POST requests.

```go
func (*Group) POST(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### PUT

PUT registers a handler for PUT requests.

```go
func (*Group) PUT(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### PrintRoutes

PrintRoutes prints all registered routes to the given writer. Routes are sorted by pattern, then by method.

```go
func (*Server) PrintRoutes(w io.Writer)
```

**Parameters:**
- `w` (io.Writer)

**Returns:**
  None

### Resource

Resource creates a new ResourceBuilder for the given pattern. The pattern should be the base path for the resource (e.g., "/users"). Optional middleware can be applied to all routes in the resource. Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.

```go
func (*Server) Resource(pattern string, mw ...any) *ResourceBuilder
```

**Parameters:**
- `pattern` (string)
- `mw` (...any)

**Returns:**
- *ResourceBuilder

### Routes

Routes returns all registered routes.

```go
func (*Server) Routes() []RouteInfo
```

**Parameters:**
  None

**Returns:**
- []RouteInfo

### Run

Run starts the server and blocks until the context is canceled or a shutdown signal is received. It performs graceful shutdown, waiting for active connections to finish within the grace period.

```go
func (*Server) Run(ctx context.Context) error
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- error

### ServeHTTP

ServeHTTP implements the http.Handler interface.

```go
func (*Server) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

**Parameters:**
- `w` (http.ResponseWriter)
- `r` (*http.Request)

**Returns:**
  None

### Shutdown

Shutdown gracefully shuts down the server without interrupting active connections. It waits for the grace period for active connections to finish.

```go
func (*Server) Shutdown(ctx context.Context) error
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- error

### Start

Start starts the server and blocks until shutdown. If an address is provided, it will be used instead of the WithAddr option. If the address is not provided and the WithAddr option is not set, it will use ":8080". This is a convenience method that calls Run with a background context.

```go
func (*Server) Start(addr ...string) error
```

**Parameters:**
- `addr` (...string)

**Returns:**
- error

### Static

Static serves static files from the given file system root.

```go
func (*Server) Static(pattern, root string)
```

**Parameters:**
- `pattern` (string)
- `root` (string)

**Returns:**
  None

### TRACE

TRACE registers a handler for TRACE requests.

```go
func (*Server) TRACE(pattern string, handler http.HandlerFunc)
```

**Parameters:**
- `pattern` (string)
- `handler` (http.HandlerFunc)

**Returns:**
  None

### Use

Use adds middleware to the server's middleware chain. Middleware is executed in the order it is added. Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.

```go
func (*Server) Use(mw ...any)
```

**Parameters:**
- `mw` (...any)

**Returns:**
  None

### TypedResourceBuilder
TypedResourceBuilder provides a fluent interface for defining typed REST resource routes.

#### Example Usage

```go
// Create a new TypedResourceBuilder
typedresourcebuilder := TypedResourceBuilder{

}
```

#### Type Definition

```go
type TypedResourceBuilder struct {
}
```

### Constructor Functions

### TypedResource

TypedResource creates a typed resource builder for the given entity type. This provides a fluent API for registering typed handlers for CRUD operations. Example: helix.TypedResource[User](s, "/users"). List(listUsers). Create(createUser). Get(getUser). Update(updateUser). Delete(deleteUser)

```go
func TypedResource(s *Server, pattern string, mw ...any) **ast.IndexExpr
```

**Parameters:**
- `s` (*Server)
- `pattern` (string)
- `mw` (...any)

**Returns:**
- **ast.IndexExpr

### TypedResourceForGroup

TypedResourceForGroup creates a typed resource builder within a group.

```go
func TypedResourceForGroup(g *Group, pattern string, mw ...any) **ast.IndexExpr
```

**Parameters:**
- `g` (*Group)
- `pattern` (string)
- `mw` (...any)

**Returns:**
- **ast.IndexExpr

## Methods

### Create

Create registers a typed POST handler for creating resources. Handler signature: func(ctx, CreateReq) (Entity, error)

```go
func (**ast.IndexExpr) Create(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### Custom

Custom registers a handler with a custom method and path suffix.

```go
func (**ast.IndexExpr) Custom(method, suffix string, handler http.HandlerFunc) **ast.IndexExpr
```

**Parameters:**
- `method` (string)
- `suffix` (string)
- `handler` (http.HandlerFunc)

**Returns:**
- **ast.IndexExpr

### Delete

Delete registers a typed DELETE handler for deleting a resource. Handler signature: func(ctx, IDRequest) error

```go
func (**ast.IndexExpr) Delete(h *ast.IndexExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexExpr)

**Returns:**
- **ast.IndexExpr

### Get

Get registers a typed GET handler for a single resource. Handler signature: func(ctx, IDRequest) (Entity, error)

```go
func (*Ctx) Get(key string) (any, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- any
- bool

### List

List registers a typed GET handler for the collection. Handler signature: func(ctx, ListReq) (ListResponse[Entity], error)

```go
func (**ast.IndexExpr) List(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### Patch

Patch registers a typed PATCH handler for partial updates.

```go
func (**ast.IndexExpr) Patch(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### Update

Update registers a typed PUT handler for updating a resource. The request type should include the ID from path and the update data.

```go
func (**ast.IndexExpr) Update(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
- `h` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexExpr

### Validatable
Validatable is an interface for types that can validate themselves.

#### Example Usage

```go
// Example implementation of Validatable
type MyValidatable struct {
    // Add your fields here
}

func (m MyValidatable) Validate() error {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Validatable interface {
    Validate() error
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### ValidationErrors
ValidationErrors collects multiple validation errors for RFC 7807 response. Implements the error interface and can be returned from Validate() methods.

#### Example Usage

```go
// Create a new ValidationErrors
validationerrors := ValidationErrors{

}
```

#### Type Definition

```go
type ValidationErrors struct {
}
```

### Constructor Functions

### NewValidationErrors

NewValidationErrors creates a new empty ValidationErrors collector.

```go
func NewValidationErrors() *ValidationErrors
```

**Parameters:**
  None

**Returns:**
- *ValidationErrors

## Methods

### Add

Add adds a validation error for a specific field.

```go
func (*ValidationErrors) Add(field, message string)
```

**Parameters:**
- `field` (string)
- `message` (string)

**Returns:**
  None

### Addf

Addf adds a validation error for a specific field with a formatted message.

```go
func (*ValidationErrors) Addf(field, format string, args ...any)
```

**Parameters:**
- `field` (string)
- `format` (string)
- `args` (...any)

**Returns:**
  None

### Err

Err returns nil if there are no errors, otherwise returns the ValidationErrors. This is useful for the common pattern: return v.Err()

```go
func (*ValidationErrors) Err() error
```

**Parameters:**
  None

**Returns:**
- error

### Error

Error implements the error interface.

```go
func (*ValidationErrors) Error() string
```

**Parameters:**
  None

**Returns:**
- string

### Errors

Errors returns the list of field errors.

```go
func (*ValidationErrors) Errors() []FieldError
```

**Parameters:**
  None

**Returns:**
- []FieldError

### HasErrors

HasErrors returns true if there are any validation errors.

```go
func (*ValidationErrors) HasErrors() bool
```

**Parameters:**
  None

**Returns:**
- bool

### Len

Len returns the number of validation errors.

```go
func (*ValidationErrors) Len() int
```

**Parameters:**
  None

**Returns:**
- int

### ToProblem

ToProblem converts ValidationErrors to a ValidationProblem for RFC 7807 response.

```go
func (*ValidationErrors) ToProblem() ValidationProblem
```

**Parameters:**
  None

**Returns:**
- ValidationProblem

### ValidationProblem
ValidationProblem is an RFC 7807 Problem with validation errors extension.

#### Example Usage

```go
// Create a new ValidationProblem
validationproblem := ValidationProblem{
    Errors: [],
}
```

#### Type Definition

```go
type ValidationProblem struct {
    Problem
    Errors []FieldError `json:"errors,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| *Problem | `Problem` |  |
| Errors | `[]FieldError` |  |

## Functions

### Accepted
Accepted writes a 202 Accepted JSON response.

```go
func (*Ctx) Accepted(v any) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `v` | `any` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of Accepted
result := Accepted(/* parameters */)
```

### Attachment
Attachment sets the Content-Disposition header to attachment with the given filename.

```go
func Attachment(w http.ResponseWriter, filename string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `filename` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of Attachment
result := Attachment(/* parameters */)
```

### BadRequest
BadRequest writes a 400 Bad Request error response.

```go
func BadRequest(w http.ResponseWriter, message string) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `message` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of BadRequest
result := BadRequest(/* parameters */)
```

### Bind
Bind binds path parameters, query parameters, headers, and JSON body to a struct. The binding sources are determined by struct tags: - `path:"name"` - binds from URL path parameters - `query:"name"` - binds from URL query parameters - `header:"name"` - binds from HTTP headers - `json:"name"` - binds from JSON body - `form:"name"` - binds from form data

```go
func Bind(r *http.Request) (T, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |
| `error` | |

**Example:**

```go
// Example usage of Bind
result := Bind(/* parameters */)
```

### BindAndValidate
BindAndValidate binds and validates a request. If the bound type implements Validatable, Validate() is called after binding.

```go
func BindAndValidate(r *http.Request) (T, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |
| `error` | |

**Example:**

```go
// Example usage of BindAndValidate
result := BindAndValidate(/* parameters */)
```

### BindHeader
BindHeader binds HTTP headers to a struct. Uses the `header` struct tag to determine field names.

```go
func BindHeader(r *http.Request) (T, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |
| `error` | |

**Example:**

```go
// Example usage of BindHeader
result := BindHeader(/* parameters */)
```

### BindJSON
BindJSON binds the JSON request body to a struct.

```go
func BindJSON(r *http.Request) (T, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |
| `error` | |

**Example:**

```go
// Example usage of BindJSON
result := BindJSON(/* parameters */)
```

### BindPath
BindPath binds URL path parameters to a struct. Uses the `path` struct tag to determine field names.

```go
func BindPath(r *http.Request) (T, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |
| `error` | |

**Example:**

```go
// Example usage of BindPath
result := BindPath(/* parameters */)
```

### BindQuery
BindQuery binds URL query parameters to a struct. Uses the `query` struct tag to determine field names.

```go
func BindQuery(r *http.Request) (T, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |
| `error` | |

**Example:**

```go
// Example usage of BindQuery
result := BindQuery(/* parameters */)
```

### Blob
Blob writes binary data with the given content type.

```go
func (*Ctx) Blob(status int, contentType string, data []byte) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | `int` | |
| `contentType` | `string` | |
| `data` | `[]byte` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of Blob
result := Blob(/* parameters */)
```

### Created
Created writes a 201 Created JSON response.

```go
func Created(w http.ResponseWriter, v any) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `v` | `any` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of Created
result := Created(/* parameters */)
```

### Error
Error writes an error response with the given status code and message.

```go
func (*ValidationErrors) Error() string
```

**Parameters:**
None

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of Error
result := Error(/* parameters */)
```

### File
File serves a file with the given content type.

```go
func File(w http.ResponseWriter, r *http.Request, path string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `r` | `*http.Request` | |
| `path` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of File
result := File(/* parameters */)
```

### Forbidden
Forbidden writes a 403 Forbidden error response.

```go
func (*Ctx) Forbidden(message string) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `message` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of Forbidden
result := Forbidden(/* parameters */)
```

### FromContext
FromContext retrieves a service from the context or falls back to global registry.

```go
func FromContext(ctx context.Context) (T, bool)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |
| `bool` | |

**Example:**

```go
// Example usage of FromContext
result := FromContext(/* parameters */)
```

### Get
Get retrieves a service from the global registry. Returns the zero value and false if not found.

```go
func (**ast.IndexExpr) Get(h *ast.IndexListExpr) **ast.IndexExpr
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `h` | `*ast.IndexListExpr` | |

**Returns:**
| Type | Description |
|------|-------------|
| `**ast.IndexExpr` | |

**Example:**

```go
// Example usage of Get
result := Get(/* parameters */)
```

### HTML
HTML writes an HTML response with the given status code.

```go
func HTML(w http.ResponseWriter, status int, html string) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `status` | `int` | |
| `html` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of HTML
result := HTML(/* parameters */)
```

### Handle
Handle wraps a generic Handler into an http.HandlerFunc. It automatically: - Binds the request to the Req type - Calls the handler with the context and request - Encodes the response as JSON - Handles errors using RFC 7807 Problem Details

```go
func (*Group) Handle(method, pattern string, handler http.HandlerFunc)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `method` | `string` | |
| `pattern` | `string` | |
| `handler` | `http.HandlerFunc` | |

**Returns:**
None

**Example:**

```go
// Example usage of Handle
result := Handle(/* parameters */)
```

### HandleAccepted
HandleAccepted wraps a generic Handler into an http.HandlerFunc that returns 202 Accepted. Useful for async operations where processing happens in the background.

```go
func HandleAccepted(h *ast.IndexListExpr) http.HandlerFunc
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `h` | `*ast.IndexListExpr` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of HandleAccepted
result := HandleAccepted(/* parameters */)
```

### HandleCreated
HandleCreated wraps a generic Handler into an http.HandlerFunc that returns 201 Created. This is a convenience wrapper for HandleWithStatus(http.StatusCreated, h).

```go
func HandleCreated(h *ast.IndexListExpr) http.HandlerFunc
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `h` | `*ast.IndexListExpr` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of HandleCreated
result := HandleCreated(/* parameters */)
```

### HandleCtx
HandleCtx wraps a CtxHandler into an http.HandlerFunc. Errors returned from the handler are automatically converted to RFC 7807 responses.

```go
func HandleCtx(h CtxHandler) http.HandlerFunc
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `h` | `CtxHandler` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of HandleCtx
result := HandleCtx(/* parameters */)
```

### HandleEmpty
HandleEmpty wraps an EmptyHandler into an http.HandlerFunc. Returns 204 No Content on success.

```go
func HandleEmpty(h EmptyHandler) http.HandlerFunc
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `h` | `EmptyHandler` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of HandleEmpty
result := HandleEmpty(/* parameters */)
```

### HandleErrorDefault
HandleErrorDefault provides the default error handling logic. This can be called from custom error handlers to fall back to default behavior.

```go
func HandleErrorDefault(w http.ResponseWriter, r *http.Request, err error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `r` | `*http.Request` | |
| `err` | `error` | |

**Returns:**
None

**Example:**

```go
// Example usage of HandleErrorDefault
result := HandleErrorDefault(/* parameters */)
```

### HandleNoRequest
HandleNoRequest wraps a NoRequestHandler into an http.HandlerFunc. Useful for endpoints that don't need request binding (e.g., GET /users).

```go
func HandleNoRequest(h *ast.IndexExpr) http.HandlerFunc
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `h` | `*ast.IndexExpr` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of HandleNoRequest
result := HandleNoRequest(/* parameters */)
```

### HandleNoResponse
HandleNoResponse wraps a NoResponseHandler into an http.HandlerFunc. Returns 204 No Content on success.

```go
func HandleNoResponse(h *ast.IndexExpr) http.HandlerFunc
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `h` | `*ast.IndexExpr` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of HandleNoResponse
result := HandleNoResponse(/* parameters */)
```

### HandleWithStatus
HandleWithStatus wraps a generic Handler into an http.HandlerFunc with a custom success status code.

```go
func HandleWithStatus(status int, h *ast.IndexListExpr) http.HandlerFunc
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | `int` | |
| `h` | `*ast.IndexListExpr` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of HandleWithStatus
result := HandleWithStatus(/* parameters */)
```

### Inline
Inline sets the Content-Disposition header to inline with the given filename.

```go
func Inline(w http.ResponseWriter, filename string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `filename` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of Inline
result := Inline(/* parameters */)
```

### InternalServerError
InternalServerError writes a 500 Internal Server Error response.

```go
func (*Ctx) InternalServerError(message string) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `message` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of InternalServerError
result := InternalServerError(/* parameters */)
```

### JSON
JSON writes a JSON response with the given status code. Uses a pooled buffer for zero-allocation in the hot path.

```go
func JSON(w http.ResponseWriter, status int, v any) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `status` | `int` | |
| `v` | `any` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of JSON
result := JSON(/* parameters */)
```

### JSONPretty
JSONPretty writes a pretty-printed JSON response with the given status code.

```go
func JSONPretty(w http.ResponseWriter, status int, v any, indent string) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `status` | `int` | |
| `v` | `any` | |
| `indent` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of JSONPretty
result := JSONPretty(/* parameters */)
```

### LivenessHandler
LivenessHandler returns a simple liveness probe handler. Returns 200 OK if the server is running.

```go
func LivenessHandler() http.HandlerFunc
```

**Parameters:**
None

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of LivenessHandler
result := LivenessHandler(/* parameters */)
```

### MustFromContext
MustFromContext retrieves a service from context or panics.

```go
func MustFromContext(ctx context.Context) T
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |

**Example:**

```go
// Example usage of MustFromContext
result := MustFromContext(/* parameters */)
```

### MustGet
MustGet retrieves a service from the global registry or panics.

```go
func (*Ctx) MustGet(key string) any
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `key` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `any` | |

**Example:**

```go
// Example usage of MustGet
result := MustGet(/* parameters */)
```

### NoContent
NoContent writes a 204 No Content response.

```go
func NoContent(w http.ResponseWriter) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of NoContent
result := NoContent(/* parameters */)
```

### NotFound
NotFound writes a 404 Not Found error response.

```go
func NotFound(w http.ResponseWriter, message string) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `message` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of NotFound
result := NotFound(/* parameters */)
```

### OK
OK writes a 200 OK JSON response.

```go
func OK(w http.ResponseWriter, v any) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `v` | `any` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of OK
result := OK(/* parameters */)
```

### Param
Param returns the value of a path parameter. Returns an empty string if the parameter does not exist.

```go
func Param(r *http.Request, name string) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |
| `name` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of Param
result := Param(/* parameters */)
```

### ParamInt
ParamInt returns the value of a path parameter as an int. Returns an error if the parameter does not exist or cannot be parsed.

```go
func (*Ctx) ParamInt(name string) (int, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `int` | |
| `error` | |

**Example:**

```go
// Example usage of ParamInt
result := ParamInt(/* parameters */)
```

### ParamInt64
ParamInt64 returns the value of a path parameter as an int64. Returns an error if the parameter does not exist or cannot be parsed.

```go
func ParamInt64(r *http.Request, name string) (int64, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |
| `name` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `int64` | |
| `error` | |

**Example:**

```go
// Example usage of ParamInt64
result := ParamInt64(/* parameters */)
```

### ParamUUID
ParamUUID returns the value of a path parameter validated as a UUID. Returns an error if the parameter does not exist or is not a valid UUID format.

```go
func (*Ctx) ParamUUID(name string) (string, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |
| `error` | |

**Example:**

```go
// Example usage of ParamUUID
result := ParamUUID(/* parameters */)
```

### Query
Query returns the first value of a query parameter. Returns an empty string if the parameter does not exist.

```go
func (*Ctx) Query(name string) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of Query
result := Query(/* parameters */)
```

### QueryBool
QueryBool returns the first value of a query parameter as a bool. Returns false if the parameter does not exist or cannot be parsed. Accepts "1", "t", "T", "true", "TRUE", "True" as true. Accepts "0", "f", "F", "false", "FALSE", "False" as false.

```go
func (*Ctx) QueryBool(name string) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of QueryBool
result := QueryBool(/* parameters */)
```

### QueryDefault
QueryDefault returns the first value of a query parameter or a default value.

```go
func QueryDefault(r *http.Request, name, defaultVal string) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |
| `name` | `string` | |
| `defaultVal` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of QueryDefault
result := QueryDefault(/* parameters */)
```

### QueryFloat64
QueryFloat64 returns the first value of a query parameter as a float64. Returns the default value if the parameter does not exist or cannot be parsed.

```go
func (*Ctx) QueryFloat64(name string, defaultVal float64) float64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | `string` | |
| `defaultVal` | `float64` | |

**Returns:**
| Type | Description |
|------|-------------|
| `float64` | |

**Example:**

```go
// Example usage of QueryFloat64
result := QueryFloat64(/* parameters */)
```

### QueryInt
QueryInt returns the first value of a query parameter as an int. Returns the default value if the parameter does not exist or cannot be parsed.

```go
func QueryInt(r *http.Request, name string, defaultVal int) int
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |
| `name` | `string` | |
| `defaultVal` | `int` | |

**Returns:**
| Type | Description |
|------|-------------|
| `int` | |

**Example:**

```go
// Example usage of QueryInt
result := QueryInt(/* parameters */)
```

### QueryInt64
QueryInt64 returns the first value of a query parameter as an int64. Returns the default value if the parameter does not exist or cannot be parsed.

```go
func QueryInt64(r *http.Request, name string, defaultVal int64) int64
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |
| `name` | `string` | |
| `defaultVal` | `int64` | |

**Returns:**
| Type | Description |
|------|-------------|
| `int64` | |

**Example:**

```go
// Example usage of QueryInt64
result := QueryInt64(/* parameters */)
```

### QuerySlice
QuerySlice returns all values of a query parameter as a string slice. Returns nil if the parameter does not exist.

```go
func QuerySlice(r *http.Request, name string) []string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |
| `name` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `[]string` | |

**Example:**

```go
// Example usage of QuerySlice
result := QuerySlice(/* parameters */)
```

### ReadinessHandler
ReadinessHandler returns a simple readiness probe handler. Uses the provided checks to determine readiness.

```go
func ReadinessHandler(checks ...func(ctx context.Context) error) http.HandlerFunc
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `checks` | `...func(ctx context.Context) error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `http.HandlerFunc` | |

**Example:**

```go
// Example usage of ReadinessHandler
result := ReadinessHandler(/* parameters */)
```

### Redirect
Redirect redirects the request to the given URL.

```go
func Redirect(w http.ResponseWriter, r *http.Request, url string, code int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `r` | `*http.Request` | |
| `url` | `string` | |
| `code` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of Redirect
result := Redirect(/* parameters */)
```

### Register
Register registers a service in the global registry by its type.

```go
func (ModuleFunc) Register(r RouteRegistrar)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `RouteRegistrar` | |

**Returns:**
None

**Example:**

```go
// Example usage of Register
result := Register(/* parameters */)
```

### Stream
Stream streams the content from the reader to the response.

```go
func Stream(w http.ResponseWriter, contentType string, reader io.Reader) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `contentType` | `string` | |
| `reader` | `io.Reader` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of Stream
result := Stream(/* parameters */)
```

### Text
Text writes a plain text response with the given status code.

```go
func Text(w http.ResponseWriter, status int, text string) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `status` | `int` | |
| `text` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of Text
result := Text(/* parameters */)
```

### Unauthorized
Unauthorized writes a 401 Unauthorized error response.

```go
func Unauthorized(w http.ResponseWriter, message string) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `message` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of Unauthorized
result := Unauthorized(/* parameters */)
```

### UnmarshalJSON
UnmarshalJSON unmarshals a JSON request body into a struct. Returns an error if the request body is not a valid JSON or if the struct cannot be unmarshalled. The error will contain the stack trace of the error.

```go
func UnmarshalJSON(reader io.Reader) (*T, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `reader` | `io.Reader` | |

**Returns:**
| Type | Description |
|------|-------------|
| `*T` | |
| `error` | |

**Example:**

```go
// Example usage of UnmarshalJSON
result := UnmarshalJSON(/* parameters */)
```

### WithService
WithService returns a new context with the service added. This is useful for request-scoped services like database transactions.

```go
func WithService(ctx context.Context, service T) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `service` | `T` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithService
result := WithService(/* parameters */)
```

### WriteProblem
WriteProblem writes a Problem response to the http.ResponseWriter.

```go
func WriteProblem(w http.ResponseWriter, p Problem) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `p` | `Problem` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of WriteProblem
result := WriteProblem(/* parameters */)
```

### WriteValidationProblem
WriteValidationProblem writes a ValidationProblem response to the http.ResponseWriter.

```go
func WriteValidationProblem(w http.ResponseWriter, v *ValidationErrors) error
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `v` | `*ValidationErrors` | |

**Returns:**
| Type | Description |
|------|-------------|
| `error` | |

**Example:**

```go
// Example usage of WriteValidationProblem
result := WriteValidationProblem(/* parameters */)
```

## External Links

- [Package Overview](../packages/helix.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/helix)
- [Source Code](https://github.com/kolosys/helix/tree/main/github.com/kolosys/helix)
