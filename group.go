package helix

import "net/http"

// Group represents a group of routes with a common prefix and middleware.
type Group struct {
	prefix     string
	middleware []Middleware
	server     *Server
	parent     *Group
}

// toMiddleware converts any middleware type to Middleware.
func toMiddleware(mw []any) []Middleware {
	result := make([]Middleware, 0, len(mw))
	for _, m := range mw {
		converted, err := convertToMiddleware(m)
		if err != nil {
			panic(err)
		}
		result = append(result, converted)
	}
	return result
}

// Group creates a new route group with the given prefix.
// The prefix is prepended to all routes registered on the group.
// Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.
func (s *Server) Group(prefix string, mw ...any) *Group {
	return &Group{
		prefix:     prefix,
		middleware: toMiddleware(mw),
		server:     s,
	}
}

// Group creates a nested group with the given prefix.
// The prefix is appended to the parent group's prefix.
// Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.
func (g *Group) Group(prefix string, mw ...any) *Group {
	return &Group{
		prefix:     g.fullPrefix() + prefix,
		middleware: toMiddleware(mw),
		server:     g.server,
		parent:     g,
	}
}

// Use adds middleware to the group.
// Middleware is applied to all routes registered on this group.
// Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.
func (g *Group) Use(mw ...any) {
	g.middleware = append(g.middleware, toMiddleware(mw)...)
}

// fullPrefix returns the complete prefix including parent prefixes.
func (g *Group) fullPrefix() string {
	return g.prefix
}

// allMiddleware returns all middleware for this group, including parent middleware.
func (g *Group) allMiddleware() []Middleware {
	var all []Middleware

	// Collect parent middleware first
	if g.parent != nil {
		all = append(all, g.parent.allMiddleware()...)
	}

	// Then add this group's middleware
	all = append(all, g.middleware...)

	return all
}

// wrapHandler wraps a handler with the group's middleware.
func (g *Group) wrapHandler(handler http.HandlerFunc) http.HandlerFunc {
	middleware := g.allMiddleware()
	if len(middleware) == 0 {
		return handler
	}

	// Build the handler chain
	var h http.Handler = handler
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

// Handle registers a handler for the given method and pattern.
func (g *Group) Handle(method, pattern string, handler http.HandlerFunc) {
	fullPattern := g.fullPrefix() + pattern
	// Prepend base path if set
	fullPattern = g.server.prependBasePath(fullPattern)
	g.server.router.Handle(method, fullPattern, g.wrapHandler(handler))
}

// GET registers a handler for GET requests.
func (g *Group) GET(pattern string, handler http.HandlerFunc) {
	g.Handle(http.MethodGet, pattern, handler)
}

// POST registers a handler for POST requests.
func (g *Group) POST(pattern string, handler http.HandlerFunc) {
	g.Handle(http.MethodPost, pattern, handler)
}

// PUT registers a handler for PUT requests.
func (g *Group) PUT(pattern string, handler http.HandlerFunc) {
	g.Handle(http.MethodPut, pattern, handler)
}

// PATCH registers a handler for PATCH requests.
func (g *Group) PATCH(pattern string, handler http.HandlerFunc) {
	g.Handle(http.MethodPatch, pattern, handler)
}

// DELETE registers a handler for DELETE requests.
func (g *Group) DELETE(pattern string, handler http.HandlerFunc) {
	g.Handle(http.MethodDelete, pattern, handler)
}

// OPTIONS registers a handler for OPTIONS requests.
func (g *Group) OPTIONS(pattern string, handler http.HandlerFunc) {
	g.Handle(http.MethodOptions, pattern, handler)
}

// HEAD registers a handler for HEAD requests.
func (g *Group) HEAD(pattern string, handler http.HandlerFunc) {
	g.Handle(http.MethodHead, pattern, handler)
}

// Any registers a handler for all HTTP methods.
func (g *Group) Any(pattern string, handler http.HandlerFunc) {
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
		g.Handle(method, pattern, handler)
	}
}

// Static serves static files from the given file system root.
func (g *Group) Static(pattern, root string) {
	if pattern == "" {
		panic("helix: pattern must not be empty")
	}
	if pattern[len(pattern)-1] != '/' {
		pattern += "/"
	}

	// Add catch-all pattern
	fullPattern := pattern + "{filepath...}"

	fs := http.Dir(root)
	fileServer := http.StripPrefix(g.fullPrefix()+pattern, http.FileServer(fs))

	g.GET(fullPattern, func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})
}

// Resource creates a new ResourceBuilder for the given pattern within this group.
// The pattern is relative to the group's prefix.
// Optional middleware can be applied to all routes in the resource.
// Accepts helix.Middleware, middleware.Middleware, or func(http.Handler) http.Handler.
func (g *Group) Resource(pattern string, mw ...any) *ResourceBuilder {
	// Combine group middleware with resource middleware
	converted := toMiddleware(mw)
	allMW := make([]Middleware, 0, len(g.allMiddleware())+len(converted))
	allMW = append(allMW, g.allMiddleware()...)
	allMW = append(allMW, converted...)

	return &ResourceBuilder{
		server:     g.server,
		group:      g,
		pattern:    pattern,
		middleware: allMW,
	}
}
