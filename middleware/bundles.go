package middleware

// API returns a middleware bundle suitable for JSON API servers.
// Includes: RequestID, Recover, and CORS.
// Add logging via helix.LoggerMiddleware with your preferred RequestLogger.
func API() []Middleware {
	return []Middleware{
		RequestID(),
		Recover(),
		CORS(),
	}
}

// APIWithCORS returns a middleware bundle suitable for JSON API servers
// with a custom CORS configuration.
// Includes: RequestID, Recover, and CORS with config.
// Add logging via helix.LoggerMiddleware with your preferred RequestLogger.
func APIWithCORS(cors CORSConfig) []Middleware {
	return []Middleware{
		RequestID(),
		Recover(),
		CORSWithConfig(cors),
	}
}

// Web returns a middleware bundle suitable for web applications.
// Includes: RequestID, Recover, and Compress.
// Add logging via helix.LoggerMiddleware with your preferred RequestLogger.
func Web() []Middleware {
	return []Middleware{
		RequestID(),
		Recover(),
		Compress(),
	}
}

// Minimal returns a minimal middleware bundle with only essential middleware.
// Includes: Recover.
func Minimal() []Middleware {
	return []Middleware{
		Recover(),
	}
}

// Production returns a middleware bundle suitable for production environments.
// Includes: RequestID, Recover.
// Add logging via helix.LoggerMiddleware with your preferred RequestLogger.
func Production() []Middleware {
	return []Middleware{
		RequestID(),
		Recover(),
	}
}

// Development returns a middleware bundle suitable for development.
// Includes: RequestID, Recover.
// Add logging via helix.LoggerMiddleware with your preferred RequestLogger.
// This is the same as what helix.Default() uses (plus logging).
func Development() []Middleware {
	return []Middleware{
		RequestID(),
		Recover(),
	}
}

// Secure returns a middleware bundle with security-focused middleware.
// Includes: RequestID, Recover, RateLimit.
// Add logging via helix.LoggerMiddleware with your preferred RequestLogger.
// Note: You should also add CORS and authentication middleware as needed.
// Parameters:
//   - rate: requests per second allowed
//   - burst: maximum burst size
func Secure(rate float64, burst int) []Middleware {
	return []Middleware{
		RequestID(),
		Recover(),
		RateLimit(rate, burst),
	}
}
