package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Problem represents an RFC 7807 Problem Details for HTTP APIs.
// See: https://tools.ietf.org/html/rfc7807
type Problem struct {
	// Type is a URI reference that identifies the problem type.
	Type string `json:"type"`

	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title"`

	// Status is the HTTP status code for this problem.
	Status int `json:"status"`

	// Detail is a human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// Instance is a URI reference that identifies the specific occurrence of the problem.
	Instance string `json:"instance,omitempty"`

	// Err is the error that caused the problem.
	Err error `json:"-"`
}

// NewProblem creates a new Problem with the given status, type, and title.
func NewProblem(status int, problemType, title string) Problem {
	return Problem{
		Type:   "about:blank#" + problemType,
		Title:  title,
		Status: status,
	}
}

// Error implements the error interface.
func (p Problem) Error() string {
	if p.Detail != "" {
		return fmt.Sprintf("%s: %s", p.Title, p.Detail)
	}
	return p.Title
}

func (p Problem) WithDetail(detail string) Problem {
	newProblem := p
	newProblem.Detail = detail
	return newProblem
}

// WithDetailf returns a copy of the Problem with the given detail message.
func (p Problem) WithDetailf(format string, args ...any) Problem {
	newProblem := p
	if len(args) > 0 {
		newProblem.Detail = fmt.Sprintf(format, args...)
	} else {
		newProblem.Detail = format
	}
	return newProblem
}

// WithInstance returns a copy of the Problem with the given instance URI.
func (p Problem) WithInstance(instance string) Problem {
	newProblem := p
	newProblem.Instance = instance
	return newProblem
}

// WithType returns a copy of the Problem with the given type URI.
func (p Problem) WithType(problemType string) Problem {
	newProblem := p
	newProblem.Type = problemType
	return newProblem
}

// WithStack returns a copy of the Problem with the given stack trace.
func (p Problem) WithErr(err error) Problem {
	newProblem := p
	newProblem.Err = err
	return newProblem
}

// Sentinel errors for common HTTP error responses.
var (
	// ErrBadRequest represents a 400 Bad Request error.
	ErrBadRequest = NewProblem(http.StatusBadRequest, "bad_request", "Bad Request")

	// ErrUnauthorized represents a 401 Unauthorized error.
	ErrUnauthorized = NewProblem(http.StatusUnauthorized, "unauthorized", "Unauthorized")

	// ErrForbidden represents a 403 Forbidden error.
	ErrForbidden = NewProblem(http.StatusForbidden, "forbidden", "Forbidden")

	// ErrNotFound represents a 404 Not Found error.
	ErrNotFound = NewProblem(http.StatusNotFound, "not_found", "Not Found")

	// ErrMethodNotAllowed represents a 405 Method Not Allowed error.
	ErrMethodNotAllowed = NewProblem(http.StatusMethodNotAllowed, "method_not_allowed", "Method Not Allowed")

	// ErrConflict represents a 409 Conflict error.
	ErrConflict = NewProblem(http.StatusConflict, "conflict", "Conflict")

	// ErrGone represents a 410 Gone error.
	ErrGone = NewProblem(http.StatusGone, "gone", "Gone")

	// ErrUnprocessableEntity represents a 422 Unprocessable Entity error.
	ErrUnprocessableEntity = NewProblem(http.StatusUnprocessableEntity, "unprocessable_entity", "Unprocessable Entity")

	// ErrTooManyRequests represents a 429 Too Many Requests error.
	ErrTooManyRequests = NewProblem(http.StatusTooManyRequests, "too_many_requests", "Too Many Requests")

	// ErrInternal represents a 500 Internal Server Error.
	ErrInternal = NewProblem(http.StatusInternalServerError, "internal_error", "Internal Server Error")

	// ErrNotImplemented represents a 501 Not Implemented error.
	ErrNotImplemented = NewProblem(http.StatusNotImplemented, "not_implemented", "Not Implemented")

	// ErrBadGateway represents a 502 Bad Gateway error.
	ErrBadGateway = NewProblem(http.StatusBadGateway, "bad_gateway", "Bad Gateway")

	// ErrServiceUnavailable represents a 503 Service Unavailable error.
	ErrServiceUnavailable = NewProblem(http.StatusServiceUnavailable, "service_unavailable", "Service Unavailable")

	// ErrGatewayTimeout represents a 504 Gateway Timeout error.
	ErrGatewayTimeout = NewProblem(http.StatusGatewayTimeout, "gateway_timeout", "Gateway Timeout")
)

// ProblemFromStatus creates a Problem from an HTTP status code.
func ProblemFromStatus(status int) Problem {
	return NewProblem(status, http.StatusText(status), http.StatusText(status))
}

// WriteProblem writes a Problem response to the http.ResponseWriter.
func WriteProblem(w http.ResponseWriter, p Problem) error {
	w.Header().Set("Content-Type", MIMEApplicationProblemJSON)
	w.WriteHeader(p.Status)
	return jsonEncode(w, p)
}

// jsonEncode encodes value to JSON without modifying Content-Type.
func jsonEncode(w http.ResponseWriter, v any) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// -----------------------------------------------------------------------------
// Convenience Error Functions
// -----------------------------------------------------------------------------

// BadRequestf creates a 400 Bad Request Problem with a formatted detail message.
func BadRequestf(format string, args ...any) Problem {
	return ErrBadRequest.WithDetailf(format, args...)
}

