# Logger Middleware

Logs HTTP requests with configurable formats. Supports multiple output formats including JSON, Apache combined/common, and development formats.

## Basic Usage

```go
// Development format (colorized)
s.Use(middleware.Logger(middleware.LogFormatDev))

// JSON format
s.Use(middleware.Logger(middleware.LogFormatJSON))

// Apache combined format
s.Use(middleware.Logger(middleware.LogFormatCombined))
```

## Available Formats

- `LogFormatDev` - Colorized development format
- `LogFormatJSON` - JSON format for structured logging
- `LogFormatCombined` - Apache combined log format
- `LogFormatCommon` - Apache common log format
- `LogFormatShort` - Shorter format
- `LogFormatTiny` - Minimal format

## Advanced Configuration

```go
s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format:     middleware.LogFormatJSON,
    Output:     os.Stdout,
    Skip:       func(r *http.Request) bool { return r.URL.Path == "/health" },
    TimeFormat: time.RFC3339,
    Fields: map[string]string{
        "api_version": "header:X-API-Version",
        "user_id":     "query:user_id",
    },
    CustomTokens: map[string]middleware.TokenExtractor{
        "user_id": func(r *http.Request, body []byte) string {
            // Extract from request body
            return extractUserID(body)
        },
    },
    CaptureBody: true,
    MaxBodySize: 64 * 1024, // 64KB
}))
```

## Log Tokens

Standard tokens available in log formats:

- `:method` - HTTP method
- `:url` - Request URL
- `:status` - Response status code
- `:response-time` - Response time in milliseconds
- `:remote-addr` - Client IP address
- `:res-length` - Response length in bytes
- `:referrer` - Referrer header
- `:user-agent` - User-Agent header
- `:date` - Request date/time

## Custom Fields

Extract custom fields from headers, query parameters, or cookies:

```go
s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: middleware.LogFormatJSON,
    Fields: map[string]string{
        "api_version": "header:X-API-Version",
        "user_id":     "query:user_id",
        "session":     "cookie:session_id",
    },
}))
```

## Body Capture

Capture request body for custom token extraction:

```go
s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format:      middleware.LogFormatJSON,
    CaptureBody: true,
    MaxBodySize: 64 * 1024, // 64KB
    CustomTokens: map[string]middleware.TokenExtractor{
        "user_id": func(r *http.Request, body []byte) string {
            // Parse JSON body and extract user_id
            var req struct {
                UserID string `json:"user_id"`
            }
            json.Unmarshal(body, &req)
            return req.UserID
        },
    },
}))
```

**Note**: When body capture is enabled, the body is read and stored. Handlers should not read the body again unless you handle it manually.

## Skipping Logs

Skip logging for specific requests:

```go
s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: middleware.LogFormatJSON,
    Skip: func(r *http.Request) bool {
        // Don't log health checks
        return r.URL.Path == "/health" || r.URL.Path == "/metrics"
    },
}))
```

## JSON Output

Configure JSON output fields:

```go
s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: middleware.LogFormatJSON,
    JSONFields: []string{
        "method",
        "url",
        "status",
        "response_time",
        "remote_addr",
    },
    JSONPretty: false, // Pretty print JSON
}))
```
