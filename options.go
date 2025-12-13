package helix

import (
	"crypto/tls"
	"time"
)

// Option configures a Server.
type Option func(*Server)

// WithAddr sets the address the server listens on.
// Default is ":8080".
func WithAddr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithReadTimeout sets the maximum duration for reading the entire request.
func WithReadTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = d
	}
}

// WithWriteTimeout sets the maximum duration before timing out writes of the response.
func WithWriteTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = d
	}
}

// WithIdleTimeout sets the maximum amount of time to wait for the next request
// when keep-alives are enabled.
func WithIdleTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.idleTimeout = d
	}
}

// WithGracePeriod sets the maximum duration to wait for active connections
// to finish during graceful shutdown.
func WithGracePeriod(d time.Duration) Option {
	return func(s *Server) {
		s.gracePeriod = d
	}
}

// WithTLS configures the server to use TLS with the provided certificate and key files.
func WithTLS(certFile, keyFile string) Option {
	return func(s *Server) {
		s.tlsCertFile = certFile
		s.tlsKeyFile = keyFile
	}
}

// WithTLSConfig sets a custom TLS configuration for the server.
func WithTLSConfig(config *tls.Config) Option {
	return func(s *Server) {
		s.tlsConfig = config
	}
}

// WithMaxHeaderBytes sets the maximum size of request headers.
func WithMaxHeaderBytes(n int) Option {
	return func(s *Server) {
		s.maxHeaderBytes = n
	}
}

// HideBanner hides the banner and sets the banner to an empty string.
func HideBanner() Option {
	return func(s *Server) {
		s.hideBanner = true
		s.banner = ""
	}
}

// WithCustomBanner sets a custom banner for the server.
func WithCustomBanner(banner string) Option {
	return func(s *Server) {
		s.banner = banner
		s.hideBanner = false
	}
}

// WithErrorHandler sets a custom error handler for the server.
// The error handler will be called whenever an error occurs in a handler.
// If not set, the default error handling (RFC 7807 Problem Details) is used.
func WithErrorHandler(handler ErrorHandler) Option {
	return func(s *Server) {
		s.errorHandler = handler
	}
}

// WithBasePath sets a base path prefix for all routes.
// All registered routes will be prefixed with this path.
// For example, with base path "/api/v1", a route "/users" becomes "/api/v1/users".
// The base path should start with "/" but should not end with "/" (it will be normalized).
func WithBasePath(path string) Option {
	return func(s *Server) {
		s.basePath = path
	}
}