// Unauthorizedf creates a 401 Unauthorized Problem with a formatted detail message.
func Unauthorizedf(format string, args ...any) Problem {
	return ErrUnauthorized.WithDetailf(format, args...)
}

// Forbiddenf creates a 403 Forbidden Problem with a formatted detail message.
func Forbiddenf(format string, args ...any) Problem {
	return ErrForbidden.WithDetailf(format, args...)
}

// NotFoundf creates a 404 Not Found Problem with a formatted detail message.
func NotFoundf(format string, args ...any) Problem {
	return ErrNotFound.WithDetailf(format, args...)
}

// MethodNotAllowedf creates a 405 Method Not Allowed Problem with a formatted detail message.
func MethodNotAllowedf(format string, args ...any) Problem {
	return ErrMethodNotAllowed.WithDetailf(format, args...)
}

// Conflictf creates a 409 Conflict Problem with a formatted detail message.
func Conflictf(format string, args ...any) Problem {
	return ErrConflict.WithDetailf(format, args...)
}

// Gonef creates a 410 Gone Problem with a formatted detail message.
func Gonef(format string, args ...any) Problem {
	return ErrGone.WithDetailf(format, args...)
}

// UnprocessableEntityf creates a 422 Unprocessable Entity Problem with a formatted detail message.
func UnprocessableEntityf(format string, args ...any) Problem {
	return ErrUnprocessableEntity.WithDetailf(format, args...)
}

// TooManyRequestsf creates a 429 Too Many Requests Problem with a formatted detail message.
func TooManyRequestsf(format string, args ...any) Problem {
	return ErrTooManyRequests.WithDetailf(format, args...)
}

// Internalf creates a 500 Internal Server Error Problem with a formatted detail message.
func Internalf(format string, args ...any) Problem {
	return ErrInternal.WithDetailf(format, args...)
}

// NotImplementedf creates a 501 Not Implemented Problem with a formatted detail message.
func NotImplementedf(format string, args ...any) Problem {
	return ErrNotImplemented.WithDetailf(format, args...)
}

// BadGatewayf creates a 502 Bad Gateway Problem with a formatted detail message.
func BadGatewayf(format string, args ...any) Problem {
	return ErrBadGateway.WithDetailf(format, args...)
}

// ServiceUnavailablef creates a 503 Service Unavailable Problem with a formatted detail message.
func ServiceUnavailablef(format string, args ...any) Problem {
	return ErrServiceUnavailable.WithDetailf(format, args...)
}

// GatewayTimeoutf creates a 504 Gateway Timeout Problem with a formatted detail message.
func GatewayTimeoutf(format string, args ...any) Problem {
	return ErrGatewayTimeout.WithDetailf(format, args...)
}

// -----------------------------------------------------------------------------
// Validation Errors
// -----------------------------------------------------------------------------

// FieldError represents a validation error for a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors collects multiple validation errors for RFC 7807 response.
// Implements the error interface and can be returned from Validate() methods.
type ValidationErrors struct {
	errors []FieldError
}

// NewValidationErrors creates a new empty ValidationErrors collector.
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		errors: make([]FieldError, 0),
	}
}

// Add adds a validation error for a specific field.
func (v *ValidationErrors) Add(field, message string) {
	v.errors = append(v.errors, FieldError{Field: field, Message: message})
}

// Addf adds a validation error for a specific field with a formatted message.
func (v *ValidationErrors) Addf(field, format string, args ...any) {
	v.errors = append(v.errors, FieldError{
		Field:   field,
		Message: fmt.Sprintf(format, args...),
	})
}

// HasErrors returns true if there are any validation errors.
func (v *ValidationErrors) HasErrors() bool {
	return len(v.errors) > 0
}

// Len returns the number of validation errors.
func (v *ValidationErrors) Len() int {
	return len(v.errors)
}

// Errors returns the list of field errors.
func (v *ValidationErrors) Errors() []FieldError {
	return v.errors
}

// Error implements the error interface.
func (v *ValidationErrors) Error() string {
	if len(v.errors) == 0 {
		return "validation failed"
	}
	if len(v.errors) == 1 {
		return fmt.Sprintf("%s: %s", v.errors[0].Field, v.errors[0].Message)
	}
	return fmt.Sprintf("validation failed: %d errors", len(v.errors))
}

// Err returns nil if there are no errors, otherwise returns the ValidationErrors.
// This is useful for the common pattern: return v.Err()
func (v *ValidationErrors) Err() error {
	if len(v.errors) == 0 {
		return nil
	}
	return v
}

// ValidationProblem is an RFC 7807 Problem with validation errors extension.
type ValidationProblem struct {
	Problem
	Errors []FieldError `json:"errors,omitempty"`
}

// ToProblem converts ValidationErrors to a ValidationProblem for RFC 7807 response.
func (v *ValidationErrors) ToProblem() ValidationProblem {
	return ValidationProblem{
		Problem: ErrUnprocessableEntity.WithDetail("One or more validation errors occurred"),
		Errors:  v.errors,
	}
}

// WriteValidationProblem writes a ValidationProblem response to the http.ResponseWriter.
func WriteValidationProblem(w http.ResponseWriter, v *ValidationErrors) error {
	p := v.ToProblem()
	w.Header().Set("Content-Type", MIMEApplicationProblemJSON)
	w.WriteHeader(p.Status)
	return jsonEncode(w, p)
}
