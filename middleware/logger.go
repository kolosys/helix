package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LogFormat represents a predefined log format.
type LogFormat string

// Predefined log formats matching Morgan.js formats.
const (
	// LogFormatCombined is the Apache combined log format.
	// :remote-addr - :remote-user [:date] ":method :url HTTP/:http-version" :status :res-length ":referrer" ":user-agent"
	LogFormatCombined LogFormat = "combined"

	// LogFormatCommon is the Apache common log format.
	// :remote-addr - :remote-user [:date] ":method :url HTTP/:http-version" :status :res-length
	LogFormatCommon LogFormat = "common"

	// LogFormatDev is a colorized development format.
	// :method :url :status :response-time ms - :res-length
	LogFormatDev LogFormat = "dev"

	// LogFormatShort is a shorter format.
	// :remote-addr :method :url HTTP/:http-version :status :res-length - :response-time ms
	LogFormatShort LogFormat = "short"

	// LogFormatTiny is the minimal format.
	// :method :url :status :res-length - :response-time ms
	LogFormatTiny LogFormat = "tiny"

	// LogFormatJSON outputs logs in JSON format.
	LogFormatJSON LogFormat = "json"
)

// TokenExtractor is a function that extracts a value from the request.
// It receives the request and the captured request body (if body capture is enabled).
type TokenExtractor func(r *http.Request, body []byte) string

// LoggerConfig configures the Logger middleware.
type LoggerConfig struct {
	// Format is the log format to use.
	// Default: LogFormatDev
	Format LogFormat

	// CustomFormat is a custom format string using tokens.
	// If set, Format is ignored (unless Format is LogFormatJSON).
	CustomFormat string

	// Output is the writer to output logs to.
	// Default: os.Stdout
	Output io.Writer

	// Skip is a function that determines if logging should be skipped.
	// If it returns true, the request is not logged.
	Skip func(r *http.Request) bool

	// TimeFormat is the time format for the :date token.
	// Default: time.RFC1123
	TimeFormat string

	// Fields maps custom field names to their sources.
	// Sources can be:
	//   - "header:X-Header-Name" - extracts from request header
	//   - "query:param_name" - extracts from query parameter
	//   - "cookie:cookie_name" - extracts from cookie
	// Example: {"api_version": "header:X-API-Version", "page": "query:page"}
	Fields map[string]string

	// CustomTokens maps token names to extractor functions.
	// These can extract data from the request body or perform custom logic.
	// Token names should not include the leading colon.
	// Example: {"user_id": func(r, body) string { ... }}
	CustomTokens map[string]TokenExtractor

	// CaptureBody enables capturing the request body for custom token extraction.
	// When enabled, the request body is read and stored for token extractors.
	// Default: false (only enable if you need body-based custom tokens)
	CaptureBody bool

	// MaxBodySize is the maximum size of the request body to capture.
	// Default: 64KB
	MaxBodySize int64

	// JSONFields specifies which fields to include in JSON output.
	// If empty, a default set of fields is used.
	// Fields can be standard tokens (without colon) or custom field names.
	JSONFields []string

	// JSONPretty enables pretty-printing for JSON output.
	// Default: false
	JSONPretty bool

	// DisableColors disables ANSI color codes in output.
	// Default: false
	DisableColors bool
}

// DefaultLoggerConfig returns the default configuration for Logger.
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Format:      LogFormatDev,
		Output:      os.Stdout,
		TimeFormat:  time.RFC1123,
		MaxBodySize: 64 << 10, // 64KB
	}
}

// Logger returns a middleware that logs HTTP requests.
// Uses the dev format by default.
func Logger(format LogFormat) Middleware {
	config := DefaultLoggerConfig()
	config.Format = format
	return LoggerWithConfig(config)
}

// LoggerJSON returns a middleware that logs HTTP requests in JSON format.
func LoggerJSON() Middleware {
	config := DefaultLoggerConfig()
	config.Format = LogFormatJSON
	return LoggerWithConfig(config)
}

