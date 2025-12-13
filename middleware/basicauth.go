package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// BasicAuthConfig configures the BasicAuth middleware.
type BasicAuthConfig struct {
	// Validator is a function that validates the username and password.
	// Return true if the credentials are valid.
	Validator func(username, password string) bool

	// Realm is the authentication realm displayed in the browser.
	// Default: "Restricted"
	Realm string

	// SkipFunc determines if authentication should be skipped.
	SkipFunc func(r *http.Request) bool
}

// BasicAuth returns a BasicAuth middleware with the given username and password.
// Uses constant-time comparison to prevent timing attacks.
func BasicAuth(username, password string) Middleware {
	return BasicAuthWithConfig(BasicAuthConfig{
		Validator: func(u, p string) bool {
			return secureCompare(u, username) && secureCompare(p, password)
		},
		Realm: "Restricted",
	})
}

// BasicAuthWithValidator returns a BasicAuth middleware with a custom validator.
func BasicAuthWithValidator(validator func(username, password string) bool) Middleware {
	return BasicAuthWithConfig(BasicAuthConfig{
		Validator: validator,
		Realm:     "Restricted",
	})
}

// BasicAuthWithConfig returns a BasicAuth middleware with the given configuration.
func BasicAuthWithConfig(config BasicAuthConfig) Middleware {
	if config.Validator == nil {
		panic("helix: BasicAuth validator is required")
	}
	if config.Realm == "" {
		config.Realm = "Restricted"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if configured
			if config.SkipFunc != nil && config.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Get credentials from request
			username, password, ok := r.BasicAuth()
			if !ok {
				unauthorized(w, config.Realm)
				return
			}

			// Validate credentials
			if !config.Validator(username, password) {
				unauthorized(w, config.Realm)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// BasicAuthUsers returns a BasicAuth middleware that validates against a map of users.
// The map key is the username and the value is the password.
func BasicAuthUsers(users map[string]string) Middleware {
	return BasicAuthWithConfig(BasicAuthConfig{
		Validator: func(username, password string) bool {
			expectedPassword, ok := users[username]
			if !ok {
				return false
			}
			return secureCompare(password, expectedPassword)
		},
		Realm: "Restricted",
	})
}

// unauthorized sends a 401 Unauthorized response with WWW-Authenticate header.
func unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized"))
}

// secureCompare performs a constant-time comparison of two strings.
// This prevents timing attacks.
func secureCompare(a, b string) bool {
	// Hash both strings to ensure constant time comparison regardless of length
	aHash := sha256.Sum256([]byte(a))
	bHash := sha256.Sum256([]byte(b))
	return subtle.ConstantTimeCompare(aHash[:], bHash[:]) == 1
}
