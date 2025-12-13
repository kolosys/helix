package helix

import (
	"context"
	"net/http"
	"reflect"
	"sync"
)

// Services provides a type-safe service registry for dependency injection.
// Register services once at startup, then access them in handlers.
//
// Example:
//
//	// At startup
//	helix.Register[*UserService](userService)
//	helix.Register[*EmailService](emailService)
//
//	// In handler
//	s.GET("/users", helix.HandleCtx(func(c *helix.Ctx) error {
//	    userSvc := helix.Get[*UserService]()
//	    users, _ := userSvc.List(c.Context())
//	    return c.OK(users)
//	}))

// serviceRegistry is the global service registry.
var (
	servicesMu sync.RWMutex
	services   = make(map[reflect.Type]any)
)

// Register registers a service in the global registry by its type.
func Register[T any](service T) {
	servicesMu.Lock()
	defer servicesMu.Unlock()

	t := reflect.TypeOf((*T)(nil)).Elem()
	services[t] = service
}

// Get retrieves a service from the global registry.
// Returns the zero value and false if not found.
func Get[T any]() (T, bool) {
	servicesMu.RLock()
	defer servicesMu.RUnlock()

	t := reflect.TypeOf((*T)(nil)).Elem()
	if v, ok := services[t]; ok {
		return v.(T), true
	}
	var zero T
	return zero, false
}

// MustGet retrieves a service from the global registry or panics.
func MustGet[T any]() T {
	svc, ok := Get[T]()
	if !ok {
		t := reflect.TypeOf((*T)(nil)).Elem()
		panic("helix: service not registered: " + t.String())
	}
	return svc
}

// contextServices holds services for a specific request context.
type contextServices struct {
	mu       sync.RWMutex
	services map[reflect.Type]any
}

// WithService returns a new context with the service added.
// This is useful for request-scoped services like database transactions.
func WithService[T any](ctx context.Context, service T) context.Context {
	cs := getContextServices(ctx)
	if cs == nil {
		cs = &contextServices{services: make(map[reflect.Type]any)}
		ctx = context.WithValue(ctx, servicesCtxKey, cs)
	}

	cs.mu.Lock()
	t := reflect.TypeOf((*T)(nil)).Elem()
	cs.services[t] = service
	cs.mu.Unlock()

	return ctx
}

// FromContext retrieves a service from the context or falls back to global registry.
func FromContext[T any](ctx context.Context) (T, bool) {
	// First check context-specific services
	cs := getContextServices(ctx)
	if cs != nil {
		cs.mu.RLock()
		t := reflect.TypeOf((*T)(nil)).Elem()
		if v, ok := cs.services[t]; ok {
			cs.mu.RUnlock()
			return v.(T), true
		}
		cs.mu.RUnlock()
	}

	// Fall back to global registry
	return Get[T]()
}

// MustFromContext retrieves a service from context or panics.
func MustFromContext[T any](ctx context.Context) T {
	svc, ok := FromContext[T](ctx)
	if !ok {
		t := reflect.TypeOf((*T)(nil)).Elem()
		panic("helix: service not found in context or registry: " + t.String())
	}
	return svc
}

// getContextServices retrieves the services from context.
func getContextServices(ctx context.Context) *contextServices {
	if v := ctx.Value(servicesCtxKey); v != nil {
		return v.(*contextServices)
	}
	return nil
}

// ProvideMiddleware returns middleware that injects services into the request context.
// Services added this way are request-scoped.
func ProvideMiddleware[T any](factory func(r *http.Request) T) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			service := factory(r)
			ctx := WithService(r.Context(), service)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
