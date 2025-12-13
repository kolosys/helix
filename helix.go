// Package helix provides a zero-dependency, context-aware, high-performance
// HTTP web framework for Go with stdlib compatibility.
package helix

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kolosys/helix/middleware"
)

const (
	// Version of Hexix
	Version = "0.1.0"
	website = "https://github.com/kolosys/helix"
	// http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Echo
	banner = `
    __ __    ___     
   / // /__ / (_)_ __
  / _  / -_) / /\ \ /
 /_//_/\__/_/_//_\_\ v%s
 Developer friendly HTTP framework
 %s
______________________________________
`
)

// Middleware is a function that wraps an http.Handler to provide additional functionality.
// This is an alias to middleware.Middleware for convenience.
type Middleware = middleware.Middleware

// Server is the main HTTP server for the Helix framework.
type Server struct {
	router     *Router
	middleware []Middleware
	httpServer *http.Server

	// Configuration
	addr           string
	readTimeout    time.Duration
	writeTimeout   time.Duration
	idleTimeout    time.Duration
	gracePeriod    time.Duration
	maxHeaderBytes int
	tlsCertFile    string
	tlsKeyFile     string
	tlsConfig      *tls.Config
	hideBanner     bool
	banner         string

	// Lifecycle hooks
	onStart []func(s *Server)
	onStop  []func(ctx context.Context, s *Server)

	// Error handling
	errorHandler ErrorHandler

	// Routing
	basePath string // Base path prefix for all routes

	// State
	once    sync.Once
	handler http.Handler // Pre-compiled middleware chain
	built   bool         // Whether the handler chain has been built

	// Object pools for zero-allocation hot path
	ctxPool sync.Pool
}

// New creates a new Server with the provided options.
func New(opts ...Option) *Server {
	s := &Server{
		router:       newRouter(),
		addr:         ":8080",
		readTimeout:  30 * time.Second,
		writeTimeout: 30 * time.Second,
		idleTimeout:  120 * time.Second,
		gracePeriod:  30 * time.Second,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.banner == "" && !s.hideBanner {
		s.banner = fmt.Sprintf(banner, Version, website)
	}

	return s
}

// Default creates a new Server with sensible defaults for development.
// It includes RequestID, Logger (dev format), and Recover middleware.
func Default(opts ...Option) *Server {
	s := New(opts...)
	s.Use(middleware.RequestID())
	s.Use(middleware.Logger(middleware.LogFormatDev))
	s.Use(middleware.Recover())
	return s
}

// Start starts the server and blocks until shutdown.
// If an address is provided, it will be used instead of the WithAddr option.
// If the address is not provided and the WithAddr option is not set, it will use ":8080".
// This is a convenience method that calls Run with a background context.
func (s *Server) Start(addr ...string) error {
	if len(addr) > 0 {
		s.addr = addr[0]
	} else if s.addr == "" {
		s.addr = ":8080"
	}
	log.Printf("Server starting on %s", s.addr)
	return s.Run(context.Background())
}

// convertToMiddleware converts any middleware type to Middleware.
// Returns an error with detailed type information if conversion fails.
func convertToMiddleware(m any) (Middleware, error) {
	switch v := m.(type) {
	case Middleware:
		return v, nil
	case func(http.Handler) http.Handler:
		return v, nil
	default:
		return nil, fmt.Errorf("helix: middleware must be Middleware or func(http.Handler) http.Handler, got %T", m)
	}
}

// Use adds middleware to the server's middleware chain.
// Middleware is executed in the order it is added.
// Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.
func (s *Server) Use(mw ...any) {
	for _, m := range mw {
		converted, err := convertToMiddleware(m)
		if err != nil {
			panic(err)
		}
		s.middleware = append(s.middleware, converted)
	}
}

// Build pre-compiles the middleware chain for optimal performance.
// This is called automatically before the server starts, but can be called
// manually after all routes and middleware are registered.
func (s *Server) Build() {
	if s.built {
		return
	}

	// Build the handler chain with middleware
	var handler http.Handler = s.router

	// If a base path is set, validate that incoming requests start with it
	// Routes are registered with the base path, so the router will match the full path
	if s.basePath != "" {
		handler = s.basePathMiddleware(handler)
	}

	// If a custom error handler is set, inject it into the request context
	// This must be done before other middleware so handlers can access it
	if s.errorHandler != nil {
		handler = s.errorHandlerMiddleware(handler)
	}

	// Apply middleware in reverse order so first added is outermost
	for i := len(s.middleware) - 1; i >= 0; i-- {
		handler = s.middleware[i](handler)
	}

	s.handler = handler
	s.built = true

	// Initialize context pool
	s.ctxPool = sync.Pool{
		New: func() any {
			return &Ctx{}
		},
	}
}

// basePathMiddleware validates that incoming requests start with the base path.
// Routes are registered with the base path prepended, so the router will match the full path.
// This middleware only validates and rejects requests that don't start with the base path.
func (s *Server) basePathMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		basePath := s.basePath

		// Normalize base path - ensure it starts with / and doesn't end with /
		if !strings.HasPrefix(basePath, "/") {
			basePath = "/" + basePath
		}
		basePath = strings.TrimSuffix(basePath, "/")

		// Check if path starts with base path
		if strings.HasPrefix(path, basePath) {
			// Path is valid, let the router handle it (routes are registered with base path)
			next.ServeHTTP(w, r)
		} else if path == "/" && basePath != "/" {
			// Root path doesn't match base path, return 404
			http.NotFound(w, r)
		} else {
			// Path doesn't match base path, return 404
			http.NotFound(w, r)
		}
	})
}

