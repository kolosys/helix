package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CacheConfig configures the Cache middleware.
type CacheConfig struct {
	// MaxAge sets the max-age directive in seconds.
	// Default: 0 (not set)
	MaxAge int

	// SMaxAge sets the s-maxage directive in seconds (for shared caches).
	// Default: 0 (not set)
	SMaxAge int

	// Public indicates the response can be cached by any cache.
	// Default: false
	Public bool

	// Private indicates the response is for a single user.
	// Default: false
	Private bool

	// NoCache indicates the response must be revalidated before use.
	// Default: false
	NoCache bool

	// NoStore indicates the response must not be stored.
	// Default: false
	NoStore bool

	// NoTransform indicates the response must not be transformed.
	// Default: false
	NoTransform bool

	// MustRevalidate indicates stale responses must be revalidated.
	// Default: false
	MustRevalidate bool

	// ProxyRevalidate is like MustRevalidate but for shared caches.
	// Default: false
	ProxyRevalidate bool

	// Immutable indicates the response body will not change.
	// Default: false
	Immutable bool

	// StaleWhileRevalidate allows serving stale content while revalidating.
	// Value is in seconds.
	// Default: 0 (not set)
	StaleWhileRevalidate int

	// StaleIfError allows serving stale content if there's an error.
	// Value is in seconds.
	// Default: 0 (not set)
	StaleIfError int

	// SkipFunc determines if cache headers should be skipped.
	SkipFunc func(r *http.Request) bool

	// VaryHeaders is a list of headers to include in the Vary header.
	VaryHeaders []string
}

// DefaultCacheConfig returns the default Cache configuration.
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{}
}

// Cache returns a Cache middleware with the given max-age in seconds.
func Cache(maxAge int) Middleware {
	return CacheWithConfig(CacheConfig{
		MaxAge: maxAge,
	})
}

// CachePublic returns a Cache middleware with public caching.
func CachePublic(maxAge int) Middleware {
	return CacheWithConfig(CacheConfig{
		MaxAge: maxAge,
		Public: true,
	})
}

// CachePrivate returns a Cache middleware with private caching.
func CachePrivate(maxAge int) Middleware {
	return CacheWithConfig(CacheConfig{
		MaxAge:  maxAge,
		Private: true,
	})
}

// CacheImmutable returns a Cache middleware for immutable content.
func CacheImmutable(maxAge int) Middleware {
	return CacheWithConfig(CacheConfig{
		MaxAge:    maxAge,
		Public:    true,
		Immutable: true,
	})
}

// NoCache returns a middleware that disables caching.
func NoCache() Middleware {
	return CacheWithConfig(CacheConfig{
		NoCache:        true,
		NoStore:        true,
		MustRevalidate: true,
	})
}

// CacheWithConfig returns a Cache middleware with the given configuration.
func CacheWithConfig(config CacheConfig) Middleware {
	// Pre-build the Cache-Control header value
	cacheControl := buildCacheControl(config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if configured
			if config.SkipFunc != nil && config.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Only apply to GET and HEAD requests
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			// Set Cache-Control header
			if cacheControl != "" {
				w.Header().Set("Cache-Control", cacheControl)
			}

			// Set Vary header
			if len(config.VaryHeaders) > 0 {
				w.Header().Set("Vary", strings.Join(config.VaryHeaders, ", "))
			}

			// Set Expires header if MaxAge is set
			if config.MaxAge > 0 && !config.NoCache && !config.NoStore {
				expires := time.Now().Add(time.Duration(config.MaxAge) * time.Second)
				w.Header().Set("Expires", expires.Format(http.TimeFormat))
			}

			next.ServeHTTP(w, r)
		})
	}
}

// buildCacheControl builds the Cache-Control header value from config.
func buildCacheControl(config CacheConfig) string {
	var directives []string

	if config.Public {
		directives = append(directives, "public")
	}
	if config.Private {
		directives = append(directives, "private")
	}
	if config.NoCache {
		directives = append(directives, "no-cache")
	}
	if config.NoStore {
		directives = append(directives, "no-store")
	}
	if config.NoTransform {
		directives = append(directives, "no-transform")
	}
	if config.MustRevalidate {
		directives = append(directives, "must-revalidate")
	}
	if config.ProxyRevalidate {
		directives = append(directives, "proxy-revalidate")
	}
	if config.Immutable {
		directives = append(directives, "immutable")
	}
	if config.MaxAge > 0 {
		directives = append(directives, fmt.Sprintf("max-age=%d", config.MaxAge))
	}
	if config.SMaxAge > 0 {
		directives = append(directives, fmt.Sprintf("s-maxage=%d", config.SMaxAge))
	}
	if config.StaleWhileRevalidate > 0 {
		directives = append(directives, fmt.Sprintf("stale-while-revalidate=%d", config.StaleWhileRevalidate))
	}
	if config.StaleIfError > 0 {
		directives = append(directives, fmt.Sprintf("stale-if-error=%d", config.StaleIfError))
	}

	return strings.Join(directives, ", ")
}

// SetCacheControl sets the Cache-Control header on the response.
func SetCacheControl(w http.ResponseWriter, value string) {
	w.Header().Set("Cache-Control", value)
}

// SetExpires sets the Expires header on the response.
func SetExpires(w http.ResponseWriter, t time.Time) {
	w.Header().Set("Expires", t.Format(http.TimeFormat))
}

// SetLastModified sets the Last-Modified header on the response.
func SetLastModified(w http.ResponseWriter, t time.Time) {
	w.Header().Set("Last-Modified", t.Format(http.TimeFormat))
}
