package middleware

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

// RecoverConfig configures the Recover middleware.
type RecoverConfig struct {
	// PrintStack enables printing the stack trace when a panic occurs.
	// Default: true
	PrintStack bool

	// StackSize is the maximum size of the stack trace buffer.
	// Default: 4KB
	StackSize int

	// Output is the writer to output the panic message to.
	// Default: os.Stderr
	Output io.Writer

	// Handler is a custom function to handle panics.
	// If set, it will be called instead of the default behavior.
	// The handler should write the response and return.
	Handler func(w http.ResponseWriter, r *http.Request, err any)
}

// DefaultRecoverConfig returns the default configuration for Recover.
func DefaultRecoverConfig() RecoverConfig {
	return RecoverConfig{
		PrintStack: true,
		StackSize:  4 << 10, // 4KB
		Output:     os.Stderr,
	}
}

// Recover returns a middleware that recovers from panics.
// It logs the panic and stack trace, then returns a 500 Internal Server Error.
func Recover() Middleware {
	return RecoverWithConfig(DefaultRecoverConfig())
}

// RecoverWithConfig returns a Recover middleware with the given configuration.
func RecoverWithConfig(config RecoverConfig) Middleware {
	if config.StackSize == 0 {
		config.StackSize = 4 << 10
	}
	if config.Output == nil {
		config.Output = os.Stderr
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Handle custom recovery handler
					if config.Handler != nil {
						config.Handler(w, r, err)
						return
					}

					// Get stack trace
					var stack []byte
					if config.PrintStack {
						stack = make([]byte, config.StackSize)
						length := runtime.Stack(stack, false)
						stack = stack[:length]
					}

					// Log the error
					if config.PrintStack {
						fmt.Fprintf(config.Output, "[PANIC RECOVER] %v\n%s\n", err, stack)
					} else {
						fmt.Fprintf(config.Output, "[PANIC RECOVER] %v\n", err)
					}

					// Return 500 Internal Server Error
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
