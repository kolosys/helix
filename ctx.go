package helix

import (
	"context"
	"encoding/json"
	"net/http"
)

// Ctx provides a unified context for HTTP handlers with fluent accessors
// for request data and response methods.
type Ctx struct {
	Request  *http.Request
	Response http.ResponseWriter

	// status holds the pending status code for chained responses
	status int

	// store holds request-scoped values for dependency injection
	store map[string]any
}

// NewCtx creates a new Ctx from an http.Request and http.ResponseWriter.
func NewCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		Request:  r,
		Response: w,
	}
}

// Reset resets the Ctx for reuse from a pool.
func (c *Ctx) Reset(w http.ResponseWriter, r *http.Request) {
	c.Request = r
	c.Response = w
	c.status = 0
	c.store = nil
}

// Context returns the request's context.Context.
func (c *Ctx) Context() context.Context {
	return c.Request.Context()
}

// -----------------------------------------------------------------------------
// Request-Scoped Storage (Dependency Injection)
// -----------------------------------------------------------------------------

// Set stores a value in the request-scoped store.
func (c *Ctx) Set(key string, value any) {
	if c.store == nil {
		c.store = make(map[string]any)
	}
	c.store[key] = value
}

// Get retrieves a value from the request-scoped store.
func (c *Ctx) Get(key string) (any, bool) {
	if c.store == nil {
		return nil, false
	}
	v, ok := c.store[key]
	return v, ok
}

// MustGet retrieves a value from the request-scoped store or panics if not found.
func (c *Ctx) MustGet(key string) any {
	v, ok := c.Get(key)
	if !ok {
		panic("helix: key not found in context: " + key)
	}
	return v
}

