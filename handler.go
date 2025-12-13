package helix

import (
	"context"
	"net/http"
)

// ErrorHandler is a function that handles errors from handlers.
// It receives the response writer, request, and error, and is responsible
// for writing an appropriate error response.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// errorHandlerKey is the context key for storing the error handler.
type errorHandlerKey struct{}

// withErrorHandler stores the error handler in the request context.
func withErrorHandler(r *http.Request, handler ErrorHandler) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), errorHandlerKey{}, handler))
}

// getErrorHandler retrieves the error handler from the request context.
func getErrorHandler(r *http.Request) (ErrorHandler, bool) {
	handler, ok := r.Context().Value(errorHandlerKey{}).(ErrorHandler)
	return handler, ok
}

// Handler is a generic handler function that accepts a typed request and returns a typed response.
// The request type is automatically bound from path parameters, query parameters, headers, and JSON body.
// The response is automatically encoded as JSON.
type Handler[Req, Res any] func(ctx context.Context, req Req) (Res, error)

// Handle wraps a generic Handler into an http.HandlerFunc.
// It automatically:
//   - Binds the request to the Req type
//   - Calls the handler with the context and request
//   - Encodes the response as JSON
//   - Handles errors using RFC 7807 Problem Details
func Handle[Req, Res any](h Handler[Req, Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Bind request
		req, err := Bind[Req](r)
		if err != nil {
			handleError(w, r, err)
			return
		}

		// Check if request is validatable
		if v, ok := any(&req).(Validatable); ok {
			if err := v.Validate(); err != nil {
				handleError(w, r, err)
				return
			}
		}

		// Call handler
		res, err := h(r.Context(), req)
		if err != nil {
			handleError(w, r, err)
			return
		}

		// Encode response
		if err := JSON(w, http.StatusOK, res); err != nil {
			handleError(w, r, err)
			return
		}
	}
}

// HandleWithStatus wraps a generic Handler into an http.HandlerFunc with a custom success status code.
func HandleWithStatus[Req, Res any](status int, h Handler[Req, Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Bind request
		req, err := Bind[Req](r)
		if err != nil {
			handleError(w, r, err)
			return
		}

		// Check if request is validatable
		if v, ok := any(&req).(Validatable); ok {
			if err := v.Validate(); err != nil {
				handleError(w, r, err)
				return
			}
		}

		// Call handler
		res, err := h(r.Context(), req)
		if err != nil {
			handleError(w, r, err)
			return
		}

		// Encode response
		if err := JSON(w, status, res); err != nil {
			handleError(w, r, err)
			return
		}
	}
}

// HandleCreated wraps a generic Handler into an http.HandlerFunc that returns 201 Created.
// This is a convenience wrapper for HandleWithStatus(http.StatusCreated, h).
func HandleCreated[Req, Res any](h Handler[Req, Res]) http.HandlerFunc {
	return HandleWithStatus(http.StatusCreated, h)
}

// HandleAccepted wraps a generic Handler into an http.HandlerFunc that returns 202 Accepted.
// Useful for async operations where processing happens in the background.
func HandleAccepted[Req, Res any](h Handler[Req, Res]) http.HandlerFunc {
	return HandleWithStatus(http.StatusAccepted, h)
}

// NoRequestHandler is a handler that takes no request body, only context.
type NoRequestHandler[Res any] func(ctx context.Context) (Res, error)

// HandleNoRequest wraps a NoRequestHandler into an http.HandlerFunc.
// Useful for endpoints that don't need request binding (e.g., GET /users).
func HandleNoRequest[Res any](h NoRequestHandler[Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := h(r.Context())
		if err != nil {
			handleError(w, r, err)
			return
		}

		if err := JSON(w, http.StatusOK, res); err != nil {
			handleError(w, r, err)
			return
		}
	}
}

// NoResponseHandler is a handler that returns no response body.
type NoResponseHandler[Req any] func(ctx context.Context, req Req) error

// HandleNoResponse wraps a NoResponseHandler into an http.HandlerFunc.
// Returns 204 No Content on success.
func HandleNoResponse[Req any](h NoResponseHandler[Req]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := Bind[Req](r)
		if err != nil {
			handleError(w, r, err)
			return
		}

		if v, ok := any(&req).(Validatable); ok {
			if err := v.Validate(); err != nil {
				handleError(w, r, err)
				return
			}
		}

		if err := h(r.Context(), req); err != nil {
			handleError(w, r, err)
			return
		}

		NoContent(w)
	}
}

// EmptyHandler is a handler that takes no request and returns no response.
type EmptyHandler func(ctx context.Context) error

// HandleEmpty wraps an EmptyHandler into an http.HandlerFunc.
// Returns 204 No Content on success.
func HandleEmpty(h EmptyHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(r.Context()); err != nil {
			handleError(w, r, err)
			return
		}

		NoContent(w)
	}
}

// handleError handles errors from handlers.
// If a custom error handler is set in the request context, it is used.
// Otherwise, the default error handling is used:
//   - If the error is a Problem, it is encoded as RFC 7807.
//   - If the error is ValidationErrors, it is encoded with field-level errors.
//   - Otherwise, a generic 500 Internal Server Error is returned.
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	// Check for custom error handler in context
	if handler, ok := getErrorHandler(r); ok {
		handler(w, r, err)
		return
	}

	// Default error handling
	HandleErrorDefault(w, r, err)
}

// HandleErrorDefault provides the default error handling logic.
// This can be called from custom error handlers to fall back to default behavior.
func HandleErrorDefault(w http.ResponseWriter, r *http.Request, err error) {
	// Check if it's a ValidationErrors
	if verrs, ok := err.(*ValidationErrors); ok {
		p := verrs.ToProblem()
		p.Instance = r.URL.RequestURI()
		w.Header().Set("Content-Type", MIMEApplicationProblemJSON)
		w.WriteHeader(p.Status)
		jsonEncode(w, p)
		return
	}

	// Check if it's a Problem error
	if problem, ok := err.(Problem); ok {
		// Set the instance to the request URI if not set
		if problem.Instance == "" {
			problem.Instance = r.URL.RequestURI()
		}
		WriteProblem(w, problem)
		return
	}

	// Check for binding errors
	if isBindingError(err) {
		problem := ErrBadRequest.WithErr(err)
		problem.Instance = r.URL.RequestURI()
		WriteProblem(w, problem)
		return
	}

	// Default to internal server error
	problem := ErrInternal.WithErr(err)
	problem.Instance = r.URL.RequestURI()
	WriteProblem(w, problem)
}

// isBindingError checks if an error is a binding error.
func isBindingError(err error) bool {
	switch err {
	case ErrBindingFailed, ErrUnsupportedType, ErrInvalidJSON, ErrRequiredField, ErrInvalidFieldValue:
		return true
	}

	// Check if error message contains binding error prefix
	errStr := err.Error()
	return len(errStr) > 6 && errStr[:6] == "helix:"
}
