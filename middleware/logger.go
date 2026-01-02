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
	LogFormatCombined LogFormat = "combined"
	LogFormatCommon   LogFormat = "common"
	LogFormatDev      LogFormat = "dev"
	LogFormatShort    LogFormat = "short"
	LogFormatTiny     LogFormat = "tiny"
	LogFormatJSON     LogFormat = "json"
)

// LogValues contains all extracted request/response data for logging.
type LogValues struct {
	Method        string
	Path          string
	URI           string
	Host          string
	Protocol      string
	RemoteIP      string
	UserAgent     string
	Referer       string
	ContentLength int64
	ContentType   string
	Status        int
	ResponseSize  int
	Latency       time.Duration
	Error         error
	RequestID     string
	StartTime     time.Time
	Headers       map[string]string
	QueryParams   map[string]string
	FormValues    map[string]string
	CustomFields  map[string]string
}

// LogOutputFunc is a callback that receives log values and outputs them.
// This is the single output mechanism - use helpers for common formats.
type LogOutputFunc func(v LogValues)

// TokenExtractor extracts a custom value from the request.
type TokenExtractor func(r *http.Request, body []byte) string

// LoggerConfig configures the Logger middleware.
type LoggerConfig struct {
	// Output is the callback that receives log values. Required.
	// Use TextOutput() for Morgan.js-style formatting.
	// Use helix.StructuredOutput() for logs package integration.
	// Or provide your own function for custom logging.
	Output LogOutputFunc

	// Skip determines if logging should be skipped for a request.
	Skip func(r *http.Request) bool

	// Fields maps custom field names to their sources.
	// Sources: "header:Name", "query:param", "cookie:name"
	Fields map[string]string

	// CustomTokens maps names to extractor functions for body/context data.
	CustomTokens map[string]TokenExtractor

	// LogHeaders specifies request headers to extract into Headers map.
	LogHeaders []string

	// LogQueryParams specifies query parameters to extract.
	LogQueryParams []string

	// LogFormValues specifies form values to extract.
	LogFormValues []string

	// CaptureBody enables request body capture for CustomTokens.
	CaptureBody bool

	// MaxBodySize limits captured body size. Default: 64KB.
	MaxBodySize int64
}

// Logger returns a middleware with dev format text output.
func Logger() Middleware {
	return LoggerWithConfig(LoggerConfig{
		Output: TextOutput(os.Stdout, LogFormatDev),
	})
}

// LoggerWithConfig returns a Logger middleware with the given configuration.
func LoggerWithConfig(config LoggerConfig) Middleware {
	if config.Output == nil {
		config.Output = TextOutput(os.Stdout, LogFormatDev)
	}
	if config.MaxBodySize == 0 {
		config.MaxBodySize = 64 << 10
	}

	// Precompile field extractors
	fieldExtractors := make(map[string]fieldExtractor)
	for name, source := range config.Fields {
		fieldExtractors[name] = parseFieldSource(source)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skip != nil && config.Skip(r) {
				next.ServeHTTP(w, r)
				return
			}

			var capturedBody []byte
			if config.CaptureBody && r.Body != nil && r.ContentLength > 0 {
				capturedBody = captureRequestBody(r, config.MaxBodySize)
			}

			start := time.Now()
			rw := newResponseWriter(w)

			next.ServeHTTP(rw, r)

			v := LogValues{
				Method:        r.Method,
				Path:          r.URL.Path,
				URI:           r.URL.RequestURI(),
				Host:          r.Host,
				Protocol:      r.Proto,
				RemoteIP:      getRemoteAddr(r),
				UserAgent:     r.UserAgent(),
				Referer:       r.Referer(),
				ContentLength: r.ContentLength,
				ContentType:   r.Header.Get("Content-Type"),
				Status:        rw.Status(),
				ResponseSize:  rw.Size(),
				Latency:       time.Since(start),
				RequestID:     r.Header.Get(RequestIDHeader),
				StartTime:     start,
			}

			// Extract headers
			if len(config.LogHeaders) > 0 {
				v.Headers = make(map[string]string, len(config.LogHeaders))
				for _, h := range config.LogHeaders {
					v.Headers[h] = r.Header.Get(h)
				}
			}

			// Extract query params
			if len(config.LogQueryParams) > 0 {
				v.QueryParams = make(map[string]string, len(config.LogQueryParams))
				query := r.URL.Query()
				for _, p := range config.LogQueryParams {
					v.QueryParams[p] = query.Get(p)
				}
			}

			// Extract form values
			if len(config.LogFormValues) > 0 {
				v.FormValues = make(map[string]string, len(config.LogFormValues))
				for _, f := range config.LogFormValues {
					v.FormValues[f] = r.FormValue(f)
				}
			}

			// Extract custom fields
			if len(fieldExtractors) > 0 || len(config.CustomTokens) > 0 {
				v.CustomFields = make(map[string]string)
				for name, ext := range fieldExtractors {
					if val := ext.extract(r); val != "" {
						v.CustomFields[name] = val
					}
				}
				for name, ext := range config.CustomTokens {
					if val := ext(r, capturedBody); val != "" {
						v.CustomFields[name] = val
					}
				}
			}

			config.Output(v)
		})
	}
}

