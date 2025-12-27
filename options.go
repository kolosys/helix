package helix

import (
	"crypto/tls"
	"net"
	"strconv"
	"strings"
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

	// AutoPort enables automatic port selection when the configured port is in use.
	// When enabled, the server will try incrementing ports until it finds an available one.
	// This is primarily useful for development environments.
	// Default is false.
	AutoPort bool

	// MaxPortAttempts is the maximum number of ports to try when AutoPort is enabled.
	// Default is 10.
	MaxPortAttempts int
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
	if o.MaxPortAttempts == 0 {
		o.MaxPortAttempts = 10
	}
}

// parseAddr parses an address into host and port components.
func parseAddr(addr string) (host string, port int, err error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		// Handle addresses like ":8080" where host is empty
		if strings.HasPrefix(addr, ":") {
			host = ""
			portStr = addr[1:]
		} else {
			return "", 0, err
		}
	}
	port, err = strconv.Atoi(portStr)
	if err != nil {
		return "", 0, err
	}
	return host, port, nil
}

// isPortAvailable checks if a port is available for listening.
func isPortAvailable(host string, port int) bool {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// findAvailableAddr finds an available address starting from the given address.
// If the port is in use, it increments the port and tries again up to maxAttempts times.
func findAvailableAddr(addr string, maxAttempts int) (string, error) {
	host, port, err := parseAddr(addr)
	if err != nil {
		return "", err
	}

	for i := 0; i < maxAttempts; i++ {
		if isPortAvailable(host, port+i) {
			return net.JoinHostPort(host, strconv.Itoa(port+i)), nil
		}
	}

	return "", &net.AddrError{Err: "no available port found", Addr: addr}
}
