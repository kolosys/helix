package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// TimeoutConfig configures the Timeout middleware.
type TimeoutConfig struct {
	// Timeout is the maximum duration for the request.
	// Default: 30 seconds
	Timeout time.Duration

	// Handler is called when the request times out.
	// If nil, a default 503 Service Unavailable response is sent.
	Handler http.HandlerFunc

	// SkipFunc is a function that determines if timeout should be skipped.
	// If it returns true, no timeout is applied.
	SkipFunc func(r *http.Request) bool
}

// DefaultTimeoutConfig returns the default Timeout configuration.
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout: 30 * time.Second,
	}
}

// Timeout returns a middleware that adds a timeout to requests.
func Timeout(timeout time.Duration) Middleware {
	config := DefaultTimeoutConfig()
	config.Timeout = timeout
	return TimeoutWithConfig(config)
}

// TimeoutWithConfig returns a Timeout middleware with the given configuration.
func TimeoutWithConfig(config TimeoutConfig) Middleware {
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if configured
			if config.SkipFunc != nil && config.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Create context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), config.Timeout)
			defer cancel()

			// Replace request context
			r = r.WithContext(ctx)

			// Channel to signal completion
			done := make(chan struct{})

			// Create a response writer wrapper that can detect writes
			tw := &timeoutWriter{
				ResponseWriter: w,
				done:           done,
			}

			// Process request in goroutine
			go func() {
				next.ServeHTTP(tw, r)
				close(done)
			}()

			// Wait for completion or timeout
			select {
			case <-done:
				// Request completed successfully
				return
			case <-ctx.Done():
				// Timeout occurred
				tw.mu.Lock()
				defer tw.mu.Unlock()

				if tw.written {
					// Headers already written, can't send timeout response
					return
				}

				tw.timedOut = true

				if config.Handler != nil {
					config.Handler(w, r)
				} else {
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte("Service Unavailable: request timeout"))
				}
			}
		})
	}
}

// timeoutWriter wraps http.ResponseWriter to track if headers have been written.
type timeoutWriter struct {
	http.ResponseWriter
	done     chan struct{}
	mu       sync.Mutex
	written  bool
	timedOut bool
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut {
		return
	}
	tw.written = true
	tw.ResponseWriter.WriteHeader(code)
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut {
		return 0, context.DeadlineExceeded
	}
	tw.written = true
	return tw.ResponseWriter.Write(b)
}