// --- Text Output Helpers (Morgan.js style) ---

// TextOutputOptions configures text output formatting.
type TextOutputOptions struct {
	TimeFormat    string
	DisableColors bool
	JSONPretty    bool // for LogFormatJSON
}

// TextOutput returns a LogOutputFunc that writes Morgan.js-style formatted logs.
func TextOutput(w io.Writer, format LogFormat) LogOutputFunc {
	return TextOutputWithOptions(w, format, TextOutputOptions{})
}

// TextOutputCustom returns a LogOutputFunc with a custom format string.
func TextOutputCustom(w io.Writer, format string, opts ...TextOutputOptions) LogOutputFunc {
	var opt TextOutputOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.TimeFormat == "" {
		opt.TimeFormat = time.RFC1123
	}
	return textOutputFunc(w, format, opt)
}

// TextOutputWithOptions returns a LogOutputFunc with custom options.
func TextOutputWithOptions(w io.Writer, format LogFormat, opts TextOutputOptions) LogOutputFunc {
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.RFC1123
	}
	if format == LogFormatJSON {
		return jsonOutputFunc(w, opts)
	}
	return textOutputFunc(w, getFormatString(format), opts)
}

func textOutputFunc(w io.Writer, formatStr string, opts TextOutputOptions) LogOutputFunc {
	pattern := regexp.MustCompile(`:(\w+)\[([^\]]+)\]`)

	return func(v LogValues) {
		line := formatStr

		status := strconv.Itoa(v.Status)
		method := v.Method
		if !opts.DisableColors {
			status = colorizeStatus(v.Status)
			method = colorizeMethod(v.Method)
		}

		replacements := map[string]string{
			":method":         method,
			":url":            v.URI,
			":path":           v.Path,
			":status":         status,
			":response-time":  formatDuration(v.Latency),
			":latency":        formatDuration(v.Latency),
			":res-length":     formatSize(v.ResponseSize),
			":remote-addr":    v.RemoteIP,
			":remote-user":    "-",
			":date":           v.StartTime.Format(opts.TimeFormat),
			":referrer":       v.Referer,
			":user-agent":     v.UserAgent,
			":http-version":   formatHTTPVersion(v.Protocol),
			":request-id":     v.RequestID,
			":content-type":   v.ContentType,
			":content-length": strconv.FormatInt(v.ContentLength, 10),
		}

		for token, val := range replacements {
			line = strings.ReplaceAll(line, token, val)
		}

		// Dynamic tokens like :header[X-Name]
		line = pattern.ReplaceAllStringFunc(line, func(match string) string {
			m := pattern.FindStringSubmatch(match)
			if len(m) != 3 {
				return match
			}
			switch m[1] {
			case "header":
				return v.Headers[m[2]]
			case "query":
				return v.QueryParams[m[2]]
			default:
				return match
			}
		})

		// Custom fields as tokens
		for name, val := range v.CustomFields {
			line = strings.ReplaceAll(line, ":"+name, val)
		}

		fmt.Fprintln(w, line)
	}
}