// LoggerWithConfig returns a Logger middleware with the given configuration.
func LoggerWithConfig(config LoggerConfig) Middleware {
	if config.Output == nil {
		config.Output = os.Stdout
	}
	if config.TimeFormat == "" {
		config.TimeFormat = time.RFC1123
	}
	if config.Format == "" && config.CustomFormat == "" {
		config.Format = LogFormatDev
	}
	if config.MaxBodySize == 0 {
		config.MaxBodySize = 64 << 10
	}

	// Get the format string for non-JSON formats
	formatStr := config.CustomFormat
	if formatStr == "" && config.Format != LogFormatJSON {
		formatStr = getFormatString(config.Format)
	}

	// Precompile field extractors
	fieldExtractors := make(map[string]fieldExtractor)
	for name, source := range config.Fields {
		fieldExtractors[name] = parseFieldSource(source)
	}

	// Buffer pool for JSON encoding
	bufPool := &sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 512))
		},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if configured
			if config.Skip != nil && config.Skip(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Capture request body if needed
			var capturedBody []byte
			if config.CaptureBody && r.Body != nil && r.ContentLength > 0 {
				capturedBody = captureRequestBody(r, config.MaxBodySize)
			}

			// Record start time
			start := time.Now()

			// Wrap response writer
			rw := newResponseWriter(w)

			// Process request
			next.ServeHTTP(rw, r)

			// Calculate duration
			duration := time.Since(start)

			// Create log context
			ctx := &logContext{
				request:         r,
				responseWriter:  rw,
				duration:        duration,
				config:          config,
				capturedBody:    capturedBody,
				fieldExtractors: fieldExtractors,
			}

			// Output log
			if config.Format == LogFormatJSON {
				writeJSONLog(config.Output, ctx, bufPool)
			} else {
				line := formatLogLine(formatStr, ctx)
				fmt.Fprintln(config.Output, line)
			}
		})
	}
}

// LoggerWithFormat returns a Logger middleware with a custom format string.
func LoggerWithFormat(format string) Middleware {
	config := DefaultLoggerConfig()
	config.CustomFormat = format
	return LoggerWithConfig(config)
}

// LoggerWithFields returns a Logger middleware with custom fields.
func LoggerWithFields(fields map[string]string) Middleware {
	config := DefaultLoggerConfig()
	config.Fields = fields
	return LoggerWithConfig(config)
}

// logContext holds all the data needed for logging.
type logContext struct {
	request         *http.Request
	responseWriter  *responseWriter
	duration        time.Duration
	config          LoggerConfig
	capturedBody    []byte
	fieldExtractors map[string]fieldExtractor
}

// fieldExtractor extracts a value from a request.
type fieldExtractor struct {
	source string // "header", "query", "cookie"
	key    string // the header/query/cookie name
}

// parseFieldSource parses a field source string like "header:X-API-Version".
func parseFieldSource(source string) fieldExtractor {
	parts := strings.SplitN(source, ":", 2)
	if len(parts) != 2 {
		return fieldExtractor{source: "literal", key: source}
	}
	return fieldExtractor{source: parts[0], key: parts[1]}
}

// extract extracts the value from the request.
func (f fieldExtractor) extract(r *http.Request) string {
	switch f.source {
	case "header":
		return r.Header.Get(f.key)
	case "query":
		return r.URL.Query().Get(f.key)
	case "cookie":
		if c, err := r.Cookie(f.key); err == nil {
			return c.Value
		}
		return ""
	case "literal":
		return f.key
	default:
		return ""
	}
}

