package middleware

import "net/http"

// NewResponseWriter exports newResponseWriter for testing.
func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return newResponseWriter(w)
}
