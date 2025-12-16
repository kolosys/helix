package helix

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

// Buffer pool for JSON encoding to reduce allocations.
var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

// JSON writes a JSON response with the given status code.
// Uses pooled buffer for zero-allocation in the hot path.
func JSON(w http.ResponseWriter, status int, v any) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return err
	}

	w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	return err
}

// JSONPretty writes a pretty-printed JSON response with the given status code.
func JSONPretty(w http.ResponseWriter, status int, v any, indent string) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", indent)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return err
	}

	w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	return err
}

// Text writes a plain text response with the given status code.
func Text(w http.ResponseWriter, status int, text string) error {
	w.Header().Set("Content-Type", MIMETextPlainCharsetUTF8)
	w.WriteHeader(status)
	_, err := io.WriteString(w, text)
	return err
}

// HTML writes an HTML response with the given status code.
func HTML(w http.ResponseWriter, status int, html string) error {
	w.Header().Set("Content-Type", MIMETextHTMLCharsetUTF8)
	w.WriteHeader(status)
	_, err := io.WriteString(w, html)
	return err
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Redirect redirects the request to the given URL.
func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

// Stream streams the content from the reader to the response.
func Stream(w http.ResponseWriter, contentType string, reader io.Reader) error {
	w.Header().Set("Content-Type", contentType)
	_, err := io.Copy(w, reader)
	return err
}

// Blob writes binary data with the given content type.
func Blob(w http.ResponseWriter, status int, contentType string, data []byte) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, err := w.Write(data)
	return err
}

// File serves a file with the given content type.
func File(w http.ResponseWriter, r *http.Request, path string) {
	http.ServeFile(w, r, path)
}

// Attachment sets the Content-Disposition header to attachment with the given filename.
func Attachment(w http.ResponseWriter, filename string) {
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
}

// Inline sets the Content-Disposition header to inline with the given filename.
func Inline(w http.ResponseWriter, filename string) {
	w.Header().Set("Content-Disposition", "inline; filename=\""+filename+"\"")
}

// Error writes an error response with the given status code and message.
func Error(w http.ResponseWriter, status int, message string) error {
	return JSON(w, status, map[string]string{"error": message})
}

// OK writes a 200 OK JSON response.
func OK(w http.ResponseWriter, v any) error {
	return JSON(w, http.StatusOK, v)
}

// Created writes a 201 Created JSON response.
func Created(w http.ResponseWriter, v any) error {
	return JSON(w, http.StatusCreated, v)
}

// Accepted writes a 202 Accepted JSON response.
func Accepted(w http.ResponseWriter, v any) error {
	return JSON(w, http.StatusAccepted, v)
}

// BadRequest writes a 400 Bad Request error response.
func BadRequest(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusBadRequest, message)
}

// Unauthorized writes a 401 Unauthorized error response.
func Unauthorized(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusUnauthorized, message)
}

// Forbidden writes a 403 Forbidden error response.
func Forbidden(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusForbidden, message)
}

// NotFound writes a 404 Not Found error response.
func NotFound(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusNotFound, message)
}

// InternalServerError writes a 500 Internal Server Error response.
func InternalServerError(w http.ResponseWriter, message string) error {
	return Error(w, http.StatusInternalServerError, message)
}