// GetString retrieves a string value from the request-scoped store.
func (c *Ctx) GetString(key string) string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetInt retrieves an int value from the request-scoped store.
func (c *Ctx) GetInt(key string) int {
	if v, ok := c.Get(key); ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

// -----------------------------------------------------------------------------
// Path Parameter Accessors
// -----------------------------------------------------------------------------

// Param returns the value of a path parameter.
func (c *Ctx) Param(name string) string {
	return Param(c.Request, name)
}

// ParamInt returns the value of a path parameter as an int.
func (c *Ctx) ParamInt(name string) (int, error) {
	return ParamInt(c.Request, name)
}

// ParamInt64 returns the value of a path parameter as an int64.
func (c *Ctx) ParamInt64(name string) (int64, error) {
	return ParamInt64(c.Request, name)
}

// ParamUUID returns the value of a path parameter validated as a UUID.
func (c *Ctx) ParamUUID(name string) (string, error) {
	return ParamUUID(c.Request, name)
}

// -----------------------------------------------------------------------------
// Query Parameter Accessors
// -----------------------------------------------------------------------------

// Query returns the first value of a query parameter.
func (c *Ctx) Query(name string) string {
	return Query(c.Request, name)
}

// QueryDefault returns the first value of a query parameter or a default value.
func (c *Ctx) QueryDefault(name, defaultVal string) string {
	return QueryDefault(c.Request, name, defaultVal)
}

// QueryInt returns the first value of a query parameter as an int.
func (c *Ctx) QueryInt(name string, defaultVal int) int {
	return QueryInt(c.Request, name, defaultVal)
}

// QueryInt64 returns the first value of a query parameter as an int64.
func (c *Ctx) QueryInt64(name string, defaultVal int64) int64 {
	return QueryInt64(c.Request, name, defaultVal)
}

// QueryFloat64 returns the first value of a query parameter as a float64.
func (c *Ctx) QueryFloat64(name string, defaultVal float64) float64 {
	return QueryFloat64(c.Request, name, defaultVal)
}

// QueryBool returns the first value of a query parameter as a bool.
func (c *Ctx) QueryBool(name string) bool {
	return QueryBool(c.Request, name)
}

// QuerySlice returns all values of a query parameter as a string slice.
func (c *Ctx) QuerySlice(name string) []string {
	return QuerySlice(c.Request, name)
}

// -----------------------------------------------------------------------------
// Header Accessors
// -----------------------------------------------------------------------------

// Header returns the value of a request header.
func (c *Ctx) Header(name string) string {
	return c.Request.Header.Get(name)
}

// -----------------------------------------------------------------------------
// Request Body Binding
// -----------------------------------------------------------------------------

// Bind binds the request body to the given struct using JSON decoding.
func (c *Ctx) Bind(v any) error {
	if c.Request.Body == nil {
		return ErrInvalidJSON
	}
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// BindJSON is an alias for Bind.
func (c *Ctx) BindJSON(v any) error {
	return c.Bind(v)
}

// -----------------------------------------------------------------------------
// Response Methods (Chainable)
// -----------------------------------------------------------------------------

// SetHeader sets a response header and returns the Ctx for chaining.
func (c *Ctx) SetHeader(key, value string) *Ctx {
	c.Response.Header().Set(key, value)
	return c
}

// AddHeader adds a response header value and returns the Ctx for chaining.
func (c *Ctx) AddHeader(key, value string) *Ctx {
	c.Response.Header().Add(key, value)
	return c
}

// SetCookie sets a cookie on the response and returns the Ctx for chaining.
func (c *Ctx) SetCookie(cookie *http.Cookie) *Ctx {
	http.SetCookie(c.Response, cookie)
	return c
}

// Status sets the pending status code for the response and returns the Ctx for chaining.
// The status is applied when a response body is written.
func (c *Ctx) Status(code int) *Ctx {
	c.status = code
	return c
}

// -----------------------------------------------------------------------------
// Response Writers
// -----------------------------------------------------------------------------

// JSON writes a JSON response with the given status code.
func (c *Ctx) JSON(status int, v any) error {
	return JSON(c.Response, status, v)
}

// OK writes a 200 OK JSON response.
func (c *Ctx) OK(v any) error {
	return OK(c.Response, v)
}

// Created writes a 201 Created JSON response.
func (c *Ctx) Created(v any) error {
	return Created(c.Response, v)
}

// Accepted writes a 202 Accepted JSON response.
func (c *Ctx) Accepted(v any) error {
	return Accepted(c.Response, v)
}

// NoContent writes a 204 No Content response.
func (c *Ctx) NoContent() error {
	return NoContent(c.Response)
}

// Text writes a plain text response with the given status code.
func (c *Ctx) Text(status int, text string) error {
	return Text(c.Response, status, text)
}

// HTML writes an HTML response with the given status code.
func (c *Ctx) HTML(status int, html string) error {
	return HTML(c.Response, status, html)
}

// Blob writes binary data with the given content type.
func (c *Ctx) Blob(status int, contentType string, data []byte) error {
	return Blob(c.Response, status, contentType, data)
}

// Problem writes an RFC 7807 Problem response.
func (c *Ctx) Problem(p Problem) error {
	if p.Instance == "" {
		p.Instance = c.Request.URL.RequestURI()
	}
	return WriteProblem(c.Response, p)
}

// Redirect redirects the request to the given URL.
func (c *Ctx) Redirect(url string, code int) {
	Redirect(c.Response, c.Request, url, code)
}

// File serves a file.
func (c *Ctx) File(path string) {
	File(c.Response, c.Request, path)
}

// Attachment sets the Content-Disposition header to attachment.
func (c *Ctx) Attachment(filename string) *Ctx {
	Attachment(c.Response, filename)
	return c
}

// Inline sets the Content-Disposition header to inline.
func (c *Ctx) Inline(filename string) *Ctx {
	Inline(c.Response, filename)
	return c
}

// -----------------------------------------------------------------------------
// Error Response Helpers
// -----------------------------------------------------------------------------

// BadRequest writes a 400 Bad Request error response.
func (c *Ctx) BadRequest(message string) error {
	return c.Problem(ErrBadRequest.WithDetailf("%s", message))
}

// Unauthorized writes a 401 Unauthorized error response.
func (c *Ctx) Unauthorized(message string) error {
	return c.Problem(ErrUnauthorized.WithDetailf("%s", message))
}

// Forbidden writes a 403 Forbidden error response.
func (c *Ctx) Forbidden(message string) error {
	return c.Problem(ErrForbidden.WithDetailf("%s", message))
}

// NotFound writes a 404 Not Found error response.
func (c *Ctx) NotFound(message string) error {
	return c.Problem(ErrNotFound.WithDetailf("%s", message))
}

// InternalServerError writes a 500 Internal Server Error response.
func (c *Ctx) InternalServerError(message string) error {
	return c.Problem(ErrInternal.WithDetailf("%s", message))
}

// -----------------------------------------------------------------------------
// Convenience Response Helpers
// -----------------------------------------------------------------------------

// SendJSON writes a JSON response with the pending status code (or 200 if not set).
func (c *Ctx) SendJSON(v any) error {
	status := c.status
	if status == 0 {
		status = 200
	}
	return JSON(c.Response, status, v)
}

// OKMessage writes a 200 OK response with a message.
func (c *Ctx) OKMessage(message string) error {
	return c.OK(map[string]string{"message": message})
}

// CreatedMessage writes a 201 Created response with a message and ID.
func (c *Ctx) CreatedMessage(message string, id any) error {
	return c.Created(map[string]any{"message": message, "id": id})
}

// DeletedMessage writes a 200 OK response indicating deletion.
func (c *Ctx) DeletedMessage(message string) error {
	return c.OK(map[string]string{"message": message, "deleted": "true"})
}

// Paginated writes a paginated JSON response.
func (c *Ctx) Paginated(items any, total, page, limit int) error {
	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}

	return c.OK(map[string]any{
		"items":       items,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"has_more":    page < totalPages,
	})
}

// -----------------------------------------------------------------------------
// CtxHandler and HandleCtx
// -----------------------------------------------------------------------------

// CtxHandler is a handler function that uses the unified Ctx type.
type CtxHandler func(c *Ctx) error

// HandleCtx wraps a CtxHandler into an http.HandlerFunc.
// Errors returned from the handler are automatically converted to RFC 7807 responses.
func HandleCtx(h CtxHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := NewCtx(w, r)
		if err := h(c); err != nil {
			handleError(w, r, err)
		}
	}
}
