package middleware_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/kolosys/helix/middleware"
)

func TestChain(t *testing.T) {
	var order []string

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}

	chain := Chain(mw1, mw2)

	handler := chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("expected order[%d] = %s, got %s", i, v, order[i])
		}
	}
}

func TestResponseWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	// Test WriteHeader
	rw.WriteHeader(http.StatusCreated)
	if rw.Status() != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rw.Status())
	}

	// Test Write
	n, err := rw.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 bytes, got %d", n)
	}
	if rw.Size() != 5 {
		t.Errorf("expected size 5, got %d", rw.Size())
	}

	// Test duplicate WriteHeader
	rw.WriteHeader(http.StatusNotFound)
	if rw.Status() != http.StatusCreated {
		t.Errorf("status should not change after first write, got %d", rw.Status())
	}
}

func TestResponseWriterDefaultStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	// Write without calling WriteHeader
	rw.Write([]byte("test"))

	if rw.Status() != http.StatusOK {
		t.Errorf("expected default status 200, got %d", rw.Status())
	}
}

func TestRecover(t *testing.T) {
	output := &bytes.Buffer{}
	mw := RecoverWithConfig(RecoverConfig{
		PrintStack: false,
		Output:     output,
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Should not panic
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

	if !strings.Contains(output.String(), "test panic") {
		t.Errorf("expected output to contain panic message, got '%s'", output.String())
	}
}

func TestRecoverWithStack(t *testing.T) {
	output := &bytes.Buffer{}
	mw := RecoverWithConfig(RecoverConfig{
		PrintStack: true,
		StackSize:  4096,
		Output:     output,
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic with stack")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	out := output.String()
	if !strings.Contains(out, "goroutine") {
		t.Errorf("expected stack trace, got '%s'", out)
	}
}

func TestRecoverWithCustomHandler(t *testing.T) {
	customCalled := false
	mw := RecoverWithConfig(RecoverConfig{
		Handler: func(w http.ResponseWriter, r *http.Request, err any) {
			customCalled = true
			w.WriteHeader(http.StatusServiceUnavailable)
		},
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !customCalled {
		t.Error("custom handler was not called")
	}
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rec.Code)
	}
}

func TestRecoverDefault(t *testing.T) {
	mw := Recover()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestRequestID(t *testing.T) {
	mw := RequestID()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestIDFromRequest(r)
		if id == "" {
			t.Error("expected request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	respID := rec.Header().Get(RequestIDHeader)
	if respID == "" {
		t.Error("expected X-Request-ID in response header")
	}
	if len(respID) != 32 {
		t.Errorf("expected 32 char ID, got %d chars", len(respID))
	}
}

func TestRequestIDPropagation(t *testing.T) {
	mw := RequestID()

	existingID := "existing-request-id-12345"

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestIDFromRequest(r)
		if id != existingID {
			t.Errorf("expected propagated ID '%s', got '%s'", existingID, id)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(RequestIDHeader, existingID)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	respID := rec.Header().Get(RequestIDHeader)
	if respID != existingID {
		t.Errorf("expected response ID '%s', got '%s'", existingID, respID)
	}
}

func TestRequestIDCustomGenerator(t *testing.T) {
	customID := "custom-id-123"
	mw := RequestIDWithConfig(RequestIDConfig{
		Generator: func() string { return customID },
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	respID := rec.Header().Get(RequestIDHeader)
	if respID != customID {
		t.Errorf("expected '%s', got '%s'", customID, respID)
	}
}

func TestCORS(t *testing.T) {
	mw := CORS()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" {
		t.Errorf("expected '*', got '%s'", allowOrigin)
	}
}

func TestCORSPreflight(t *testing.T) {
	mw := CORS()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for preflight")
	}))

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}

	allowMethods := rec.Header().Get("Access-Control-Allow-Methods")
	if allowMethods == "" {
		t.Error("expected Access-Control-Allow-Methods header")
	}
}

func TestCORSWithConfig(t *testing.T) {
	mw := CORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"http://allowed.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"X-Custom"},
		AllowCredentials: true,
		MaxAge:           3600,
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test allowed origin
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://allowed.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "http://allowed.com" {
		t.Errorf("expected 'http://allowed.com', got '%s'", allowOrigin)
	}

	credentials := rec.Header().Get("Access-Control-Allow-Credentials")
	if credentials != "true" {
		t.Errorf("expected 'true', got '%s'", credentials)
	}

	// Test disallowed origin
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://notallowed.com")
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	allowOrigin = rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "" {
		t.Errorf("expected no CORS header, got '%s'", allowOrigin)
	}
}

func TestCORSNoOrigin(t *testing.T) {
	mw := CORS()

	handlerCalled := false
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Error("handler should be called for non-CORS request")
	}
}

func TestCORSAllowAll(t *testing.T) {
	mw := CORSAllowAll()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://any.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" {
		t.Errorf("expected '*', got '%s'", allowOrigin)
	}
}

func TestTimeout(t *testing.T) {
	mw := Timeout(100 * time.Millisecond)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rec.Code)
	}
}

