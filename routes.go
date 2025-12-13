package helix

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

// prependBasePath prepends the base path to a route pattern if a base path is set.
func (s *Server) prependBasePath(pattern string) string {
	if s.basePath == "" {
		return pattern
	}

	// Normalize base path - ensure it starts with / and doesn't end with /
	basePath := s.basePath
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	basePath = strings.TrimSuffix(basePath, "/")

	// Normalize pattern - ensure it starts with /
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	// Combine base path and pattern
	return basePath + pattern
}

// Handle registers a handler for the given method and pattern.
func (s *Server) Handle(method, pattern string, handler http.HandlerFunc) {
	s.router.Handle(method, s.prependBasePath(pattern), handler)
}

// GET registers a handler for GET requests.
func (s *Server) GET(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodGet, pattern, handler)
}

// POST registers a handler for POST requests.
func (s *Server) POST(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodPost, pattern, handler)
}

// PUT registers a handler for PUT requests.
func (s *Server) PUT(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodPut, pattern, handler)
}

// PATCH registers a handler for PATCH requests.
func (s *Server) PATCH(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodPatch, pattern, handler)
}

// DELETE registers a handler for DELETE requests.
func (s *Server) DELETE(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodDelete, pattern, handler)
}

// OPTIONS registers a handler for OPTIONS requests.
func (s *Server) OPTIONS(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodOptions, pattern, handler)
}

// HEAD registers a handler for HEAD requests.
func (s *Server) HEAD(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodHead, pattern, handler)
}

// CONNECT registers a handler for CONNECT requests.
func (s *Server) CONNECT(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodConnect, pattern, handler)
}

// TRACE registers a handler for TRACE requests.
func (s *Server) TRACE(pattern string, handler http.HandlerFunc) {
	s.router.Handle(http.MethodTrace, pattern, handler)
}

// Any registers a handler for all HTTP methods.
func (s *Server) Any(pattern string, handler http.HandlerFunc) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	}
	for _, method := range methods {
		s.router.Handle(method, pattern, handler)
	}
}

// Static serves static files from the given file system root.
func (s *Server) Static(pattern, root string) {
	if pattern == "" {
		panic("helix: pattern must not be empty")
	}
	if pattern[len(pattern)-1] != '/' {
		pattern += "/"
	}

	// Add catch-all pattern
	fullPattern := pattern + "{filepath...}"

	fs := http.Dir(root)
	// Strip the original pattern (without base path) for file serving
	fileServer := http.StripPrefix(pattern, http.FileServer(fs))

	s.GET(fullPattern, func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})
}

// Routes returns all registered routes.
func (s *Server) Routes() []RouteInfo {
	return s.router.Routes()
}

// PrintRoutes prints all registered routes to the given writer.
// Routes are sorted by pattern, then by method.
func (s *Server) PrintRoutes(w io.Writer) {
	routes := s.Routes()

	// Sort routes by pattern, then by method
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Pattern != routes[j].Pattern {
			return routes[i].Pattern < routes[j].Pattern
		}
		return routes[i].Method < routes[j].Method
	})

	// Find max method length for alignment
	maxMethodLen := 0
	for _, r := range routes {
		if len(r.Method) > maxMethodLen {
			maxMethodLen = len(r.Method)
		}
	}

	for _, r := range routes {
		fmt.Fprintf(w, "%-*s  %s\n", maxMethodLen, r.Method, r.Pattern)
	}
}