// errorHandlerMiddleware injects the error handler into the request context.
func (s *Server) errorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = withErrorHandler(r, s.errorHandler)
		next.ServeHTTP(w, r)
	})
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure handler chain is built (lazy initialization)
	if !s.built {
		s.Build()
	}

	s.handler.ServeHTTP(w, r)
}

// Run starts the server and blocks until the context is canceled or a shutdown
// signal is received. It performs graceful shutdown, waiting for active connections
// to finish within the grace period.
func (s *Server) Run(ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:           s.addr,
		Handler:        s,
		ReadTimeout:    s.readTimeout,
		WriteTimeout:   s.writeTimeout,
		IdleTimeout:    s.idleTimeout,
		MaxHeaderBytes: s.maxHeaderBytes,
		TLSConfig:      s.tlsConfig,
	}

	if !s.hideBanner {
		fmt.Println(strings.ReplaceAll(s.banner, "{version}", Version))
	}

	// Call onStart hooks
	for _, fn := range s.onStart {
		fn(s)
	}

	// Channel to receive server errors
	errCh := make(chan error, 1)

	// Start server in goroutine
	go func() {
		var err error
		if s.tlsCertFile != "" && s.tlsKeyFile != "" {
			err = s.httpServer.ListenAndServeTLS(s.tlsCertFile, s.tlsKeyFile)
		} else if s.tlsConfig != nil {
			err = s.httpServer.ListenAndServeTLS("", "")
		} else {
			err = s.httpServer.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	// Wait for shutdown signal or context cancellation
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
		// Received shutdown signal
	case <-ctx.Done():
		// Context canceled
	}

	// Perform graceful shutdown
	return s.Shutdown(context.Background())
}

// Shutdown gracefully shuts down the server without interrupting active connections.
// It waits for the grace period for active connections to finish.
func (s *Server) Shutdown(ctx context.Context) error {
	var err error
	s.once.Do(func() {
		// Create shutdown context with grace period
		shutdownCtx, cancel := context.WithTimeout(ctx, s.gracePeriod)
		defer cancel()

		// Call onStop hooks
		for _, fn := range s.onStop {
			fn(shutdownCtx, s)
		}

		if s.httpServer == nil {
			return
		}

		err = s.httpServer.Shutdown(shutdownCtx)
	})
	return err
}

// Addr returns the address the server is configured to listen on.
func (s *Server) Addr() string {
	return s.addr
}

// OnStart registers a function to be called when the server starts.
// Multiple functions can be registered and will be called in order.
func (s *Server) OnStart(fn func(s *Server)) {
	s.onStart = append(s.onStart, fn)
}

// OnStop registers a function to be called when the server stops.
// Multiple functions can be registered and will be called in order.
// The context passed to the function has the grace period as its deadline.
func (s *Server) OnStop(fn func(ctx context.Context, s *Server)) {
	s.onStop = append(s.onStop, fn)
}
