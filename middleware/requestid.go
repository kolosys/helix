package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// RequestIDHeader is the default header name for the request ID.
const RequestIDHeader = "X-Request-ID"

// requestIDKey is the context key for the request ID.
type requestIDKey struct{}

// RequestIDConfig configures the RequestID middleware.
type RequestIDConfig struct {
	// Header is the name of the header to read/write the request ID.
	// Default: "X-Request-ID"
	Header string

	// Generator is a function that generates a new request ID.
	// Default: generates a random 16-byte hex string
	Generator func() string

	// TargetHeader is the header name to set on the response.
	// Default: same as Header
	TargetHeader string
}

// DefaultRequestIDConfig returns the default configuration for RequestID.
func DefaultRequestIDConfig() RequestIDConfig {
	return RequestIDConfig{
		Header:       RequestIDHeader,
		Generator:    generateRequestID,
		TargetHeader: RequestIDHeader,
	}
}

// RequestID returns a middleware that generates or propagates a request ID.
// The request ID is stored in the request context and the response header.
func RequestID() Middleware {
	return RequestIDWithConfig(DefaultRequestIDConfig())
}

// RequestIDWithConfig returns a RequestID middleware with the given configuration.
func RequestIDWithConfig(config RequestIDConfig) Middleware {
	if config.Header == "" {
		config.Header = RequestIDHeader
	}
	if config.Generator == nil {
		config.Generator = generateRequestID
	}
	if config.TargetHeader == "" {
		config.TargetHeader = config.Header
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for existing request ID
			id := r.Header.Get(config.Header)
			if id == "" {
				id = config.Generator()
			}

			// Set response header
			w.Header().Set(config.TargetHeader, id)

			// Store in context
			ctx := context.WithValue(r.Context(), requestIDKey{}, id)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetRequestID retrieves the request ID from the context.
// Returns an empty string if no request ID is set.
func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// GetRequestIDFromRequest retrieves the request ID from the request context.
func GetRequestIDFromRequest(r *http.Request) string {
	return GetRequestID(r.Context())
}

// generateRequestID generates a random 16-byte hex string (32 characters).
func generateRequestID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to a timestamp-based ID if random fails
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(b)
}