// captureRequestBody reads and stores the request body, then replaces it.
func captureRequestBody(r *http.Request, maxSize int64) []byte {
	if r.Body == nil {
		return nil
	}

	limitReader := io.LimitReader(r.Body, maxSize)
	body, err := io.ReadAll(limitReader)
	if err != nil {
		return nil
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	return body
}

// getFormatString returns the format string for a predefined format.
func getFormatString(format LogFormat) string {
	switch format {
	case LogFormatCombined:
		return `:remote-addr - :remote-user [:date] ":method :url HTTP/:http-version" :status :res-length ":referrer" ":user-agent"`
	case LogFormatCommon:
		return `:remote-addr - :remote-user [:date] ":method :url HTTP/:http-version" :status :res-length`
	case LogFormatDev:
		return `:method :url :status :response-time - :res-length`
	case LogFormatShort:
		return `:remote-addr :method :url HTTP/:http-version :status :res-length - :response-time`
	case LogFormatTiny:
		return `:method :url :status :res-length - :response-time`
	default:
		return `:method :url :status :response-time - :res-length`
	}
}

// Token pattern for dynamic tokens like :header[X-API-Version] or :query[page]
var dynamicTokenPattern = regexp.MustCompile(`:(\w+)\[([^\]]+)\]`)

// formatLogLine replaces tokens in the format string with actual values.
func formatLogLine(format string, ctx *logContext) string {
	r := ctx.request
	rw := ctx.responseWriter
	config := ctx.config

	line := format

	status := rw.Status()
	statusStr := strconv.Itoa(status)
	if config.Format == LogFormatDev && !config.DisableColors {
		statusStr = colorizeStatus(status)
	}

	method := r.Method
	if config.Format == LogFormatDev && !config.DisableColors {
		method = colorizeMethod(r.Method)
	}

	responseTime := formatDuration(ctx.duration)

	replacements := []struct {
		token string
		value string
	}{
		{":method", method},
		{":url", r.URL.RequestURI()},
		{":path", r.URL.Path},
		{":status", statusStr},
		{":response-time", responseTime},
		{":latency", responseTime},
		{":res-length", formatSize(rw.Size())},
		{":remote-addr", getRemoteAddr(r)},
		{":remote-user", getRemoteUser(r)},
		{":date", time.Now().Format(config.TimeFormat)},
		{":referrer", r.Referer()},
		{":user-agent", r.UserAgent()},
		{":http-version", fmt.Sprintf("%d.%d", r.ProtoMajor, r.ProtoMinor)},
		{":request-id", r.Header.Get(RequestIDHeader)},
		{":content-type", r.Header.Get("Content-Type")},
		{":content-length", strconv.FormatInt(r.ContentLength, 10)},
	}

	for _, rep := range replacements {
		line = strings.ReplaceAll(line, rep.token, rep.value)
	}

	// Replace dynamic tokens like :header[X-API-Version] or :query[page]
	line = dynamicTokenPattern.ReplaceAllStringFunc(line, func(match string) string {
		matches := dynamicTokenPattern.FindStringSubmatch(match)
		if len(matches) != 3 {
			return match
		}
		tokenType := matches[1]
		tokenKey := matches[2]

		switch tokenType {
		case "header":
			return r.Header.Get(tokenKey)
		case "query":
			return r.URL.Query().Get(tokenKey)
		case "cookie":
			if c, err := r.Cookie(tokenKey); err == nil {
				return c.Value
			}
			return ""
		default:
			return match
		}
	})

	// Replace custom field tokens (from Fields config)
	for name, extractor := range ctx.fieldExtractors {
		token := ":" + name
		line = strings.ReplaceAll(line, token, extractor.extract(r))
	}

	// Replace custom token extractors
	for name, extractor := range config.CustomTokens {
		token := ":" + name
		if strings.Contains(line, token) {
			value := extractor(r, ctx.capturedBody)
			line = strings.ReplaceAll(line, token, value)
		}
	}

	return line
}

// LogEntry represents a JSON log entry.
type LogEntry struct {
	Timestamp    string            `json:"timestamp"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	URL          string            `json:"url,omitempty"`
	Status       int               `json:"status"`
	Latency      string            `json:"latency"`
	LatencyMs    float64           `json:"latency_ms"`
	Size         int               `json:"size"`
	RemoteAddr   string            `json:"remote_addr"`
	UserAgent    string            `json:"user_agent,omitempty"`
	Referer      string            `json:"referer,omitempty"`
	RequestID    string            `json:"request_id,omitempty"`
	Error        string            `json:"error,omitempty"`
	CustomFields map[string]string `json:"custom,omitempty"`
}

// writeJSONLog writes a JSON-formatted log entry.
func writeJSONLog(w io.Writer, ctx *logContext, bufPool *sync.Pool) {
	r := ctx.request
	rw := ctx.responseWriter

	entry := LogEntry{
		Timestamp:  time.Now().Format(time.RFC3339),
		Method:     r.Method,
		Path:       r.URL.Path,
		URL:        r.URL.RequestURI(),
		Status:     rw.Status(),
		Latency:    formatDuration(ctx.duration),
		LatencyMs:  float64(ctx.duration.Microseconds()) / 1000.0,
		Size:       rw.Size(),
		RemoteAddr: getRemoteAddr(r),
		UserAgent:  r.UserAgent(),
		Referer:    r.Referer(),
		RequestID:  r.Header.Get(RequestIDHeader),
	}

	// Add custom fields from Fields config
	if len(ctx.fieldExtractors) > 0 || len(ctx.config.CustomTokens) > 0 {
		entry.CustomFields = make(map[string]string)

		for name, extractor := range ctx.fieldExtractors {
			if value := extractor.extract(r); value != "" {
				entry.CustomFields[name] = value
			}
		}

		for name, extractor := range ctx.config.CustomTokens {
			if value := extractor(r, ctx.capturedBody); value != "" {
				entry.CustomFields[name] = value
			}
		}
	}

	// Get buffer from pool
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	// Encode JSON
	encoder := json.NewEncoder(buf)
	if ctx.config.JSONPretty {
		encoder.SetIndent("", "  ")
	}
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(entry); err != nil {
		fmt.Fprintf(w, `{"error":"failed to encode log: %s"}`+"\n", err.Error())
		return
	}

	w.Write(buf.Bytes())
}

// ANSI color codes
const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

// colorizeStatus returns a colorized status code.
func colorizeStatus(status int) string {
	switch {
	case status >= 500:
		return red + strconv.Itoa(status) + reset
	case status >= 400:
		return yellow + strconv.Itoa(status) + reset
	case status >= 300:
		return cyan + strconv.Itoa(status) + reset
	case status >= 200:
		return green + strconv.Itoa(status) + reset
	default:
		return strconv.Itoa(status)
	}
}

// colorizeMethod returns a colorized HTTP method.
func colorizeMethod(method string) string {
	switch method {
	case http.MethodGet:
		return blue + method + reset
	case http.MethodPost:
		return cyan + method + reset
	case http.MethodPut:
		return yellow + method + reset
	case http.MethodDelete:
		return red + method + reset
	case http.MethodPatch:
		return green + method + reset
	case http.MethodHead:
		return magenta + method + reset
	case http.MethodOptions:
		return white + method + reset
	default:
		return method
	}
}

// formatDuration formats a duration for logging.
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d.Microseconds()))
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// formatSize formats a byte size for logging.
func formatSize(size int) string {
	if size < 0 {
		return "-"
	}
	return strconv.Itoa(size)
}

// getRemoteAddr gets the remote address from the request.
func getRemoteAddr(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Get first IP from comma-separated list
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	addr := r.RemoteAddr
	// Strip port if present
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		// Check if it's an IPv6 address
		if strings.Count(addr, ":") > 1 {
			// IPv6: [::1]:8080
			if bracket := strings.LastIndex(addr, "]"); bracket != -1 && idx > bracket {
				return addr[:idx]
			}
		} else {
			return addr[:idx]
		}
	}
	return addr
}

// getRemoteUser gets the remote user from Basic Auth (if present).
func getRemoteUser(r *http.Request) string {
	user, _, ok := r.BasicAuth()
	if ok {
		return user
	}
	return "-"
}

// JSONBodyExtractor creates a token extractor that extracts a field from JSON body.
// The path can be a simple field name like "user_id" or a nested path like "user.id".
func JSONBodyExtractor(path string) TokenExtractor {
	parts := strings.Split(path, ".")

	return func(r *http.Request, body []byte) string {
		if len(body) == 0 {
			return ""
		}

		var data map[string]any
		if err := json.Unmarshal(body, &data); err != nil {
			return ""
		}

		var current any = data
		for _, part := range parts {
			if m, ok := current.(map[string]any); ok {
				current = m[part]
			} else {
				return ""
			}
		}

		switch v := current.(type) {
		case string:
			return v
		case float64:
			// Check if it's an integer
			if v == float64(int64(v)) {
				return strconv.FormatInt(int64(v), 10)
			}
			return strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(v)
		case nil:
			return ""
		default:
			// For complex types, return JSON representation
			if b, err := json.Marshal(v); err == nil {
				return string(b)
			}
			return fmt.Sprintf("%v", v)
		}
	}
}

// FormValueExtractor creates a token extractor that extracts a form field.
func FormValueExtractor(field string) TokenExtractor {
	return func(r *http.Request, body []byte) string {
		// Parse form if not already parsed
		if r.Form == nil {
			r.ParseForm()
		}
		return r.FormValue(field)
	}
}

// ContextValueExtractor creates a token extractor that extracts a value from request context.
func ContextValueExtractor(key any) TokenExtractor {
	return func(r *http.Request, body []byte) string {
		if v := r.Context().Value(key); v != nil {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}
}