func jsonOutputFunc(w io.Writer, opts TextOutputOptions) LogOutputFunc {
	var mu sync.Mutex
	bufPool := sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 0, 512)) }}

	return func(v LogValues) {
		entry := map[string]any{
			"timestamp":  v.StartTime.Format(time.RFC3339),
			"method":     v.Method,
			"path":       v.Path,
			"status":     v.Status,
			"latency":    formatDuration(v.Latency),
			"latency_ms": float64(v.Latency.Microseconds()) / 1000.0,
			"size":       v.ResponseSize,
			"remote_ip":  v.RemoteIP,
		}
		if v.RequestID != "" {
			entry["request_id"] = v.RequestID
		}
		if v.UserAgent != "" {
			entry["user_agent"] = v.UserAgent
		}
		if len(v.CustomFields) > 0 {
			entry["custom"] = v.CustomFields
		}

		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufPool.Put(buf)

		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if opts.JSONPretty {
			enc.SetIndent("", "  ")
		}
		enc.Encode(entry)

		mu.Lock()
		w.Write(buf.Bytes())
		mu.Unlock()
	}
}

// --- Format Helpers ---

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

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return strconv.FormatFloat(float64(d.Microseconds()), 'f', 2, 64) + "µs"
	}
	if d < time.Second {
		return strconv.FormatFloat(float64(d.Microseconds())/1000, 'f', 2, 64) + "ms"
	}
	return strconv.FormatFloat(d.Seconds(), 'f', 2, 64) + "s"
}

func formatSize(size int) string {
	if size < 0 {
		return "-"
	}
	return strconv.Itoa(size)
}

func formatHTTPVersion(proto string) string {
	if strings.HasPrefix(proto, "HTTP/") {
		return proto[5:]
	}
	return proto
}

// --- Color Helpers ---

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
)

func colorizeStatus(status int) string {
	s := strconv.Itoa(status)
	switch {
	case status >= 500:
		return red + s + reset
	case status >= 400:
		return yellow + s + reset
	case status >= 300:
		return cyan + s + reset
	case status >= 200:
		return green + s + reset
	default:
		return s
	}
}

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
	default:
		return method
	}
}

// --- Field Extraction ---

type fieldExtractor struct {
	source, key string
}

func parseFieldSource(source string) fieldExtractor {
	parts := strings.SplitN(source, ":", 2)
	if len(parts) != 2 {
		return fieldExtractor{source: "literal", key: source}
	}
	return fieldExtractor{source: parts[0], key: parts[1]}
}

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
	case "literal":
		return f.key
	}
	return ""
}

// --- Body Capture ---

func captureRequestBody(r *http.Request, maxSize int64) []byte {
	if r.Body == nil {
		return nil
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxSize))
	if err != nil {
		return nil
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

// --- Request Helpers ---

func getRemoteAddr(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

// --- Token Extractors ---

// JSONBodyExtractor creates a TokenExtractor for JSON body fields.
func JSONBodyExtractor(path string) TokenExtractor {
	parts := strings.Split(path, ".")
	return func(r *http.Request, body []byte) string {
		if len(body) == 0 {
			return ""
		}
		var data map[string]any
		if json.Unmarshal(body, &data) != nil {
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
			if v == float64(int64(v)) {
				return strconv.FormatInt(int64(v), 10)
			}
			return strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(v)
		case nil:
			return ""
		default:
			if b, err := json.Marshal(v); err == nil {
				return string(b)
			}
			return fmt.Sprintf("%v", v)
		}
	}
}

// FormValueExtractor creates a TokenExtractor for form fields.
func FormValueExtractor(field string) TokenExtractor {
	return func(r *http.Request, body []byte) string {
		return r.FormValue(field)
	}
}

// ContextValueExtractor creates a TokenExtractor for context values.
func ContextValueExtractor(key any) TokenExtractor {
	return func(r *http.Request, body []byte) string {
		if v := r.Context().Value(key); v != nil {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}
}
