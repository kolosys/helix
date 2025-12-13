package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
)

// ETagConfig configures the ETag middleware.
type ETagConfig struct {
	// Weak indicates whether to generate weak ETags.
	// Weak ETags are prefixed with W/ and indicate semantic equivalence.
	// Default: false (strong ETags)
	Weak bool

	// SkipFunc determines if ETag generation should be skipped.
	SkipFunc func(r *http.Request) bool
}

// DefaultETagConfig returns the default ETag configuration.
func DefaultETagConfig() ETagConfig {
	return ETagConfig{
		Weak: false,
	}
}

// ETag returns an ETag middleware with default configuration.
func ETag() Middleware {
	return ETagWithConfig(DefaultETagConfig())
}

// ETagWeak returns an ETag middleware that generates weak ETags.
func ETagWeak() Middleware {
	return ETagWithConfig(ETagConfig{Weak: true})
}

// ETagWithConfig returns an ETag middleware with the given configuration.
func ETagWithConfig(config ETagConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if configured
			if config.SkipFunc != nil && config.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Only process GET and HEAD requests
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			// Create a response buffer
			ew := &etagWriter{
				ResponseWriter: w,
				buffer:         bytes.NewBuffer(nil),
				weak:           config.Weak,
			}

			// Process request
			next.ServeHTTP(ew, r)

			// Generate ETag from response body
			if ew.buffer.Len() > 0 {
				hash := sha256.Sum256(ew.buffer.Bytes())
				etag := `"` + hex.EncodeToString(hash[:8]) + `"`
				if config.Weak {
					etag = "W/" + etag
				}

				// Check If-None-Match header
				ifNoneMatch := r.Header.Get("If-None-Match")
				if ifNoneMatch != "" && matchETag(ifNoneMatch, etag) {
					w.Header().Set("ETag", etag)
					w.WriteHeader(http.StatusNotModified)
					return
				}

				// Set ETag header
				w.Header().Set("ETag", etag)
			}

			// Write response
			if !ew.headerWritten {
				w.WriteHeader(ew.status)
			}
			if r.Method != http.MethodHead {
				w.Write(ew.buffer.Bytes())
			}
		})
	}
}

// etagWriter buffers the response to compute ETag.
type etagWriter struct {
	http.ResponseWriter
	buffer        *bytes.Buffer
	status        int
	headerWritten bool
	weak          bool
}

func (ew *etagWriter) WriteHeader(code int) {
	ew.status = code
	// Don't write header yet - we need to compute ETag first
}

func (ew *etagWriter) Write(b []byte) (int, error) {
	if ew.status == 0 {
		ew.status = http.StatusOK
	}
	return ew.buffer.Write(b)
}

// matchETag checks if an ETag matches the If-None-Match header.
func matchETag(ifNoneMatch, etag string) bool {
	// Handle wildcard
	if ifNoneMatch == "*" {
		return true
	}

	// Parse If-None-Match (can be comma-separated list)
	for _, match := range strings.Split(ifNoneMatch, ",") {
		match = strings.TrimSpace(match)

		// Remove weak indicator for comparison
		matchValue := strings.TrimPrefix(match, "W/")
		etagValue := strings.TrimPrefix(etag, "W/")

		// Remove quotes
		matchValue = strings.Trim(matchValue, `"`)
		etagValue = strings.Trim(etagValue, `"`)

		if matchValue == etagValue {
			return true
		}
	}

	return false
}

// ETagFromContent generates an ETag from content.
func ETagFromContent(content []byte, weak bool) string {
	hash := sha256.Sum256(content)
	etag := `"` + hex.EncodeToString(hash[:8]) + `"`
	if weak {
		etag = "W/" + etag
	}
	return etag
}

// ETagFromString generates an ETag from a string.
func ETagFromString(s string, weak bool) string {
	return ETagFromContent([]byte(s), weak)
}

// ETagFromVersion generates an ETag from a version number.
func ETagFromVersion(version int64, weak bool) string {
	etag := `"` + strconv.FormatInt(version, 10) + `"`
	if weak {
		etag = "W/" + etag
	}
	return etag
}