func TestTimeoutNoTimeout(t *testing.T) {
	mw := Timeout(1 * time.Second)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestTimeoutSkip(t *testing.T) {
	mw := TimeoutWithConfig(TimeoutConfig{
		Timeout: 10 * time.Millisecond,
		SkipFunc: func(r *http.Request) bool {
			return r.URL.Path == "/skip"
		},
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	// Should skip timeout
	req := httptest.NewRequest(http.MethodGet, "/skip", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 (skipped), got %d", rec.Code)
	}
}

func TestCompress(t *testing.T) {
	mw := Compress()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Write enough data to trigger compression
		data := strings.Repeat(`{"key":"value"}`, 200)
		w.Write([]byte(data))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	encoding := rec.Header().Get("Content-Encoding")
	if encoding != "gzip" {
		t.Errorf("expected gzip encoding, got '%s'", encoding)
	}

	// Verify we can decompress
	reader, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read gzip body: %v", err)
	}

	if !strings.Contains(string(body), "key") {
		t.Errorf("expected decompressed body to contain 'key'")
	}
}

func TestCompressNoAcceptEncoding(t *testing.T) {
	mw := Compress()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key":"value"}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	encoding := rec.Header().Get("Content-Encoding")
	if encoding != "" {
		t.Errorf("expected no encoding, got '%s'", encoding)
	}
}

func TestCompressSmallResponse(t *testing.T) {
	mw := CompressWithConfig(CompressConfig{
		MinSize: 1024,
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"small":"data"}`)) // Less than MinSize
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Small responses should not be compressed
	encoding := rec.Header().Get("Content-Encoding")
	if encoding == "gzip" {
		t.Error("small response should not be compressed")
	}
}

func TestRateLimit(t *testing.T) {
	mw := RateLimit(2, 2) // 2 requests per second, burst of 2

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First 2 requests should succeed
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}

	// Third request should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", rec.Code)
	}

	// Check rate limit headers
	limit := rec.Header().Get("X-RateLimit-Limit")
	if limit == "" {
		t.Error("expected X-RateLimit-Limit header")
	}
}

func TestRateLimitDifferentClients(t *testing.T) {
	mw := RateLimit(1, 1) // 1 request per second, burst of 1

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First client uses their quota
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("client 1: expected status 200, got %d", rec.Code)
	}

	// Second client should have their own quota
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("client 2: expected status 200, got %d", rec.Code)
	}
}

func TestRateLimitSkip(t *testing.T) {
	mw := RateLimitWithConfig(RateLimitConfig{
		Rate:  1,
		Burst: 1,
		SkipFunc: func(r *http.Request) bool {
			return r.URL.Path == "/health"
		},
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Make many requests to /health - all should succeed
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}
}

func TestBasicAuth(t *testing.T) {
	mw := BasicAuth("admin", "secret")

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	}))

	// Without credentials
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}

	authHeader := rec.Header().Get("WWW-Authenticate")
	if !strings.Contains(authHeader, "Basic") {
		t.Errorf("expected WWW-Authenticate header, got '%s'", authHeader)
	}

	// With correct credentials
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("admin", "secret")
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	// With incorrect credentials
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("admin", "wrong")
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestBasicAuthUsers(t *testing.T) {
	users := map[string]string{
		"user1": "pass1",
		"user2": "pass2",
	}
	mw := BasicAuthUsers(users)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		user     string
		pass     string
		expected int
	}{
		{"user1", "pass1", 200},
		{"user2", "pass2", 200},
		{"user1", "wrong", 401},
		{"unknown", "pass", 401},
	}

	for _, tc := range tests {
		t.Run(tc.user, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.SetBasicAuth(tc.user, tc.pass)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.expected {
				t.Errorf("expected status %d, got %d", tc.expected, rec.Code)
			}
		})
	}
}

func TestBasicAuthSkip(t *testing.T) {
	mw := BasicAuthWithConfig(BasicAuthConfig{
		Validator: func(u, p string) bool { return u == "admin" && p == "secret" },
		Realm:     "Test",
		SkipFunc: func(r *http.Request) bool {
			return r.URL.Path == "/public"
		},
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Public path - should skip auth
	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestETag(t *testing.T) {
	mw := ETag()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":"test"}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	etag := rec.Header().Get("ETag")
	if etag == "" {
		t.Error("expected ETag header")
	}
	if !strings.HasPrefix(etag, `"`) || !strings.HasSuffix(etag, `"`) {
		t.Errorf("ETag should be quoted, got '%s'", etag)
	}
}

func TestETagNotModified(t *testing.T) {
	mw := ETag()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":"test"}`))
	}))

	// First request - get ETag
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	etag := rec.Header().Get("ETag")

	// Second request with If-None-Match
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("If-None-Match", etag)
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotModified {
		t.Errorf("expected status 304, got %d", rec.Code)
	}

	if rec.Body.Len() != 0 {
		t.Error("expected empty body for 304")
	}
}

