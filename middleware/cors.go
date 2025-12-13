package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig configures the CORS middleware.
type CORSConfig struct {
	// AllowOrigins is a list of origins that are allowed.
	// Use "*" to allow all origins.
	// Default: []
	AllowOrigins []string

	// AllowOriginFunc is a custom function to validate the origin.
	// If set, AllowOrigins is ignored.
	AllowOriginFunc func(origin string) bool

	// AllowMethods is a list of methods that are allowed.
	// Default: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
	AllowMethods []string

	// AllowHeaders is a list of headers that are allowed.
	// Default: Origin, Content-Type, Accept, Authorization
	AllowHeaders []string

	// ExposeHeaders is a list of headers that are exposed to the client.
	// Default: []
	ExposeHeaders []string

	// AllowCredentials indicates whether credentials are allowed.
	// Default: false
	AllowCredentials bool

	// MaxAge is the maximum age (in seconds) of the preflight cache.
	// Default: 0 (no caching)
	MaxAge int
}

// DefaultCORSConfig returns the default CORS configuration.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
		},
	}
}

// CORS returns a CORS middleware with default configuration.
func CORS() Middleware {
	return CORSWithConfig(DefaultCORSConfig())
}

// CORSWithConfig returns a CORS middleware with the given configuration.
func CORSWithConfig(config CORSConfig) Middleware {
	// Set defaults
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig().AllowMethods
	}
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = DefaultCORSConfig().AllowHeaders
	}

	// Precompute header values
	allowMethods := strings.Join(config.AllowMethods, ", ")
	allowHeaders := strings.Join(config.AllowHeaders, ", ")
	exposeHeaders := strings.Join(config.ExposeHeaders, ", ")
	maxAge := strconv.Itoa(config.MaxAge)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// If no origin header, not a CORS request
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if origin is allowed
			allowed := false
			if config.AllowOriginFunc != nil {
				allowed = config.AllowOriginFunc(origin)
			} else {
				for _, o := range config.AllowOrigins {
					if o == "*" || o == origin {
						allowed = true
						break
					}
				}
			}

			if !allowed {
				next.ServeHTTP(w, r)
				return
			}

			// Set CORS headers
			if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" && !config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if exposeHeaders != "" {
				w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
			}

			// Handle preflight request
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", allowMethods)
				w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", maxAge)
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSAllowAll returns a CORS middleware that allows all origins, methods, and headers.
func CORSAllowAll() Middleware {
	return CORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
	})
}
