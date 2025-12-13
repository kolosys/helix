package helix

import "net/http"

// Module is an interface for modular route definitions.
// Modules allow you to organize routes into separate files or packages
// and mount them onto a server or group.
//
// Example:
//
//	type UserModule struct {
//	    store *UserStore
//	}
//
//	func (m *UserModule) Register(r RouteRegistrar) {
//	    r.GET("/", m.list)
//	    r.POST("/", m.create)
//	    r.GET("/{id}", m.get)
//	}
//
//	// Mount the module
//	s.Mount("/users", &UserModule{store: store})
type Module interface {
	Register(r RouteRegistrar)
}

// ModuleFunc is a function that implements Module.
type ModuleFunc func(r RouteRegistrar)

// Register implements Module.
func (f ModuleFunc) Register(r RouteRegistrar) {
	f(r)
}

// RouteRegistrar is an interface for registering routes.
// Both Server and Group implement this interface.
type RouteRegistrar interface {
	GET(pattern string, handler http.HandlerFunc)
	POST(pattern string, handler http.HandlerFunc)
	PUT(pattern string, handler http.HandlerFunc)
	PATCH(pattern string, handler http.HandlerFunc)
	DELETE(pattern string, handler http.HandlerFunc)
	OPTIONS(pattern string, handler http.HandlerFunc)
	HEAD(pattern string, handler http.HandlerFunc)
	Handle(method, pattern string, handler http.HandlerFunc)
	Group(prefix string, mw ...any) *Group
	Resource(pattern string, mw ...any) *ResourceBuilder
}

// Mount mounts a module at the given prefix.
// The module's routes are prefixed with the given path.
func (s *Server) Mount(prefix string, m Module, mw ...any) {
	g := s.Group(prefix, mw...)
	m.Register(g)
}

// Mount mounts a module at the given prefix within a group.
func (g *Group) Mount(prefix string, m Module, mw ...any) {
	sub := g.Group(prefix, mw...)
	m.Register(sub)
}

// MountFunc mounts a function as a module at the given prefix.
func (s *Server) MountFunc(prefix string, fn func(r RouteRegistrar), mw ...any) {
	s.Mount(prefix, ModuleFunc(fn), mw...)
}

// MountFunc mounts a function as a module at the given prefix within a group.
func (g *Group) MountFunc(prefix string, fn func(r RouteRegistrar), mw ...any) {
	g.Mount(prefix, ModuleFunc(fn), mw...)
}