func TestETagWeak(t *testing.T) {
	mw := ETagWeak()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":"test"}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	etag := rec.Header().Get("ETag")
	if !strings.HasPrefix(etag, "W/") {
		t.Errorf("expected weak ETag, got '%s'", etag)
	}
}

func TestETagSkipNonGET(t *testing.T) {
	mw := ETag()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":"test"}`))
	}))

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	etag := rec.Header().Get("ETag")
	if etag != "" {
		t.Errorf("expected no ETag for POST, got '%s'", etag)
	}
}

func TestETagHelpers(t *testing.T) {
	// Test ETagFromContent
	etag := ETagFromContent([]byte("test"), false)
	if etag == "" {
		t.Error("expected ETag")
	}

	// Test weak ETag
	weakEtag := ETagFromString("test", true)
	if !strings.HasPrefix(weakEtag, "W/") {
		t.Errorf("expected weak ETag, got '%s'", weakEtag)
	}

	// Test version ETag
	versionEtag := ETagFromVersion(123, false)
	if versionEtag != `"123"` {
		t.Errorf("expected '\"123\"', got '%s'", versionEtag)
	}
}

func TestCache(t *testing.T) {
	mw := Cache(3600)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	cacheControl := rec.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "max-age=3600") {
		t.Errorf("expected max-age=3600, got '%s'", cacheControl)
	}
}

func TestCachePublic(t *testing.T) {
	mw := CachePublic(3600)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	cacheControl := rec.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "public") {
		t.Errorf("expected public, got '%s'", cacheControl)
	}
}

func TestCachePrivate(t *testing.T) {
	mw := CachePrivate(600)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	cacheControl := rec.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "private") {
		t.Errorf("expected private, got '%s'", cacheControl)
	}
}

func TestCacheImmutable(t *testing.T) {
	mw := CacheImmutable(31536000)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	cacheControl := rec.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "immutable") {
		t.Errorf("expected immutable, got '%s'", cacheControl)
	}
}

func TestNoCache(t *testing.T) {
	mw := NoCache()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	cacheControl := rec.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "no-cache") {
		t.Errorf("expected no-cache, got '%s'", cacheControl)
	}
	if !strings.Contains(cacheControl, "no-store") {
		t.Errorf("expected no-store, got '%s'", cacheControl)
	}
	if !strings.Contains(cacheControl, "must-revalidate") {
		t.Errorf("expected must-revalidate, got '%s'", cacheControl)
	}
}

func TestCacheSkipNonGET(t *testing.T) {
	mw := Cache(3600)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	cacheControl := rec.Header().Get("Cache-Control")
	if cacheControl != "" {
		t.Errorf("expected no Cache-Control for POST, got '%s'", cacheControl)
	}
}

func TestCacheWithVary(t *testing.T) {
	mw := CacheWithConfig(CacheConfig{
		MaxAge:      3600,
		VaryHeaders: []string{"Accept", "Accept-Encoding"},
	})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	vary := rec.Header().Get("Vary")
	if !strings.Contains(vary, "Accept") {
		t.Errorf("expected Vary to contain Accept, got '%s'", vary)
	}
}

func TestCacheHelpers(t *testing.T) {
	rec := httptest.NewRecorder()

	SetCacheControl(rec, "max-age=3600")
	if rec.Header().Get("Cache-Control") != "max-age=3600" {
		t.Error("SetCacheControl failed")
	}

	now := time.Now()
	SetExpires(rec, now)
	if rec.Header().Get("Expires") == "" {
		t.Error("SetExpires failed")
	}

	SetLastModified(rec, now)
	if rec.Header().Get("Last-Modified") == "" {
		t.Error("SetLastModified failed")
	}
}

func BenchmarkRecoverMiddleware(b *testing.B) {
	mw := Recover()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkRequestIDMiddleware(b *testing.B) {
	mw := RequestID()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkCORSMiddleware(b *testing.B) {
	mw := CORS()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkRateLimitMiddleware(b *testing.B) {
	mw := RateLimit(1000000, 1000000) // High limits to avoid throttling
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkBasicAuthMiddleware(b *testing.B) {
	mw := BasicAuth("admin", "secret")
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("admin", "secret")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkCacheMiddleware(b *testing.B) {
	mw := Cache(3600)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
