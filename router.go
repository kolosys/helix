package helix

import (
	"net/http"
	"strings"
	"sync"
)

// RouteInfo contains information about a registered route.
type RouteInfo struct {
	Method  string
	Pattern string
}

// Router handles HTTP request routing.
type Router struct {
	trees      map[string]*routeNode // method -> root
	routes     []RouteInfo           // registered routes for introspection
	mu         sync.RWMutex
	paramsPool sync.Pool
}

// routeNode represents a node in the routing tree.
type routeNode struct {
	path     string           // static path segment
	children []*routeNode     // child nodes
	param    *routeNode       // parameter child node
	paramKey string           // parameter name if this is a param node
	catchAll *routeNode       // catch-all child node
	handler  http.HandlerFunc // handler for this route
}

// params holds path parameters extracted from a route.
type params struct {
	keys   []string
	values []string
}

func (p *params) reset() {
	p.keys = p.keys[:0]
	p.values = p.values[:0]
}

func (p *params) add(key, value string) {
	p.keys = append(p.keys, key)
	p.values = append(p.values, value)
}

func (p *params) get(key string) string {
	for i, k := range p.keys {
		if k == key {
			return p.values[i]
		}
	}
	return ""
}

// newRouter creates a new Router.
func newRouter() *Router {
	return &Router{
		trees: make(map[string]*routeNode),
		paramsPool: sync.Pool{
			New: func() any {
				return &params{
					keys:   make([]string, 0, 8),
					values: make([]string, 0, 8),
				}
			},
		},
	}
}

// Handle registers a new route with the given method and pattern.
func (r *Router) Handle(method, pattern string, handler http.HandlerFunc) {
	if pattern == "" {
		panic("helix: pattern must not be empty")
	}
	if pattern[0] != '/' {
		panic("helix: pattern must begin with '/'")
	}
	if handler == nil {
		panic("helix: handler must not be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	root := r.trees[method]
	if root == nil {
		root = &routeNode{}
		r.trees[method] = root
	}

	// Track the route for introspection
	r.routes = append(r.routes, RouteInfo{
		Method:  method,
		Pattern: pattern,
	})

	// Parse pattern into segments
	segments := parsePattern(pattern)
	r.addRoute(root, segments, handler)
}

// Routes returns all registered routes.
func (r *Router) Routes() []RouteInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routes := make([]RouteInfo, len(r.routes))
	copy(routes, r.routes)
	return routes
}

// segment represents a path segment.
type segment struct {
	value    string // segment value (static text or param name)
	isParam  bool   // is this a parameter?
	catchAll bool   // is this a catch-all?
}

// parsePattern parses a pattern into segments.
func parsePattern(pattern string) []segment {
	// Remove leading slash
	if len(pattern) > 0 && pattern[0] == '/' {
		pattern = pattern[1:]
	}

	if pattern == "" {
		return nil
	}

	parts := strings.Split(pattern, "/")
	segments := make([]segment, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		if len(part) > 2 && part[0] == '{' && part[len(part)-1] == '}' {
			paramName := part[1 : len(part)-1]
			if strings.HasSuffix(paramName, "...") {
				segments = append(segments, segment{
					value:    paramName[:len(paramName)-3],
					isParam:  true,
					catchAll: true,
				})
			} else {
				segments = append(segments, segment{
					value:   paramName,
					isParam: true,
				})
			}
		} else {
			segments = append(segments, segment{value: part})
		}
	}

	return segments
}

// addRoute adds a route to the tree.
func (r *Router) addRoute(n *routeNode, segments []segment, handler http.HandlerFunc) {
	if len(segments) == 0 {
		if n.handler != nil {
			panic("helix: route already registered")
		}
		n.handler = handler
		return
	}

	seg := segments[0]
	remaining := segments[1:]

	if seg.catchAll {
		if n.catchAll == nil {
			n.catchAll = &routeNode{paramKey: seg.value}
		}
		n.catchAll.handler = handler
		return
	}

	if seg.isParam {
		if n.param == nil {
			n.param = &routeNode{paramKey: seg.value}
		}
		r.addRoute(n.param, remaining, handler)
		return
	}

	for _, child := range n.children {
		if child.path == seg.value {
			r.addRoute(child, remaining, handler)
			return
		}
	}

	child := &routeNode{path: seg.value}
	n.children = append(n.children, child)
	r.addRoute(child, remaining, handler)
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	r.mu.RLock()
	root := r.trees[req.Method]
	r.mu.RUnlock()

	if root == nil {
		http.NotFound(w, req)
		return
	}

	ps := r.paramsPool.Get().(*params)
	ps.reset()

	handler := r.lookup(root, path, ps)

	if handler == nil {
		r.paramsPool.Put(ps)
		http.NotFound(w, req)
		return
	}

	if len(ps.keys) > 0 {
		ctx := setParams(req.Context(), ps)
		req = req.WithContext(ctx)
	}

	handler(w, req)

	r.paramsPool.Put(ps)
}

// lookup finds a handler for the given path.
func (r *Router) lookup(n *routeNode, path string, ps *params) http.HandlerFunc {
	// Remove leading slash
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	return r.lookupRecursive(n, path, ps)
}

// lookupRecursive recursively searches for a matching route.
func (r *Router) lookupRecursive(n *routeNode, path string, ps *params) http.HandlerFunc {
	if path == "" {
		return n.handler
	}

	before, after, ok := strings.Cut(path, "/")
	var segment, remaining string
	if !ok {
		segment = path
		remaining = ""
	} else {
		segment = before
		remaining = after
	}

	for _, child := range n.children {
		if child.path == segment {
			if handler := r.lookupRecursive(child, remaining, ps); handler != nil {
				return handler
			}
		}
	}

	if n.param != nil {
		ps.add(n.param.paramKey, segment)
		if handler := r.lookupRecursive(n.param, remaining, ps); handler != nil {
			return handler
		}
		ps.keys = ps.keys[:len(ps.keys)-1]
		ps.values = ps.values[:len(ps.values)-1]
	}

	if n.catchAll != nil {
		fullPath := segment
		if remaining != "" {
			fullPath = segment + "/" + remaining
		}
		ps.add(n.catchAll.paramKey, fullPath)
		return n.catchAll.handler
	}

	return nil
}
