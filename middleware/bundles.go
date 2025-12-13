package middleware

// API returns a middleware bundle suitable for JSON API servers.
// Includes: RequestID, Logger (JSON format), Recover, and CORS.
func API() []Middleware {
	return []Middleware{
		RequestID(),
		Logger(LogFormatJSON),
		Recover(),
		CORS(),
	}
}

// APIWithCORS returns a middleware bundle suitable for JSON API servers
// with a custom CORS configuration.
// Includes: RequestID, Logger (JSON format), Recover, and CORS with config.
func APIWithCORS(cors CORSConfig) []Middleware {
	return []Middleware{
		RequestID(),
		Logger(LogFormatJSON),
		Recover(),
		CORSWithConfig(cors),
	}
}

// Web returns a middleware bundle suitable for web applications.
// Includes: RequestID, Logger (dev format), Recover, and Compress.
func Web() []Middleware {
	return []Middleware{
		RequestID(),
		Logger(LogFormatDev),
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
// Includes: RequestID, Logger (combined format), Recover.
func Production() []Middleware {
	return []Middleware{
		RequestID(),
		Logger(LogFormatCombined),
		Recover(),
	}
}

// Development returns a middleware bundle suitable for development.
// Includes: RequestID, Logger (dev format), Recover.
// This is the same as what helix.Default() uses.
func Development() []Middleware {
	return []Middleware{
		RequestID(),
		Logger(LogFormatDev),
		Recover(),
	}
}

// Secure returns a middleware bundle with security-focused middleware.
// Includes: RequestID, Logger (JSON format), Recover, RateLimit.
// Note: You should also add CORS and authentication middleware as needed.
// Parameters:
//   - rate: requests per second allowed
//   - burst: maximum burst size
func Secure(rate float64, burst int) []Middleware {
	return []Middleware{
		RequestID(),
		Logger(LogFormatJSON),
		Recover(),
		RateLimit(rate, burst),
	}
}
