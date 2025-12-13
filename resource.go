package helix

import "net/http"

// ResourceBuilder provides a fluent interface for defining REST resource routes.
type ResourceBuilder struct {
	server     *Server
	group      *Group
	pattern    string
	middleware []Middleware
}

// Resource creates a new ResourceBuilder for the given pattern.
// The pattern should be the base path for the resource (e.g., "/users").
// Optional middleware can be applied to all routes in the resource.
// Accepts Middleware (helix.Middleware is an alias for middleware.Middleware) or func(http.Handler) http.Handler.
func (s *Server) Resource(pattern string, mw ...any) *ResourceBuilder {
	return &ResourceBuilder{
		server:     s,
		pattern:    pattern,
		middleware: toMiddleware(mw),
	}
}

// wrapHandler wraps a handler with the resource's middleware.
func (rb *ResourceBuilder) wrapHandler(handler http.HandlerFunc) http.HandlerFunc {
	if len(rb.middleware) == 0 {
		return handler
	}

	var h http.Handler = handler
	for i := len(rb.middleware) - 1; i >= 0; i-- {
		h = rb.middleware[i](h)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

// handle registers a route using either the server or group.
func (rb *ResourceBuilder) handle(method, pattern string, handler http.HandlerFunc) {
	wrapped := rb.wrapHandler(handler)
	if rb.group != nil {
		rb.group.Handle(method, pattern, wrapped)
	} else {
		rb.server.Handle(method, pattern, wrapped)
	}
}

// List registers a GET handler for the collection (e.g., GET /users).
func (rb *ResourceBuilder) List(handler http.HandlerFunc) *ResourceBuilder {
	rb.handle(http.MethodGet, rb.pattern, handler)
	return rb
}

// Create registers a POST handler for creating resources (e.g., POST /users).
func (rb *ResourceBuilder) Create(handler http.HandlerFunc) *ResourceBuilder {
	rb.handle(http.MethodPost, rb.pattern, handler)
	return rb
}

// Get registers a GET handler for a single resource (e.g., GET /users/{id}).
func (rb *ResourceBuilder) Get(handler http.HandlerFunc) *ResourceBuilder {
	rb.handle(http.MethodGet, rb.pattern+"/{id}", handler)
	return rb
}

// Update registers a PUT handler for updating a resource (e.g., PUT /users/{id}).
func (rb *ResourceBuilder) Update(handler http.HandlerFunc) *ResourceBuilder {
	rb.handle(http.MethodPut, rb.pattern+"/{id}", handler)
	return rb
}

// Patch registers a PATCH handler for partial updates (e.g., PATCH /users/{id}).
func (rb *ResourceBuilder) Patch(handler http.HandlerFunc) *ResourceBuilder {
	rb.handle(http.MethodPatch, rb.pattern+"/{id}", handler)
	return rb
}

// Delete registers a DELETE handler for deleting a resource (e.g., DELETE /users/{id}).
func (rb *ResourceBuilder) Delete(handler http.HandlerFunc) *ResourceBuilder {
	rb.handle(http.MethodDelete, rb.pattern+"/{id}", handler)
	return rb
}

// Custom registers a handler with a custom method and path suffix.
// The suffix is appended to the base pattern.
// Example: Custom("POST", "/{id}/archive", archiveHandler) for POST /users/{id}/archive
func (rb *ResourceBuilder) Custom(method, suffix string, handler http.HandlerFunc) *ResourceBuilder {
	rb.handle(method, rb.pattern+suffix, handler)
	return rb
}

// Index is an alias for List.
func (rb *ResourceBuilder) Index(handler http.HandlerFunc) *ResourceBuilder {
	return rb.List(handler)
}

// Store is an alias for Create.
func (rb *ResourceBuilder) Store(handler http.HandlerFunc) *ResourceBuilder {
	return rb.Create(handler)
}

// Show is an alias for Get.
func (rb *ResourceBuilder) Show(handler http.HandlerFunc) *ResourceBuilder {
	return rb.Get(handler)
}

// Destroy is an alias for Delete.
func (rb *ResourceBuilder) Destroy(handler http.HandlerFunc) *ResourceBuilder {
	return rb.Delete(handler)
}

// CRUD registers all standard CRUD handlers in one call.
// Handlers: list, create, get, update, delete
func (rb *ResourceBuilder) CRUD(list, create, get, update, delete http.HandlerFunc) *ResourceBuilder {
	if list != nil {
		rb.List(list)
	}
	if create != nil {
		rb.Create(create)
	}
	if get != nil {
		rb.Get(get)
	}
	if update != nil {
		rb.Update(update)
	}
	if delete != nil {
		rb.Delete(delete)
	}
	return rb
}

// ReadOnly registers only read handlers (list and get).
func (rb *ResourceBuilder) ReadOnly(list, get http.HandlerFunc) *ResourceBuilder {
	if list != nil {
		rb.List(list)
	}
	if get != nil {
		rb.Get(get)
	}
	return rb
}

// -----------------------------------------------------------------------------
// Typed Resource Builder
// -----------------------------------------------------------------------------

// TypedResource creates a typed resource builder for the given entity type.
// This provides a fluent API for registering typed handlers for CRUD operations.
//
// Example:
//
//	helix.TypedResource[User](s, "/users").
//	    List(listUsers).
//	    Create(createUser).
//	    Get(getUser).
//	    Update(updateUser).
//	    Delete(deleteUser)
func TypedResource[Entity any](s *Server, pattern string, mw ...any) *TypedResourceBuilder[Entity] {
	return &TypedResourceBuilder[Entity]{
		server:     s,
		pattern:    pattern,
		middleware: toMiddleware(mw),
	}
}

// TypedResourceBuilder provides a fluent interface for defining typed REST resource routes.
type TypedResourceBuilder[Entity any] struct {
	server     *Server
	group      *Group
	pattern    string
	middleware []Middleware
}

// wrapHandler wraps a handler with the resource's middleware.
func (rb *TypedResourceBuilder[Entity]) wrapHandler(handler http.HandlerFunc) http.HandlerFunc {
	if len(rb.middleware) == 0 {
		return handler
	}

	var h http.Handler = handler
	for i := len(rb.middleware) - 1; i >= 0; i-- {
		h = rb.middleware[i](h)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

// handle registers a route using either the server or group.
func (rb *TypedResourceBuilder[Entity]) handle(method, pattern string, handler http.HandlerFunc) {
	wrapped := rb.wrapHandler(handler)
	if rb.group != nil {
		rb.group.Handle(method, pattern, wrapped)
	} else {
		rb.server.Handle(method, pattern, wrapped)
	}
}

// ListRequest is a common request type for list operations.
type ListRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
	Order  string `query:"order"`
	Search string `query:"search"`
}

// ListResponse wraps a list of entities with pagination metadata.
type ListResponse[Entity any] struct {
	Items []Entity `json:"items"`
	Total int      `json:"total"`
	Page  int      `json:"page,omitempty"`
	Limit int      `json:"limit,omitempty"`
}

// IDRequest is a common request type for single-entity operations.
type IDRequest struct {
	ID int `path:"id"`
}

// List registers a typed GET handler for the collection.
// Handler signature: func(ctx, ListReq) (ListResponse[Entity], error)
func (rb *TypedResourceBuilder[Entity]) List(h Handler[ListRequest, ListResponse[Entity]]) *TypedResourceBuilder[Entity] {
	rb.handle(http.MethodGet, rb.pattern, Handle(h))
	return rb
}

// Create registers a typed POST handler for creating resources.
// Handler signature: func(ctx, CreateReq) (Entity, error)
func (rb *TypedResourceBuilder[Entity]) Create(h Handler[Entity, Entity]) *TypedResourceBuilder[Entity] {
	rb.handle(http.MethodPost, rb.pattern, HandleCreated(h))
	return rb
}

// Get registers a typed GET handler for a single resource.
// Handler signature: func(ctx, IDRequest) (Entity, error)
func (rb *TypedResourceBuilder[Entity]) Get(h Handler[IDRequest, Entity]) *TypedResourceBuilder[Entity] {
	rb.handle(http.MethodGet, rb.pattern+"/{id}", Handle(h))
	return rb
}

// Update registers a typed PUT handler for updating a resource.
// The request type should include the ID from path and the update data.
func (rb *TypedResourceBuilder[Entity]) Update(h Handler[Entity, Entity]) *TypedResourceBuilder[Entity] {
	rb.handle(http.MethodPut, rb.pattern+"/{id}", Handle(h))
	return rb
}

// Patch registers a typed PATCH handler for partial updates.
func (rb *TypedResourceBuilder[Entity]) Patch(h Handler[Entity, Entity]) *TypedResourceBuilder[Entity] {
	rb.handle(http.MethodPatch, rb.pattern+"/{id}", Handle(h))
	return rb
}

// Delete registers a typed DELETE handler for deleting a resource.
// Handler signature: func(ctx, IDRequest) error
func (rb *TypedResourceBuilder[Entity]) Delete(h NoResponseHandler[IDRequest]) *TypedResourceBuilder[Entity] {
	rb.handle(http.MethodDelete, rb.pattern+"/{id}", HandleNoResponse(h))
	return rb
}

// Custom registers a handler with a custom method and path suffix.
func (rb *TypedResourceBuilder[Entity]) Custom(method, suffix string, handler http.HandlerFunc) *TypedResourceBuilder[Entity] {
	rb.handle(method, rb.pattern+suffix, handler)
	return rb
}

// TypedResourceForGroup creates a typed resource builder within a group.
func TypedResourceForGroup[Entity any](g *Group, pattern string, mw ...any) *TypedResourceBuilder[Entity] {
	converted := toMiddleware(mw)
	allMW := make([]Middleware, 0, len(g.allMiddleware())+len(converted))
	allMW = append(allMW, g.allMiddleware()...)
	allMW = append(allMW, converted...)

	return &TypedResourceBuilder[Entity]{
		server:     g.server,
		group:      g,
		pattern:    pattern,
		middleware: allMW,
	}
}
