# middleware API

Complete API documentation for the middleware package.

**Import Path:** `github.com/kolosys/helix/middleware`

## Package Documentation

Package middleware provides HTTP middleware for the Helix framework.


## Constants

### RequestIDHeader

RequestIDHeader is the default header name for the request ID.


```go
&{<nil> [RequestIDHeader] <nil> [0xc0002db600] <nil>}
```

## Types

### BasicAuthConfig
BasicAuthConfig configures the BasicAuth middleware.

#### Example Usage

```go
// Create a new BasicAuthConfig
basicauthconfig := BasicAuthConfig{
    Validator: /* value */,
    Realm: "example",
    SkipFunc: /* value */,
}
```

#### Type Definition

```go
type BasicAuthConfig struct {
    Validator func(username, password string) bool
    Realm string
    SkipFunc func(r *http.Request) bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Validator | `func(username, password string) bool` | Validator is a function that validates the username and password. Return true if the credentials are valid. |
| Realm | `string` | Realm is the authentication realm displayed in the browser. Default: "Restricted" |
| SkipFunc | `func(r *http.Request) bool` | SkipFunc determines if authentication should be skipped. |

### CORSConfig
CORSConfig configures the CORS middleware.

#### Example Usage

```go
// Create a new CORSConfig
corsconfig := CORSConfig{
    AllowOrigins: [],
    AllowOriginFunc: /* value */,
    AllowMethods: [],
    AllowHeaders: [],
    ExposeHeaders: [],
    AllowCredentials: true,
    MaxAge: 42,
}
```

#### Type Definition

```go
type CORSConfig struct {
    AllowOrigins []string
    AllowOriginFunc func(origin string) bool
    AllowMethods []string
    AllowHeaders []string
    ExposeHeaders []string
    AllowCredentials bool
    MaxAge int
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| AllowOrigins | `[]string` | AllowOrigins is a list of origins that are allowed. Use "*" to allow all origins. Default: [] |
| AllowOriginFunc | `func(origin string) bool` | AllowOriginFunc is a custom function to validate the origin. If set, AllowOrigins is ignored. |
| AllowMethods | `[]string` | AllowMethods is a list of methods that are allowed. Default: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS |
| AllowHeaders | `[]string` | AllowHeaders is a list of headers that are allowed. Default: Origin, Content-Type, Accept, Authorization |
| ExposeHeaders | `[]string` | ExposeHeaders is a list of headers that are exposed to the client. Default: [] |
| AllowCredentials | `bool` | AllowCredentials indicates whether credentials are allowed. Default: false |
| MaxAge | `int` | MaxAge is the maximum age (in seconds) of the preflight cache. Default: 0 (no caching) |

### Constructor Functions

### DefaultCORSConfig

DefaultCORSConfig returns the default CORS configuration.

```go
func DefaultCORSConfig() CORSConfig
```

**Parameters:**
  None

**Returns:**
- CORSConfig

### CacheConfig
CacheConfig configures the Cache middleware.

#### Example Usage

```go
// Create a new CacheConfig
cacheconfig := CacheConfig{
    MaxAge: 42,
    SMaxAge: 42,
    Public: true,
    Private: true,
    NoCache: true,
    NoStore: true,
    NoTransform: true,
    MustRevalidate: true,
    ProxyRevalidate: true,
    Immutable: true,
    StaleWhileRevalidate: 42,
    StaleIfError: 42,
    SkipFunc: /* value */,
    VaryHeaders: [],
}
```

#### Type Definition

```go
type CacheConfig struct {
    MaxAge int
    SMaxAge int
    Public bool
    Private bool
    NoCache bool
    NoStore bool
    NoTransform bool
    MustRevalidate bool
    ProxyRevalidate bool
    Immutable bool
    StaleWhileRevalidate int
    StaleIfError int
    SkipFunc func(r *http.Request) bool
    VaryHeaders []string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| MaxAge | `int` | MaxAge sets the max-age directive in seconds. Default: 0 (not set) |
| SMaxAge | `int` | SMaxAge sets the s-maxage directive in seconds (for shared caches). Default: 0 (not set) |
| Public | `bool` | Public indicates the response can be cached by any cache. Default: false |
| Private | `bool` | Private indicates the response is for a single user. Default: false |
| NoCache | `bool` | NoCache indicates the response must be revalidated before use. Default: false |
| NoStore | `bool` | NoStore indicates the response must not be stored. Default: false |
| NoTransform | `bool` | NoTransform indicates the response must not be transformed. Default: false |
| MustRevalidate | `bool` | MustRevalidate indicates stale responses must be revalidated. Default: false |
| ProxyRevalidate | `bool` | ProxyRevalidate is like MustRevalidate but for shared caches. Default: false |
| Immutable | `bool` | Immutable indicates the response body will not change. Default: false |
| StaleWhileRevalidate | `int` | StaleWhileRevalidate allows serving stale content while revalidating. Value is in seconds. Default: 0 (not set) |
| StaleIfError | `int` | StaleIfError allows serving stale content if there's an error. Value is in seconds. Default: 0 (not set) |
| SkipFunc | `func(r *http.Request) bool` | SkipFunc determines if cache headers should be skipped. |
| VaryHeaders | `[]string` | VaryHeaders is a list of headers to include in the Vary header. |

### Constructor Functions

### DefaultCacheConfig

DefaultCacheConfig returns the default Cache configuration.

```go
func DefaultCacheConfig() CacheConfig
```

**Parameters:**
  None

**Returns:**
- CacheConfig

### CompressConfig
CompressConfig configures the Compress middleware.

#### Example Usage

```go
// Create a new CompressConfig
compressconfig := CompressConfig{
    Level: 42,
    MinSize: 42,
    Types: [],
    SkipFunc: /* value */,
}
```

#### Type Definition

```go
type CompressConfig struct {
    Level int
    MinSize int
    Types []string
    SkipFunc func(r *http.Request) bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Level | `int` | Level is the compression level. Valid levels: -1 (default), 0 (no compression), 1 (best speed) to 9 (best compression) Default: -1 (gzip.DefaultCompression) |
| MinSize | `int` | MinSize is the minimum size in bytes to trigger compression. Default: 1024 (1KB) |
| Types | `[]string` | Types is a list of content types to compress. Default: text/*, application/json, application/javascript, application/xml |
| SkipFunc | `func(r *http.Request) bool` | SkipFunc is a function that determines if compression should be skipped. |

### Constructor Functions

### DefaultCompressConfig

DefaultCompressConfig returns the default Compress configuration.

```go
func DefaultCompressConfig() CompressConfig
```

**Parameters:**
  None

**Returns:**
- CompressConfig

### ETagConfig
ETagConfig configures the ETag middleware.

#### Example Usage

```go
// Create a new ETagConfig
etagconfig := ETagConfig{
    Weak: true,
    SkipFunc: /* value */,
}
```

#### Type Definition

```go
type ETagConfig struct {
    Weak bool
    SkipFunc func(r *http.Request) bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Weak | `bool` | Weak indicates whether to generate weak ETags. Weak ETags are prefixed with W/ and indicate semantic equivalence. Default: false (strong ETags) |
| SkipFunc | `func(r *http.Request) bool` | SkipFunc determines if ETag generation should be skipped. |

### Constructor Functions

### DefaultETagConfig

DefaultETagConfig returns the default ETag configuration.

```go
func DefaultETagConfig() ETagConfig
```

**Parameters:**
  None

**Returns:**
- ETagConfig

### LogEntry
LogEntry represents a JSON log entry.

#### Example Usage

```go
// Create a new LogEntry
logentry := LogEntry{
    Timestamp: "example",
    Method: "example",
    Path: "example",
    URL: "example",
    Status: 42,
    Latency: "example",
    LatencyMs: 3.14,
    Size: 42,
    RemoteAddr: "example",
    UserAgent: "example",
    Referer: "example",
    RequestID: "example",
    Error: "example",
    CustomFields: map[],
}
```

#### Type Definition

```go
type LogEntry struct {
    Timestamp string `json:"timestamp"`
    Method string `json:"method"`
    Path string `json:"path"`
    URL string `json:"url,omitempty"`
    Status int `json:"status"`
    Latency string `json:"latency"`
    LatencyMs float64 `json:"latency_ms"`
    Size int `json:"size"`
    RemoteAddr string `json:"remote_addr"`
    UserAgent string `json:"user_agent,omitempty"`
    Referer string `json:"referer,omitempty"`
    RequestID string `json:"request_id,omitempty"`
    Error string `json:"error,omitempty"`
    CustomFields map[string]string `json:"custom,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Timestamp | `string` |  |
| Method | `string` |  |
| Path | `string` |  |
| URL | `string` |  |
| Status | `int` |  |
| Latency | `string` |  |
| LatencyMs | `float64` |  |
| Size | `int` |  |
| RemoteAddr | `string` |  |
| UserAgent | `string` |  |
| Referer | `string` |  |
| RequestID | `string` |  |
| Error | `string` |  |
| CustomFields | `map[string]string` |  |

### LogFormat
LogFormat represents a predefined log format.

#### Example Usage

```go
// Example usage of LogFormat
var value LogFormat
// Initialize with appropriate value
```

#### Type Definition

```go
type LogFormat string
```

### LoggerConfig
LoggerConfig configures the Logger middleware.

#### Example Usage

```go
// Create a new LoggerConfig
loggerconfig := LoggerConfig{
    Format: LogFormat{},
    CustomFormat: "example",
    Output: /* value */,
    Skip: /* value */,
    TimeFormat: "example",
    Fields: map[],
    CustomTokens: map[],
    CaptureBody: true,
    MaxBodySize: 42,
    JSONFields: [],
    JSONPretty: true,
    DisableColors: true,
}
```

#### Type Definition

```go
type LoggerConfig struct {
    Format LogFormat
    CustomFormat string
    Output io.Writer
    Skip func(r *http.Request) bool
    TimeFormat string
    Fields map[string]string
    CustomTokens map[string]TokenExtractor
    CaptureBody bool
    MaxBodySize int64
    JSONFields []string
    JSONPretty bool
    DisableColors bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Format | `LogFormat` | Format is the log format to use. Default: LogFormatDev |
| CustomFormat | `string` | CustomFormat is a custom format string using tokens. If set, Format is ignored (unless Format is LogFormatJSON). |
| Output | `io.Writer` | Output is the writer to output logs to. Default: os.Stdout |
| Skip | `func(r *http.Request) bool` | Skip is a function that determines if logging should be skipped. If it returns true, the request is not logged. |
| TimeFormat | `string` | TimeFormat is the time format for the :date token. Default: time.RFC1123 |
| Fields | `map[string]string` | Fields maps custom field names to their sources. Sources can be: - "header:X-Header-Name" - extracts from request header - "query:param_name" - extracts from query parameter - "cookie:cookie_name" - extracts from cookie Example: {"api_version": "header:X-API-Version", "page": "query:page"} |
| CustomTokens | `map[string]TokenExtractor` | CustomTokens maps token names to extractor functions. These can extract data from the request body or perform custom logic. Token names should not include the leading colon. Example: {"user_id": func(r, body) string { ... }} |
| CaptureBody | `bool` | CaptureBody enables capturing the request body for custom token extraction. When enabled, the request body is read and stored for token extractors. Default: false (only enable if you need body-based custom tokens) |
| MaxBodySize | `int64` | MaxBodySize is the maximum size of the request body to capture. Default: 64KB |
| JSONFields | `[]string` | JSONFields specifies which fields to include in JSON output. If empty, a default set of fields is used. Fields can be standard tokens (without colon) or custom field names. |
| JSONPretty | `bool` | JSONPretty enables pretty-printing for JSON output. Default: false |
| DisableColors | `bool` | DisableColors disables ANSI color codes in output. Default: false |

### Constructor Functions

### DefaultLoggerConfig

DefaultLoggerConfig returns the default configuration for Logger.

```go
func DefaultLoggerConfig() LoggerConfig
```

**Parameters:**
  None

**Returns:**
- LoggerConfig

### Middleware
Middleware is a function that wraps an http.Handler to provide additional functionality.

#### Example Usage

```go
// Example usage of Middleware
var value Middleware
// Initialize with appropriate value
```

#### Type Definition

```go
type Middleware func(next http.Handler) http.Handler
```

### Constructor Functions

### API

API returns a middleware bundle suitable for JSON API servers. Includes: RequestID, Logger (JSON format), Recover, and CORS.

```go
func API() []Middleware
```

**Parameters:**
  None

**Returns:**
- []Middleware

### APIWithCORS

APIWithCORS returns a middleware bundle suitable for JSON API servers with a custom CORS configuration. Includes: RequestID, Logger (JSON format), Recover, and CORS with config.

```go
func APIWithCORS(cors CORSConfig) []Middleware
```

**Parameters:**
- `cors` (CORSConfig)

**Returns:**
- []Middleware

### BasicAuth

BasicAuth returns a BasicAuth middleware with the given username and password. Uses constant-time comparison to prevent timing attacks.

```go
func BasicAuth(username, password string) Middleware
```

**Parameters:**
- `username` (string)
- `password` (string)

**Returns:**
- Middleware

### BasicAuthUsers

BasicAuthUsers returns a BasicAuth middleware that validates against a map of users. The map key is the username and the value is the password.

```go
func BasicAuthUsers(users map[string]string) Middleware
```

**Parameters:**
- `users` (map[string]string)

**Returns:**
- Middleware

### BasicAuthWithConfig

BasicAuthWithConfig returns a BasicAuth middleware with the given configuration.

```go
func BasicAuthWithConfig(config BasicAuthConfig) Middleware
```

**Parameters:**
- `config` (BasicAuthConfig)

**Returns:**
- Middleware

### BasicAuthWithValidator

BasicAuthWithValidator returns a BasicAuth middleware with a custom validator.

```go
func BasicAuthWithValidator(validator func(username, password string) bool) Middleware
```

**Parameters:**
- `validator` (func(username, password string) bool)

**Returns:**
- Middleware

### CORS

CORS returns a CORS middleware with default configuration.

```go
func CORS() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### CORSAllowAll

CORSAllowAll returns a CORS middleware that allows all origins, methods, and headers.

```go
func CORSAllowAll() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### CORSWithConfig

CORSWithConfig returns a CORS middleware with the given configuration.

```go
func CORSWithConfig(config CORSConfig) Middleware
```

**Parameters:**
- `config` (CORSConfig)

**Returns:**
- Middleware

### Cache

Cache returns a Cache middleware with the given max-age in seconds.

```go
func Cache(maxAge int) Middleware
```

**Parameters:**
- `maxAge` (int)

**Returns:**
- Middleware

### CacheImmutable

CacheImmutable returns a Cache middleware for immutable content.

```go
func CacheImmutable(maxAge int) Middleware
```

**Parameters:**
- `maxAge` (int)

**Returns:**
- Middleware

### CachePrivate

CachePrivate returns a Cache middleware with private caching.

```go
func CachePrivate(maxAge int) Middleware
```

**Parameters:**
- `maxAge` (int)

**Returns:**
- Middleware

### CachePublic

CachePublic returns a Cache middleware with public caching.

```go
func CachePublic(maxAge int) Middleware
```

**Parameters:**
- `maxAge` (int)

**Returns:**
- Middleware

### CacheWithConfig

CacheWithConfig returns a Cache middleware with the given configuration.

```go
func CacheWithConfig(config CacheConfig) Middleware
```

**Parameters:**
- `config` (CacheConfig)

**Returns:**
- Middleware

### Chain

Chain creates a new middleware chain from the given middlewares. The first middleware in the chain is the outermost (executed first on request, last on response).

```go
func Chain(middlewares ...Middleware) Middleware
```

**Parameters:**
- `middlewares` (...Middleware)

**Returns:**
- Middleware

### Compress

Compress returns a middleware that compresses responses using gzip or deflate.

```go
func Compress() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### CompressWithConfig

CompressWithConfig returns a Compress middleware with the given configuration.

```go
func CompressWithConfig(config CompressConfig) Middleware
```

**Parameters:**
- `config` (CompressConfig)

**Returns:**
- Middleware

### CompressWithLevel

CompressWithLevel returns a Compress middleware with the given compression level.

```go
func CompressWithLevel(level int) Middleware
```

**Parameters:**
- `level` (int)

**Returns:**
- Middleware

### Development

Development returns a middleware bundle suitable for development. Includes: RequestID, Logger (dev format), Recover. This is the same as what helix.Default() uses.

```go
func Development() []Middleware
```

**Parameters:**
  None

**Returns:**
- []Middleware

### ETag

ETag returns an ETag middleware with default configuration.

```go
func ETag() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### ETagWeak

ETagWeak returns an ETag middleware that generates weak ETags.

```go
func ETagWeak() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### ETagWithConfig

ETagWithConfig returns an ETag middleware with the given configuration.

```go
func ETagWithConfig(config ETagConfig) Middleware
```

**Parameters:**
- `config` (ETagConfig)

**Returns:**
- Middleware

### Logger

Logger returns a middleware that logs HTTP requests. Uses the dev format by default.

```go
func Logger(format LogFormat) Middleware
```

**Parameters:**
- `format` (LogFormat)

**Returns:**
- Middleware

### LoggerJSON

LoggerJSON returns a middleware that logs HTTP requests in JSON format.

```go
func LoggerJSON() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### LoggerWithConfig

LoggerWithConfig returns a Logger middleware with the given configuration.

```go
func LoggerWithConfig(config LoggerConfig) Middleware
```

**Parameters:**
- `config` (LoggerConfig)

**Returns:**
- Middleware

### LoggerWithFields

LoggerWithFields returns a Logger middleware with custom fields.

```go
func LoggerWithFields(fields map[string]string) Middleware
```

**Parameters:**
- `fields` (map[string]string)

**Returns:**
- Middleware

### LoggerWithFormat

LoggerWithFormat returns a Logger middleware with a custom format string.

```go
func LoggerWithFormat(format string) Middleware
```

**Parameters:**
- `format` (string)

**Returns:**
- Middleware

### Minimal

Minimal returns a minimal middleware bundle with only essential middleware. Includes: Recover.

```go
func Minimal() []Middleware
```

**Parameters:**
  None

**Returns:**
- []Middleware

### NoCache

NoCache returns a middleware that disables caching.

```go
func NoCache() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### Production

Production returns a middleware bundle suitable for production environments. Includes: RequestID, Logger (combined format), Recover.

```go
func Production() []Middleware
```

**Parameters:**
  None

**Returns:**
- []Middleware

### RateLimit

RateLimit returns a rate limiting middleware with the given rate and burst.

```go
func RateLimit(rate float64, burst int) Middleware
```

**Parameters:**
- `rate` (float64)
- `burst` (int)

**Returns:**
- Middleware

### RateLimitWithConfig

RateLimitWithConfig returns a RateLimit middleware with the given configuration.

```go
func RateLimitWithConfig(config RateLimitConfig) Middleware
```

**Parameters:**
- `config` (RateLimitConfig)

**Returns:**
- Middleware

### Recover

Recover returns a middleware that recovers from panics. It logs the panic and stack trace, then returns a 500 Internal Server Error.

```go
func Recover() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### RecoverWithConfig

RecoverWithConfig returns a Recover middleware with the given configuration.

```go
func RecoverWithConfig(config RecoverConfig) Middleware
```

**Parameters:**
- `config` (RecoverConfig)

**Returns:**
- Middleware

### RequestID

RequestID returns a middleware that generates or propagates a request ID. The request ID is stored in the request context and the response header.

```go
func RequestID() Middleware
```

**Parameters:**
  None

**Returns:**
- Middleware

### RequestIDWithConfig

RequestIDWithConfig returns a RequestID middleware with the given configuration.

```go
func RequestIDWithConfig(config RequestIDConfig) Middleware
```

**Parameters:**
- `config` (RequestIDConfig)

**Returns:**
- Middleware

### Secure

Secure returns a middleware bundle with security-focused middleware. Includes: RequestID, Logger (JSON format), Recover, RateLimit. Note: You should also add CORS and authentication middleware as needed.

```go
func Secure(rate float64, burst int) []Middleware
```

**Parameters:**

- `rate` (float64) - requests per second allowed

- `burst` (int) - maximum burst size

**Returns:**
- []Middleware

### Timeout

Timeout returns a middleware that adds a timeout to requests.

```go
func Timeout(timeout time.Duration) Middleware
```

**Parameters:**
- `timeout` (time.Duration)

**Returns:**
- Middleware

### TimeoutWithConfig

TimeoutWithConfig returns a Timeout middleware with the given configuration.

```go
func TimeoutWithConfig(config TimeoutConfig) Middleware
```

**Parameters:**
- `config` (TimeoutConfig)

**Returns:**
- Middleware

### Web

Web returns a middleware bundle suitable for web applications. Includes: RequestID, Logger (dev format), Recover, and Compress.

```go
func Web() []Middleware
```

**Parameters:**
  None

**Returns:**
- []Middleware

### RateLimitConfig
RateLimitConfig configures the RateLimit middleware.

#### Example Usage

```go
// Create a new RateLimitConfig
ratelimitconfig := RateLimitConfig{
    Rate: 3.14,
    Burst: 42,
    KeyFunc: /* value */,
    Handler: /* value */,
    SkipFunc: /* value */,
    CleanupInterval: /* value */,
    ExpirationTime: /* value */,
}
```

#### Type Definition

```go
type RateLimitConfig struct {
    Rate float64
    Burst int
    KeyFunc func(r *http.Request) string
    Handler http.HandlerFunc
    SkipFunc func(r *http.Request) bool
    CleanupInterval time.Duration
    ExpirationTime time.Duration
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Rate | `float64` | Rate is the number of requests allowed per second. Default: 100 |
| Burst | `int` | Burst is the maximum number of requests allowed in a burst. Default: 10 |
| KeyFunc | `func(r *http.Request) string` | KeyFunc extracts the rate limit key from the request. Default: uses client IP address |
| Handler | `http.HandlerFunc` | Handler is called when the rate limit is exceeded. If nil, a default 429 Too Many Requests response is sent. |
| SkipFunc | `func(r *http.Request) bool` | SkipFunc determines if rate limiting should be skipped. |
| CleanupInterval | `time.Duration` | CleanupInterval is the interval for cleaning up expired entries. Default: 1 minute |
| ExpirationTime | `time.Duration` | ExpirationTime is how long to keep entries after last access. Default: 5 minutes |

### Constructor Functions

### DefaultRateLimitConfig

DefaultRateLimitConfig returns the default RateLimit configuration.

```go
func DefaultRateLimitConfig() RateLimitConfig
```

**Parameters:**
  None

**Returns:**
- RateLimitConfig

### RecoverConfig
RecoverConfig configures the Recover middleware.

#### Example Usage

```go
// Create a new RecoverConfig
recoverconfig := RecoverConfig{
    PrintStack: true,
    StackSize: 42,
    Output: /* value */,
    Handler: /* value */,
}
```

#### Type Definition

```go
type RecoverConfig struct {
    PrintStack bool
    StackSize int
    Output io.Writer
    Handler func(w http.ResponseWriter, r *http.Request, err any)
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| PrintStack | `bool` | PrintStack enables printing the stack trace when a panic occurs. Default: true |
| StackSize | `int` | StackSize is the maximum size of the stack trace buffer. Default: 4KB |
| Output | `io.Writer` | Output is the writer to output the panic message to. Default: os.Stderr |
| Handler | `func(w http.ResponseWriter, r *http.Request, err any)` | Handler is a custom function to handle panics. If set, it will be called instead of the default behavior. The handler should write the response and return. |

### Constructor Functions

### DefaultRecoverConfig

DefaultRecoverConfig returns the default configuration for Recover.

```go
func DefaultRecoverConfig() RecoverConfig
```

**Parameters:**
  None

**Returns:**
- RecoverConfig

### RequestIDConfig
RequestIDConfig configures the RequestID middleware.

#### Example Usage

```go
// Create a new RequestIDConfig
requestidconfig := RequestIDConfig{
    Header: "example",
    Generator: /* value */,
    TargetHeader: "example",
}
```

#### Type Definition

```go
type RequestIDConfig struct {
    Header string
    Generator func() string
    TargetHeader string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Header | `string` | Header is the name of the header to read/write the request ID. Default: "X-Request-ID" |
| Generator | `func() string` | Generator is a function that generates a new request ID. Default: generates a random 16-byte hex string |
| TargetHeader | `string` | TargetHeader is the header name to set on the response. Default: same as Header |

### Constructor Functions

### DefaultRequestIDConfig

DefaultRequestIDConfig returns the default configuration for RequestID.

```go
func DefaultRequestIDConfig() RequestIDConfig
```

**Parameters:**
  None

**Returns:**
- RequestIDConfig

### TimeoutConfig
TimeoutConfig configures the Timeout middleware.

#### Example Usage

```go
// Create a new TimeoutConfig
timeoutconfig := TimeoutConfig{
    Timeout: /* value */,
    Handler: /* value */,
    SkipFunc: /* value */,
}
```

#### Type Definition

```go
type TimeoutConfig struct {
    Timeout time.Duration
    Handler http.HandlerFunc
    SkipFunc func(r *http.Request) bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Timeout | `time.Duration` | Timeout is the maximum duration for the request. Default: 30 seconds |
| Handler | `http.HandlerFunc` | Handler is called when the request times out. If nil, a default 503 Service Unavailable response is sent. |
| SkipFunc | `func(r *http.Request) bool` | SkipFunc is a function that determines if timeout should be skipped. If it returns true, no timeout is applied. |

### Constructor Functions

### DefaultTimeoutConfig

DefaultTimeoutConfig returns the default Timeout configuration.

```go
func DefaultTimeoutConfig() TimeoutConfig
```

**Parameters:**
  None

**Returns:**
- TimeoutConfig

### TokenExtractor
TokenExtractor is a function that extracts a value from the request. It receives the request and the captured request body (if body capture is enabled).

#### Example Usage

```go
// Example usage of TokenExtractor
var value TokenExtractor
// Initialize with appropriate value
```

#### Type Definition

```go
type TokenExtractor func(r *http.Request, body []byte) string
```

### Constructor Functions

### ContextValueExtractor

ContextValueExtractor creates a token extractor that extracts a value from request context.

```go
func ContextValueExtractor(key any) TokenExtractor
```

**Parameters:**
- `key` (any)

**Returns:**
- TokenExtractor

### FormValueExtractor

FormValueExtractor creates a token extractor that extracts a form field.

```go
func FormValueExtractor(field string) TokenExtractor
```

**Parameters:**
- `field` (string)

**Returns:**
- TokenExtractor

### JSONBodyExtractor

JSONBodyExtractor creates a token extractor that extracts a field from JSON body. The path can be a simple field name like "user_id" or a nested path like "user.id".

```go
func JSONBodyExtractor(path string) TokenExtractor
```

**Parameters:**
- `path` (string)

**Returns:**
- TokenExtractor

## Functions

### ETagFromContent
ETagFromContent generates an ETag from content.

```go
func ETagFromContent(content []byte, weak bool) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `content` | `[]byte` | |
| `weak` | `bool` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of ETagFromContent
result := ETagFromContent(/* parameters */)
```

### ETagFromString
ETagFromString generates an ETag from a string.

```go
func ETagFromString(s string, weak bool) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `s` | `string` | |
| `weak` | `bool` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of ETagFromString
result := ETagFromString(/* parameters */)
```

### ETagFromVersion
ETagFromVersion generates an ETag from a version number.

```go
func ETagFromVersion(version int64, weak bool) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `version` | `int64` | |
| `weak` | `bool` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of ETagFromVersion
result := ETagFromVersion(/* parameters */)
```

### GetRequestID
GetRequestID retrieves the request ID from the context. Returns an empty string if no request ID is set.

```go
func GetRequestID(ctx context.Context) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of GetRequestID
result := GetRequestID(/* parameters */)
```

### GetRequestIDFromRequest
GetRequestIDFromRequest retrieves the request ID from the request context.

```go
func GetRequestIDFromRequest(r *http.Request) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `r` | `*http.Request` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of GetRequestIDFromRequest
result := GetRequestIDFromRequest(/* parameters */)
```

### SetCacheControl
SetCacheControl sets the Cache-Control header on the response.

```go
func SetCacheControl(w http.ResponseWriter, value string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `value` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetCacheControl
result := SetCacheControl(/* parameters */)
```

### SetExpires
SetExpires sets the Expires header on the response.

```go
func SetExpires(w http.ResponseWriter, t time.Time)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `t` | `time.Time` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetExpires
result := SetExpires(/* parameters */)
```

### SetLastModified
SetLastModified sets the Last-Modified header on the response.

```go
func SetLastModified(w http.ResponseWriter, t time.Time)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `w` | `http.ResponseWriter` | |
| `t` | `time.Time` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetLastModified
result := SetLastModified(/* parameters */)
```

## External Links

- [Package Overview](../packages/middleware.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/helix/middleware)
- [Source Code](https://github.com/kolosys/helix/tree/main/middleware)
