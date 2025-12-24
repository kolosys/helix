package helix

import (
	"crypto/tls"
	"time"
)

// Options configures a Server.
type Options struct {
	// Addr is the address the server listens on.
	// Default is ":8080".
	Addr string

	// ReadTimeout is the maximum duration for reading the entire request.
	// Default is 30 seconds.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out writes of the response.
	// Default is 30 seconds.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the next request
	// when keep-alives are enabled.
	// Default is 120 seconds.
	IdleTimeout time.Duration

	// GracePeriod is the maximum duration to wait for active connections
	// to finish during graceful shutdown.
	// Default is 30 seconds.
	GracePeriod time.Duration

	// TLSCertFile is the path to the TLS certificate file.
	// If set along with TLSKeyFile, the server will use TLS.
	TLSCertFile string

	// TLSKeyFile is the path to the TLS key file.
	// If set along with TLSCertFile, the server will use TLS.
	TLSKeyFile string

	// TLSConfig is a custom TLS configuration for the server.
	// If set, TLSCertFile and TLSKeyFile are ignored.
	TLSConfig *tls.Config

	// MaxHeaderBytes is the maximum size of request headers.
	// Default is 0 (no limit).
	MaxHeaderBytes int

	// HideBanner hides the banner on startup.
	// Default is false.
	HideBanner bool

	// Banner is a custom banner to display on startup.
	// If set, HideBanner is ignored.
	Banner string

	// ErrorHandler is a custom error handler for the server.
	// If not set, the default error handling (RFC 7807 Problem Details) is used.
	ErrorHandler ErrorHandler

	// BasePath is a base path prefix for all routes.
	// All registered routes will be prefixed with this path.
	// For example, with base path "/api/v1", a route "/users" becomes "/api/v1/users".
	// The base path should start with "/" but should not end with "/" (it will be normalized).
	BasePath string
}

// applyDefaults applies default values to nil or zero-valued options.
func (o *Options) applyDefaults() {
	if o.Addr == "" {
		o.Addr = ":8080"
	}
	if o.ReadTimeout == 0 {
		o.ReadTimeout = 30 * time.Second
	}
	if o.WriteTimeout == 0 {
		o.WriteTimeout = 30 * time.Second
	}
	if o.IdleTimeout == 0 {
		o.IdleTimeout = 120 * time.Second
	}
	if o.GracePeriod == 0 {
		o.GracePeriod = 30 * time.Second
	}
}
